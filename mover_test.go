package main

import (
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
	if err := os.WriteFile(dir+"/main.tf", []byte(newTf), 0644); err != nil {
		t.Fatal(err)
	}

	if err := move(config{}, Resource{"null_resource.old", "null_resource"}, Resource{"null_resource.new", "null_resource"}); err != nil {
		t.Fatal(err)
	}

	var want []ResChange
	isPre012, err := isPre012()
	if err != nil {
		t.Fatal(err)
	}
	if !isPre012 {
		want = []ResChange{
			{"null_resource.new", "null_resource", Change{[]changeAction{noOp}}},
			{"null_resource.second", "null_resource", Change{[]changeAction{noOp}}},
		}
	}

	if got, err := changes(config{}, []string{}); err != nil && !reflect.DeepEqual(got, want) {
		t.Errorf("changes() = %q, want %q", got, want)
	}
}
