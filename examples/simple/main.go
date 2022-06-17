package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mikelorant/easyredir-cli/pkg/easyredir"
	"github.com/mikelorant/easyredir-cli/pkg/easyredir/rule"
)

func main() {
	// New client
	e := easyredir.New(
		easyredir.WithAPIKey(os.Getenv("EASYREDIR_API_KEY")),
		easyredir.WithAPISecret(os.Getenv("EASYREDIR_API_SECRET")),
	)

	// Create rule
	rattr := rule.Attributes{
		ForwardParams: ref(true),
		ForwardPath:   ref(true),
		ResponseType:  ref(rule.ResponseMovedPermanently),
		SourceURLs: []string{
			"source.nineexample.com",
		},
		TargetURL: ref("target.nineexample.com"),
	}

	cr, err := e.CreateRule(rattr)
	if err != nil {
		log.Fatalf("unable to create rule: %v\n", err)
	}

	fmt.Println("Create rule output:")
	fmt.Println(cr)

	// List rule
	lr, err := e.ListRules(
		easyredir.WithSourceFilter("source.nineexample.com"),
		easyredir.WithTargetFilter("target.nineexample.com"),
	)
	if err != nil {
		log.Fatalf("unable to list rule: %v\n", err)
	}

	fmt.Println("List rule output:")
	fmt.Println(lr)

	// Update rule
	rattr = rule.Attributes{
		SourceURLs: []string{
			"sourceupdated.ninenineexample.com",
		},
		TargetURL: ref("targetupdated.nineexample.com"),
	}
	ur, err := e.UpdateRule(cr.Data.ID, rattr)
	if err != nil {
		log.Fatalf("unable to update rule: %v\n", err)
		return
	}

	fmt.Println("Update rule output:")
	fmt.Println(ur)

	// Remove rule
	res, err := e.RemoveRule(cr.Data.ID)
	if err != nil {
		log.Fatalf("unable to remove rule: %v\n", err)
		return
	}

	fmt.Printf("Result of remove rule for %v: %v\n", cr.Data.ID, res)
}

func ref[T any](x T) *T {
	return &x
}
