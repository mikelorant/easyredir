package main

import (
	"fmt"
	"log"

	"github.com/alexflint/go-arg"
	"github.com/mikelorant/easyredir-cli/pkg/easyredir"
	"github.com/mikelorant/easyredir-cli/pkg/easyredir/host"
	"github.com/mikelorant/easyredir-cli/pkg/easyredir/rule"
)

type CreateCmd struct {
	Rule *CreateRuleCmd `arg:"subcommand:rule"`
}

type CreateRuleCmd struct {
	ForwardParams *bool              `arg:"--forward-params,required"`
	ForwardPath   *bool              `arg:"--forward-path,required"`
	ResponseType  *rule.ResponseType `arg:"--response-type,required"`
	SourceURLs    []string           `arg:"--source-url,required"`
	TargetURL     *string            `arg:"--target-url,required"`
}

type GetCmd struct {
	Host *GetHostCmd `arg:"subcommand:host"`
}

type GetHostCmd struct {
	ID string `arg:"positional"`
}

type ListCmd struct {
	Host *ListHostsCmd `arg:"subcommand:hosts"`
	Rule *ListRulesCmd `arg:"subcommand:rules"`
}

type ListHostsCmd struct{}

type ListRulesCmd struct {
	SourceFilter string `arg:"--source-filter"`
	TargetFilter string `arg:"--target-filter"`
}

type RemoveCmd struct {
	Rule *RemoveRuleCmd `arg:"subcommand:rule"`
}

type RemoveRuleCmd struct {
	ID string `arg:"positional"`
}

type UpdateCmd struct {
	Host *UpdateHostCmd `arg:"subcommand:host"`
	Rule *UpdateRuleCmd `arg:"subcommand:rule"`
}

type UpdateHostCmd struct {
	ID                      string             `arg:"positional"`
	CaseInsensitive         *bool              `arg:"--case-insenstive"`
	SlashInsensitive        *bool              `arg:"--slash-insensitive"`
	ForwardParams           *bool              `arg:"--forward-params"`
	ForwardPath             *bool              `arg:"--forward-path"`
	Custom404Body           *string            `arg:"--custom-404-body"`
	ResponseCode            *host.ResponseCode `arg:"--response-code"`
	ResponseURL             *string            `arg:"--response-url"`
	HTTPSUpgrade            *bool              `arg:"--https-upgrade"`
	PreventForeignEmbedding *bool              `arg:"--prevent-foreign-embedding"`
	HSTSIncludeSubDomains   *bool              `arg:"--hsts-include-sub-domains"`
	HSTSMaxAge              *int               `arg:"--hsts-max-age"`
	HSTSPreload             *bool              `arg:"--hsts-preload"`
}

type UpdateRuleCmd struct {
	ID            string             `arg:"positional"`
	ForwardParams *bool              `arg:"--forward-params"`
	ForwardPath   *bool              `arg:"--forward-path"`
	ResponseType  *rule.ResponseType `arg:"--response-type"`
	SourceURLs    []string           `arg:"--source-url"`
	TargetURL     *string            `arg:"--target-url"`
}

var args struct {
	APIKey    string     `arg:"env:EASYREDIR_API_KEY"`
	APISecret string     `arg:"env:EASYREDIR_API_SECRET"`
	Create    *CreateCmd `arg:"subcommand:create"`
	Get       *GetCmd    `arg:"subcommand:get"`
	List      *ListCmd   `arg:"subcommand:list"`
	Remove    *RemoveCmd `arg:"subcommand:remove"`
	Update    *UpdateCmd `arg:"subcommand:update"`
}

func main() {
	log.SetFlags(0)
	arg.MustParse(&args)

	e := easyredir.New(
		easyredir.WithAPIKey(args.APIKey),
		easyredir.WithAPISecret(args.APISecret),
	)

	switch {
	case args.Create != nil:
		switch {
		case args.Create.Rule != nil:
			r, err := e.CreateRule(rule.Attributes{
				ForwardParams: args.Create.Rule.ForwardParams,
				ForwardPath:   args.Create.Rule.ForwardPath,
				ResponseType:  args.Create.Rule.ResponseType,
				SourceURLs:    args.Create.Rule.SourceURLs,
				TargetURL:     args.Create.Rule.TargetURL,
			})
			if err != nil {
				log.Fatalf("unable to create rule: %v\n", err)
			}
			fmt.Print(r)
		}

	case args.Get != nil:
		switch {
		case args.Get.Host != nil:
			h, err := e.GetHost(args.Get.Host.ID)
			if err != nil {
				log.Fatalf("unable to get host: %v: %v\n", args.Get.Host.ID, err)
			}
			fmt.Print(h)
		}

	case args.List != nil:
		switch {
		case args.List.Host != nil:
			r, err := e.ListHosts()
			if err != nil {
				log.Fatalf("unable to list hosts: %v\n", err)
			}
			log.Print(r)

		case args.List.Rule != nil:
			r, err := e.ListRules()
			if err != nil {
				log.Fatalf("unable to list rules: %v\n", err)
			}
			log.Print(r)
		}

	case args.Remove != nil:
		switch {
		case args.Remove.Rule != nil:
			res, err := e.RemoveRule(args.Remove.Rule.ID)
			if err != nil {
				log.Fatalf("unable to remove rule: %v\n", err)
			}
			log.Printf("Result of remove rule for %v: %v\n", args.Remove.Rule.ID, res)
		}

	case args.Update != nil:
		switch {
		case args.Update.Host != nil:
			r, err := e.UpdateHost(args.Update.Host.ID, host.Attributes{
				MatchOptions: host.MatchOptions{
					CaseInsensitive:  args.Update.Host.CaseInsensitive,
					SlashInsensitive: args.Update.Host.SlashInsensitive,
				},
				NotFoundAction: host.NotFoundAction{
					ForwardParams: args.Update.Host.ForwardParams,
					ForwardPath:   args.Update.Host.ForwardPath,
					Custom404Body: args.Update.Host.Custom404Body,
					ResponseCode:  args.Update.Host.ResponseCode,
					ResponseURL:   args.Update.Host.ResponseURL,
				},
				Security: host.Security{
					HTTPSUpgrade:            args.Update.Host.HTTPSUpgrade,
					PreventForeignEmbedding: args.Update.Host.PreventForeignEmbedding,
					HSTSIncludeSubDomains:   args.Update.Host.HSTSIncludeSubDomains,
					HSTSMaxAge:              args.Update.Host.HSTSMaxAge,
					HSTSPreload:             args.Update.Host.HSTSPreload,
				},
			})
			if err != nil {
				log.Fatalf("unable to update host: %v\n", err)
			}
			fmt.Print(r)

		case args.Update.Rule != nil:
			r, err := e.UpdateRule(args.Update.Rule.ID, rule.Attributes{
				ForwardParams: args.Update.Rule.ForwardParams,
				ForwardPath:   args.Update.Rule.ForwardPath,
				ResponseType:  args.Update.Rule.ResponseType,
				SourceURLs:    args.Update.Rule.SourceURLs,
				TargetURL:     args.Update.Rule.TargetURL,
			})
			if err != nil {
				log.Fatalf("unable to update rule: %v\n", err)
			}
			fmt.Print(r)
		}
	}
}
