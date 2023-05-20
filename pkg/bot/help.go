package bot

import tele "gopkg.in/telebot.v3"

func HandleHelpCommand(c tele.Context) error {
	return c.Reply(`/help - show this help message`)
}
