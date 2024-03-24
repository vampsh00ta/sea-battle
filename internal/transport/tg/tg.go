package tg

import (
	"context"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram/bot"
	tgmodels "github.com/go-telegram/bot/models"
	"seabattle/internal/repository/models"
	"seabattle/internal/service"
	"seabattle/internal/service/rules"
	"seabattle/internal/transport/tg/keyboard"
	"seabattle/internal/transport/tg/request"
	"strings"

	"strconv"
)

type Transport struct {
	srvc service.Service
}

func New(s service.Service) *Transport {
	return &Transport{
		s,
	}
}

// func (t Transport) Test(ctx context.Context, bot *tgbotapi.Bot, update *models.Update) {
//
//		_, err := bot.SendMessage(ctx, &tgbotapi.SendMessageParams{
//			ChatID: update.Message.Chat.ID,
//			Text:   "Доступные файлы",
//		})
//	}
const (
	inviteLink = "http://telegram.me/OceanBattle_bot?start="
)

func (t Transport) CreateFight(ctx context.Context, bot *tgbotapi.Bot, update *tgmodels.Update) {
	tgChatId := update.Message.Chat.ID
	inviteCode, err := t.srvc.CreateFight(ctx, strconv.Itoa(int(tgChatId)))
	if err != nil {
		fmt.Println(err)
		bot.SendMessage(ctx, &tgbotapi.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   err.Error(),
		})
		return
	}

	bot.SendMessage(ctx, &tgbotapi.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   fmt.Sprintf("Отравь эту ссылку своему другу \n ``` %s ``` ", inviteLink+inviteCode),
	})

}

