package cmd

import (
	"os"

	"github.com/pkg/errors"
	"github.com/sujamess/fastgql/api"
	"github.com/sujamess/fastgql/codegen/config"
	"github.com/urfave/cli/v2"
)

var genCmd = &cli.Command{
	Name:  "generate",
	Usage: "generate a graphql server based on schema",
	Flags: []cli.Flag{
		&cli.BoolFlag{Name: "verbose, v", Usage: "show logs"},
		&cli.StringFlag{Name: "config, c", Usage: "the config filename"},
	},
	Action: func(ctx *cli.Context) error {
		var cfg *config.Config
		var err error
		if configFilename := ctx.String("config"); configFilename != "" {
			cfg, err = config.LoadConfig(configFilename)
			if err != nil {
				return err
			}
		} else {
			cfg, err = config.LoadConfigFromDefaultLocations()
			if os.IsNotExist(errors.Cause(err)) {
				cfg, err = config.LoadDefaultConfig()
			}

			if err != nil {
				return err
			}
		}

		if err = api.Generate(cfg); err != nil {
			return err
		}
		return nil
	},
}
