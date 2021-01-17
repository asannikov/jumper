package command

import (
	"github.com/urfave/cli/v2"
)

type getProjectPathDialog interface {
}

type pathOptions interface {
	GetInitFuntion() func(bool) string
}

// GetProjectPath gets absolute path to the project without doing anything with docker
func GetProjectPath(d getProjectPathDialog, options pathOptions) *cli.Command {
	initf := options.GetInitFuntion()

	return &cli.Command{
		Name:        "path",
		Aliases:     []string{},
		Usage:       `Gets project path`,
		Description: ``,
		//SkipFlagParsing: true,
		Action: func(c *cli.Context) (err error) {
			initf(true)
			return nil
		},
	}
}
