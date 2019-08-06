package main

import "log"

func move(dir string, from Resource, to Resource) {
	log.Printf("Moving %v to %v.\n", from.Address, to.Address)

	if err := terraformExec(dir, "state", "mv", from.Address, to.Address); err != nil {
		log.Panicf("terraform state mv failed with %s\n", err)
	}
}
