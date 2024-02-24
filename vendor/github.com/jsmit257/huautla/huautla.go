package huautla

import (
	"fmt"

	"github.com/jsmit257/huautla/internal/data"
	"github.com/jsmit257/huautla/types"

	log "github.com/sirupsen/logrus"
)

// var _ = db.(interface {
// 	GetLifecycleEvents(context.Context, *types.Lifecycle, types.CID) error
// })

func New(cfg *types.Config, log *log.Entry) (types.DB, error) {
	var cnxFmt = "host=%s port=%d user=%s password=%s dbname=huautla sslmode=%s"
	var cnxInfo string

	if host := cfg.PGHost; host == "" {
		return nil, fmt.Errorf("postgres connection needs hostname attribute")
	} else if user := cfg.PGUser; user == "" {
		return nil, fmt.Errorf("postgres connection needs username attribute")
	} else if pass := cfg.PGPass; pass == "" {
		return nil, fmt.Errorf("postgres connection needs password attribute")
	} else if port := cfg.PGPort; port == 0 {
		return nil, fmt.Errorf("postgres connection needs port attribute")
	} else {
		cnxInfo = fmt.Sprintf(cnxFmt, host, port, user, pass, cfg.PGSSL)
	}

	return data.New(cnxInfo, log)
}
