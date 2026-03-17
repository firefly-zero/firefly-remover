#![no_std]
#![no_main]
extern crate alloc;

mod state;
mod translations;

use firefly_rust::*;
use firefly_ui::{Input, InputManager, Translate};
use state::*;
use translations::*;

#[unsafe(no_mangle)]
extern "C" fn boot() {
    load_state();
}

#[unsafe(no_mangle)]
extern "C" fn update() {
    let state = get_state();
    state.input.update();
    // ...
}

#[unsafe(no_mangle)]
extern "C" fn render() {
    let state = get_state();
    let theme = state.settings.theme;
    let lang = state.settings.language;
    let font = state.font.as_font();
    firefly_ui::draw_bg(theme);

    if let Some(msg) = state.msg {
        let msg = msg.translate(lang);
        render_message(state, msg);
        return;
    }

    let title = Message::WhatToDelete.translate(lang);
    firefly_ui::draw_title(title, false, &font, theme.accent);
    let pressed = state.input.pressed();
    firefly_ui::draw_cursor((state.cursor + 1).into(), theme, &font, pressed, 0);

    // Draw switches.
    for (switch, i) in state.switches.iter().zip(1u8..) {
        let pressed = pressed && i == state.cursor + 1;
        firefly_ui::draw_switch(i32::from(i), switch.selected, pressed, &font, theme);

        let point = Point::new(20, 25 + 13 * i);
        let name = match switch.kind {
            Kind::Rom => Message::Rom,
            Kind::Data => Message::Data,
            Kind::Shots => Message::Shots,
        };
        let name = name.translate(lang);
        draw_text(name, &font, point, theme.primary);
    }

    // Draw button.
    let any_selected = state.switches.iter().any(|s| s.selected);
    let msg = if any_selected {
        Message::Remove
    } else {
        Message::Cancel
    };
    let msg = msg.translate(lang);
    let point = Point::new(20, 25 + 13 * (state.switches.len() as i32 + 1));
    draw_text(msg, &font, point, theme.accent);
}

fn render_message(state: &State, msg: &str) {
    let theme = state.settings.theme;
    let font = state.font.as_font();
    let x = (WIDTH - font.line_width_utf8(msg) as i32) / 2;
    let y = (HEIGHT - i32::from(font.char_height())) / 2;
    let point = Point::new(x, y);
    draw_text(msg, &font, point, theme.accent);
}
