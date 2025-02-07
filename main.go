package main

import (
	"context"
	"io"
	"log"
	"os"
	"time"

	"github.com/adroll/ecs-ship/clients"
	"github.com/adroll/ecs-ship/models"
	"github.com/adroll/ecs-ship/services"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/fatih/color"
	"github.com/urfave/cli/v3"
	"gopkg.in/yaml.v3"
)

func main() {
	initLogger()

	app := &cli.Command{
		Name:                   "ecs-ship",
		Usage:                  "Deploy your aws ecs services.",
		Version:                "2.0.0",
		UseShortOptionHandling: true,
		ArgsUsage:              "<cluster> <service>",
		UsageText:              "ecs-deploy [options] <cluster> <service>",
		HideHelpCommand:        true,
		Flags: []cli.Flag{
			&cli.StringFlag{
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
		Action: func(ctx context.Context, cmd *cli.Command) error {
			color.NoColor = color.NoColor || cmd.Bool("no-color")

			if cmd.NArg() != 2 {
				ec := cli.Exit(color.RedString("Please specify a cluster and service to update"), 2)
				if err := cli.ShowAppHelp(cmd); err != nil {
					log.Println(color.RedString("failed to show help: %s"), err.Error())
				}
				return ec
			}
			args := cmd.Args()
			cluster := args.Get(0)
			service := args.Get(1)

			data, err := readConfigPayload(cmd.String("updates"))
			if err != nil {
				ec := cli.Exit(color.RedString("Unable to read input file"), 3)
				return ec
			}

			var cfg models.TaskConfig
			if err := yaml.Unmarshal(data, &cfg); err != nil {
				return err
			}

			ecsConfig, err := config.LoadDefaultConfig(ctx)
			if err != nil {
				return err
			}

			ecsClient := ecs.NewFromConfig(ecsConfig)

			client := clients.NewECSClient(ecsClient)
			svc := services.NewDeployerService(client)
			return svc.Deploy(ctx, &services.DeployInput{
				Cluster:   cluster,
				Service:   service,
				NewConfig: cfg,
				DryRun:    cmd.Bool("dry"),
				Timeout:   cmd.Duration("timeout"),
				NoWait:    cmd.Bool("no-wait"),
			})
		},
	}

	err := app.Run(context.Background(), os.Args)
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
