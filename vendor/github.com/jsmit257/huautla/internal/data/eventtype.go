package data

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"

	"github.com/jsmit257/huautla/types"
)

func (db *Conn) SelectAllEventTypes(ctx context.Context, cid types.CID) ([]types.EventType, error) {
	var err error

	deferred, start, l := initAccessFuncs("SelectAllEventTypes", db.logger, "nil", cid)
	defer deferred(start, err, l)

	var rows *sql.Rows

	result := make([]types.EventType, 0, 100)

	rows, err = db.query.QueryContext(ctx, psqls["eventtype"]["select-all"])
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		row := types.EventType{}
		if err = rows.Scan(
			&row.UUID,
			&row.Name,
			&row.Severity,
			&row.Stage.UUID,
			&row.Stage.Name); err != nil {

			return result, err
		}
		result = append(result, row)
	}

	return result, err
}

func (db *Conn) SelectEventType(ctx context.Context, id types.UUID, cid types.CID) (types.EventType, error) {
	var err error

	deferred, start, l := initAccessFuncs("SelectEventType", db.logger, id, cid)
	defer deferred(start, err, l)

	result := types.EventType{}

	err = db.
		QueryRowContext(ctx, psqls["eventtype"]["select"], id).
		Scan(
			&result.UUID,
			&result.Name,
			&result.Severity,
			&result.Stage.UUID,
			&result.Stage.Name)

	return result, err
}

func (db *Conn) InsertEventType(ctx context.Context, e types.EventType, cid types.CID) (types.EventType, error) {
	var err error

	e.UUID = types.UUID(db.generateUUID().String())

	deferred, start, l := initAccessFuncs("InsertEventType", db.logger, e.UUID, cid)
	defer deferred(start, err, l)

	result, err := db.ExecContext(ctx, psqls["eventtype"]["insert"], e.UUID, e.Name, e.Severity, e.Stage.UUID)
	if err != nil {
		return e, err
	} else if rows, err := result.RowsAffected(); err != nil {
		return e, err
	} else if rows != 1 { // most likely cause is a bad vendor.uuid
		return e, fmt.Errorf("eventtype was not added")
	}

	return e, err
}

func (db *Conn) UpdateEventType(ctx context.Context, id types.UUID, e types.EventType, cid types.CID) error {
	var err error

	deferred, start, l := initAccessFuncs("UpdateEventType", db.logger, id, cid)
	defer deferred(start, err, l)

	result, err := db.ExecContext(ctx, psqls["eventtype"]["update"], e.Name, e.Severity, e.Stage.UUID, id)
	if err != nil {
		return err
	} else if rows, err := result.RowsAffected(); err != nil {
		return err
	} else if rows != 1 {
		return fmt.Errorf("eventtype was not updated: '%s'", id)
	}
	return nil
}

func (db *Conn) DeleteEventType(ctx context.Context, id types.UUID, cid types.CID) error {
	return db.deleteByUUID(ctx, id, cid, "DeleteEventType", "eventtype", db.logger)
}

func (e eventtype) children(db *Conn, ctx context.Context, cid types.CID, p *rpttree) error {
	var err error

	deferred, start, l := initAccessFuncs("eventtype::children", db.logger, types.UUID(e.UUID), cid)
	defer deferred(start, err, l)

	param, _ := types.NewReportAttrs(url.Values{"eventtype-id": {string(e.UUID)}})

	lcs, err := db.lifecycleReport(ctx, param, cid, p)
	if err != nil {
		return err
	} else if len(lcs) != 0 {
		p.data["lifecycles"] = lcs
	}

	gens, err := db.generationReport(ctx, param, cid, p)
	if err != nil {
		return err
	} else if len(gens) != 0 {
		p.data["generations"] = gens
	}

	return nil
}

func (db *Conn) EventTypeReport(ctx context.Context, id types.UUID, cid types.CID) (types.Entity, error) {
	var err error
	var rpt rpt

	deferred, start, l := initAccessFuncs("EventTypeReport", db.logger, id, cid)
	defer deferred(start, err, l)

	result, err := db.SelectEventType(ctx, id, cid)
	if err != nil {
		return nil, err
	} else if rpt, err = db.newRpt(ctx, eventtype(result), cid, nil); err != nil {
		return nil, err
	} else if rpt == nil {
		return nil, nil
	}

	return rpt.Data(), nil

}
