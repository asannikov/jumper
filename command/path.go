package command

import (
	"github.com/urfave/cli/v2"
)

// GetProjectPath gets absolute path to the project without doing anything with docker
func GetProjectPath(initf func(), cfg projectConfig, d dialog) *cli.Command {
	return &cli.Command{
		Name:        "path",
		Aliases:     []string{},
		Usage:       `gets project path`,
		Description: ``,
		//SkipFlagParsing: true,
		Action: func(c *cli.Context) (err error) {
			initf()
			return nil
		},
	}
}
