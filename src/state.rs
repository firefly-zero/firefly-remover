use crate::*;
use alloc::{string::String, vec, vec::Vec};
use core::mem::MaybeUninit;
use firefly_rust::*;
use firefly_types::{Encode, Stats};

static mut STATE: MaybeUninit<State> = MaybeUninit::uninit();

#[derive(Clone, Copy, PartialEq)]
pub enum Kind {
    Rom,
    Data,
    Shots,
    Badges,
    Scores,
}

#[derive(Clone, Copy)]
pub struct Switch {
    pub kind: Kind,
    pub selected: bool,
}

impl Switch {
    fn new(kind: Kind) -> Self {
        Self {
            kind,
            selected: false,
        }
    }
}

pub struct State {
    pub font: FileBuf,
    pub target: Option<(String, String)>,
    pub settings: Settings,
    pub msg: Option<Message>,
    pub switches: Vec<Switch>,
    pub cursor: u8,
    pub input: firefly_ui::InputManager,
    pub fire: Option<Fire>,
}

pub fn get_state() -> &'static mut State {
    #[allow(static_mut_refs)]
    unsafe {
        STATE.assume_init_mut()
    }
}

pub fn load_state() {
    let settings = get_settings(get_me());
    let encoding = settings.language.encoding();
    let font = load_file_buf(encoding).unwrap_or_else(|| load_file_buf("ascii").unwrap());
    let target = load_target();
    let mut msg = None;
    if target.is_none() {
        msg = Some(Message::NoAppSelected);
    }
    let switches = if let Some((author_id, app_id)) = &target {
        detect_switches(author_id, app_id)
    } else {
        Vec::new()
    };
    if switches.is_empty() && msg.is_none() {
        msg = Some(Message::AppAlreadyRemoved);
    }
    let state = State {
        font,
        target,
        settings,
        msg,
        switches,
        cursor: 0,
        input: firefly_ui::InputManager::new(),
        fire: Some(Fire::new()),
    };
    #[allow(static_mut_refs)]
    unsafe {
        STATE.write(state)
    };
}

fn detect_switches(author_id: &str, app_id: &str) -> Vec<Switch> {
    let mut switches = Vec::new();
    let bin_path = alloc::format!("roms/{author_id}/{app_id}/_bin");
    if sudo::get_file_size(&bin_path) != 0 {
        switches.push(Switch::new(Kind::Rom));
    }
    if has_data(author_id, app_id) {
        switches.push(Switch::new(Kind::Data));
    }
    if has_shots(author_id, app_id) {
        switches.push(Switch::new(Kind::Shots));
    }
    let (has_badges, has_scores) = has_stats(author_id, app_id);
    if has_badges {
        switches.push(Switch::new(Kind::Badges));
    }
    if has_scores {
        switches.push(Switch::new(Kind::Scores));
    }
    switches
}

fn has_data(author_id: &str, app_id: &str) -> bool {
    let data_path = alloc::format!("data/{author_id}/{app_id}");
    let stash_path = alloc::format!("{data_path}/stash");
    if sudo::get_file_size(&stash_path) != 0 {
        return true;
    }
    let etc_path = alloc::format!("{data_path}/etc");
    let data_files = sudo::DirBuf::list_files(&etc_path);
    data_files.iter().next().is_some()
}

fn has_shots(author_id: &str, app_id: &str) -> bool {
    let shots_path = alloc::format!("data/{author_id}/{app_id}/shots");
    let files = sudo::DirBuf::list_files(&shots_path);
    files.iter().next().is_some()
}

fn has_stats(author_id: &str, app_id: &str) -> (bool, bool) {
    let path = alloc::format!("data/{author_id}/{app_id}/stats");
    let size = sudo::get_file_size(&path);
    if size == 0 {
        log_error("app has no stats file");
        return (false, false);
    }
    let mut raw = vec![0; size];
    sudo::load_file(&path, &mut raw);
    let Ok(stats) = Stats::decode(&raw) else {
        log_error("app stats cannot be decoded");
        return (false, false);
    };
    let has_badges = stats.badges.iter().any(|b| b.done != 0);
    let has_scores = stats.scores.iter().any(|board| {
        board.me.iter().any(|s| *s != 0) || board.friends.iter().any(|s| s.score != 0)
    });
    (has_badges, has_scores)
}

/// Read the ID of the app to be removed.
fn load_target() -> Option<(String, String)> {
    let raw = load_file_buf("target")?;
    let raw = raw.as_bytes();
    let raw = raw.trim_ascii();
    let raw = alloc::str::from_utf8(raw).ok()?;
    let (author, app) = split_by(raw, '.')?;
    let target = (String::from(author), String::from(&app[1..]));
    Some(target)
}

/// Split the string once at the given character.
fn split_by(input: &str, sep: char) -> Option<(&str, &str)> {
    let mut split_at = None;
    let sep: u8 = sep.try_into().unwrap();
    for (i, ch) in input.bytes().enumerate() {
        if ch == sep {
            split_at = Some(i);
            break;
        }
    }
    let split_at = split_at?;
    Some(input.split_at(split_at))
}
