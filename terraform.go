package main

import (
	"bytes"
	"github.com/hashicorp/go-version"
	"os"
	"os/exec"
	"regexp"
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

func terraformExec(args []string, extraArgs ...string) error {
	args = append(extraArgs, args...)
	cmd := exec.Command("terraform", args...)
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(),
		"TF_INPUT=false",
	)
	return cmd.Run()
}

func isPre012() (bool, error) {
	cmd := exec.Command("terraform", "version")
	cmdOutput := &bytes.Buffer{}
	cmd.Stdout = cmdOutput
	err := cmd.Run()
	if err != nil {
		return false, err
	}
	output := cmdOutput.Bytes()
	var ver = regexp.MustCompile(`Terraform v(\d+\.\d+\.\d+)`)
	result := ver.FindStringSubmatch(string(output))
	v012, err := version.NewVersion("0.12")
	if err != nil {
		return false, err
	}
	current, err := version.NewVersion(result[1])
	if err != nil {
		return false, err
	}
	return current.LessThan(v012), nil
}
