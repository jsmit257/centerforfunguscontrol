package main

import (
	"os"
	"sync"
	"syscall"

	"github.com/jsmit257/centerforfunguscontrol/internal/config"
	"github.com/jsmit257/centerforfunguscontrol/internal/data/huautla"
	"github.com/jsmit257/huautla/types"

	log "github.com/sirupsen/logrus"
)

var traps = []os.Signal{
	syscall.SIGINT,
	syscall.SIGHUP,
	syscall.SIGTERM,
	syscall.SIGQUIT,
}

func main() {
	cfg := config.NewConfig()
	log.SetLevel(log.DebugLevel) // TODO: grab this from the config
	log.SetFormatter(&log.JSONFormatter{})

	log := log.WithFields(log.Fields{
		"app":     "cffc",
		"ingress": "http",
	})

	wg := &sync.WaitGroup{}

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

	startServer(cfg, r, wg, log).Wait()

	log.Info("done")

	os.Exit(0)
}
