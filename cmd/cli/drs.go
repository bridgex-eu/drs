package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/bridgex-eu/drs/internal/command"
	"github.com/bridgex-eu/drs/internal/command/app"
	"github.com/bridgex-eu/drs/internal/command/docker"
	"github.com/bridgex-eu/drs/internal/command/key"
	"github.com/bridgex-eu/drs/internal/command/machine"
	"github.com/bridgex-eu/drs/internal/command/profile"
	"github.com/bridgex-eu/drs/internal/config"
	"github.com/lmittmann/tint"
	"github.com/mattn/go-colorable"
	"github.com/mattn/go-isatty"
	"github.com/urfave/cli/v2"
)

func main() {
	w := os.Stderr

	slog.SetDefault(slog.New(
		tint.NewHandler(colorable.NewColorable(w), &tint.Options{
			Level:      slog.LevelDebug,
			TimeFormat: time.Kitchen,
			NoColor:    !isatty.IsTerminal(w.Fd()),
		}),
	))

	cfg, err := config.ParseDefaultConfig()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	drsCli := &command.Cli{
		Config: cfg,
		In:     os.Stdin,
		Out:    os.Stdout,
		Err:    os.Stderr,
	}

	app := &cli.App{
		Name:    "drs",
		Usage:   "Run multiple apps on single machine.",
		Version: "v1.0",
		Commands: []*cli.Command{
			app.NewDeployCmd(drsCli),
			docker.NewDockerCmd(drsCli),
			machine.NewCmd(drsCli),
			key.NewCmd(drsCli),
			profile.NewCmd(drsCli),
		},
		Suggest:   true,
		Reader:    drsCli.In,
		Writer:    drsCli.Out,
		ErrWriter: drsCli.Err,
		ExitErrHandler: func(ctx *cli.Context, err error) {
			if err != nil {
				fmt.Fprintln(ctx.App.Writer, err)
				os.Exit(0)
			}
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
