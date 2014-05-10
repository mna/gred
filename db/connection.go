package db

import "errors"

var (
	pong    = []byte("+PONG\r\n")
	errPong = errors.New("pong")
)

var cmdEcho = CheckArgCount(
	func(ctx *Ctx) (interface{}, error) {
		return ctx.s0, nil
	}, 1, 1)

var cmdPing = CheckArgCount(
	func(ctx *Ctx) (interface{}, error) {
		return nil, errPong
	}, 0, 0)

var cmdQuit = CheckArgCount(
	func(ctx *Ctx) (interface{}, error) {
		return nil, nil
	}, 0, 0)
