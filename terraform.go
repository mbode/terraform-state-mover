package main

import (
	"os"
	"os/exec"
)

type resChanges struct {
	ResChanges []ResChange `json:"resource_changes"`
}

// ResChange represents a resource change in a Terraform plan.
type ResChange struct {
	Address string
	Type    string
	Change  Change
}

// Change represents a list of actions to one resource in a Terraform plan.
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

// Resource represents a Terraform resource and consists of a type and an address.
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
