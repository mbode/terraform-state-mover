// Terraform-state-mover helps refactoring terraform code by offering an interactive prompt for the `terraform state mv` command.
package main

import (
	"log"
	"os"
	"reflect"
)

func main() {
	args := os.Args[1:]

	changes, err := changes(args)
	if err != nil {
		log.Panicf("Detecting changes failed %v\n", err)
	}
	srcs := filter(changes, del)
	dests := filter(changes, create)

	moves := make(map[Resource]Resource)
	for len(srcs) > 0 && len(dests) > 0 {
		src, dest, err := prompt(srcs, dests)
		if err != nil {
			log.Panicf("Prompting failed %v\n", err)
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
			log.Panicf("Moving resource failed %v\n", err)
		}
	}
}
