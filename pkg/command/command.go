package command

import (
	"github.com/h3th-IV/climateer/pkg/runner"
	"github.com/h3th-IV/climateer/pkg/utils"
	"github.com/urfave/cli/v2"
)

func StartCommand() *cli.Command {
	var (
		startRunner = &runner.StartRunner{}
	)

	cmd := &cli.Command{
		Name:  "start",
		Usage: "starts the server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "listen-addr",
				EnvVars:     []string{"LISTEN_ADDR"},
				Usage:       "the address that the server will listen for request on",
				Destination: &startRunner.ListenAddr,
				Value:       ":8080", // TODO: check that this is correct port to serve on
			},
			&cli.StringFlag{
				Name:        "mysql-database-name",
				EnvVars:     []string{"CL_DBNAME"},
				Usage:       "Sample database name",
				Destination: &startRunner.MySQLDatabaseName,
				Value:       "",
			},
			&cli.StringFlag{
				Name:        "mysql-database-password",
				EnvVars:     []string{"CL_PASSWORD"},
				Usage:       "Sample database password",
				Destination: &startRunner.MySQLDatabasePassword,
				Value:       "",
			},
			&cli.StringFlag{
				Name:        "mysql-database-User",
				EnvVars:     []string{"CL_USER"},
				Usage:       "Sample database user",
				Destination: &startRunner.MySQLDatabaseUser,
				Value:       "",
			},
			&cli.StringFlag{
				Name:        "mysql-database-Host",
				EnvVars:     []string{"CL_HOST"},
				Usage:       "Sample database host",
				Destination: &startRunner.MySQLDatabaseHost,
				Value:       "",
			},
			&cli.StringFlag{
				Name:        "mysql-database-Port",
				EnvVars:     []string{"CL_PORT"},
				Usage:       "Sample database port",
				Destination: &startRunner.MySQLDatabasePort,
				Value:       "",
			},
			&cli.StringFlag{
				Name:        "jwt_issuer",
				EnvVars:     []string{"JWTISSUER"},
				Usage:       "Sample database port",
				Destination: &utils.JWTISSUER,
				Value:       "",
			},
			&cli.StringFlag{
				Name:        "secret",
				EnvVars:     []string{"MYSTIC"},
				Usage:       "Sample database port",
				Destination: &utils.MYSTIC,
				Value:       "",
			},
		},

		Action: startRunner.Run,
	}
	return cmd
}
