package main

import (
	"os"
	"sync"
	"syscall"

	"github.com/jsmit257/centerforfunguscontrol/internal/config"
	"github.com/jsmit257/centerforfunguscontrol/internal/data/huautla"
	"github.com/jsmit257/huautla/types"

	"github.com/sirupsen/logrus"
)

var traps = []os.Signal{
	syscall.SIGPIPE, // why not this?
	syscall.SIGINT,
	syscall.SIGHUP,
	syscall.SIGTERM,
	syscall.SIGQUIT,
}

func main() {
	cfg := config.NewConfig()
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel) // TODO: grab this from the config
	logger.SetFormatter(&logrus.JSONFormatter{})

	log := logger.WithFields(logrus.Fields{
		"app":     "cffc",
		"ingress": "http",
	})

	ha, err := huautla.New(&types.Config{
		PGHost: cfg.HuautlaHost,
		PGPort: cfg.HuautlaPort,
		PGUser: cfg.HuautlaUser,
		PGPass: cfg.HuautlaPass,
		PGSSL:  cfg.HuautlaSSL,
	},
		log)
	if err != nil {
		panic(err)
	}

	r := newHuautla(cfg, ha, log)
	newHC(r)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go startServer(cfg, r, wg, log)
	wg.Wait()

	log.Info("done, really?")

	os.Exit(0)
}
