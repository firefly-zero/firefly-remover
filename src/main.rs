#![no_std]
#![no_main]
extern crate alloc;

mod state;
mod translations;

use firefly_rust::*;
use firefly_ui::{Input, Translate};
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

    let input = state.input.get();
    if state.msg.is_some() {
        if input == Input::Select {
            quit();
        }
        return;
    }

    match input {
        Input::Up => {
            if state.cursor > 0 {
                state.cursor -= 1;
            }
        }
        Input::Down => {
            if usize::from(state.cursor) < state.switches.len() {
                state.cursor += 1;
            }
        }
        Input::Left => state.cursor = 0,
        Input::Right => state.cursor = state.switches.len() as u8,
        Input::Select => {
            let cursor = usize::from(state.cursor);
            if let Some(switch) = state.switches.get_mut(cursor) {
                switch.selected = !switch.selected;
            } else {
                let any_selected = state.switches.iter().any(|s| s.selected);
                if any_selected {
                    remove_app(state);
                } else {
                    quit();
                }
            }
        }
        Input::Back => quit(),
        Input::None => {}
    }
}

fn remove_app(state: &mut State) {
    let (author_id, app_id) = state.target.as_ref().unwrap();
    let id = alloc::format!("{author_id}/{app_id}");
    // Delete ROM and detect if we can remove the whole data dir
    // or only some subdirs.
    let mut del_all_data = true;
    let mut del_etc = false;
    let mut del_shots = false;
    for option in &state.switches {
        if option.kind == Kind::Rom {
            if option.selected {
                sudo::remove_dir(&alloc::format!("roms/{id}"));
                // Reset launcher cache.
                sudo::remove_file("data/sys/launcher/etc/metas");
            }
            continue;
        }
        if !option.selected {
            del_all_data = false;
            continue;
        }
        if option.kind == Kind::Data {
            del_etc = true
        }
        if option.kind == Kind::Shots {
            del_shots = true
        }
    }

    state.msg = Some(Message::Removed);
    if del_all_data {
        sudo::remove_dir(&alloc::format!("data/{id}"));
        return;
    }

    // Delete the app data and stash.
    if del_etc {
        sudo::remove_dir(&alloc::format!("data/{id}/etc"));
        let stash_path = &alloc::format!("roms/{id}/stash");
        if sudo::get_file_size(stash_path) != 0 {
            sudo::remove_file(stash_path)
        }
    }

    // Delete screenshots.
    if del_shots {
        sudo::remove_dir(&alloc::format!("roms/{id}/shots"))
    }
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

        let mut point = Point::new(20, 25 + 13 * i);
        if pressed {
            point.x += 1;
            point.y += 1;
        }
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
    let pressed = pressed && usize::from(state.cursor) == state.switches.len();
    let mut point = Point::new(20, 25 + 13 * (state.switches.len() as i32 + 1));
    if pressed {
        point.x += 1;
        point.y += 1;
    }
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
