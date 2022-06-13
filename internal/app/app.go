package app

import (
	"fmt"
	"os"

	"github.com/mikelorant/easyredir-cli/pkg/easyredir"
)

func Run() error {
	e := easyredir.New(&easyredir.Config{
		APIKey:    os.Getenv("EASYREDIR_API_KEY"),
		APISecret: os.Getenv("EASYREDIR_API_SECRET"),
	})

	if len(os.Args) <= 1 {
		return nil
	}

	switch os.Args[1] {
	case "rules":
		if err := listRules(e); err != nil {
			return fmt.Errorf("unable to list rules: %w", err)
		}
	case "hosts":
		if err := listHosts(e); err != nil {
			return fmt.Errorf("unable to list hosts: %w", err)
		}
	}

	return nil
}

func listRules(e *easyredir.Easyredir) error {
	rules, err := e.ListRules(easyredir.WithLimit(100))
	if err != nil {
		return fmt.Errorf("unable to list rules: %w", err)
	}

	fmt.Print(rules)

	return nil
}

func listHosts(e *easyredir.Easyredir) error {
	hosts, err := e.ListHosts(easyredir.WithHostsLimit(100))
	if err != nil {
		return fmt.Errorf("unable to list hosts: %w", err)
	}

	fmt.Print(hosts)

	return nil
}
