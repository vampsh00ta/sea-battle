package routes

import (
	"context"
	"encoding/json"
	tgbotapi "github.com/go-telegram/bot"
	tgmodels "github.com/go-telegram/bot/models"
	"seabattle/internal/entity"
	"seabattle/internal/pb"
	"seabattle/internal/service/rules"
	"seabattle/internal/transport/tg/keyboard"
	"strconv"
	"strings"
)

func (t router) GameAction(ctx context.Context, bot *tgbotapi.Bot, update *tgmodels.Update) {
	tgId := update.CallbackQuery.Message.Chat.ID
	dataStr := strings.Split(update.CallbackQuery.Data, "#")[1]
	var req entity.Shoot
	if err := json.Unmarshal([]byte(dataStr), &req); err != nil {
		_, _ = bot.SendMessage(ctx, &tgbotapi.SendMessageParams{
			ChatID: update.CallbackQuery.Message.Chat.ID,
			Text:   err.Error(),
		})
		return
	}
	req.TgId = strconv.Itoa(int(tgId))
	fight, res, err := t.battleAction.Shoot(ctx, req)
	if err != nil {
		_, _ = bot.SendMessage(ctx, &tgbotapi.SendMessageParams{
			ChatID: update.CallbackQuery.Message.Chat.ID,
			Text:   err.Error(),
		})
		return
	}
	var end bool
	switch res {
	case rules.Killed, rules.Missed:
		end = false

	case rules.Lost:
		end = true

	}
	toSend := make(map[entity.User][]*tgmodels.InlineKeyboardMarkup)
	user_my1, user_enemy1 := keyboard.BattlefieldAction(&fight.Users[0], fight.Turn, req.Code, end)
	toSend[fight.Users[0]] = []*tgmodels.InlineKeyboardMarkup{user_my1, user_enemy1}

	user_my2, user_enemy2 := keyboard.BattlefieldAction(&fight.Users[1], fight.Turn, req.Code, end)
	toSend[fight.Users[1]] = []*tgmodels.InlineKeyboardMarkup{user_my2, user_enemy2}
	var loser, winner int64
	for user, k := range toSend {
		myQueryIdInt, _ := strconv.Atoi(user.MyFieldQueryId)
		enemyQueryIdInt, _ := strconv.Atoi(user.EnemyFieldQueryId)

		_, _ = bot.EditMessageReplyMarkup(ctx, &tgbotapi.EditMessageReplyMarkupParams{

			ChatID:      user.TgId,
			MessageID:   myQueryIdInt,
			ReplyMarkup: k[0],
		})
		_, _ = bot.EditMessageReplyMarkup(ctx, &tgbotapi.EditMessageReplyMarkupParams{

			ChatID:      user.TgId,
			MessageID:   enemyQueryIdInt,
			ReplyMarkup: k[1],
		})

		if end {
			var endText string
			var tgID int
			tgID, _ = strconv.Atoi(user.TgId)

			if user.TgId == fight.Turn {
				endText = "Ты выиграл"
				winner = int64(tgID)
			} else {
				endText = "Ты проиграл"
				loser = int64(tgID)
			}
			_, _ = bot.SendMessage(ctx, &tgbotapi.SendMessageParams{
				ChatID: user.TgId,
				Text:   endText,
			})
		}

	}
	if end {
		go func() {
			_, _ = t.gc.MatchResult(ctx, &pb.MatchResultRequest{TgIDLoser: loser, TgIDWinner: winner})
		}()

	}
}
