package routes

import (
	"context"
	"encoding/json"
	tgbotapi "github.com/go-telegram/bot"
	tgmodels "github.com/go-telegram/bot/models"
	"seabattle/internal/entity"
	"seabattle/internal/service/rules"
	"seabattle/internal/transport/tg/keyboard"
	"seabattle/internal/transport/tg/request"
	"strconv"
	"strings"
)

func (t router) SetShip(ctx context.Context, bot *tgbotapi.Bot, update *tgmodels.Update) {
	dataStr := strings.Split(update.CallbackQuery.Data, "#")[1]
	var req request.SetShip
	if err := json.Unmarshal([]byte(dataStr), &req); err != nil {
		bot.SendMessage(ctx, &tgbotapi.SendMessageParams{
			ChatID: update.CallbackQuery.Message.Chat.ID,
			Text:   err.Error(),
		})
		return
	}
	tgId := strconv.Itoa(int(update.CallbackQuery.Message.Chat.ID))
	setShipReq := entity.SetShip{TgId: tgId, Point: req.Point, Code: req.Code}
	b, statusCode, err := t.battlePreparation.SetShip(ctx, setShipReq)

	if err != nil && !checkError(err) {
		bot.SendMessage(ctx, &tgbotapi.SendMessageParams{
			ChatID: update.CallbackQuery.Message.Chat.ID,
			Text:   err.Error(),
		})
		_, _ = bot.AnswerCallbackQuery(ctx, &tgbotapi.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
			ShowAlert:       false,
		})

		return
	}
	var text string
	var kb *tgmodels.InlineKeyboardMarkup
	switch statusCode {
	case rules.ShipFirstPoint:
		text = keyboard.FirstPoint
		kb = keyboard.SetBattlefield(b, req.Code, text)
	case rules.ShipSecondPoint:
		text = keyboard.SecondPoint

		kb = keyboard.SetBattlefield(b, req.Code, text)

	case rules.PersonReady:
		text = keyboard.SettingReady
		kb = keyboard.SetBattlefieldWaiting(b)

	case rules.PersonsReady:

		t.createGameAction(ctx, bot, update, req.Code)
		return

	}
	if checkError(err) {
		kb = keyboard.SetBattlefieldWithError(b, req.Code, text, err.Error())

	}
	_, _ = bot.EditMessageReplyMarkup(ctx, &tgbotapi.EditMessageReplyMarkupParams{
		ChatID:    update.CallbackQuery.Message.Chat.ID,
		MessageID: update.CallbackQuery.Message.ID,

		ReplyMarkup: kb,
	})

}
