package metrics

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

type (
	logwriter struct {
		msgs []string
	}
)

func (l *logwriter) Write(b []byte) (int, error) {
	l.msgs = append(l.msgs, string(b))
	return len(l.msgs), nil // doesn't match the semantics of normal write
}

func (l *logwriter) String() string {
	return strings.Join(l.msgs, ",")
}

func (l *logwriter) Audit() ([]map[string]interface{}, error) {
	result := make([]map[string]interface{}, 0, len(l.msgs))
	for _, v := range l.msgs {
		tmp := map[string]interface{}{}
		if err := json.Unmarshal([]byte(v), &tmp); err != nil {
			return nil, err
		}
		delete(tmp, "cid")
		delete(tmp, "duration")
		delete(tmp, "function")
		delete(tmp, "test")
		delete(tmp, "time")
		result = append(result, tmp)
	}
	return result, nil
}

func (l *logwriter) Clear() {
	l.msgs = nil
}

func Test_NewTracker(t *testing.T) {
	w := &logwriter{}

	log := logrus.New()
	log.SetLevel(logrus.DebugLevel)
	log.SetOutput(w)
	log.SetFormatter(&logrus.JSONFormatter{})

	ctx := context.WithValue(
		context.WithValue(
			context.TODO(),
			Log,
			log.WithField("test", "Test_NewTracker"),
		),
		Metrics,
		DataMetrics.MustCurryWith(prometheus.Labels{
			"pkg": "test",
			"db":  "test",
		}),
	)
	tracker := NewDataTracker(ctx, "Test_NewTracker")
	w.Clear()

	err := tracker.Lap().Done("testing Lap() (and Field(), Done(), Err())").Err()
	require.Nil(t, err)
	audit, err := w.Audit()
	require.Nil(t, err)
	require.Equal(t, 1, len(audit))
	result := map[string]interface{}{
		"level": "info",
		"msg":   "testing Lap() (and Field(), Done(), Err())",
	}
	require.Equal(t, result, audit[0])
	w.Clear()

	err = tracker.Fields(logrus.Fields{"field": "test"}).Done("testing Fields()").Err()
	require.Nil(t, err)
	audit, err = w.Audit()
	require.Nil(t, err)
	require.Equal(t, 1, len(audit))
	result = map[string]interface{}{
		"field": "test",
		"level": "info",
		"msg":   "testing Fields()",
	}
	require.Equal(t, result, audit[0])
	w.Clear()

	sc := tracker.SC(http.StatusTeapot).Done("testing both SC()s").SC()
	require.Equal(t, http.StatusTeapot, sc)
	audit, err = w.Audit()
	require.Nil(t, err)
	require.Equal(t, 1, len(audit))
	result = map[string]interface{}{
		"level":      "info",
		"msg":        "testing both SC()s",
		"statuscode": float64(418),
	}
	require.Equal(t, result, audit[0])
	w.Clear()

	err = tracker.Err(fmt.Errorf("some error")).Done("testing Err(error)").Err()
	require.Equal(t, fmt.Errorf("some error"), err)
	audit, err = w.Audit()
	require.Nil(t, err)
	require.Equal(t, 1, len(audit))
	result = map[string]interface{}{
		"error": "some error",
		"level": "error",
		"msg":   "testing Err(error)",
	}
	require.Equal(t, result, audit[0])
	w.Clear()

	sc = tracker.OK().SC()
	require.Equal(t, 0, sc)
	audit, err = w.Audit()
	require.Nil(t, err)
	require.Equal(t, 1, len(audit))
	result = map[string]interface{}{
		"level": "info",
		"msg":   "finished work",
	}
	require.Equal(t, result, audit[0], w.String())
	w.Clear()

	func() {
		defer func() {
			require.NotNil(t, recover())
		}()
		tracker := tracker.Err(fmt.Errorf("some error"))
		_ = tracker.Err(fmt.Errorf("run for it"))
	}()

	tracker = tracker.Debug("debug '%s'", "args")
	require.NotNil(t, tracker)
	audit, err = w.Audit()
	require.Nil(t, err)
	require.Equal(t, 1, len(audit))
	result = map[string]interface{}{
		"level": "debug",
		"msg":   "debug 'args'",
	}
	require.Equal(t, result, audit[0])
	w.Clear()

	tracker = tracker.Warn("warn '%s'", "args")
	require.NotNil(t, tracker)
	audit, err = w.Audit()
	require.Nil(t, err)
	require.Equal(t, 1, len(audit))
	result = map[string]interface{}{
		"level": "warning",
		"msg":   "warn 'args'",
	}
	require.Equal(t, result, audit[0])
	w.Clear()
}
