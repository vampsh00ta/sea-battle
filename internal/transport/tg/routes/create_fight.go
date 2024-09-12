package routes

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram/bot"
	tgmodels "github.com/go-telegram/bot/models"
	"strconv"
)

func (t router) CreateFight(ctx context.Context, bot *tgbotapi.Bot, update *tgmodels.Update) {

	tgChatId := update.Message.Chat.ID

	inviteCode, err := t.battlePreparation.CreateFight(ctx, strconv.Itoa(int(tgChatId)))

	if err != nil {
		_, _ = bot.SendMessage(ctx, &tgbotapi.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   err.Error(),
		})
		return
	}
	_, _ = bot.SendMessage(ctx, &tgbotapi.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   fmt.Sprintf("Отравь эту ссылку своему другу \n ```%s``` ", inviteLink+inviteCode),
	})

}
