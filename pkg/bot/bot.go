package bot

import (
	"time"

	"github.com/codingconcepts/env"
	"github.com/joho/godotenv"

	log "github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

const defaultEnvFile = ".env"
const defaultTimeout = 10 * time.Second

type BotConfig struct {
	// Telegram bot token
	TelegramBotToken string `env:"TELEGRAM_BOT_TOKEN" required:"true"`

	// whitelist (chat ID)
	Whitelist []int64 `env:"BOT_WHITELIST"`
}

func Execute() {
	err := godotenv.Load(defaultEnvFile)
	if err != nil {
		log.Warnf("failed to load .env file: %+v", err)
	}

	var config BotConfig
	if err := env.Set(&config); err != nil {
		log.Fatal(err)
	}

	pref := tele.Settings{
		Token:  config.TelegramBotToken,
		Poller: &tele.LongPoller{Timeout: defaultTimeout},
	}

	bot, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}
	if len(config.Whitelist) > 0 {
		bot.Use(whitelist(config.Whitelist...))
	}

	bot.Use(responseTimer)
	bot.Use(messageLogger)
	bot.Handle("/help", HandleHelpCommand)

	log.Infof("Starting bot")
	bot.Start()
}
