package db

import "errors"

var (
	// pong is the standard PONG serialized response, to avoid allocations
	// for this common case.
	pong = []byte("+PONG\r\n")

	// errPong is a sentinel error value to indicate that the standard PONG
	// response should be returned.
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
