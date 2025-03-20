package main

import (
	"github.com/savageking-io/nec-lib/conf"
	"github.com/savageking-io/nec-user/user_client"
	"time"

	//user_client "github.com/savageking-io/nec-user/user_client"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "necrest"
	app.Version = AppVersion
	app.Description = "Smart backend service for smart game developers"
	app.Usage = "REST Microservice of NoErrorCode ecosystem"

	app.Authors = []cli.Author{
		{
			Name:  "savageking.io",
			Email: "i@savageking.io",
		},
		{
			Name:  "Mike Savochkin (crioto)",
			Email: "mike@crioto.com",
		},
	}

	app.Copyright = "2025 (c) savageking.io. All Rights Reserved"

	app.Commands = []cli.Command{
		{
			Name:  "serve",
			Usage: "Start REST",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "config",
					Usage:       "Configuration filepath",
					Value:       ConfigFilepath,
					Destination: &ConfigFilepath,
				},
				cli.StringFlag{
					Name:        "log",
					Usage:       "Specify logging level",
					Value:       LogLevel,
					Destination: &LogLevel,
				},
			},
			Action: Serve,
		},
	}

	_ = app.Run(os.Args)
}
func Serve(c *cli.Context) error {
	if err := SetLogLevel(LogLevel); err != nil {
		log.Errorf("Failed to set log level to %s: %s. Using INFO", LogLevel, err)
	}

	log.Infof("Starting REST service")
	dir, file, err := conf.ExtractDirectoryAndFilenameFromPath(ConfigFilepath)
	if err != nil {
		log.Error("Bad configuration file: %s", err.Error())
		return err
	}

	config := new(conf.Config)
	if err := config.Init(dir); err != nil {
		log.Errorf("Unrecoverable error: %s", err.Error())
		return err
	}

	restConfig := new(Config)
	fs := os.DirFS(dir)
	if err := config.ReadConfig(fs, file, restConfig); err != nil {
		log.Errorf("Failed to read REST configuration: %s", err.Error())
		return err
	}

	userService := new(user_client.UserClient)
	log.Infof("Connecting to user service at %s:%d", restConfig.UserService.Hostname, restConfig.UserService.Port)
	for {
		if err := userService.Connect(restConfig.UserService.Hostname, uint32(restConfig.UserService.Port)); err != nil {
			log.Errorf("Failed to connect to user service: %s", err.Error())
			wait(3 * time.Second)
			continue
		}
		log.Infof("Connected to user service")
		break
	}

	service := new(REST)
	if err := service.Init(restConfig); err != nil {
		log.Errorf("Failed to init REST service: %s", err.Error())
		return err
	}

	return service.Start()
}

func SetLogLevel(level string) error {
	lvl, err := log.ParseLevel(level)
	if err != nil {
		log.SetLevel(log.InfoLevel)
		return err
	}
	log.SetLevel(lvl)
	return nil
}

func wait(seconds time.Duration) {
	started := time.Now()
	for time.Since(started) < seconds {
		time.Sleep(100 * time.Millisecond)
	}
}
