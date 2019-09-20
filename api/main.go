package main

import (
	"github.com/trinet-fasie/unity-prototype/api/v1"
	"database/sql"
	"github.com/NeowayLabs/wabbit/amqp"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/urfave/cli"
	"github.com/urfave/negroni"
	"log"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "One More World API Server"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "db-connection, db",
			Value:  "postgres://tm:Skdj18jaKppkdk17jjskh@db:5432/tm?sslmode=disable",
			Usage:  "database connection string",
			EnvVar: "OMW_DB_CONNECTION",
		},
		cli.StringFlag{
			Name:   "amqp-connection, amqp",
			Value:  "amqp://tm:Alsjqskj39aknkjls2e@rabbitmq:5672/",
			Usage:  "rabbitmq connection string",
			EnvVar: "OMW_AMQP_CONNECTION",
		},
		cli.StringFlag{
			Name:   "port, p",
			Value:  "3000",
			Usage:  "api port",
			EnvVar: "OMW_PORT",
		},
	}

	app.Action = func(c *cli.Context) error {
		rabbitMq, err := amqp.Dial(c.GlobalString("amqp-connection"))
		if err != nil {
			log.Fatalf("Failed to connect to message broker: %s", err)
		}
		defer rabbitMq.Close()

		db, err := sql.Open("postgres", c.GlobalString("db-connection"))
		if err != nil {
			log.Fatalf("Failed to connect to database: %s", err)
		}
		defer db.Close()

		r := mux.NewRouter()
		v1.New(r, rabbitMq, db)

		n := negroni.Classic()
		n.UseHandler(r)
		n.Run(":" + c.GlobalString("port"))

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
