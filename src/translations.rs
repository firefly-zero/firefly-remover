#[derive(Clone, Copy)]
pub enum Message {
    WhatToDelete,
    NoAppSelected,
    AppAlreadyRemoved,
    Rom,
    Data,
    Shots,
    Remove,
    Cancel,
    Removed,
}

impl firefly_ui::Translate<'static> for Message {
    fn translate_english(&self) -> &'static str {
        match self {
            Message::WhatToDelete => "What do you want to delete?",
            Message::NoAppSelected => "no app selected",
            Message::AppAlreadyRemoved => "app already removed",
            Message::Rom => "app ROM",
            Message::Data => "data and save files",
            Message::Shots => "screenshots",
            Message::Remove => "remove",
            Message::Cancel => "cancel",
            Message::Removed => "app is removed",
        }
    }

    fn translate_dutch(&self) -> &'static str {
        self.translate_english()
    }

    fn translate_french(&self) -> &'static str {
        self.translate_english()
    }

    fn translate_german(&self) -> &'static str {
        self.translate_english()
    }

    fn translate_italian(&self) -> &'static str {
        self.translate_english()
    }

    fn translate_polish(&self) -> &'static str {
        self.translate_english()
    }

    fn translate_romanian(&self) -> &'static str {
        self.translate_english()
    }

    fn translate_russian(&self) -> &'static str {
        self.translate_english()
    }

    fn translate_spanish(&self) -> &'static str {
        self.translate_english()
    }

    fn translate_swedish(&self) -> &'static str {
        self.translate_english()
    }

    fn translate_turkish(&self) -> &'static str {
        self.translate_english()
    }

    fn translate_ukrainian(&self) -> &'static str {
        self.translate_english()
    }

    fn translate_toki_pona(&self) -> &'static str {
        self.translate_english()
    }
}
