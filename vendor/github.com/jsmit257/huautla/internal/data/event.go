package data

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jsmit257/huautla/types"
)

func (db *Conn) GetLifecycleEvents(ctx context.Context, lc *types.Lifecycle, cid types.CID) error {
	var err error

	deferred, start, l := initAccessFuncs("GetLifecycleEvents", db.logger, lc.UUID, cid)
	defer deferred(start, err, l)

	lc.Events, err = db.selectEventsList(ctx, psqls["event"]["all-by-lifecycle"], lc.UUID)

	return err
}

func (db *Conn) SelectByEventType(ctx context.Context, et types.EventType, cid types.CID) ([]types.Event, error) {
	var err error

	deferred, start, l := initAccessFuncs("SelectByEventType", db.logger, et.UUID, cid)
	defer deferred(start, err, l)

	return db.selectEventsList(ctx, psqls["event"]["all-by-eventtype"], et.UUID)
}

func (db *Conn) selectEventsList(ctx context.Context, query string, id types.UUID) ([]types.Event, error) {
	var err error
	var rows *sql.Rows

	result := make([]types.Event, 0, 1000)

	rows, err = db.query.QueryContext(ctx, query, id)
	if err != nil {
		return result, err
	}

	defer rows.Close()

	for rows.Next() {
		row := types.Event{}
		if err = rows.Scan(
			&row.UUID,
			&row.Temperature,
			&row.Humidity,
			&row.MTime,
			&row.CTime,
			&row.EventType.UUID,
			&row.EventType.Name,
			&row.EventType.Severity,
			&row.EventType.Stage.UUID,
			&row.EventType.Stage.Name); err != nil {

			return result, err
		}
		result = append(result, row)
	}

	return result, err
}

// XXX: is this even useful??
func (db *Conn) SelectEvent(ctx context.Context, id types.UUID, cid types.CID) (types.Event, error) {
	var err error

	deferred, start, l := initAccessFuncs("SelectEvent", db.logger, id, cid)
	defer deferred(start, err, l)

	result := types.Event{UUID: id}

	if err = db.
		QueryRowContext(ctx, psqls["event"]["select"], id).
		Scan(
			&result.Temperature,
			&result.Humidity,
			&result.MTime,
			&result.CTime,
			&result.EventType.UUID,
			&result.EventType.Name,
			&result.EventType.Severity,
			&result.EventType.Stage.UUID,
			&result.EventType.Stage.Name); err != nil {

		return result, err
	}

	return result, err
}

func (db *Conn) AddEvent(ctx context.Context, lc *types.Lifecycle, e types.Event, cid types.CID) error {
	var err error
	var result sql.Result

	deferred, start, l := initAccessFuncs("AddEvent", db.logger, lc.UUID, cid)
	defer deferred(start, err, l)

	e.UUID = types.UUID(db.generateUUID().String())
	e.MTime = time.Now().UTC()
	e.CTime = e.MTime

	result, err = db.ExecContext(ctx, psqls["event"]["add"],
		e.UUID,
		e.Temperature,
		e.Humidity,
		e.MTime,
		e.CTime,
		lc.UUID,
		e.EventType.UUID)
	if err != nil {
		if isPrimaryKeyViolation(err) {
			return db.AddEvent(ctx, lc, e, cid) // FIXME: infinite loop?
		}
		return err
	} else if rows, err := result.RowsAffected(); err != nil {
		return err
	} else if rows != 1 { // most likely cause is a bad vendor.uuid
		return fmt.Errorf("event was not added")
	}

	lc.Events = append(lc.Events, e)

	return err
}

func (db *Conn) ChangeEvent(ctx context.Context, lc *types.Lifecycle, e types.Event, cid types.CID) error {
	var err error
	var result sql.Result

	deferred, start, l := initAccessFuncs("ChangeEvent", db.logger, lc.UUID, cid)
	defer deferred(start, err, l)

	e.MTime = time.Now().UTC()

	result, err = db.ExecContext(ctx, psqls["event"]["change"],
		e.Temperature,
		e.Humidity,
		e.MTime,
		e.UUID,
		e.EventType.UUID)
	if err != nil {
		return err
	} else if rows, err := result.RowsAffected(); err != nil {
		return err
	} else if rows != 1 { // most likely cause is a bad vendor.uuid
		return fmt.Errorf("event was not changed")
	}

	i, j := 0, len(lc.Events)
	for i < j && lc.Events[i].UUID != e.UUID {
		i++
	}
	lc.Events[i] = e

	return err
}

func (db *Conn) RemoveEvent(ctx context.Context, lc *types.Lifecycle, id types.UUID, cid types.CID) error {

	var err error

	deferred, start, l := initAccessFuncs("RemoveEvent", db.logger, lc.UUID, cid)
	defer deferred(start, err, l)

	result, err := db.ExecContext(ctx, psqls["event"]["remove"], id)

	if err != nil {
		return err
	} else if rows, err := result.RowsAffected(); err != nil {
		return err
	} else if rows != 1 { // most likely cause is a bad vendor.uuid
		return fmt.Errorf("event could not be removed")
	}

	i, j := 0, len(lc.Events)
	for i < j && lc.Events[i].UUID != id {
		i++
	}
	lc.Events = append(lc.Events[:i], lc.Events[i+1:]...)

	return nil
}
