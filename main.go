// Terraform-state-mover helps refactoring terraform code by offering an interactive prompt for the `terraform state mv` command.
package main

import (
	"gopkg.in/urfave/cli.v2"
	"log"
	"os"
	"reflect"
)

var (
	version = "dev"
)

func main() {
	app := &cli.App{
		Name:      "terraform-state-mover",
		Usage:     "refactoring Terraform code has never been easier",
		Authors:   []*cli.Author{{Name: "Maximilian Bode", Email: "maxbode@gmail.com"}},
		Action:    action,
		UsageText: "terraform-state-mover [-- <terraform args>]",
		Version:   version,
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func action(c *cli.Context) error {
	var args []string
	if len(os.Args) >= 2 && os.Args[1] == "--" {
		args = os.Args[2:]
	}

	changes, err := changes(args)
	if err != nil {
		return err
	}
	dests := filterByAction(changes, create)
	srcs := filterByDestinationResourceTypes(filterByAction(changes, del), dests)

	moves := make(map[Resource]Resource)
	for len(srcs) > 0 && len(dests) > 0 {
		src, dest, err := prompt(srcs, dests)
		if err != nil {
			return err
		}
		if reflect.DeepEqual(src, Resource{}) {
			break
		}
		moves[src] = dest
		delete(srcs, src)
		delete(dests, dest)
	}

	for src, dest := range moves {
		if err := move(src, dest); err != nil {
			return err
		}
	}
	return nil
}
