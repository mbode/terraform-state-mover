package main

import (
	"os"
	"os/exec"
)

type resChanges struct {
	ResChanges []ResChange `json:"resource_changes"`
}

type ResChange struct {
	Address string
	Type    string
	Change  Change
}

type Change struct {
	Actions []changeAction
}

type changeAction string

const (
	noOp   changeAction = "no-op"
	create changeAction = "create"
	read   changeAction = "read"
	update changeAction = "update"
	del    changeAction = "delete"
)

type Resource struct {
	Address string
	Type    string
}

func terraformExec(dir string, args ...string) error {
	cmd := exec.Command("terraform", args...)
	cmd.Dir = dir
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
