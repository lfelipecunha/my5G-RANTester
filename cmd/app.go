package main

import (
	"my5G-RANTester/config"
	// "fmt"
	"github.com/davecgh/go-spew/spew"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"my5G-RANTester/internal/templates"
	"os"
)

const version = "0.1"

func init() {
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)
	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)
	spew.Config.Indent = "\t"

	log.Info("my5G-RANTester version " + version)

}

func execLoadTest(name string, numberUes int) {
	switch name {
	case "tnla":
		templates.TestMultiAttachUesInConcurrencyWithTNLAs(numberUes)
	case "gnb":
		templates.TestMultiAttachUesInConcurrencyWithGNBs(numberUes)
	default:
		templates.TestMultiAttachUesInQueue(numberUes)
	}
}

func main() {

	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:    "load-test",
				Aliases: []string{"load-test"},
				Usage: "\nLoad endurance stress tests.\n" +
					"Example for ues in queue: load-test -n 5 \n" +
					"Example for concurrency testing with different GNBs: load-test -g -n 10\n" +
					"Example for concurrency testing with some TNLAs: load-test -t -n 10\n",
				Flags: []cli.Flag{
					&cli.IntFlag{Name: "number-of-ues", Value: 1, Aliases: []string{"n"}},
					&cli.BoolFlag{Name: "gnb", Aliases: []string{"g"}},
					&cli.BoolFlag{Name: "tnla", Aliases: []string{"t"}},
				},
				Action: func(c *cli.Context) error {
					var execName string
					var name string
					var numUes int
					cfg := config.Data

					if c.IsSet("number-of-ues") {
						execName = "queue"
						name = "Testing multiple UEs attached in queue"
						numUes = c.Int("number-of-ues")
					} else {
						log.Info(c.Command.Usage)
						return nil
					}

					if c.Bool("tnla") {
						execName = "tnla"
						name = "Testing multiple UEs attached in concurrency with TNLAs"
					} else if c.Bool("gnb") {
						execName = "gnb"
						name = "Testing multiple UEs attached in concurrency with different GNBs"
					}
					log.Info("---------------------------------------")
					log.Info("Starting test function: ", name)
					log.Info("Number of UEs: ", numUes)
					log.Info("gNodeB control interface IP/Port: ", cfg.GNodeB.ControlIF.Ip, "/", cfg.GNodeB.ControlIF.Port)
					log.Info("gNodeB data interface IP/Port: ", cfg.GNodeB.DataIF.Ip, "/", cfg.GNodeB.DataIF.Port)
					log.Info("AMF IP/Port: ", cfg.AMF.Ip, "/", cfg.AMF.Port)
					log.Info("UPF IP/Port: ", cfg.UPF.Ip, "/", cfg.UPF.Port)
					log.Info("---------------------------------------")
					execLoadTest(execName, numUes)

					return nil
				},
			},
			{
				Name:    "stress-tests",
				Aliases: []string{"stress-test"},
				Usage: "\nLoad endurance stress tests.\n" +
					"Example for increase test: stress-tests -start 10 -step 5 -end 100 -interval 10 \n" +
					"Example for decrease test: stress-tests -start 100 -step 5 -end 10 -interval 10 \n" +
					"Example for constant test: stress-tests -start 100 -step 10 -end 100 -interval 10 \n\t\t In this case step is times of loop",
				Flags: []cli.Flag{
					&cli.IntFlag{Name: "start", Aliases: []string{"st"}},
					&cli.IntFlag{Name: "step", Aliases: []string{"sp"}},
					&cli.IntFlag{Name: "end", Aliases: []string{"e"}},
					&cli.IntFlag{Name: "interval", Aliases: []string{"i"}},
				},
				Action: func(c *cli.Context) error {
					cfg := config.Data

					if !c.IsSet("start") || !c.IsSet("step") || !c.IsSet("end") || !c.IsSet("interval") {
						log.Info(c.Command.Usage)
						return nil
					}

					log.Info("gNodeB control interface IP/Port: ", cfg.GNodeB.ControlIF.Ip, "/", cfg.GNodeB.ControlIF.Port)
					log.Info("gNodeB data interface IP/Port: ", cfg.GNodeB.DataIF.Ip, "/", cfg.GNodeB.DataIF.Port)
					log.Info("AMF IP/Port: ", cfg.AMF.Ip, "/", cfg.AMF.Port)
					log.Info("UPF IP/Port: ", cfg.UPF.Ip, "/", cfg.UPF.Port)
					log.Info("---------------------------------------")
					templates.TestMultiAttachUesLoadStress(c.Int("start"), c.Int("step"), c.Int("end"), c.Int("interval"))

					return nil
				},
			},
			{
				Name:    "ue",
				Aliases: []string{"ue"},
				Usage:   "Testing an ue attached with configuration",
				Action: func(c *cli.Context) error {
					name := "Testing an ue attached with configuration"
					cfg := config.Data

					log.Info("---------------------------------------")
					log.Info("Starting test function: ", name)
					log.Info("Number of UEs: ", 1)
					log.Info("gNodeB control interface IP/Port: ", cfg.GNodeB.ControlIF.Ip, "/", cfg.GNodeB.ControlIF.Port)
					log.Info("gNodeB data interface IP/Port: ", cfg.GNodeB.DataIF.Ip, "/", cfg.GNodeB.DataIF.Port)
					log.Info("AMF IP/Port: ", cfg.AMF.Ip, "/", cfg.AMF.Port)
					log.Info("UPF IP/Port: ", cfg.UPF.Ip, "/", cfg.UPF.Port)
					log.Info("---------------------------------------")
					templates.TestAttachUeWithConfiguration()
					return nil
				},
			},
			{
				Name:    "gnb",
				Aliases: []string{"gnb"},
				Usage: "Testing multiple GNBs attached.\n" +
					"Example for testing attached gnbs: gnb -n 5",
				Flags: []cli.Flag{
					&cli.IntFlag{Name: "number-of-gnbs", Value: 1, Aliases: []string{"n"}},
				},
				Action: func(c *cli.Context) error {
					var numGnbs int

					if c.IsSet("number-of-gnbs") {
						numGnbs = c.Int("number-of-gnbs")
					} else {
						log.Info(c.Command.Usage)
						return nil
					}

					name := "Testing multiple GNBs attached"
					cfg := config.Data

					log.Info("---------------------------------------")
					log.Info("Starting test function: ", name)
					log.Info("Number of GNBs: ", numGnbs)
					log.Info("gNodeB control interface IP/Port: ", cfg.GNodeB.ControlIF.Ip, "/", cfg.GNodeB.ControlIF.Port)
					log.Info("gNodeB data interface IP/Port: ", cfg.GNodeB.DataIF.Ip, "/", cfg.GNodeB.DataIF.Port)
					log.Info("AMF IP/Port: ", cfg.AMF.Ip, "/", cfg.AMF.Port)
					log.Info("UPF IP/Port: ", cfg.UPF.Ip, "/", cfg.UPF.Port)
					log.Info("---------------------------------------")
					templates.TestMultiAttachGnbInConcurrency(numGnbs)
					return nil
				},
			},
			{
				Name:    "nlinear-tests",
				Aliases: []string{"nlinear-tests"},
				Usage: "Testing multiple UEs attached but with samples generated by Poisson and Exponential Distribution .\n" +
					"Example for testing attached UEs with samples(s) generated by Poisson with mean(mu) : nlinear-tests -s 5 -mu 2 -se 14",
				Flags: []cli.Flag{
					&cli.IntFlag{Name: "number-of-samples", Value: 1, Aliases: []string{"s"}},
					&cli.IntFlag{Name: "mean", Value: 1, Aliases: []string{"mu"}},
					&cli.IntFlag{Name: "seed", Value: 1, Aliases: []string{"se"}},
				},
				Action: func(c *cli.Context) error {
					var numSamples int
					var mean float64
					var seed int

					numSamples = c.Int("number-of-samples")
					mean = c.Float64("mean")
					seed = c.Int("seed")

					name := "Testing multiple UE attached"
					cfg := config.Data

					log.Info("---------------------------------------")
					log.Info("Starting test function: ", name)
					log.Info("Number of Samples: ", numSamples)
					log.Info("Poisson and Exponential distribution with mean: ", mean)
					log.Info("Poisson and Exponential distribution with seed: ", seed)
					log.Info("gNodeB control interface IP/Port: ", cfg.GNodeB.ControlIF.Ip, "/", cfg.GNodeB.ControlIF.Port)
					log.Info("gNodeB data interface IP/Port: ", cfg.GNodeB.DataIF.Ip, "/", cfg.GNodeB.DataIF.Port)
					log.Info("AMF IP/Port: ", cfg.AMF.Ip, "/", cfg.AMF.Port)
					log.Info("UPF IP/Port: ", cfg.UPF.Ip, "/", cfg.UPF.Port)
					log.Info("---------------------------------------")
					templates.TestMultiAttachUesInConcurrencyWithGNBsUsingPoissonAndExponential(numSamples, mean, seed)
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
