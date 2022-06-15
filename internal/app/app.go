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
		r, err := e.ListRules(easyredir.WithLimit(100))
		if err != nil {
			return fmt.Errorf("unable to list rules: %w", err)
		}

		fmt.Print(r)

		return nil

	case "hosts":
		if len(os.Args) > 2 {
			h, err := e.GetHost(os.Args[2])
			if err != nil {
				return fmt.Errorf("unable to get host: %v: %w", os.Args[2], err)
			}

			fmt.Print(h)

			return nil
		}

		h, err := e.ListHosts(easyredir.WithLimit(100))
		if err != nil {
			return fmt.Errorf("unable to list hosts: %w", err)
		}

		fmt.Print(h)

		return nil
	}

	return nil
}
