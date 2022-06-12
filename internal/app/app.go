package app

import (
	"fmt"
	"os"

	"github.com/gotidy/ptr"
	"github.com/mikelorant/easyredir-cli/pkg/easyredir"
)

func Run() error {
	e := easyredir.New(&easyredir.Config{
		APIKey:    os.Getenv("EASYREDIR_API_KEY"),
		APISecret: os.Getenv("EASYREDIR_API_SECRET"),
	})
	e.Ping()
	rules, err := e.ListRulesPaginator()
	if err != nil {
		return fmt.Errorf("unable to list rules: %w", err)
	}

	for _, r := range rules.Data {
		fmt.Println(ptr.ToString(r.ID))
	}

	return nil
}
