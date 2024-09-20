package main

import (
	"io"
	"log"
	"os"
	"time"

	"github.com/adroll/ecs-ship/action"
	"github.com/adroll/ecs-ship/ecs"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
)

func main() {
	initLogger()

	app := &cli.App{
		Name:                   "ecs-ship",
		Usage:                  "Deploy your aws ecs services.",
		Version:                "1.2.0",
		UseShortOptionHandling: true,
		ArgsUsage:              "<cluster> <service>",
		UsageText:              "ecs-deploy [options] <cluster> <service>",
		HideHelpCommand:        true,
		Flags: []cli.Flag{
			&cli.PathFlag{
				Name:        "updates",
				Aliases:     []string{"u"},
				Usage:       "Use an input `FILE` to describe service updates",
				Value:       "-",
				DefaultText: "stdin",
			},
			&cli.DurationFlag{
				Name:     "timeout",
				Aliases:  []string{"t"},
				Usage:    "Wait this `DURATION` for the service to be correctly updated",
				Value:    time.Minute * 5,
				Required: false,
			},
			&cli.BoolFlag{
				Name:     "no-color",
				Aliases:  []string{"n"},
				Usage:    "Disable colored output",
				Required: false,
			},
			&cli.BoolFlag{
				Name:     "no-wait",
				Aliases:  []string{"w"},
				Usage:    "Disable waiting for updates to be completed.",
				Required: false,
			},
			&cli.BoolFlag{
				Name:     "dry",
				Aliases:  []string{"d"},
				Usage:    "Don't deploy just show what would change in the remote service",
				Required: false,
			},
		},
		Action: func(ctx *cli.Context) error {
			color.NoColor = color.NoColor || ctx.Bool("no-color")

			if ctx.NArg() != 2 {
				ec := cli.Exit(color.RedString("Please specify a cluster and service to update"), 2)
				if err := cli.ShowAppHelp(ctx); err != nil {
					log.Println(color.RedString("failed to show help: %s"), err.Error())
				}
				return ec
			}
			args := ctx.Args()
			cluster := args.Get(0)
			service := args.Get(1)

			data, err := readConfigPayload(ctx.Path("updates"))
			if err != nil {
				ec := cli.Exit(color.RedString("Unable to read input file"), 3)
				return ec
			}

			var cfg ecs.TaskConfig
			if err := yaml.Unmarshal(data, &cfg); err != nil {
				return err
			}

			client, err := ecs.BuildDefaultClient(ctx.Context)
			if err != nil {
				return err
			}
			return action.ECSDeploy(ctx.Context, cluster, service, client, ctx.Duration("timeout"), &cfg, ctx.Bool("dry"), ctx.Bool("no-wait"))
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(color.RedString("%s", err))
	}
}

func initLogger() {
	log.SetFlags(log.Lmsgprefix)
	log.SetPrefix("")
}

func readConfigPayload(inputName string) ([]byte, error) {
	if inputName == "-" {
		return io.ReadAll(os.Stdin)
	}
	return os.ReadFile(inputName)
}
