package main

import (
	"embed"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
	"html/template"
	"io/fs"
	"io/ioutil"
	"libra-go/config"
	"libra-go/discovery"
	"libra-go/monitor"
	"libra-go/rocketmq"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

//go:embed templates/*.html
var templates embed.FS

//go:embed templates/assets/*
var assets embed.FS

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
			&cli.IntFlag{
				Name:    "port",
				Aliases: []string{"p"},
				Value:   6161,
				Usage:   "配置 with ui 使用，配置 web 的端口号",
			},
			&cli.BoolFlag{
				Name:  "loop",
				Usage: "是否循环执行",
				Value: true,
			},
			&cli.Int64Flag{
				Name:  "intervalMS",
				Value: 20_000,
				Usage: "循环间隔时间，毫秒",
			},
			&cli.Int64Flag{
				Name:    "loopCount",
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
			loopCount := c.Int("loopCount")
			libraConfigFilePath := c.String("config")
			intervalMS := c.Int64("intervalMS")
			withUI := c.Bool("ui")
			webPort := c.Int("port")
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
			rocketmq.Setup(libraConfig.RocketMQClients)

			if len(libraConfig.Services) == 0 {
				return nil
			}

			if withUI {
				go runningMonitor(libraConfig, intervalMS)
				r := gin.Default()

				fp, _ := fs.Sub(assets, "templates/assets")
				r.StaticFS("/assets", http.FS(fp))
				t, _ := template.ParseFS(templates, "templates/*.html")
				r.SetHTMLTemplate(t)
				r.GET("/", func(c *gin.Context) {
					c.HTML(http.StatusOK, "index.html", nil)
				})
				r.GET("/ping", func(c *gin.Context) {
					c.JSON(200, gin.H{
						"message": "pong",
					})
				})
				r.GET("/monitor", func(c *gin.Context) {
					monitoringResult := monitor.Monitoring(libraConfig)
					c.JSON(http.StatusOK, monitoringResult)

					//c.JSON(200, gin.H{
					//	"monitoringResult": monitoringResult,
					//})
				})
				r.Run(":" + strconv.Itoa(webPort))
			} else {
				if isNeedLoop {
					foreverLoop := loopCount < 1
					for i := 0; foreverLoop || i <= loopCount; i++ {
						runningMonitor(libraConfig, intervalMS)
					}
				} else {
					runningMonitor(libraConfig, 0)
				}
			}
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func runningMonitor(libra config.LibraConfig, intervalMS int64) {
	monitor.Monitoring(libra)
	if intervalMS > 0 {
		time.Sleep(time.Duration(intervalMS) * time.Millisecond)
	}
}
