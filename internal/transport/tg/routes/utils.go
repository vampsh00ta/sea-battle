package routes

import (
	"context"
	tgbotapi "github.com/go-telegram/bot"
	tgmodels "github.com/go-telegram/bot/models"
	"seabattle/internal/entity"
	"seabattle/internal/transport/tg/keyboard"
	"strconv"
)

func (t router) createGameAction(ctx context.Context, bot *tgbotapi.Bot, update *tgmodels.Update, token string) {
	fight, err := t.battlePreparation.InitFightAction(ctx, token)

	if err != nil {
		bot.SendMessage(ctx, &tgbotapi.SendMessageParams{
			ChatID: update.CallbackQuery.Message.Chat.ID,
			Text:   err.Error(),
		})
		return
	}
	toSend := make(map[entity.User][]*tgmodels.InlineKeyboardMarkup)

	user_my1, user_enemy1 := keyboard.BattlefieldAction(&fight.Users[0], fight.Turn, token, false)
	toSend[fight.Users[0]] = []*tgmodels.InlineKeyboardMarkup{user_my1, user_enemy1}

	user_my2, user_enemy2 := keyboard.BattlefieldAction(&fight.Users[1], fight.Turn, token, false)
	toSend[fight.Users[1]] = []*tgmodels.InlineKeyboardMarkup{user_my2, user_enemy2}

	for user, keyboard := range toSend {
		idQueryInt, _ := strconv.Atoi(user.MyFieldQueryId)
		_, _ = bot.EditMessageReplyMarkup(ctx, &tgbotapi.EditMessageReplyMarkupParams{

			ChatID:      user.TgId,
			MessageID:   idQueryInt,
			ReplyMarkup: keyboard[0],
		})

		field, _ := bot.SendMessage(ctx, &tgbotapi.SendMessageParams{
			ChatID:      user.TgId,
			Text:        "Бой начался",
			ReplyMarkup: keyboard[1],
		})
		if err := t.battlePreparation.SetFieldQueryId(ctx, token, user.TgId, strconv.Itoa(field.ID), false); err != nil {
			bot.SendMessage(ctx, &tgbotapi.SendMessageParams{
				ChatID: user.TgId,
				Text:   err.Error(),
			})
			return
		}

	}

}
