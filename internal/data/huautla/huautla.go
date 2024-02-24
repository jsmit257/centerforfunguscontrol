package huautla

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"

	log "github.com/sirupsen/logrus"

	"github.com/jsmit257/huautla"
	"github.com/jsmit257/huautla/types"
)

type (
	HuautlaAdaptor struct {
		db    types.DB
		log   *log.Entry
		mtrcs interface{}
	}

	methodStats struct {
		cid types.CID
		l   *log.Entry
		m   interface{}
		s   time.Time
	}
)

func New(cfg *types.Config, log *log.Entry, mtrcs interface{}) (*HuautlaAdaptor, error) {
	if db, err := huautla.New(cfg, log); err != nil {
		return nil, err
	} else {
		log.Info("connected to database")
		return &HuautlaAdaptor{
			db:    db,
			log:   log,
			mtrcs: mtrcs,
		}, nil
	}
}

// helper function adds fields `method` and `cid` to all subsequent logs; returns an object
// that encapsulates various success/error events with appropriate logging/metrics/responses
func (ha *HuautlaAdaptor) start(method string) *methodStats {
	result := methodStats{
		cid: cid(),
		m:   nil, // later
		s:   time.Now().UTC(),
	}
	result.l = ha.log.WithFields(log.Fields{
		"method": method,
		"cid":    result.cid,
	})

	result.l.Info("starting work")

	return &result
}

func (ms *methodStats) elapsed() *log.Entry {
	return ms.l.WithField("elapsed", time.Now().UTC().Sub(ms.s))
}

// simple way to log, emit metrics and respond to the client in a regular way
func (ms *methodStats) error(w http.ResponseWriter, err error, sc int, msg string) {
	ms.elapsed().WithError(err).Error(msg)
	// ms.m.??? // fit metrics in here eventually
	w.WriteHeader(sc)
}

// assuming noone has called error() on this object, send() is the next likely step,
// to get the result data to the client
func (ms *methodStats) send(w http.ResponseWriter, i interface{}, sc int) {
	result, err := json.Marshal(i)
	if err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to marshal result")
	} else {
		w.WriteHeader(sc)
		if _, err = w.Write([]byte(result)); err != nil {
			ms.error(w, err, http.StatusInternalServerError, "failed to send result")
		}
	}
}

func (ms *methodStats) end() {
	ms.elapsed().Info("finished work")
}

func cid() types.CID {
	return types.CID(uuid.NewString())
}
