#![no_std]
#![no_main]
extern crate alloc;

mod state;

use firefly_rust::*;
use firefly_ui::{Input, InputManager};
use state::*;

#[unsafe(no_mangle)]
extern "C" fn boot() {
    load_state();
}

#[unsafe(no_mangle)]
extern "C" fn update() {
    let state = get_state();
    // ...
}

#[unsafe(no_mangle)]
extern "C" fn render() {
    let state = get_state();
    firefly_ui::draw_bg(state.settings.theme);
}
