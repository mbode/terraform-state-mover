package main

import (
	"os"
	"reflect"
	"testing"
)

func TestCreate(t *testing.T) {
	dir := createDir(t)
	defer os.RemoveAll(dir)

	content := `resource "null_resource" "first" {}
resource "null_resource" "second" {}`
	if err := os.WriteFile(dir+"/main.tf", []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	if err := terraformExec(config{}, true, []string{}, "init"); err != nil {
		t.Fatalf("terraform init failed with %s\n", err)
	}

	want := []ResChange{
		{"null_resource.first", "null_resource", Change{[]changeAction{create}}},
		{"null_resource.second", "null_resource", Change{[]changeAction{create}}},
	}
	got, err := changes(config{}, []string{})
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

	if err := os.WriteFile(dir+"/main.tf", []byte("\n"), 0644); err != nil {
		t.Fatal(err)
	}

	want := []ResChange{
		{"null_resource.first", "null_resource", Change{[]changeAction{del}}},
		{"null_resource.second", "null_resource", Change{[]changeAction{del}}},
	}
	got, err := changes(config{}, []string{})
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

	if err := os.WriteFile(dir+"/main.tf", []byte(content), 0644); err != nil {
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

	got, err := changes(config{}, []string{})
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

	if got := filterByAction(resources, create); !reflect.DeepEqual(got, want) {
		t.Errorf("changes() = %v, want %v", got, want)
	}

	want = make(map[Resource]bool)
	want[Resource{"null_resource.delete", "null_resource"}] = true
	if got := filterByAction(resources, del); !reflect.DeepEqual(got, want) {
		t.Errorf("changes() = %v, want %v", got, want)
	}
}

func createDir(t *testing.T) string {
	dir, err := os.MkdirTemp("", t.Name())
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
	if err := os.WriteFile(dir+"/main.tf", []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	if err := terraformExec(config{}, true, []string{}, "init"); err != nil {
		t.Fatal(err)
	}
	if err := terraformExec(config{}, true, []string{}, "apply", "-auto-approve"); err != nil {
		t.Fatal(err)
	}
}

func TestFilterByDestinationResourceTypes(t *testing.T) {
	resSrc1 := Resource{Address: "null_resource.resource_alpha", Type: "null_resource"}
	resSrc2 := Resource{Address: "null_resource.resource_beta", Type: "another_type"}
	resDest1 := Resource{Address: "null_resource.resource_gamma", Type: "null_resource"}
	resDest2 := Resource{Address: "null_resource.resource_delta", Type: "new_type"}

	sourceResources := make(map[Resource]bool)
	sourceResources[resSrc1] = true
	sourceResources[resSrc2] = true
	destResources := make(map[Resource]bool)
	destResources[resDest1] = true
	destResources[resDest2] = true

	want := make(map[Resource]bool)
	want[resSrc1] = true

	if got := filterByDestinationResourceTypes(sourceResources, destResources); !reflect.DeepEqual(got, want) {
		t.Errorf("filterByDestinationResourceTypes() = %v, want %v", got, want)
	}
}
