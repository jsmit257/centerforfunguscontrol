package huautla

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	log "github.com/sirupsen/logrus"

	"github.com/jsmit257/huautla"
	"github.com/jsmit257/huautla/types"
)

type (
	HuautlaAdaptor struct {
		db    types.DB
		log   *log.Entry
		filer func(string, []byte, fs.FileMode) error
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
			filer: os.WriteFile,
		}, nil
	}
}

func getUUIDByName(name string, w http.ResponseWriter, r *http.Request, ms *methodStats) (uuid types.UUID, err error) {
	if id := chi.URLParam(r, name); id == "" {
		err = fmt.Errorf("missing required id parameter")
	} else if id, err = url.QueryUnescape(id); err != nil {
		err = fmt.Errorf("malformed id parameter")
	} else {
		uuid = types.UUID(id)
	}
	return uuid, err
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
	// ms.m.??? // fit metrics in here eventually
	w.Header().Add("cid", string(ms.cid))
	w.WriteHeader(sc)
	ms.elapsed().WithError(err).Error(msg)
}

// assuming noone has called error() on this object, send() is the next likely step,
// to get the result data to the client
func (ms *methodStats) send(w http.ResponseWriter, i interface{}, sc int) {
	result, err := json.Marshal(i)
	if err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to marshal result")
	} else {
		w.Header().Add("Content-type", "application/json")
		w.WriteHeader(sc)
		if sc == http.StatusNoContent {
			return
		}
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
