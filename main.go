package main

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
)

var (
	version = "0.0.0"
	build   = "0"
)

func main() {
	app := cli.NewApp()
	app.Name = "rancher publish"
	app.Usage = "rancher publish"
	app.Action = run
	app.Version = fmt.Sprintf("%s+%s", version, build)

	app.Flags = []cli.Flag{

		cli.StringFlag{
			Name:   "url",
			Usage:  "url to the rancher api",
			EnvVar: "PLUGIN_URL, RANCHER_URL",
		},
		cli.StringFlag{
			Name:   "token",
			Usage:  "rancher token",
			EnvVar: "PLUGIN_TOKEN, RANCHER_TOKEN",
		},
		cli.StringFlag{
			Name:   "project",
			Usage:  "rancher project id (cxyz5:p-a1bc2)",
			EnvVar: "PLUGIN_PROJECT, RANCHER_PROJECT",
		},
		cli.StringFlag{
			Name:   "deployment",
			Usage:  "Name of the kubernetes deployment",
			EnvVar: "PLUGIN_DEPLOYMENT",
		},
		cli.StringFlag{
			Name:   "namespace",
			Usage:  "Namespace of the kubernetes deployment",
			EnvVar: "PLUGIN_NAMESPACE",
		},
		cli.StringFlag{
			Name:   "docker-image",
			Usage:  "image to use",
			EnvVar: "PLUGIN_DOCKER_IMAGE",
		},
		cli.BoolTFlag{
			Name:   "start-first",
			Usage:  "Start new container before stopping old",
			EnvVar: "PLUGIN_START_FIRST",
		},
		cli.BoolFlag{
			Name:   "confirm",
			Usage:  "auto confirm the service upgrade if successful",
			EnvVar: "PLUGIN_CONFIRM",
		},
		cli.IntFlag{
			Name:   "timeout",
			Usage:  "the maximum wait time in seconds for the service to upgrade",
			Value:  30,
			EnvVar: "PLUGIN_TIMEOUT",
		},
		cli.Int64Flag{
			Name:   "interval-millis",
			Usage:  "The interval for batch size upgrade",
			Value:  1000,
			EnvVar: "PLUGIN_INTERVAL_MILLIS",
		},
		cli.Int64Flag{
			Name:   "batch-size",
			Usage:  "The upgrade batch size",
			Value:  1,
			EnvVar: "PLUGIN_BATCH_SIZE",
		},
		cli.BoolTFlag{
			Name:   "yaml-verified",
			Usage:  "Ensure the yaml was signed",
			EnvVar: "DRONE_YAML_VERIFIED",
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	plugin := Plugin{
		URL:            c.String("url"),
		Token:          c.String("token"),
		Project:        c.String("project"),
		Deployment:     c.String("deployment"),
		Namespace:      c.String("namespace"),
		DockerImage:    c.String("docker-image"),
		StartFirst:     c.BoolT("start-first"),
		Confirm:        c.Bool("confirm"),
		Timeout:        c.Int("timeout"),
		IntervalMillis: c.Int64("interval-millis"),
		BatchSize:      c.Int64("batch-size"),
		YamlVerified:   c.BoolT("yaml-verified"),
	}
	return plugin.Exec()
}
