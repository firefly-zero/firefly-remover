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
    let font = state.font.as_font();
    firefly_ui::draw_bg(theme);

    if let Some(msg) = state.msg {
        render_message(state, msg);
        return;
    }

    let title = "What do you want to delete?";
    firefly_ui::draw_title(title, false, &font, theme.accent);
    let pressed = state.input.pressed();
    firefly_ui::draw_cursor((state.cursor + 1).into(), theme, &font, pressed, 0);
    for (switch, i) in state.switches.iter().zip(1u8..) {
        let pressed = pressed && i == state.cursor + 1;
        firefly_ui::draw_switch(i32::from(i), switch.selected, pressed, &font, theme);

        let point = Point::new(20, 25 + 13 * i);
        let name = match switch.kind {
            Kind::Rom => "app ROM",
            Kind::Data => "data and save files",
            Kind::Shots => "screenshots",
        };
        draw_text(name, &font, point, theme.primary);
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
