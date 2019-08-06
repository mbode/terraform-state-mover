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

func terraformExec(dir string, args []string, extraArgs ...string) error {
	args = append(extraArgs, args...)
	cmd := exec.Command("terraform", args...)
	cmd.Dir = dir
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(),
		"TF_INPUT=false",
	)
	return cmd.Run()
}
