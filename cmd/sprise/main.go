package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/nathan-osman/sprise/db"
	"github.com/nathan-osman/sprise/queue"
	"github.com/nathan-osman/sprise/server"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "sprise"
	app.Usage = "Photo gallery and organizer"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "db-args",
			Value:  "dbname=postgres user=postgres",
			EnvVar: "DB_ARGS",
			Usage:  "database arguments",
		},
		cli.StringFlag{
			Name:   "db-driver",
			Value:  "postgres",
			EnvVar: "DB_DRIVER",
			Usage:  "database driver",
		},
		cli.StringFlag{
			Name:   "server-addr",
			Value:  ":8000",
			EnvVar: "SERVER_ADDR",
			Usage:  "server address",
		},
	}
	app.Action = func(c *cli.Context) error {

		// Create the database connection
		conn, err := db.New(&db.Config{
			Driver: c.String("db-driver"),
			Args:   c.String("db-args"),
		})
		if err != nil {
			return err
		}
		defer conn.Close()

		// Initialize the upload queue
		q, err := queue.New(&queue.Config{
			Conn: conn,
		})
		if err != nil {
			return err
		}
		defer q.Close()

		// Initialize the server
		s, err := server.New(&server.Config{
			Addr: c.String("server-addr"),
			Conn: conn,
		})
		if err != nil {
			return err
		}
		defer s.Close()

		// Wait for SIGINT or SIGTERM
		sigChan := make(chan os.Signal)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		return nil
	}
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Fatal: %s\n", err.Error())
	}
}
