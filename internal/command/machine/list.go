package machine

import (
	"sort"

	"github.com/bridgex-eu/drs/internal/command"
	"github.com/dustin/go-humanize"
	"github.com/urfave/cli/v2"
)

func NewListCmd(drsCli *command.Cli) *cli.Command {
	return &cli.Command{
		Name:  "ls",
		Usage: "List machines",
		Action: func(ctx *cli.Context) error {
			return listMachines(drsCli)
		},
	}
}

func listMachines(drsCli *command.Cli) error {
	machines := drsCli.Config.Machines().All()

	sort.Slice(machines, func(i, j int) bool {
		return machines[i].CreatedAt.After(machines[j].CreatedAt)
	})

	data := [][]string{
		{
			"NAME",
			"HOST",
			"CREATED",
		},
	}

	for _, machine := range machines {
		data = append(data, []string{
			machine.Name,
			machine.Host.String(),
			humanize.Time(machine.CreatedAt),
		})
	}

	command.PrintTable(drsCli.Out, data)

	return nil
}
