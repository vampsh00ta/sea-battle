package routes

import (
	"context"
	tgbotapi "github.com/go-telegram/bot"
	tgmodels "github.com/go-telegram/bot/models"
	isrvc "seabattle/internal/app/service"
	"seabattle/internal/pb"
)

type router struct {
	battleAction      isrvc.BattleAction
	battlePreparation isrvc.BattlePreparation
	gc                pb.MatchmakingClient
}

func New(battleAction isrvc.BattleAction, battlePreparation isrvc.BattlePreparation, gc pb.MatchmakingClient) router {
	return router{battleAction: battleAction, battlePreparation: battlePreparation, gc: gc}
}

func (t router) Pass(ctx context.Context, bot *tgbotapi.Bot, update *tgmodels.Update) {
	_, _ = bot.AnswerCallbackQuery(ctx, &tgbotapi.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})

}
