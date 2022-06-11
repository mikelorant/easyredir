package app

import (
	"github.com/mikelorant/easyredir-cli/pkg/easyredir"
)

func Run() error {
	e := easyredir.New(&easyredir.Config{
		Key:    "key",
		Secret: "secret",
	})
	e.Ping()

	return nil
}
