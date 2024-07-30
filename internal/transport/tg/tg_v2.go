package tg

import (
	tgbotapi "github.com/go-telegram/bot"
	isrvc "seabattle/internal/app/service"
	"seabattle/internal/pb"
	"seabattle/internal/transport/tg/routes"
)

func New(bot *tgbotapi.Bot, s isrvc.Service, gc pb.MatchmakingClient) {
	tr := routes.New(s, gc)
	bot.RegisterHandler(tgbotapi.HandlerTypeMessageText,
		"/creategame",
		tgbotapi.MatchTypeExact, tr.CreateFight,
	)
	bot.RegisterHandler(tgbotapi.HandlerTypeMessageText,
		"/start",
		tgbotapi.MatchTypePrefix, tr.JoinFight,
	)
	bot.RegisterHandler(tgbotapi.HandlerTypeMessageText,
		"/search",
		tgbotapi.MatchTypePrefix, tr.SearchFight,
	)
	//bot.RegisterHandler(tgbotapi.HandlerTypeMessageText,
	//	"/search_battle",
	//	tgbotapi.MatchTypePrefix, tr.Search,
	//)
	bot.RegisterHandler(tgbotapi.HandlerTypeMessageText,
		"/joingame",
		tgbotapi.MatchTypePrefix, tr.JoinFight,
	)
	bot.RegisterHandler(tgbotapi.HandlerTypeCallbackQueryData,
		"set#",
		tgbotapi.MatchTypePrefix, tr.SetShip,
	)
	bot.RegisterHandler(tgbotapi.HandlerTypeCallbackQueryData,
		"pass#",
		tgbotapi.MatchTypePrefix, tr.Pass,
	)
	bot.RegisterHandler(tgbotapi.HandlerTypeCallbackQueryData,
		"shoot#",
		tgbotapi.MatchTypePrefix, tr.GameAction,
	)

}

//	func (t router) Search(ctx context.Context, bot *tgbotapi.Bot, update *tgmodels.Update) {
//		tgChatId := update.Message.Chat.ID
//		fmt.Println("ASdsadasd")
//		if err := t.srvc.SearchFight(ctx, int(tgChatId)); err != nil {
//			return
//		}
//		bot.SendMessage(ctx, &tgbotapi.SendMessageParams{
//			ChatID: update.Message.Chat.ID,
//			Text:   fmt.Sprintf("Идет поиск боя ..."),
//		})
//
// }
