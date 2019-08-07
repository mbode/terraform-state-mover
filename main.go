// Terraform-state-mover helps refactoring terraform code by offering an interactive prompt for the `terraform state mv` command.
package main

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"log"
	"os"
	"strings"
)

func main() {
	args := os.Args[1:]

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	changes := changes(dir, args)
	srcs := filter(changes, del)

	if len(srcs) == 0 {
		fmt.Println("Found no resources to move.")
		return
	}
	allDests := filter(changes, create)

	srcTempl := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "\U0000261E {{ .Address | cyan | underline }} ({{ .Type | red }})",
		Inactive: "  {{ .Address | cyan }} ({{ .Type | red }})",
		Selected: "{{ .Address }} \U00002192",
	}
	prompt := promptui.Select{Label: "Select Source", Items: srcs, Templates: srcTempl}
	i, _, err := prompt.Run()
	if err != nil {
		log.Panicf("Prompt failed %v\n", err)
	}
	src := srcs[i]

	var dests []Resource
	for _, r := range allDests {
		if r.Type == src.Type {
			dests = append(dests, r)
		}
	}

	spaces := strings.Repeat(" ", len(src.Address)+3)
	destTempl := &promptui.SelectTemplates{
		Label:    spaces + "{{ . }}",
		Active:   spaces + "\U0000261E {{ .Address | cyan | underline }} ({{ .Type | red }})",
		Inactive: spaces + "  {{ .Address | cyan }} ({{ .Type | red }})",
		Selected: spaces + "{{ .Address }}",
	}
	prompt = promptui.Select{Label: "Select Destination", Items: dests, Templates: destTempl}
	j, _, err := prompt.Run()
	if err != nil {
		log.Panicf("Prompt failed %v\n", err)
	}
	dest := dests[j]

	move(dir, src, dest)
}
