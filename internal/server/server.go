package server

import (
	"github.com/google/wire"
)

type Server interface {
	Run() error
}

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewHTTPServer, NewTaskHTTPHandler)
