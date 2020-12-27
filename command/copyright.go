package command

import (
	"github.com/urfave/cli/v2"
)

type copyRightGlobalConfig interface {
	EnableCopyright() error
	DisableCopyright() error
}

// CallCopyrightCommand runs copyright dialog
func CallCopyrightCommand(initf func(bool), cfg copyRightGlobalConfig, d dialog) *cli.Command {
	return &cli.Command{
		Name: "copyright",
		Subcommands: []*cli.Command{
			{
				Name: "enable",
				Action: func(c *cli.Context) error {
					initf(false)
					return cfg.EnableCopyright()
				},
			},
			{
				Name: "disable",
				Action: func(c *cli.Context) error {
					initf(false)
					return cfg.DisableCopyright()
				},
			},
		},
	}
}
