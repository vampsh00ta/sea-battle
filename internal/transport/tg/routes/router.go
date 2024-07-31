package routes

import (
	"context"
	tgbotapi "github.com/go-telegram/bot"
	tgmodels "github.com/go-telegram/bot/models"
	isrvc "seabattle/internal/app/service"
	"seabattle/internal/pb"
)

type router struct {
	srvc isrvc.Service
	gc   pb.MatchmakingClient
}

func New(srvc isrvc.Service, gc pb.MatchmakingClient) router {
	return router{srvc: srvc, gc: gc}
}

func (t router) Pass(ctx context.Context, bot *tgbotapi.Bot, update *tgmodels.Update) {
	_, _ = bot.AnswerCallbackQuery(ctx, &tgbotapi.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})

}
