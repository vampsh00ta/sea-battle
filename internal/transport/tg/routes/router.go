package routes

import (
	"context"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram/bot"
	tgmodels "github.com/go-telegram/bot/models"
	isrvc "seabattle/internal/app/service"
	"seabattle/internal/entity"
	"seabattle/internal/pb"
	"seabattle/internal/service/rules"
	"seabattle/internal/transport/tg/keyboard"
	"seabattle/internal/transport/tg/request"
	"strconv"
	"strings"
)

type router struct {
	srvc isrvc.Service
	gc   pb.MatchmakingClient
}

func New(srvc isrvc.Service, gc pb.MatchmakingClient) router {
	return router{srvc: srvc, gc: gc}
}

func (t router) CreateFight(ctx context.Context, bot *tgbotapi.Bot, update *tgmodels.Update) {

	tgChatId := update.Message.Chat.ID

	inviteCode, err := t.srvc.CreateFight(ctx, strconv.Itoa(int(tgChatId)))

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

func (t router) Pass(ctx context.Context, bot *tgbotapi.Bot, update *tgmodels.Update) {
	_, _ = bot.AnswerCallbackQuery(ctx, &tgbotapi.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})

}
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
	f, err := t.srvc.JoinFight(ctx, code, strconv.Itoa(int(tgId)))
	if err != nil {
		bot.SendMessage(ctx, &tgbotapi.SendMessageParams{
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
		if err := t.srvc.SetFieldQueryId(ctx, code, f.Users[i].TgId, strconv.Itoa(a.ID), true); err != nil {
			bot.SendMessage(ctx, &tgbotapi.SendMessageParams{
				ChatID: f.Users[i].TgId,
				Text:   err.Error(),
			})
			return
		}

	}

}
func (t router) GameAction(ctx context.Context, bot *tgbotapi.Bot, update *tgmodels.Update) {
	tgId := update.CallbackQuery.Message.Chat.ID
	dataStr := strings.Split(update.CallbackQuery.Data, "#")[1]
	var req entity.Shoot
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
	toSend := make(map[entity.User][]*tgmodels.InlineKeyboardMarkup)
	user_my1, user_enemy1 := keyboard.BattlefieldAction(&fight.Users[0], fight.Turn, req.Code, end)
	toSend[fight.Users[0]] = []*tgmodels.InlineKeyboardMarkup{user_my1, user_enemy1}

	user_my2, user_enemy2 := keyboard.BattlefieldAction(&fight.Users[1], fight.Turn, req.Code, end)
	toSend[fight.Users[1]] = []*tgmodels.InlineKeyboardMarkup{user_my2, user_enemy2}

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

func (t router) createGameAction(ctx context.Context, bot *tgbotapi.Bot, update *tgmodels.Update, token string) {
	fight, err := t.srvc.InitFightAction(ctx, token)

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
		if err := t.srvc.SetFieldQueryId(ctx, token, user.TgId, strconv.Itoa(field.ID), false); err != nil {
			bot.SendMessage(ctx, &tgbotapi.SendMessageParams{
				ChatID: user.TgId,
				Text:   err.Error(),
			})
			return
		}

	}

}

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
	b, statusCode, err := t.srvc.SetShip(ctx, setShipReq)

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
	_, err = bot.EditMessageReplyMarkup(ctx, &tgbotapi.EditMessageReplyMarkupParams{
		ChatID:    update.CallbackQuery.Message.Chat.ID,
		MessageID: update.CallbackQuery.Message.ID,

		ReplyMarkup: kb,
	})

}
