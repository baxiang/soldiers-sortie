package main

import (
	"fmt"

	"os"

	handlerouter "github.com/baxiang/soldiers-sortie/hello-machinery/handle"
	itasks "github.com/baxiang/soldiers-sortie/hello-machinery/tasks"

	conf "github.com/baxiang/soldiers-sortie/hello-machinery/config"

	"github.com/RichardKnop/machinery/v2"
	"github.com/RichardKnop/machinery/v2/config"
	"github.com/gin-gonic/gin"
	"github.com/urfave/cli"

	redisbackend "github.com/RichardKnop/machinery/v2/backends/redis"
	redisbroker "github.com/RichardKnop/machinery/v2/brokers/redis"
	eagerlock "github.com/RichardKnop/machinery/v2/locks/eager"
)

var (
	server *machinery.Server
	cnf    *config.Config
	app    *cli.App
	tasks  map[string]interface{}
)

func init() {
	//初始化配置
	conf.InitConfig()

	app = cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "c",
			Value: "",
			Usage: "Path to a configuration file",
		},
	}

	//
	tasks = map[string]interface{}{
		"add":               itasks.Add,
		"long_running_task": itasks.LongRunningTask,
	}


	cnf := &config.Config{
		DefaultQueue:    "machinery_tasks",
		ResultsExpireIn: 3600,
		Redis: &config.RedisConfig{
			MaxIdle:                3,
			IdleTimeout:            240,
			ReadTimeout:            15,
			WriteTimeout:           15,
			ConnectTimeout:         15,
			NormalTasksPollPeriod:  1000,
			DelayedTasksPollPeriod: 500,
		},
	}

	// Create server instance
	broker := redisbroker.NewGR(cnf, []string{"localhost:6379"}, 0)
	backend := redisbackend.NewGR(cnf, []string{"localhost:6379"}, 0)
	lock := eagerlock.New()

	server = machinery.NewServer(cnf, broker, backend, lock)

}
func main() {

	// Set the CLI app commands
	app.Commands = []cli.Command{
		{
			Name:  "worker",
			Usage: "launch machinery worker",
			Action: func(c *cli.Context) error {
				if err := runWorker(); err != nil {
					return cli.NewExitError(err.Error(), 1)
				}
				return nil
			},
		},
		{
			Name:  "send",
			Usage: "send async tasks ",
			Action: func(c *cli.Context) error {
				if err := runSender(); err != nil {
					return cli.NewExitError(err.Error(), 1)
				}
				return nil
			},
		},
	}

	// Run the CLI app
	app.Run(os.Args)

}


func runWorker() (err error) {

	err = server.RegisterTasks(tasks)
	if err != nil {
		panic(err)
		return
	}
	workers := server.NewWorker("worker_test", 10)
	err = workers.Launch()
	if err != nil {
		panic(err)
		return
	}
	return
}
// sender对外提供接口的
func runSender() (err error) {

	err = server.RegisterTasks(tasks)
	if err != nil {
		panic(err)
		return
	}
	r := gin.Default()
	r.GET("/add", func(c *gin.Context) {
		handlerouter.Add(c, server)
	})
	r.POST("/longRunningTask", func(c *gin.Context) {
		handlerouter.LongRunningTask(c, server)
	})

	err = r.Run(fmt.Sprintf(":%d", conf.Cfg.AppPort))
	return
}

