use alloc::string::String;
use core::cell::OnceCell;
use firefly_rust::*;

static mut STATE: OnceCell<State> = OnceCell::new();

pub struct State {
    pub font: FileBuf,
    pub target: Option<(String, String)>,
    pub settings: Settings,
    pub msg: Option<&'static str>,
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
        msg = Some("no app selected");
    }
    let state = State {
        font,
        target,
        settings: get_settings(get_me()),
        msg,
    };
    #[allow(static_mut_refs)]
    unsafe { STATE.set(state) }.ok().unwrap();
}

/// Read the ID of the app to be removed.
fn load_target() -> Option<(String, String)> {
    let raw = load_file_buf("target")?;
    let raw = raw.as_bytes();
    let raw = raw.trim_ascii();
    let raw = alloc::str::from_utf8(raw).ok()?;
    let (author, app) = split_by(raw, '.')?;
    let target = (String::from(author), String::from(app));
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
