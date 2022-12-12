package model

import "context"

type GlobalContext struct {
	Ctx    context.Context
	Cancel context.CancelFunc
}
