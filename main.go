// Terraform-state-mover helps refactoring terraform code by offering an interactive prompt for the `terraform state mv` command.
package main

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"reflect"
	"time"
)

var (
	version = "dev"
)

type config struct {
	delay   time.Duration
	verbose bool
	dryrun  bool
}

func main() {

	// do not use "-v" to print the version
	cli.VersionFlag = &cli.BoolFlag{
		Name: "version", Aliases: []string{},
		Usage: "print the version only",
	}

	app := &cli.App{
		Name:    "terraform-state-mover",
		Usage:   "refactoring Terraform code has never been easier",
		Authors: []*cli.Author{{Name: "Maximilian Bode", Email: "maxbode@gmail.com"}},
		Action:  action,
		Flags: []cli.Flag{
			&cli.DurationFlag{
				Name: "delay", Aliases: []string{"d"},
				Usage: "Delay between terraform state mv calls. Helps to avoid rate-limits.",
				Value: time.Second * 0,
			},
			&cli.BoolFlag{
				Name: "verbose", Aliases: []string{"v"},
				Usage: "Be more verbose - prints e.g. terraform mv calls",
				Value: false,
			},
			&cli.BoolFlag{
				Name: "dry-run", Aliases: []string{"n"},
				Usage: "Do not actually move state, enables -v",
				Value: false,
			},
		},
		UsageText: "terraform-state-mover [-v] [-d delay] [-n] [-- <terraform args>]",
		Version:   version,
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func action(ctx *cli.Context) error {
	var planArgs []string
	for i, elem := range os.Args {
		if "--" == elem {
			planArgs = os.Args[i+1:]
		}
	}
	cfg := readConfig(ctx)

	changes, err := changes(cfg, planArgs)
	if err != nil {
		return err
	}
	dests := filterByAction(changes, create)
	srcs := filterByDestinationResourceTypes(filterByAction(changes, del), dests)

	moves := make(map[Resource]Resource)
	for len(srcs) > 0 && len(dests) > 0 {
		src, dest, err := prompt(srcs, dests)
		if err != nil {
			if err == promptui.ErrInterrupt && len(moves) > 0 {
				fmt.Println("Interrupted. These moves would have been executed based on your selections:")
				for src, dest := range moves {
					fmt.Printf("  terraform state mv '%s' '%s'\n", src.Address, dest.Address)
				}
			}
			return err
		}
		if reflect.DeepEqual(src, Resource{}) {
			break
		}
		moves[src] = dest
		delete(srcs, src)
		delete(dests, dest)
	}

	if len(moves) == 0 {
		fmt.Println("Nothing to do.")
	}

	var firstEntry = true
	for src, dest := range moves {
		if firstEntry {
			firstEntry = false
		} else {
			wait(cfg)
		}
		if err := move(cfg, src, dest); err != nil {
			return err
		}
	}
	return nil
}

func readConfig(ctx *cli.Context) config {
	return config{
		delay:   ctx.Duration("delay"),
		verbose: ctx.Bool("verbose") || ctx.Bool("dry-run"),
		dryrun:  ctx.Bool("dry-run"),
	}
}

func wait(cfg config) {
	if cfg.verbose && cfg.delay > 0 {
		fmt.Println("Waiting", cfg.delay, "...")
	}
	time.Sleep(cfg.delay)
}
