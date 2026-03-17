use crate::*;
use alloc::{string::String, vec::Vec};
use core::cell::OnceCell;
use firefly_rust::*;

static mut STATE: OnceCell<State> = OnceCell::new();

#[derive(Clone, Copy)]
pub enum Kind {
    Rom,
    Data,
    Shots,
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
}

pub fn get_state() -> &'static mut State {
    #[allow(static_mut_refs)]
    unsafe { STATE.get_mut() }.unwrap()
}

pub fn load_state() {
    let font = load_file_buf("ascii").unwrap();
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
    let settings = get_settings(get_me());
    let state = State {
        font,
        target,
        settings,
        msg,
        switches,
        cursor: 0,
        input: firefly_ui::InputManager::new(),
    };
    #[allow(static_mut_refs)]
    unsafe { STATE.set(state) }.ok().unwrap();
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
        switches.push(Switch::new(Kind::Data));
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
