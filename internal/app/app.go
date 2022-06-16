package app

import (
	"fmt"
	"os"

	"github.com/gotidy/ptr"
	"github.com/mikelorant/easyredir-cli/pkg/easyredir"
	"github.com/mikelorant/easyredir-cli/pkg/easyredir/rule"
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

	case "create":
		attr := rule.Attributes{
			ForwardParams: ptr.Bool(true),
			ForwardPath:   ptr.Bool(true),
			ResponseType:  ptr.String("moved_permanently"),
			SourceURLs: []string{
				"source.example.com",
			},
			TargetURL: ptr.String("target.example.com"),
		}

		r, err := e.CreateRule(attr)
		if err != nil {
			return fmt.Errorf("unable to create rule: %w", err)
		}

		fmt.Print(r)

		return nil

	case "remove":
		res, err := e.RemoveRule(os.Args[2])
		if err != nil {
			return fmt.Errorf("unable to remove rule: %w", err)
		}

		fmt.Print(res)

		return nil
	}

	return nil
}
