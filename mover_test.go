package main

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestMove(t *testing.T) {
	dir := createDir(t)
	defer os.RemoveAll(dir)

	oldTf := `Resource "null_resource" "old" {}
Resource "null_resource" "second" {}`
	prepareState(dir, oldTf, t)

	newTf := `Resource "null_resource" "new" {}
Resource "null_resource" "second" {}`
	if err := ioutil.WriteFile(dir+"/main.tf", []byte(newTf), 0644); err != nil {
		t.Fatal(err)
	}

	move(dir, Resource{"null_resource.old", "null_resource"}, Resource{"null_resource.new", "null_resource"})

	want := []ResChange{
		{"null_resource.new", "null_resource", Change{[]changeAction{noOp}}},
		{"null_resource.second", "null_resource", Change{[]changeAction{noOp}}},
	}
	if got := changes(dir); !reflect.DeepEqual(got, want) {
		t.Errorf("changes() = %q, want %q", got, want)
	}
}
