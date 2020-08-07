package main

func move(cfg config, from Resource, to Resource) error {
	return terraformExec(cfg, false, []string{}, "state", "mv", from.Address, to.Address)
}
