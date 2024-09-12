package routes

import (
	"context"
	tgbotapi "github.com/go-telegram/bot"
	tgmodels "github.com/go-telegram/bot/models"
	"seabattle/internal/transport/tg/keyboard"
	"strconv"
)

func (t router) JoinFight(ctx context.Context, bot *tgbotapi.Bot, update *tgmodels.Update) {
	msg := update.Message.Text
	if len(msg) <= len("/start")+1 {
		_, _ = bot.SendMessage(ctx, &tgbotapi.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "err",
		})
		return
	}
	tgId := update.Message.Chat.ID
	code := msg[len("/start")+1:]
	f, err := t.battlePreparation.JoinFight(ctx, code, strconv.Itoa(int(tgId)))
	if err != nil {
		_, _ = bot.SendMessage(ctx, &tgbotapi.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   err.Error(),
		})
		return
	}
	for i := range f.Users {
		a, _ := bot.SendMessage(ctx, &tgbotapi.SendMessageParams{
			ChatID:      f.Users[i].TgId,
			Text:        "Раставь корабли",
			ReplyMarkup: keyboard.SetBattlefield(f.Users[i].MyField, code, keyboard.FirstPoint),
		})
		if err := t.battlePreparation.SetFieldQueryId(ctx, code, f.Users[i].TgId, strconv.Itoa(a.ID), true); err != nil {
			_, _ = bot.SendMessage(ctx, &tgbotapi.SendMessageParams{
				ChatID: f.Users[i].TgId,
				Text:   err.Error(),
			})
			return
		}

	}

}
