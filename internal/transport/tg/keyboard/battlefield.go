package keyboard

import (
	"encoding/json"
	tgmodels "github.com/go-telegram/bot/models"
	"seabattle/internal/repository/models"
	"seabattle/internal/service/entity"
	"seabattle/internal/transport/tg/request"
)

func handlePoint(field models.Field) string {

	if field.Dead {
		return "❌"
	}
	if field.Shooted {
		return "🔥"
	}
	if field.Ship {
		return "🚢"
	}
	if field.Marked {
		return " # "
	}
	return "  *  "
}

const (
	FirstPoint       = "Выбери первую точку"
	SecondPoint      = "Выбери вторую точку"
	SettingReady     = "Дождись своего оппонента"
	YourTurn         = "Твой ход"
	WaitOpponentTurn = "Дождись хода оппонента"
)

func Battlefield(fight *models.BattleField, sessionId string) *tgmodels.InlineKeyboardMarkup {

	kb := &tgmodels.InlineKeyboardMarkup{
		InlineKeyboard: [][]tgmodels.InlineKeyboardButton{},
	}

	for i := 0; i < 8; i++ {

		res := []tgmodels.InlineKeyboardButton{}
		for j := 0; j < 8; j++ {
			callbackData := request.SetShip{
				Point: entity.Point{
					X: j,
					Y: i,
				},
				Code: sessionId,
			}
			callbackDataBytes, err := json.Marshal(callbackData)
			if err != nil {
				return nil
			}

			res = append(res, tgmodels.InlineKeyboardButton{

				Text: handlePoint(fight.Fields[i][j]), CallbackData: "set#" + string(callbackDataBytes),
			})
		}
		kb.InlineKeyboard = append(kb.InlineKeyboard, res)
	}
	return kb

}
func BattlefieldAction(user *models.User, turn, token string, end bool) (*tgmodels.InlineKeyboardMarkup, *tgmodels.InlineKeyboardMarkup) {

	my := &tgmodels.InlineKeyboardMarkup{
		InlineKeyboard: [][]tgmodels.InlineKeyboardButton{},
	}

	my.InlineKeyboard = append(my.InlineKeyboard, []tgmodels.InlineKeyboardButton{
		tgmodels.InlineKeyboardButton{
			Text: "Твое поле", CallbackData: "pass#",
		},
	})
	for i := 0; i < 8; i++ {

		myTemp := []tgmodels.InlineKeyboardButton{}
		for j := 0; j < 8; j++ {

			myTemp = append(myTemp, tgmodels.InlineKeyboardButton{

				Text: handlePoint(user.MyField.Fields[i][j]), CallbackData: "pass#",
			})

		}
		my.InlineKeyboard = append(my.InlineKeyboard, myTemp)

	}
	enemy := &tgmodels.InlineKeyboardMarkup{
		InlineKeyboard: [][]tgmodels.InlineKeyboardButton{},
	}

	enemy.InlineKeyboard = append(enemy.InlineKeyboard,
		[]tgmodels.InlineKeyboardButton{
			tgmodels.InlineKeyboardButton{
				Text: "Поле противника", CallbackData: "pass#",
			},
		})
	var callbackData string
	for i := 0; i < 8; i++ {
		enemyTemp := []tgmodels.InlineKeyboardButton{}
		for j := 0; j < 8; j++ {
			if user.TgId == turn && !end {
				callbackDataModel := request.Shoot{
					Point: entity.Point{
						X: j,
						Y: i,
					},
					Code: token,
				}
				callbackDataBytes, err := json.Marshal(callbackDataModel)
				if err != nil {
					return nil, nil
				}
				callbackData = "shoot#" + string(callbackDataBytes)
			} else {
				callbackData = "pass#"
			}

			enemyTemp = append(enemyTemp, tgmodels.InlineKeyboardButton{

				Text: handlePoint(user.EnemyField.Fields[i][j]), CallbackData: callbackData,
			})
		}
		enemy.InlineKeyboard = append(enemy.InlineKeyboard, enemyTemp)

	}
	var enemyFieldText string

	if user.TgId == turn {
		enemyFieldText = YourTurn

	} else {
		enemyFieldText = WaitOpponentTurn
	}
	enemy.InlineKeyboard = append(enemy.InlineKeyboard, []tgmodels.InlineKeyboardButton{
		tgmodels.InlineKeyboardButton{
			Text: enemyFieldText, CallbackData: "pass#",
		},
	})

	return my, enemy

}
func SetBattlefieldWaiting(fight *models.BattleField) *tgmodels.InlineKeyboardMarkup {

	kb := &tgmodels.InlineKeyboardMarkup{
		InlineKeyboard: [][]tgmodels.InlineKeyboardButton{},
	}

	for i := 0; i < 8; i++ {

		res := []tgmodels.InlineKeyboardButton{}
		for j := 0; j < 8; j++ {

			res = append(res, tgmodels.InlineKeyboardButton{

				Text: handlePoint(fight.Fields[i][j]), CallbackData: "pass#",
			})
		}
		kb.InlineKeyboard = append(kb.InlineKeyboard, res)
	}
	kb.InlineKeyboard = append(kb.InlineKeyboard, []tgmodels.InlineKeyboardButton{
		tgmodels.InlineKeyboardButton{
			Text: SettingReady, CallbackData: "pass#",
		},
	})
	return kb

}
func SetBattlefield(fight *models.BattleField, sessionId string, text string) *tgmodels.InlineKeyboardMarkup {

	kb := Battlefield(fight, sessionId)
	kb.InlineKeyboard = append(kb.InlineKeyboard, []tgmodels.InlineKeyboardButton{
		tgmodels.InlineKeyboardButton{
			Text: text, CallbackData: "apply",
		},
	})
	return kb

}

func SetBattlefieldWithError(fight *models.BattleField, sessionId string, text, err string) *tgmodels.InlineKeyboardMarkup {

	kb := Battlefield(fight, sessionId)
	kb.InlineKeyboard = append(kb.InlineKeyboard, []tgmodels.InlineKeyboardButton{
		tgmodels.InlineKeyboardButton{
			Text: text, CallbackData: "apply",
		},
		tgmodels.InlineKeyboardButton{
			Text: err, CallbackData: "apply",
		},
	})
	return kb

}
