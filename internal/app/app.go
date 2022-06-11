package app

import (
	"github.com/mikelorant/easyredir-cli/pkg/easyredir"
)

func Run() error {
	e := easyredir.New(&easyredir.Config{})
	e.Ping()

	return nil
}
