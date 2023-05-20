package bot

type BotConfig struct {
	// Telegram bot token
	TelegramBotToken string `env:"TELEGRAM_BOT_TOKEN" required:"true"`

	// whitelist (chat ID)
	Whitelist []int64 `env:"BOT_WHITELIST"`
}
