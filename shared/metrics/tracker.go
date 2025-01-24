package metrics

import (
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

const (
	Cid     ctxkey = "cid"
	Metrics ctxkey = "metrics"
	Log     ctxkey = "log"
)

type (
	Tracker interface {
		Field(n string, v interface{}) Tracker
		Fields(f logrus.Fields) Tracker
		SC(code int) Tracker
		Err(error) Tracker
		Lap() Tracker

		Debug(string, ...any) Tracker
		Warn(string, ...any) Tracker

		Done(string) resulter
		OK() resulter
	}

	resulter interface {
		Err() error
		SC() int
	}

	trackresults struct {
		code int
		e    error
	}

	track struct {
		l *logrus.Entry
		m *prometheus.CounterVec
		r *trackresults
		s time.Time
	}

	core struct {
		metrics *prometheus.CounterVec
		log     *logrus.Entry
	}

	datatracker    track
	servicetracker track

	ctxkey string
)

func (t *track) Lap() Tracker {
	return t.Field("duration", time.Since(t.s).String())
}

func (t *track) Field(n string, v interface{}) Tracker {
	return &track{l: t.l.WithField(n, v), m: t.m, r: t.r, s: t.s}
}

func (t *track) Fields(f logrus.Fields) Tracker {
	return &track{l: t.l.WithFields(f), m: t.m, r: t.r, s: t.s}
}

func (t *track) SC(code int) Tracker {
	return &track{
		l: t.l.WithField("statuscode", code),
		m: t.m,
		r: &trackresults{code: code, e: t.r.e},
		s: t.s,
	}
}

func (t *track) Err(e error) Tracker {
	if t.r.e != nil {
		// if we don't do this MustCurryWith (probably?) will; don't really want to return an error
		panic(fmt.Errorf("error is already set for tracker (old: %w), (new: %w)", t.r.e, e))
	}

	return &track{
		l: t.l.WithError(e),
		m: t.m.MustCurryWith(prometheus.Labels{"status": e.Error()}),
		r: &trackresults{code: t.r.code, e: e},
		s: t.s,
	}
}

func (t *track) Done(msg string) resulter {
	var status []string
	var closer = t.Lap().(*track).l.Info
	// var closer = t.l.WithField("duration", time.Since(t.s).String()).Info
	if t.r.e == nil {
		status = []string{"ok"}
	} else {
		closer = t.l.Error
	}
	t.m.WithLabelValues(status...).Inc()
	closer(msg)
	return t.r
}

func (t *track) OK() resulter {
	return t.Done("finished work")
}

func (t *track) Debug(msg string, args ...any) Tracker {
	t.l.Debugf(msg, args...)
	return t
}

func (t *track) Warn(msg string, args ...any) Tracker {
	t.l.Warnf(msg, args...)
	return t
}

func (r *trackresults) Err() error {
	return r.e
}

func (r *trackresults) SC() int {
	return r.code
}
