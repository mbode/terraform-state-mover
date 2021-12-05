package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"reflect"
	"regexp"
	"strings"
)

func changes(cfg config, planArgs []string) ([]ResChange, error) {
	tfPlan, err := ioutil.TempFile("", "tfplan")
	if err != nil {
		return nil, err
	}
	tfPlanName := tfPlan.Name()
	defer os.Remove(tfPlanName)

	if err := terraformExec(cfg, true, planArgs, "plan", "-out="+tfPlanName); err != nil {
		return nil, err
	}

	isPre012, err := isPre012()
	if err != nil {
		return nil, err
	}
	if isPre012 {
		cmd := exec.Command("terraform", "show", "-no-color", tfPlanName)
		var stdout bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = os.Stderr
		if err = cmd.Run(); err != nil {
			return nil, err
		}
		stdoutBytes := stdout.Bytes()
		var changes []ResChange
		for _, res := range regexp.MustCompile(`(?m)\+ (.*)$`).FindAllSubmatch(stdoutBytes, -1) {
			address := string(res[1])
			parts := strings.Split(address, ".")
			changes = append(changes, ResChange{address, parts[len(parts)-2], Change{[]changeAction{create}}})
		}
		for _, res := range regexp.MustCompile(`(?m)- (.*)$`).FindAllSubmatch(stdoutBytes, -1) {
			address := string(res[1])
			parts := strings.Split(address, ".")
			changes = append(changes, ResChange{address, parts[len(parts)-2], Change{[]changeAction{del}}})
		}
		return changes, nil
	}
	cmd := exec.Command("terraform", "show", "-json", tfPlanName)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = os.Stderr
	if err = cmd.Run(); err != nil {
		return nil, err
	}

	changes := resChanges{}

	if err := json.Unmarshal(stdout.Bytes(), &changes); err != nil {
		return nil, err
	}
	return changes.ResChanges, nil
}

func filterByAction(resources []ResChange, action changeAction) map[Resource]bool {
	set := make(map[Resource]bool)
	for _, res := range resources {
		if reflect.DeepEqual(res.Change.Actions, []changeAction{action}) {
			set[Resource{res.Address, res.Type}] = true
		}
	}
	return set
}

func filterByDestinationResourceTypes(sourceResources map[Resource]bool, destResources map[Resource]bool) map[Resource]bool {
	set := make(map[Resource]bool)
	types := make(map[string]bool)
	for res := range destResources {
		types[res.Type] = true
	}

	for res := range sourceResources {
		if types[res.Type] || res.Address == FinishedAddress {
			set[Resource{res.Address, res.Type}] = true
		}
	}
	return set
}
