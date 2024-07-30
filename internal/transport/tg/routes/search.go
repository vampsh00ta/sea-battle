package routes

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram/bot"
	tgmodels "github.com/go-telegram/bot/models"
	"seabattle/internal/entity"
	"seabattle/internal/pb"
	"seabattle/internal/transport/tg/keyboard"
	"strconv"
)

func (t router) SearchFight(ctx context.Context, bot *tgbotapi.Bot, update *tgmodels.Update) {
	//msg := update.Message.Text
	tgId := update.Message.Chat.ID

	res, err := t.gc.FindMatch(ctx, &pb.FindMatchRequest{TgID: int64(tgId)})
	if err != nil {
		_, _ = bot.SendMessage(ctx, &tgbotapi.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   err.Error(),
		})
		return
	}
	if err != nil {
		_, _ = bot.SendMessage(ctx, &tgbotapi.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   err.Error(),
		})
		return
	}
	switch {
	case res.TgID == -1:
		bot.SendMessage(ctx, &tgbotapi.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Идет поиск оппонента",
		})
		return
	default:
		tgIDs := []string{
			strconv.Itoa(int(tgId)),
			strconv.Itoa(int(res.TgID)),
		}
		fmt.Println(tgIDs)
		f, err := t.srvc.InitFight(ctx, tgIDs...)
		if err != nil {
			_, _ = bot.SendMessage(ctx, &tgbotapi.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   err.Error(),
			})
			return
		}
		if err := t.notifyUsersAboutFight(ctx, bot, f); err != nil {
			return
		}

	}

}

func (t router) notifyUsersAboutFight(ctx context.Context, bot *tgbotapi.Bot, f entity.Fight) error {
	for i := range f.Users {
		a, _ := bot.SendMessage(ctx, &tgbotapi.SendMessageParams{
			ChatID:      f.Users[i].TgId,
			Text:        "Раставь корабли",
			ReplyMarkup: keyboard.SetBattlefield(f.Users[i].MyField, f.SessionId, keyboard.FirstPoint),
		})
		if err := t.srvc.SetFieldQueryId(ctx, f.SessionId, f.Users[i].TgId, strconv.Itoa(a.ID), true); err != nil {
			bot.SendMessage(ctx, &tgbotapi.SendMessageParams{
				ChatID: f.Users[i].TgId,
				Text:   err.Error(),
			})
			return err
		}

	}
	return nil
}
