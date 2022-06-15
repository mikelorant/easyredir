package app

import (
	"fmt"
	"os"

	"github.com/mikelorant/easyredir-cli/pkg/easyredir"
)

func Run() error {
	e := easyredir.New(
		easyredir.WithAPIKey(os.Getenv("EASYREDIR_API_KEY")),
		easyredir.WithAPISecret(os.Getenv("EASYREDIR_API_SECRET")),
	)

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
	r, err := e.ListRules()
	if err != nil {
		return fmt.Errorf("unable to list rules: %w", err)
	}

	fmt.Print(r)

	return nil
}

func listHosts(e *easyredir.Easyredir) error {
	h, err := e.ListHosts()
	if err != nil {
		return fmt.Errorf("unable to list hosts: %w", err)
	}

	fmt.Print(h)

	return nil
}

func getHost(e *easyredir.Easyredir) error {
	h, err := e.GetHost(os.Args[2])
	if err != nil {
		return fmt.Errorf("unable to get host: %v: %w", os.Args[2], err)
	}

	fmt.Print(h)

	return nil
}
