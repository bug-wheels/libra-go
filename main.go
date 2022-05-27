package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"libra-go/config"
	"libra-go/discovery"
	"libra-go/monitor"
	"log"
	"os"
	"time"
)

func main() {
	log.SetFlags(log.Lmicroseconds | log.Ldate | log.Lshortfile)
	app := &cli.App{
		Name:  "libra",
		Usage: "services monitoring",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "ui",
				Usage: "with ui",
			},
			&cli.BoolFlag{
				Name:  "loop",
				Usage: "是否循环执行",
				Value: true,
			},
			&cli.Int64Flag{
				Name:  "intervalMS",
				Value: 10_000,
				Usage: "循环间隔时间，毫秒",
			},
			&cli.Int64Flag{
				Name:    "loopcount",
				Aliases: []string{"lc"},
				Value:   0,
				Usage:   "循环次数，小于 1 则一直执行",
			},
			&cli.StringFlag{
				Name:    "config",
				Value:   "libra.yaml",
				Aliases: []string{"c"},
				Usage:   "Load configuration from `FILE`",
			},
		},
		Action: func(c *cli.Context) error {
			isNeedLoop := c.Bool("loop")
			loopcount := c.Int("loopcount")
			libraConfigFilePath := c.String("config")
			intervalMS := c.Int64("intervalMS")
			configFile, err := ioutil.ReadFile(libraConfigFilePath)
			if err != nil {
				fmt.Print(err)
			}

			//yaml文件内容影射到结构体中
			var libraConfig config.LibraConfig

			err1 := yaml.Unmarshal(configFile, &libraConfig)
			if err1 != nil {
				fmt.Println("error", err1)
			}
			fmt.Printf("config.app: %#v\n", libraConfig)

			discovery.Setup(libraConfig.RegistrationCenters)

			if len(libraConfig.Services) == 0 {
				return nil
			}

			if isNeedLoop {
				foreverLoop := loopcount < 1
				for i := 0; foreverLoop || i <= loopcount; i++ {
					monitor.MonitoringService(libraConfig.Services)
					time.Sleep(time.Duration(intervalMS) * time.Millisecond)
				}
			} else {
				monitor.MonitoringService(libraConfig.Services)
			}
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
