package main

import (
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/lncapital/torq/build"
	"github.com/lncapital/torq/migrations"
	"github.com/lncapital/torq/pkg/database"
	"github.com/lncapital/torq/pkg/lndutil"
	"github.com/lncapital/torq/server"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
	"log"
	"os"
)

func loadFlags() func(context *cli.Context) (altsrc.InputSourceContext, error) {

	return func(context *cli.Context) (altsrc.InputSourceContext, error) {
		return altsrc.NewTomlSourceFromFile(context.String("config"))
	}

}

func main() {
	app := cli.NewApp()
	app.Name = "torq"
	app.EnableBashCompletion = true
	app.Version = build.Version()

	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("error finding home directory of user: %v", err)
	}

	cmdFlags := []cli.Flag{

		// All these flags can be set though a common config file.
		&cli.StringFlag{
			Name:    "config",
			Value:   homedir + "/.torq/torq.conf",
			Aliases: []string{"c"},
			Usage:   "Path to config file",
		},

		// Torq connection details
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:  "grpc_host",
			Value: "localhost",
			Usage: "Host address for your regular grpc",
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:  "grpc_port",
			Value: "50050",
			Usage: "Port for your regular grpc",
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:  "grpc_web_port",
			Value: "50051",
			Usage: "Port for your web grpc",
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:  "cert",
			Value: "./cert.pem",
			Usage: "Path to your cert.pem file used by the GRPC server (torq)",
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:  "key",
			Value: "./key.pem",
			Usage: "Path to your key.pem file used by the GRPC server",
		}),

		// Torq database
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:  "db_name",
			Value: "torq",
			Usage: "Name of the database",
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:  "db_user",
			Usage: "Name of the postgres user with access to the database",
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:  "db_password",
			Usage: "Name of the postgres user with access to the database",
		}),

		// LND node connection details
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:    "lnd.node_address",
			Aliases: []string{"na"},
			Value:   "localhost:10009",
			Usage:   "Where to reach the lnd. Default: localhost:10009",
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:  "lnd.tls",
			Usage: "Path to your tls.cert file (LND node).",
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:  "lnd.macaroon",
			Usage: "Path to your admin.macaroon file. (LND node)",
		}),
	}

	start := &cli.Command{
		Name:  "start",
		Usage: "Starts the server, checking ",
		Action: func(c *cli.Context) error {

			// Check if the database needs to be migrated.
			err := migrations.MigrateUp()
			if err != nil && !errors.Is(err, migrate.ErrNoChange) {
				return fmt.Errorf("%", err)
			}

			conn, err := lndutil.ConnectLnd(
				c.String("lnd.node_address"),
				c.String("lnd.tls"),
				c.String("lnd.macaroon"))

			if err != nil {
				return fmt.Errorf("failed to connect to lnd: %v", err)
			}

			db, err := database.PgConnect(c.String("db_name"), c.String("db_user"), c.String("db_password"))
			if err != nil {
				return fmt.Errorf("(cmd/lnc streamHtlcCommand) error connecting to db: %v", err)
			}

			// Start the server
			fmt.Printf("Starting Torq v%s\n", build.Version())

			server.Start(conn, db)

			return nil
		},
	}

	migrate := &cli.Command{
		Name:  "migrate",
		Usage: "Migrates the database to the latest version",
		Action: func(c *cli.Context) error {
			err := migrations.MigrateUp()
			if err != nil {
				return fmt.Errorf("%v", err)
			}
			return nil
		},
	}

	migrateDown := &cli.Command{
		Name:  "migratedown",
		Usage: "Migrates the database down one step",
		Action: func(c *cli.Context) error {
			err := migrations.MigrateDown()
			if err != nil {
				return fmt.Errorf("%v", err)
			}
			return nil
		},
	}

	app.Flags = cmdFlags

	app.Before = altsrc.InitInputSourceWithContext(cmdFlags, loadFlags())

	app.Commands = cli.Commands{
		start,
		migrate,
		migrateDown,
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
