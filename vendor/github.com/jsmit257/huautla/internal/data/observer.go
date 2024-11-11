package data

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jsmit257/huautla/types"
)

type (
	nullnote struct {
		uuid         *types.UUID
		note         *string
		ctime, mtime *time.Time
	}
	nullphoto struct {
		uuid         *types.UUID
		filename     *string
		ctime, mtime *time.Time
	}
)

func (db *Conn) SelectByEventType(ctx context.Context, et types.EventType, cid types.CID) ([]types.Event, error) {
	var err error
	deferred, l := initAccessFuncs("SelectByEventType", db.logger, et.UUID, cid)
	defer deferred(&err, l)

	result, err := db.selectEventsList(ctx, psqls["event"]["all-by-eventtype"], et.UUID, cid)

	return result, err
}

func (db *Conn) selectEventsList(ctx context.Context, query string, id types.UUID, _ types.CID) ([]types.Event, error) {
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

func (db *Conn) notesAndPhotos(ctx context.Context, e []types.Event, id types.UUID, cid types.CID) error {
	var err error
	deferred, l := initAccessFuncs("notesAndPhotos", db.logger, id, cid)
	defer deferred(&err, l)

	if len(e) == 0 { // not really needed for safety, but it saves a hit to the db
		return nil
	}

	evts := make(map[types.UUID]*types.Event, len(e))
	for i, v := range e {
		evts[v.UUID] = &e[i]
	}

	rows, err := db.query.QueryContext(ctx, psqls["event"]["notes-and-photos"], id)
	if err != nil {
		return err
	}
	defer rows.Close()

	var lastnote *types.Note
	var lastphoto *types.Photo
	var eventUUID types.UUID
	for rows.Next() {
		n := nullnote{}
		p := nullphoto{}
		pn := nullnote{}
		if err = rows.Scan(
			&eventUUID,
			&n.uuid,
			&n.note,
			&n.mtime,
			&n.ctime,
			&p.uuid,
			&p.filename,
			&p.mtime,
			&p.ctime,
			&pn.uuid,
			&pn.note,
			&pn.mtime,
			&pn.ctime,
		); err != nil {
			return err
		}

		evt := evts[eventUUID]

		if n.uuid != nil {
			note := types.Note{
				UUID:  *n.uuid,
				Note:  *n.note,
				MTime: *n.mtime,
				CTime: *n.ctime,
			}

			if lastnote == nil || *lastnote != note {
				evt.Notes = append([]types.Note{note}, evt.Notes...)
			}

			lastnote = &note
		}

		if p.uuid != nil {
			photo := types.Photo{
				UUID:     *p.uuid,
				Filename: *p.filename,
				MTime:    *p.mtime,
				CTime:    *p.ctime,
			}

			if pn.uuid != nil {
				photo.Notes = []types.Note{{
					UUID:  *pn.uuid,
					Note:  *pn.note,
					MTime: *pn.mtime,
					CTime: *pn.ctime,
				}}
			}

			if lastphoto == nil || lastphoto.UUID != photo.UUID {
				evt.Photos = append([]types.Photo{photo}, evt.Photos...)
			} else if lastphoto.UUID == photo.UUID {
				evt.Photos[0].Notes = append(photo.Notes, evt.Photos[0].Notes...)
			}

			lastphoto = &photo
		}
	}

	return nil
}

func (db *Conn) SelectEvent(ctx context.Context, id types.UUID, cid types.CID) (types.Event, error) {
	var err error
	deferred, l := initAccessFuncs("SelectEvent", db.logger, id, cid)
	defer deferred(&err, l)

	result := types.Event{UUID: id}

	if err = db.
		QueryRowContext(ctx, psqls["event"]["select"], id).
		Scan(
			&result.UUID,
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

func (db *Conn) removeEvent(ctx context.Context, events []types.Event, id types.UUID, _ types.CID) ([]types.Event, error) {
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
