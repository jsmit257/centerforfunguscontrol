package data

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jsmit257/huautla/types"

	"github.com/google/uuid"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	log "github.com/sirupsen/logrus"

	_ "github.com/lib/pq"
)

type (
	query interface {
		ExecContext(context.Context, string, ...any) (sql.Result, error)
		QueryContext(context.Context, string, ...any) (*sql.Rows, error)
		QueryRowContext(context.Context, string, ...any) *sql.Row
	}

	Conn struct {
		query
		generateUUID uuidgen
		logger       *log.Entry
		// sql          map[string]map[string]string
	}

	uuidgen func() uuid.UUID

	getMockDB func(*sql.DB, sqlmock.Sqlmock, error) *sql.DB

	deferred func(*error, *log.Entry)
)

func New(cnxInfo string, log *log.Entry) (types.DB, error) {
	var err error
	var query *sql.DB

	if query, err = sql.Open("postgres", cnxInfo); err != nil {
		return nil, err
	} else if err = query.Ping(); err != nil {
		return nil, err
	}

	return &Conn{
		query:        query,
		generateUUID: uuid.New,
		logger:       log,
	}, nil
}

func (db *Conn) deleteByUUID(ctx context.Context, id types.UUID, cid types.CID, method, table string, l *log.Entry) error {
	var err error

	deferred, l := initAccessFuncs(method, l, id, cid)
	defer deferred(&err, l)

	var result sql.Result

	result, err = db.ExecContext(ctx, psqls[table]["delete"], id)
	if err != nil {
		return err
	} else if rows, err := result.RowsAffected(); err != nil {
		return err
	} else if rows != 1 {
		// this won't be reported in the WithError log in `defer ...`, b/c it's operator error
		return fmt.Errorf("%s could not be deleted: '%s'", table, id)
	}

	return err
}

func initAccessFuncs(fn string, l *log.Entry, id any, cid types.CID) (deferred, *log.Entry) {
	start := time.Now()
	l = l.WithFields(log.Fields{
		"function": fn,
		"cid":      cid,
	})
	if id != nil {
		l = l.WithField("key", id)
	}

	l.Info("starting work")

	return func(err *error, l *log.Entry) {
		duration := time.Since(start)

		if err != nil {
			l.WithError(*err)
		}
		l.WithField("duration", duration).Infof("finished work")

		// TODO: metrics
	}, l
}

func isPrimaryKeyViolation(error) bool {
	return false
}
