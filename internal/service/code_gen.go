package service

import "github.com/dchest/uniuri"

type CodeGenerator interface {
	GetInviteCode() string
}

func (s service) GetInviteCode() string {
	chars := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	str := uniuri.NewLenChars(6, chars)
	return str
}
