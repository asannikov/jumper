package command

import (
	"github.com/urfave/cli/v2"
)

type copyRightGlobalConfig interface {
	EnableCopyright() error
	DisableCopyright() error
}

type callCopyrightCommandDialog interface {
}

// CallCopyrightCommand runs copyright dialog
func CallCopyrightCommand(initf func(bool) string, cfg copyRightGlobalConfig, d callCopyrightCommandDialog) *cli.Command {
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
