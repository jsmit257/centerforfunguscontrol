package data

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jsmit257/huautla/types"
)

func (db *Conn) SelectAllStages(ctx context.Context, cid types.CID) ([]types.Stage, error) {
	var err error

	deferred, start, l := initAccessFuncs("SelectAllStages", db.logger, types.UUID("nil"), cid)
	defer deferred(start, err, l)

	var rows *sql.Rows

	result := make([]types.Stage, 0, 100)

	rows, err = db.query.QueryContext(ctx, psqls["stage"]["select-all"])
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		row := types.Stage{}
		if err = rows.Scan(&row.UUID, &row.Name); err != nil {
			return result, err
		}
		result = append(result, row)
	}

	return result, err
}

func (db *Conn) SelectStage(ctx context.Context, id types.UUID, cid types.CID) (types.Stage, error) {
	var err error

	deferred, start, l := initAccessFuncs("SelectStage", db.logger, id, cid)
	defer deferred(start, err, l)

	result := types.Stage{UUID: id}
	err = db.
		QueryRowContext(ctx, psqls["stage"]["select"], id).
		Scan(&result.Name)

	return result, err
}

func (db *Conn) InsertStage(ctx context.Context, s types.Stage, cid types.CID) (types.Stage, error) {
	var err error

	s.UUID = types.UUID(db.generateUUID().String())

	deferred, start, l := initAccessFuncs("InsertStage", db.logger, s.UUID, cid)
	defer deferred(start, err, l)

	result, err := db.ExecContext(ctx, psqls["stage"]["insert"], s.UUID, s.Name)
	if err != nil {
		// FIXME: choose what to do based on the tupe of error
		duplicatePrimaryKeyErr := false
		if duplicatePrimaryKeyErr {
			return db.InsertStage(ctx, s, cid) // FIXME: infinite loop?
		}
		return s, err
	} else if rows, err := result.RowsAffected(); err != nil {
		return s, err
	} else if rows != 1 {
		return s, fmt.Errorf("stage was not added")
	}

	return s, err
}

func (db *Conn) UpdateStage(ctx context.Context, id types.UUID, s types.Stage, cid types.CID) error {
	var err error

	deferred, start, l := initAccessFuncs("UpdateStage", db.logger, id, cid)
	defer deferred(start, err, l)

	result, err := db.ExecContext(ctx, psqls["stage"]["update"], s.Name, id)
	if err != nil {
		return err
	} else if rows, err := result.RowsAffected(); err != nil {
		return err
	} else if rows != 1 {
		return fmt.Errorf("stage was not updated: '%s'", id)
	}
	return nil
}

func (db *Conn) DeleteStage(ctx context.Context, id types.UUID, cid types.CID) error {
	return db.deleteByUUID(ctx, id, cid, "DeleteStage", "stage", db.logger)
}
