package main

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	cli "gopkg.in/urfave/cli.v1"
	"log"
	"os"
	"path/filepath"
	"time"
	//grpc lib
	pb "google.golang.org/grpc/health/grpc_health_v1"
)

const (
	VERSION = "1.0.1"
	USAGE   = "grpc health check client"
)

var app *cli.App

func init() {
	app = cli.NewApp()
	app.Name = filepath.Base(os.Args[0])
	app.Version = VERSION
	app.Usage = USAGE
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "address, a", Usage: "server addr"},
		cli.StringFlag{Name: "service, s", Usage: "service", Value: "NULL"},
	}
	app.Action = func(ctx *cli.Context) error {
		a := ctx.GlobalString("address")
		s := ctx.GlobalString("service")
		if a == "" {
			log.Fatalln("Missing address parameter! see --help")
			return errors.New("Missing address parameter! see --help")
		}
		conn, err := grpc.Dial(a, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("did not connect: %v", err)
			return err
		}
		defer conn.Close()
		f := pb.NewHealthClient(conn)
		c, cancel := context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()
		r, err := f.Check(c, &pb.HealthCheckRequest{
			Service: s,
		})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
			return err
		}
		log.Println(r)
		return nil
	}
}
func main() {
	if err := app.Run(os.Args); err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
