package huautla

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"

	"github.com/jsmit257/centerforfunguscontrol/shared/metrics"
	"github.com/jsmit257/huautla"
	"github.com/jsmit257/huautla/types"
)

type (
	HuautlaAdaptor struct {
		db types.DB
		// log   *logrus.Entry
		filer func(string, []byte, fs.FileMode) error
	}

	methodStats struct {
		cid types.CID
		l   *logrus.Entry
		m   *prometheus.CounterVec
		s   time.Time
	}

	ParamError error
)

func New(cfg *types.Config, log *logrus.Entry) (*HuautlaAdaptor, error) {
	if db, err := huautla.New(cfg, log); err != nil {
		return nil, err
	} else {
		log.Info("connected to database")
		return &HuautlaAdaptor{
			db:    db,
			filer: os.WriteFile,
		}, nil
	}
}

func getUUIDByName(name string, _ http.ResponseWriter, r *http.Request, _ *methodStats) (uuid types.UUID, err error) {
	if id := chi.URLParam(r, name); id == "" {
		err = ParamError(fmt.Errorf("missing required parameter"))
	} else if id, err = url.QueryUnescape(id); err != nil {
		err = ParamError(fmt.Errorf("malformed parameter"))
	} else {
		uuid = types.UUID(id)
	}
	return uuid, err
}

func bodyHelper(r *http.Request, box any) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, &box)
}

// helper function adds fields `method` and `cid` to all subsequent logs; returns an object
// that encapsulates various success/error events with appropriate logging/metrics/responses
func (ha *HuautlaAdaptor) start(ctx context.Context, method string) *methodStats {
	cid := ctx.Value(metrics.Cid).(types.CID)

	result := &methodStats{
		cid: cid,
		m:   metrics.GetContextMetrics(ctx),
		l: metrics.GetContextLog(ctx).WithFields(logrus.Fields{
			"method": method,
			"cid":    cid,
		}),
		s: time.Now().UTC(),
	}

	result.l.Info("starting work")

	return result
}

func (ms *methodStats) lap() *methodStats {
	return &methodStats{
		cid: ms.cid,
		l:   ms.l.WithField("elapsed", time.Now().UTC().Sub(ms.s)),
		m:   ms.m,
		s:   ms.s,
	}

}

func (ms *methodStats) err(e error) *methodStats {
	return &methodStats{
		cid: ms.cid,
		l:   ms.l.WithError(e),
		m:   ms.m,
		s:   ms.s,
	}
}

func (ms *methodStats) error(w http.ResponseWriter, err error, sc int, msg interface{}) {
	ms.err(err).send(w, sc, msg)
}

func (ms *methodStats) send(w http.ResponseWriter, sc int, i interface{}) {
	result, err := json.Marshal(i)
	if err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to marshal result")
	} else {
		w.Header().Add("Content-type", "application/json")
		w.WriteHeader(sc)
		if sc == http.StatusNoContent {
			return
		}
		_, _ = w.Write([]byte(result))
	}
	ms.m.WithLabelValues(strconv.Itoa(sc)).Inc()
	ms.lap().l.Info("finished work")
}

func (ms *methodStats) empty(w http.ResponseWriter) {
	ms.send(w, http.StatusNoContent, nil)
}

func (ms *methodStats) created(w http.ResponseWriter, i interface{}) {
	ms.send(w, http.StatusCreated, i)
}

func (ms *methodStats) ok(w http.ResponseWriter, i interface{}) {
	ms.send(w, http.StatusOK, i)
}
