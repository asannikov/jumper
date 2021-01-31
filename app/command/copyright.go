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

type copyrightOptions interface {
	GetInitFunction() func(bool) string
}

// CallCopyrightCommand runs copyright dialog
func CallCopyrightCommand(cfg copyRightGlobalConfig, options copyrightOptions) *cli.Command {
	initf := options.GetInitFunction()

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
