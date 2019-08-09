package main

import (
	"log"
)

func move(from Resource, to Resource) {
	if err := terraformExec([]string{}, "state", "mv", from.Address, to.Address); err != nil {
		log.Panicf("terraform state mv failed with %s\n", err)
	}
}
