package routes

import "seabattle/internal/service/rules"

func checkError(err error) bool {
	if err == nil {
		return false
	}
	switch err.Error() {
	case rules.WrongPlacementErr, rules.WrongLengthErr, rules.MaxShipCountErr:
		return true
	}
	return false
}
