package entity

type Fight struct {
	Attacker  string
	Defender  string
	Turn      string
	SessionId string
	Stage     int
}

const (
	SettedBeginVector = iota
	SettedEndVector
)
