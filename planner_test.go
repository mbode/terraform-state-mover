package main

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestCreate(t *testing.T) {
	dir := createDir(t)
	defer os.RemoveAll(dir)

	content := `resource "null_resource" "first" {}
resource "null_resource" "second" {}`
	if err := ioutil.WriteFile(dir+"/main.tf", []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	if err := terraformExec([]string{}, "init"); err != nil {
		t.Fatalf("terraform init failed with %s\n", err)
	}

	want := []ResChange{
		{"null_resource.first", "null_resource", Change{[]changeAction{create}}},
		{"null_resource.second", "null_resource", Change{[]changeAction{create}}},
	}
	got, err := changes([]string{})
	if err != nil {
		t.Fatalf("failed computing changes")
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("changes() = %q, want %q", got, want)
	}
}

func TestDelete(t *testing.T) {
	dir := createDir(t)
	defer os.RemoveAll(dir)

	content := `resource "null_resource" "first" {}
resource "null_resource" "second" {}`
	prepareState(dir, content, t)

	if err := ioutil.WriteFile(dir+"/main.tf", []byte("\n"), 0644); err != nil {
		t.Fatal(err)
	}

	want := []ResChange{
		{"null_resource.first", "null_resource", Change{[]changeAction{del}}},
		{"null_resource.second", "null_resource", Change{[]changeAction{del}}},
	}
	got, err := changes([]string{})
	if err != nil {
		t.Fatalf("failed computing changes")
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("changes() = %q, want %q", got, want)
	}
}

func TestNoOp(t *testing.T) {
	dir := createDir(t)
	defer os.RemoveAll(dir)

	content := `resource "null_resource" "first" {}
resource "null_resource" "second" {}`
	prepareState(dir, content, t)

	if err := ioutil.WriteFile(dir+"/main.tf", []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	var want []ResChange
	isPre012, err := isPre012()
	if err != nil {
		t.Fatal(err)
	}
	if !isPre012 {
		want = []ResChange{
			{"null_resource.first", "null_resource", Change{[]changeAction{noOp}}},
			{"null_resource.second", "null_resource", Change{[]changeAction{noOp}}},
		}
	}

	got, err := changes([]string{})
	if err != nil {
		t.Fatalf("failed computing changes")
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("changes() = %q, want %q", got, want)
	}
}

func TestFilter(t *testing.T) {
	resources := []ResChange{
		{"null_resource.create", "null_resource", Change{[]changeAction{create}}},
		{"null_resource.delete", "null_resource", Change{[]changeAction{del}}},
		{"null_resource.noop", "null_resource", Change{[]changeAction{noOp}}},
		{"null_resource.update", "null_resource", Change{[]changeAction{update}}},
	}

	want := make(map[Resource]bool)
	want[Resource{"null_resource.create", "null_resource"}] = true

	if got := filter(resources, create); !reflect.DeepEqual(got, want) {
		t.Errorf("changes() = %v, want %v", got, want)
	}

	want = make(map[Resource]bool)
	want[Resource{"null_resource.delete", "null_resource"}] = true
	if got := filter(resources, del); !reflect.DeepEqual(got, want) {
		t.Errorf("changes() = %v, want %v", got, want)
	}
}

func createDir(t *testing.T) string {
	dir, err := ioutil.TempDir("", t.Name())
	if err != nil {
		t.Fatal(err)
	}
	err = os.Chdir(dir)
	if err != nil {
		t.Fatal(err)
	}
	return dir
}

func prepareState(dir string, content string, t *testing.T) {
	if err := ioutil.WriteFile(dir+"/main.tf", []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	if err := terraformExec([]string{}, "init"); err != nil {
		t.Fatal(err)
	}
	if err := terraformExec([]string{}, "apply", "-auto-approve"); err != nil {
		t.Fatal(err)
	}
}
