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

	oldTf := `resource "null_resource" "old" {}
resource "null_resource" "second" {}`
	prepareState(dir, oldTf, t)

	newTf := `resource "null_resource" "new" {}
resource "null_resource" "second" {}`
	if err := ioutil.WriteFile(dir+"/main.tf", []byte(newTf), 0644); err != nil {
		t.Fatal(err)
	}

	move(dir, Resource{"null_resource.old", "null_resource"}, Resource{"null_resource.new", "null_resource"})

	want := []ResChange{
		{"null_resource.new", "null_resource", Change{[]changeAction{noOp}}},
		{"null_resource.second", "null_resource", Change{[]changeAction{noOp}}},
	}
	if got := changes(dir, []string{}); !reflect.DeepEqual(got, want) {
		t.Errorf("changes() = %q, want %q", got, want)
	}
}
