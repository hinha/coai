package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/urfave/cli/v2"

	"github.com/hinha/coai/cmd/server"
)

var topLevelCommands []*cli.Command

func init() {
	topLevelCommands = append(topLevelCommands, server.MakeServerCmd())
}

func main() {
	var app cli.App
	app.Name = "coai"
	app.Commands = topLevelCommands

	sort.Slice(app.Commands, func(i, j int) bool {
		return strings.Compare(app.Commands[i].Name, app.Commands[j].Name) < 0
	})

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
