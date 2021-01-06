package command

import (
	"github.com/urfave/cli/v2"
)

// GetProjectPath gets absolute path to the project without doing anything with docker
func GetProjectPath(initf func(bool) string, d dialog) *cli.Command {
	return &cli.Command{
		Name:        "path",
		Aliases:     []string{},
		Usage:       `gets project path`,
		Description: ``,
		//SkipFlagParsing: true,
		Action: func(c *cli.Context) (err error) {
			initf(true)
			return nil
		},
	}
}
