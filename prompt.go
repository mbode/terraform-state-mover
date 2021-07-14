package main

import (
	"fmt"
	"github.com/agnivade/levenshtein"
	"github.com/manifoldco/promptui"
	"sort"
	"strings"
)

func prompt(sources map[Resource]bool, destinations map[Resource]bool) (Resource, Resource, error) {
	srcTempl := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "\U0000261E {{ .Address | cyan | underline }} ({{ .Type | red }})",
		Inactive: "  {{ .Address | cyan }} ({{ .Type | red }})",
		Selected: "{{ .Address }}",
	}
	srcs := toSlice(sources)
	srcSearcher := func(input string, index int) bool {
		if index >= len(srcs) {
			return false
		}

		resource := srcs[index]
		return strings.Contains(resource.Address, input)
	}

	prompt := promptui.Select{Label: "Select Source", Items: append(srcs, Resource{"Finished", "no more resources to move"}), Templates: srcTempl, Searcher: srcSearcher, StartInSearchMode: true}
	i, _, err := prompt.Run()
	var empty Resource
	if err != nil {
		return empty, empty, err
	}
	if i == len(srcs) {
		return empty, empty, nil
	}
	src := srcs[i]

	dests := toSlice(destinations)
	var compatDests []Resource
	for _, r := range dests {
		if r.Type == src.Type {
			compatDests = append(compatDests, r)
		}
	}
	destSearcher := func(input string, index int) bool {
		resource := compatDests[index]
		return strings.Contains(resource.Address, input)
	}
	sortByLevenshteinDistance(compatDests, src)

	fmt.Println(strings.Repeat(" ", len(src.Address)), "↘")

	spaces := strings.Repeat(" ", len(src.Address)+3)
	destTempl := &promptui.SelectTemplates{
		Label:    spaces + "{{ . }}",
		Active:   spaces + "\U0000261E {{ .Address | cyan | underline }} ({{ .Type | red }})",
		Inactive: spaces + "  {{ .Address | cyan }} ({{ .Type | red }})",
		Selected: spaces + "{{ .Address }}",
	}
	prompt = promptui.Select{Label: "Select Destination", Items: compatDests, Templates: destTempl, Searcher: destSearcher, StartInSearchMode: true}
	j, _, err := prompt.Run()
	if err != nil {
		return Resource{}, Resource{}, err
	}
	return src, compatDests[j], nil
}

func sortByLevenshteinDistance(dests []Resource, src Resource) {
	sort.Slice(dests, func(i, j int) bool {
		distanceToItemI := levenshtein.ComputeDistance(src.Address, dests[i].Address)
		distanceToItemJ := levenshtein.ComputeDistance(src.Address, dests[j].Address)
		if distanceToItemI != distanceToItemJ {
			return distanceToItemI < distanceToItemJ
		}
		return dests[i].Address < dests[j].Address
	})
}

func toSlice(set map[Resource]bool) []Resource {
	result := []Resource{}
	for elem := range set {
		result = append(result, elem)
	}
	return result
}
