package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"reflect"
)

func changes(dir string) []ResChange {
	tfPlan, err := ioutil.TempFile(dir, "tfplan")
	if err != nil {
		log.Panicf("ioutil.TempFile() failed with %s\n", err)
	}
	tfPlanName := tfPlan.Name()
	defer os.Remove(tfPlanName)

	if err := terraformExec(dir, "plan", "-out="+tfPlanName); err != nil {
		log.Panicf("terraform plan failed with %s\n", err)
	}

	cmd := exec.Command("terraform", "show", "-json", tfPlanName)
	cmd.Dir = dir
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

func filter(resources []ResChange, action changeAction) []Resource {
	var list []Resource
	for _, res := range resources {
		if reflect.DeepEqual(res.Change.Actions, []changeAction{action}) {
			list = append(list, Resource{res.Address, res.Type})
		}
	}
	return list
}
