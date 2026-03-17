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
    let theme = state.settings.theme;
    firefly_ui::draw_bg(theme);

    if let Some(msg) = state.msg {
        render_message(state, msg);
        return;
    }
}

fn render_message(state: &State, msg: &str) {
    let theme = state.settings.theme;
    let font = state.font.as_font();
    let x = (WIDTH - font.line_width_utf8(msg) as i32) / 2;
    let y = (HEIGHT - i32::from(font.char_height())) / 2;
    let point = Point::new(x, y);
    draw_text(msg, &font, point, theme.accent);
}
