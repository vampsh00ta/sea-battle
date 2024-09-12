package service

import "github.com/dchest/uniuri"

type CodeGenerator interface {
	GetCode() string
}
type codeGenerator struct{}

func newCodeGenerator() CodeGenerator {
	return &codeGenerator{}
}
func (s codeGenerator) GetCode() string {
	chars := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	str := uniuri.NewLenChars(6, chars)
	return str
}