func (t Transport) Pass(ctx context.Context, bot *tgbotapi.Bot, update *tgmodels.Update) {
	bot.AnswerCallbackQuery(ctx, &tgbotapi.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})

	return

}
func (t Transport) JoinFight(ctx context.Context, bot *tgbotapi.Bot, update *tgmodels.Update) {
	msg := update.Message.Text
	if len(msg) <= len("/start")+1 {
		bot.SendMessage(ctx, &tgbotapi.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "err",
		})
		return
	}
	tgId := update.Message.Chat.ID
	code := msg[len("/start")+1:]
	f, err := t.srvc.JoinFight(ctx, code, strconv.Itoa(int(tgId)))
	if err != nil {
		bot.SendMessage(ctx, &tgbotapi.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   err.Error(),
		})
		return
	}
	a, _ := bot.SendMessage(ctx, &tgbotapi.SendMessageParams{
		ChatID:      f.User1.TgId,
		Text:        "Раставь корабли",
		ReplyMarkup: keyboard.SetBattlefield(f.User1.MyField, code, keyboard.FirstPoint),
	})
	if err := t.srvc.SetFieldQueryId(ctx, f.User1.TgId, strconv.Itoa(a.ID), true); err != nil {
		bot.SendMessage(ctx, &tgbotapi.SendMessageParams{
			ChatID: update.CallbackQuery.Message.Chat.ID,
			Text:   err.Error(),
		})
		return
	}
	b, _ := bot.SendMessage(ctx, &tgbotapi.SendMessageParams{
		ChatID:      f.User2.TgId,
		Text:        "Раставь корабли",
		ReplyMarkup: keyboard.SetBattlefield(f.User2.MyField, code, keyboard.FirstPoint),
	})
	if err := t.srvc.SetFieldQueryId(ctx, f.User2.TgId, strconv.Itoa(b.ID), true); err != nil {
		bot.SendMessage(ctx, &tgbotapi.SendMessageParams{
			ChatID: update.CallbackQuery.Message.Chat.ID,
			Text:   err.Error(),
		})
		return
	}

}
func (t Transport) GameAction(ctx context.Context, bot *tgbotapi.Bot, update *tgmodels.Update) {
	tgId := update.CallbackQuery.Message.Chat.ID
	dataStr := strings.Split(update.CallbackQuery.Data, "#")[1]
	var req request.Shoot
	if err := json.Unmarshal([]byte(dataStr), &req); err != nil {
		bot.SendMessage(ctx, &tgbotapi.SendMessageParams{
			ChatID: update.CallbackQuery.Message.Chat.ID,
			Text:   err.Error(),
		})
		return
	}
	req.TgId = strconv.Itoa(int(tgId))
	fight, res, err := t.srvc.Shoot(ctx, req)
	if err != nil {
		bot.SendMessage(ctx, &tgbotapi.SendMessageParams{
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
	toSend := make(map[models.User][]*tgmodels.InlineKeyboardMarkup)
	user_my1, user_enemy1 := keyboard.BattlefieldAction(&fight.User1, fight.Turn, req.Code, end)
	toSend[fight.User1] = []*tgmodels.InlineKeyboardMarkup{user_my1, user_enemy1}

	user_my2, user_enemy2 := keyboard.BattlefieldAction(&fight.User2, fight.Turn, req.Code, end)
	toSend[fight.User2] = []*tgmodels.InlineKeyboardMarkup{user_my2, user_enemy2}

	for user, keyboard := range toSend {
		myQueryIdInt, _ := strconv.Atoi(user.MyFieldQueryId)
		enemyQueryIdInt, _ := strconv.Atoi(user.EnemyFieldQueryId)
		fmt.Println(myQueryIdInt, enemyQueryIdInt)
		_, err := bot.EditMessageReplyMarkup(ctx, &tgbotapi.EditMessageReplyMarkupParams{

			ChatID:      user.TgId,
			MessageID:   myQueryIdInt,
			ReplyMarkup: keyboard[0],
		})
		fmt.Println(err)
		_, err = bot.EditMessageReplyMarkup(ctx, &tgbotapi.EditMessageReplyMarkupParams{

			ChatID:      user.TgId,
			MessageID:   enemyQueryIdInt,
			ReplyMarkup: keyboard[1],
		})

		if end {
			var endText string
			if user.TgId == fight.Turn {
				endText = "Ты выиграл"
			} else {
				endText = "Ты проиграл"

			}
			bot.SendMessage(ctx, &tgbotapi.SendMessageParams{
				ChatID: user.TgId,
				Text:   endText,
			})
		}

	}
}

func (t Transport) createGameAction(ctx context.Context, bot *tgbotapi.Bot, update *tgmodels.Update, token string) {
	fight, err := t.srvc.InitFightAction(ctx, token)

	if err != nil {
		bot.SendMessage(ctx, &tgbotapi.SendMessageParams{
			ChatID: update.CallbackQuery.Message.Chat.ID,
			Text:   err.Error(),
		})
		return
	}
	toSend := make(map[models.User][]*tgmodels.InlineKeyboardMarkup)

	user_my1, user_enemy1 := keyboard.BattlefieldAction(&fight.User1, fight.Turn, token, false)
	toSend[fight.User1] = []*tgmodels.InlineKeyboardMarkup{user_my1, user_enemy1}

	user_my2, user_enemy2 := keyboard.BattlefieldAction(&fight.User2, fight.Turn, token, false)
	toSend[fight.User2] = []*tgmodels.InlineKeyboardMarkup{user_my2, user_enemy2}
	for user, keyboard := range toSend {
		idQueryInt, _ := strconv.Atoi(user.MyFieldQueryId)
		bot.EditMessageReplyMarkup(ctx, &tgbotapi.EditMessageReplyMarkupParams{

			ChatID:      user.TgId,
			MessageID:   idQueryInt,
			ReplyMarkup: keyboard[0],
		})

		field, _ := bot.SendMessage(ctx, &tgbotapi.SendMessageParams{
			ChatID:      user.TgId,
			Text:        "Бой начался",
			ReplyMarkup: keyboard[1],
		})
		if err := t.srvc.SetFieldQueryId(ctx, user.TgId, strconv.Itoa(field.ID), false); err != nil {
			bot.SendMessage(ctx, &tgbotapi.SendMessageParams{
				ChatID: user.TgId,
				Text:   err.Error(),
			})
			return
		}

	}

}

func (t Transport) SetShip(ctx context.Context, bot *tgbotapi.Bot, update *tgmodels.Update) {
	dataStr := strings.Split(update.CallbackQuery.Data, "#")[1]
	var req request.SetShip
	if err := json.Unmarshal([]byte(dataStr), &req); err != nil {
		bot.SendMessage(ctx, &tgbotapi.SendMessageParams{
			ChatID: update.CallbackQuery.Message.Chat.ID,
			Text:   err.Error(),
		})
		return
	}
	tgId := update.CallbackQuery.Message.Chat.ID
	b, statusCode, err := t.srvc.SetShip(ctx, strconv.Itoa(int(tgId)), req.Point, req.Code)

	if err != nil && !checkError(err) {
		bot.SendMessage(ctx, &tgbotapi.SendMessageParams{
			ChatID: update.CallbackQuery.Message.Chat.ID,
			Text:   err.Error(),
		})
		bot.AnswerCallbackQuery(ctx, &tgbotapi.AnswerCallbackQueryParams{
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
		text = "working"

		t.createGameAction(ctx, bot, update, req.Code)
		return

	}
	if checkError(err) {
		kb = keyboard.SetBattlefieldWithError(b, req.Code, text, err.Error())

	}
	_, err = bot.EditMessageReplyMarkup(ctx, &tgbotapi.EditMessageReplyMarkupParams{
		ChatID:    update.CallbackQuery.Message.Chat.ID,
		MessageID: update.CallbackQuery.Message.ID,

		ReplyMarkup: kb,
	})
	fmt.Println(err)

}
