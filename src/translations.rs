#[derive(Clone, Copy)]
pub enum Message {
    /// Title: "What do you want to delete?".
    WhatToDelete,
    /// Error message: "no app selected".
    NoAppSelected,
    /// Error message: "app already removed".
    AppAlreadyRemoved,
    /// Success message: "app is removed".
    Removed,
    /// Success message: "data is cleared".
    Cleared,
    /// Menu option: "app ROM".
    Rom,
    /// Menu option: "data and save files".
    Data,
    /// Menu option: "screenshots".
    Shots,
    /// Menu option: "achievements".
    Badges,
    /// Menu option: "scores".
    Scores,
    /// Button: "remove".
    Remove,
    /// Button: "cancel".
    Cancel,
}

impl firefly_ui::Translate<'static> for Message {
    fn translate_english(&self) -> &'static str {
        match self {
            Message::WhatToDelete => "What do you want to delete?",
            Message::NoAppSelected => "no app selected",
            Message::AppAlreadyRemoved => "app already removed",
            Message::Removed => "app is removed",
            Message::Cleared => "data is cleared",
            Message::Rom => "app ROM",
            Message::Data => "data and save files",
            Message::Shots => "screenshots",
            Message::Badges => "achievements",
            Message::Scores => "scores",
            Message::Remove => "remove",
            Message::Cancel => "cancel",
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
        match self {
            Message::WhatToDelete => "What do you want to delete?", // TODO
            Message::NoAppSelected => "nicio aplicație selectată",
            Message::AppAlreadyRemoved => "aplicație deja dezinstalată",
            Message::Removed => "aplicația e dezinstalată",
            Message::Cleared => "data is cleared", // TODO
            Message::Rom => "ROM-ul aplicației",
            Message::Data => "date și salvează fișierele",
            Message::Shots => "capturi de ecran",
            Message::Badges => "achievements", // TODO
            Message::Scores => "scores",       // TODO
            Message::Remove => "elimină",
            Message::Cancel => "anulează",
        }
    }

    fn translate_russian(&self) -> &'static str {
        match self {
            Message::WhatToDelete => "Что удалить?",
            Message::NoAppSelected => "приложение не выбрано",
            Message::AppAlreadyRemoved => "приложение уже удалено",
            Message::Removed => "приложение удалено",
            Message::Cleared => "данные удалены",
            Message::Rom => "ROM приложения",
            Message::Data => "данные",
            Message::Shots => "скриншоты",
            Message::Badges => "достижения",
            Message::Scores => "результаты",
            Message::Remove => "удалить",
            Message::Cancel => "отменить",
        }
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
