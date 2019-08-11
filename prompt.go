package main

import (
	"fmt"
	"github.com/manifoldco/promptui"
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
	prompt := promptui.Select{Label: "Select Source", Items: append(srcs, Resource{"Finished", "no more resources to move"}), Templates: srcTempl}
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

	fmt.Println(strings.Repeat(" ", len(src.Address)), "↘")

	spaces := strings.Repeat(" ", len(src.Address)+3)
	destTempl := &promptui.SelectTemplates{
		Label:    spaces + "{{ . }}",
		Active:   spaces + "\U0000261E {{ .Address | cyan | underline }} ({{ .Type | red }})",
		Inactive: spaces + "  {{ .Address | cyan }} ({{ .Type | red }})",
		Selected: spaces + "{{ .Address }}",
	}
	prompt = promptui.Select{Label: "Select Destination", Items: compatDests, Templates: destTempl}
	j, _, err := prompt.Run()
	if err != nil {
		return Resource{}, Resource{}, err
	}
	return src, compatDests[j], nil
}

func toSlice(set map[Resource]bool) []Resource {
	result := []Resource{}
	for elem := range set {
		result = append(result, elem)
	}
	return result
}
