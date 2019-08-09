package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"reflect"
	"regexp"
	"strings"
)

func changes(args []string) []ResChange {
	tfPlan, err := ioutil.TempFile("", "tfplan")
	if err != nil {
		log.Panicf("ioutil.TempFile() failed with %s\n", err)
	}
	tfPlanName := tfPlan.Name()
	defer os.Remove(tfPlanName)

	if err := terraformExec(args, "plan", "-out="+tfPlanName); err != nil {
		log.Panicf("terraform plan failed with %s\n", err)
	}

	isPre012, err := isPre012()
	if err != nil {
		log.Panicf("terraform version check failed with %s\n", err)
	}
	if isPre012 {
		cmd := exec.Command("terraform", "show", "-no-color", tfPlanName)
		var stdout bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = os.Stderr
		if err = cmd.Run(); err != nil {
			log.Panicf("cmd.Run() failed with %s\n", err)
		}
		stdoutBytes := stdout.Bytes()
		var changes []ResChange
		for _, res := range regexp.MustCompile("(?m)\\+ (.*)$").FindAllSubmatch(stdoutBytes, -1) {
			address := string(res[1])
			parts := strings.Split(address, ".")
			changes = append(changes, ResChange{address, parts[len(parts)-2], Change{[]changeAction{create}}})
		}
		for _, res := range regexp.MustCompile("(?m)- (.*)$").FindAllSubmatch(stdoutBytes, -1) {
			address := string(res[1])
			parts := strings.Split(address, ".")
			changes = append(changes, ResChange{address, parts[len(parts)-2], Change{[]changeAction{del}}})
		}
		return changes
	} else {
		cmd := exec.Command("terraform", "show", "-json", tfPlanName)
		var stdout bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = os.Stderr
		if err = cmd.Run(); err != nil {
			log.Panicf("cmd.Run() failed with %s\n", err)
		}

		changes := resChanges{}

		if err := json.Unmarshal(stdout.Bytes(), &changes); err != nil {
			log.Panicf("json.Unmarshal() failed with %s\n", err)
		}
		return changes.ResChanges
	}
}

func filter(resources []ResChange, action changeAction) []Resource {
	var list []Resource
	for _, res := range resources {
		if reflect.DeepEqual(res.Change.Actions, []changeAction{action}) {
			list = append(list, Resource{res.Address, res.Type})
		}
	}
	return list
}
