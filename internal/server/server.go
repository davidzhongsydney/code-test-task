package server

import (
	"github.com/google/wire"
)

type IServer interface {
	Run() error
}

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewHTTPServer, NewTaskHTTPHandler)
