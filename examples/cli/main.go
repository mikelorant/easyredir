package main

import (
	"fmt"
	"log"
	"os"

	"github.com/alexflint/go-arg"
	"github.com/mikelorant/easyredir/pkg/easyredir"
	"github.com/mikelorant/easyredir/pkg/easyredir/host"
	"github.com/mikelorant/easyredir/pkg/easyredir/rule"
)

type CreateCmd struct {
	Rule *struct {
		ForwardParams *bool              `arg:"--forward-params" default:"false"`
		ForwardPath   *bool              `arg:"--forward-path" default:"false"`
		ResponseType  *rule.ResponseType `arg:"--response-type", default:"moved_permanently"`
		SourceURLs    []string           `arg:"--source-url,required"`
		TargetURL     *string            `arg:"--target-url,required"`
	} `arg:"subcommand:rule"`
}

type GetCmd struct {
	Host *struct {
		ID string `arg:"positional"`
	} `arg:"subcommand:host"`
}

type ListCmd struct {
	Host *struct{} `arg:"subcommand:hosts"`
	Rule *struct {
		SourceFilter string `arg:"--source-filter"`
		TargetFilter string `arg:"--target-filter"`
	} `arg:"subcommand:rules"`
}

type RemoveCmd struct {
	Rule *struct {
		ID string `arg:"positional"`
	} `arg:"subcommand:rule"`
}

type UpdateCmd struct {
	Host *struct {
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
	} `arg:"subcommand:host"`
	Rule *struct {
		ID            string             `arg:"positional"`
		ForwardParams *bool              `arg:"--forward-params"`
		ForwardPath   *bool              `arg:"--forward-path"`
		ResponseType  *rule.ResponseType `arg:"--response-type"`
		SourceURLs    []string           `arg:"--source-url"`
		TargetURL     *string            `arg:"--target-url"`
	} `arg:"subcommand:rule"`
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
	p := arg.MustParse(&args)

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
			r, err := e.ListRules(
				easyredir.WithSourceFilter(args.List.Rule.SourceFilter),
				easyredir.WithTargetFilter(args.List.Rule.TargetFilter),
			)
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

	default:
		p.WriteHelp(os.Stdout)
	}
}
