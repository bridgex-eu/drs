package key

import (
	"sort"

	"github.com/bridgex-eu/drs/internal/command"
	"github.com/dustin/go-humanize"
	"github.com/urfave/cli/v2"
)

func NewListCmd(drsCli *command.Cli) *cli.Command {
	return &cli.Command{
		Name:  "ls",
		Usage: "List all SSH keys",
		Action: func(ctx *cli.Context) error {
			return listKeys(drsCli)
		},
	}
}

func listKeys(drsCli *command.Cli) error {
	keys := drsCli.Config.Keys().All()

	sort.Slice(keys, func(i, j int) bool {
		return keys[i].CreatedAt.After(keys[j].CreatedAt)
	})

	data := [][]string{
		{
			"NAME",
			"CREATED",
		},
	}

	for _, machine := range keys {
		data = append(data, []string{
			machine.Name,
			humanize.Time(machine.CreatedAt),
		})
	}

	command.PrintTable(drsCli.Out, data)

	return nil
}
