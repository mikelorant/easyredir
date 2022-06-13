package app

import (
	"fmt"
	"os"

	"github.com/mikelorant/easyredir-cli/pkg/easyredir"
	"github.com/mikelorant/easyredir-cli/pkg/easyredir/host"
	"github.com/mikelorant/easyredir-cli/pkg/easyredir/rule"
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
		if len(os.Args) > 2 {
			if err := getHost(e); err != nil {
				return fmt.Errorf("unable to list hosts: %w", err)
			}
			return nil
		}

		if err := listHosts(e); err != nil {
			return fmt.Errorf("unable to list hosts: %w", err)
		}
	}

	return nil
}

func listRules(e *easyredir.Easyredir) error {
	r, err := rule.ListRules(e, rule.WithLimit(100))
	if err != nil {
		return fmt.Errorf("unable to list rules: %w", err)
	}

	fmt.Print(r)

	return nil
}

func listHosts(e *easyredir.Easyredir) error {
	h, err := host.ListHostsPaginator(e, host.WithHostsLimit(100))
	if err != nil {
		return fmt.Errorf("unable to list hosts: %w", err)
	}

	fmt.Print(h)

	return nil
}

func getHost(e *easyredir.Easyredir) error {
	h, err := host.GetHost(e, os.Args[2])
	if err != nil {
		return fmt.Errorf("unable to get host: %v: %w", os.Args[2], err)
	}

	fmt.Print(h)

	return nil
}
