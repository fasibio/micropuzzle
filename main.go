package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/urfave/cli/v2"
)

const (
	CliFallbackLoader = "fallbackloader"
	CliMicrofrontends = "microfrontends"
	CliPublicFolder   = "publicfoder"
	CliTimeout        = "timeoutms"
	CliLogLevel       = "logLevel"
	CliPort           = "port"
	CliManagementPort = "managementport"
	CliRedisAddress   = "redisaddr"
	CliRedisUser      = "redisuser"
	CliRedisPassword  = "redispassword"
	CliRedisDb        = "redisdb"
)

const (
	CliGenDestination = "destination"
	CliGenUrl         = "url"
)

func getFlagEnvByFlagName(flagName string) string {
	return strings.ToUpper(flagName)
}

func main() {
	app := cli.NewApp()
	app.Name = "Micropuzzle"
	runner := NewRunner()
	app.Commands = []*cli.Command{
		{
			Name:    "generateType",
			Aliases: []string{"gen"},
			Usage:   "Generate Typescript types",
			Action:  runner.GenerateType,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    CliGenDestination,
					EnvVars: []string{getFlagEnvByFlagName(CliGenDestination)},
					Usage:   "Destination folder for generated files",
					Aliases: []string{"dest"},
					Value:   "./micropuzzle",
				},
				&cli.StringFlag{
					Name:    CliGenUrl,
					EnvVars: []string{getFlagEnvByFlagName(CliGenUrl)},
					Usage:   "URL to the microfrontends server management port. f.e. http://localhost:3301",
				},
			},
		},
	}
	app.Description = "Application to combine Server Side Include and Afterloading"
	app.Action = runner.Run
	app.Flags = []cli.Flag{
		&cli.DurationFlag{
			Name:    CliTimeout,
			EnvVars: []string{getFlagEnvByFlagName(CliTimeout)},
			Usage:   "Timeout for loading Microfrontends (for all slower, it will be use streaming to bring it to the client)",
			Value:   45 * time.Millisecond,
		},
		&cli.StringFlag{
			Name:    CliLogLevel,
			EnvVars: []string{getFlagEnvByFlagName(CliLogLevel)},
			Usage:   "Loglevel debug, info, warn, error",
			Value:   "info",
		},
		&cli.StringFlag{
			Name:    CliPort,
			EnvVars: []string{getFlagEnvByFlagName(CliPort)},
			Usage:   "port where server will be started",
			Value:   "3300",
		},
		&cli.StringFlag{
			Name:    CliPublicFolder,
			EnvVars: []string{getFlagEnvByFlagName(CliPublicFolder)},
			Value:   "./public",
			Usage:   "Folder where all html js css from server directly will be foundable (Public folder for the web)",
		},
		&cli.StringFlag{
			Name:    CliMicrofrontends,
			EnvVars: []string{getFlagEnvByFlagName(CliMicrofrontends)},
			Value:   "./config/frontends.yaml",
			Usage:   "A yaml file to describe available Frontends",
		},
		&cli.StringFlag{
			Name:    CliFallbackLoader,
			EnvVars: []string{getFlagEnvByFlagName(CliFallbackLoader)},
			Usage:   "key of inifile where to find fallbackhtml which will shown if microfrontend is lower than timeout",
			Value:   "fallback",
		},
		&cli.StringFlag{
			Name:    CliRedisAddress,
			EnvVars: []string{getFlagEnvByFlagName(CliRedisAddress)},
			Usage:   "The domian/ip:port of redis",
			Value:   "localhost:6379",
		},
		&cli.StringFlag{
			Name:    CliRedisUser,
			EnvVars: []string{getFlagEnvByFlagName(CliRedisUser)},
			Usage:   "Username to connect to redis",
			Value:   "",
		},
		&cli.StringFlag{
			Name:    CliRedisPassword,
			EnvVars: []string{getFlagEnvByFlagName(CliRedisPassword)},
			Usage:   "Password to connect to redis",
			Value:   "",
		},
		&cli.Int64Flag{
			Name:    CliRedisDb,
			EnvVars: []string{getFlagEnvByFlagName(CliRedisDb)},
			Usage:   "Db to use by redis",
			Value:   0,
		},
		&cli.Int64Flag{
			Name:    CliManagementPort,
			EnvVars: []string{getFlagEnvByFlagName(CliManagementPort)},
			Usage:   "Port to get data not needed from client",
			Value:   3301,
		},
	}
	if err := app.Run(os.Args); err != nil {
		fmt.Println("Error: ", err)
	}
}
