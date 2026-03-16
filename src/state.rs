use core::cell::OnceCell;
use firefly_rust::*;

static mut STATE: OnceCell<State> = OnceCell::new();

pub struct State {
    pub font: FileBuf,
    pub settings: Settings,
}

pub fn get_state() -> &'static mut State {
    #[allow(static_mut_refs)]
    unsafe { STATE.get_mut() }.unwrap()
}

pub fn load_state() {
    let font = load_file_buf("ascii").unwrap();
    let state = State {
        font,
        settings: get_settings(get_me()),
    };
    #[allow(static_mut_refs)]
    unsafe { STATE.set(state) }.ok().unwrap();
}
