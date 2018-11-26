package main

import (
	"fmt"
	"os"
	"time"

	"counter/internal/counter/handlers"
	"counter/pkg/counter"

	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"github.com/urfave/cli"
)

var (
	Version   = "unknown"
	BuildDate = "unknown"
)

func main() {
	var app *cli.App
	{
		app = &cli.App{
			Name:    "counter",
			Version: fmt.Sprintf("%s+%s", Version, BuildDate),
			Action:  run,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "redis-host",
					Usage:  "redis host",
					Value:  "127.0.0.1",
					EnvVar: "COUNTER_REDIS_HOST",
				},
				cli.IntFlag{
					Name:   "redis-port",
					Usage:  "redis port",
					Value:  6379,
					EnvVar: "COUNTER_REDIS_PORT",
				},
				cli.BoolFlag{
					Name:   "debug",
					Usage:  "enable debug mode",
					EnvVar: "COUNTER_DEBUG",
				},
			},
		}
	}

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}

func run(c *cli.Context) error {
	var redisPool *redis.Pool
	{
		redisPool = &redis.Pool{
			MaxIdle:     3,
			IdleTimeout: 10 * time.Second,
			Dial: func() (redis.Conn, error) {
				return redis.Dial("tcp", fmt.Sprintf(
					"%s:%d",
					c.String("redis-host"),
					c.Int("redis-port"),
				))
			},
		}
	}

	var redisCounter counter.Counter
	{
		redisCounter = counter.NewRedisCounter(redisPool, "")
	}

	var engine *gin.Engine
	{
		if !c.Bool("debug") {
			gin.SetMode(gin.ReleaseMode)
		}

		engine = gin.New()
		engine.Use(gin.Recovery())
		engine.Use(func(context *gin.Context) {
			context.Next()
		})
		v1 := engine.Group("v1")
		handlers.BindCounterHandler(v1, "redis", redisCounter)
	}

	return engine.Run("0.0.0.0:8080")
}
