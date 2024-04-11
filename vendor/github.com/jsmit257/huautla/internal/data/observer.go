package data

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jsmit257/huautla/types"
)

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
			&row.EventType.Stage.Name,
		); err != nil {
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
			&result.EventType.Stage.Name,
		); err != nil {
		return result, err
	}

	return result, err
}

func (db *Conn) addEvent(ctx context.Context, oID types.UUID, events []types.Event, e *types.Event, cid types.CID) ([]types.Event, error) {
	var err error
	var result sql.Result

	e.UUID = types.UUID(db.generateUUID().String())
	e.MTime = time.Now().UTC()
	e.CTime = e.MTime

	if result, err = db.ExecContext(ctx, psqls["event"]["add"],
		e.UUID,
		e.Temperature,
		e.Humidity,
		e.MTime,
		e.CTime,
		oID,
		e.EventType.UUID,
	); err != nil {
		if isPrimaryKeyViolation(err) {
			return db.addEvent(ctx, oID, events, e, cid)
		}
		return events, err
	} else if rows, err := result.RowsAffected(); err != nil {
		return events, err
	} else if rows != 1 { // most likely cause is a bad eventtype.uuid
		return events, fmt.Errorf("event was not added")
	}

	if e.EventType, err = db.SelectEventType(ctx, e.EventType.UUID, cid); err != nil {
		return events, fmt.Errorf("couldn't fetch eventtype")
	}

	return append([]types.Event{*e}, events...), err
}

func (db *Conn) changeEvent(ctx context.Context, events []types.Event, e *types.Event, cid types.CID) ([]types.Event, error) {
	var err error

	e.MTime = time.Now().UTC()

	if result, err := db.ExecContext(ctx, psqls["event"]["change"],
		e.Temperature,
		e.Humidity,
		e.MTime,
		e.UUID,
		e.EventType.UUID,
	); err != nil {
		return events, err
	} else if rows, err := result.RowsAffected(); err != nil {
		return events, err
	} else if rows != 1 { // most likely cause is a bad eventtype.uuid
		return events, fmt.Errorf("event was not changed")
	}

	if e.EventType, err = db.SelectEventType(ctx, e.EventType.UUID, cid); err != nil {
		return events, fmt.Errorf("couldn't fetch eventtype")
	}

	i, j := 0, len(events)
	for i < j && events[i].UUID != e.UUID {
		i++
	}

	return append(append([]types.Event{*e}, events[:i]...), events[i+1:]...), nil
}

func (db *Conn) removeEvent(ctx context.Context, events []types.Event, id types.UUID, cid types.CID) ([]types.Event, error) {
	if result, err := db.ExecContext(ctx, psqls["event"]["remove"], id); err != nil {
		return events, err
	} else if rows, err := result.RowsAffected(); err != nil {
		return events, err
	} else if rows != 1 { // most likely cause is a bad vendor.uuid
		return events, fmt.Errorf("event could not be removed")
	}

	i, j := 0, len(events)
	for i < j && events[i].UUID != id {
		i++
	}

	return append(events[:i], events[i+1:]...), nil
}
