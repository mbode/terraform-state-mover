package main

func move(from Resource, to Resource) error {
	return terraformExec([]string{}, "state", "mv", from.Address, to.Address)
}
