package main

import (
	"reflect"
	"testing"
)

func TestSortByLevenshteinDistance(t *testing.T) {
	type args struct {
		inputDests []Resource
		inputSrc   Resource
		expect     []Resource
	}
	resSrc := Resource{Address: "module.oldname.null_resource.resource_alpha"}
	resDest1 := Resource{Address: "module.newname.null_resource.resource_alpha"}
	resDest2 := Resource{Address: "module.newname.null_resource.resource_beta"}
	resDest3 := Resource{Address: "module.oldname.null_resource.other_resourceA"}
	resDest4 := Resource{Address: "module.oldname.null_resource.other_resourceB"}

	tests := []struct {
		name string
		args args
	}{
		{name: "simple sort", args: args{inputDests: []Resource{resDest2, resDest1}, inputSrc: resSrc, expect: []Resource{resDest1, resDest2}}},
		{name: "already sorted", args: args{inputDests: []Resource{resDest1, resDest2}, inputSrc: resSrc, expect: []Resource{resDest1, resDest2}}},
		{name: "bigger test case", args: args{inputDests: []Resource{resDest3, resDest1, resDest2}, inputSrc: resSrc, expect: []Resource{resDest1, resDest2, resDest3}}},
		{name: "sorts finally alphabetically", args: args{inputDests: []Resource{resDest4, resDest3}, inputSrc: resSrc, expect: []Resource{resDest3, resDest4}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if sortByLevenshteinDistance(tt.args.inputDests, tt.args.inputSrc); !reflect.DeepEqual(tt.args.inputDests, tt.args.expect) {
				t.Errorf("sortByLevenshteinDistance() = %v, want %v", tt.args.inputDests, tt.args.expect)
			}
		})
	}
}
