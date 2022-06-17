package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mikelorant/easyredir-cli/pkg/easyredir"
	"github.com/mikelorant/easyredir-cli/pkg/easyredir/host"
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

	// Get source host for rule
	hostID := cr.Data.Relationships.SourceHosts.Data[0].ID
	gh, err := e.GetHost(hostID)
	if err != nil {
		log.Fatalf("unable to get host: %v: %v", hostID, err)
	}

	fmt.Println("Get host output:")
	fmt.Println(gh)

	// Update source host for rule
	hattr := host.Attributes{
		MatchOptions: host.MatchOptions{
			CaseInsensitive:  ref(true),
			SlashInsensitive: ref(true),
		},
	}
	uh, err := e.UpdateHost(hostID, hattr)
	if err != nil {
		log.Fatalf("unable to update host: %v: %v", hostID, err)
	}

	fmt.Println("Update host output:")
	fmt.Println(uh)

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
