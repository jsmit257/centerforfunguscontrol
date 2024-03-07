package data

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jsmit257/huautla/types"
)

func (db *Conn) SelectAllSubstrates(ctx context.Context, cid types.CID) ([]types.Substrate, error) {
	var err error

	deferred, start, l := initAccessFuncs("SelectAllSubstrates", db.logger, "nil", cid)
	defer deferred(start, err, l)

	var rows *sql.Rows

	result := make([]types.Substrate, 0, 100)

	rows, err = db.query.QueryContext(ctx, psqls["substrate"]["select-all"])
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		row := types.Substrate{}
		if err = rows.Scan(
			&row.UUID,
			&row.Name,
			&row.Type,
			&row.Vendor.UUID,
			&row.Vendor.Name,
			&row.Vendor.Website); err != nil {

			return nil, err
		} else if err = db.GetAllIngredients(ctx, &row, cid); err != nil {
			return nil, err
		}
		result = append(result, row)
	}

	return result, err
}

func (db *Conn) SelectSubstrate(ctx context.Context, id types.UUID, cid types.CID) (types.Substrate, error) {
	var err error

	deferred, start, l := initAccessFuncs("SelectSubstrate", db.logger, id, cid)
	defer deferred(start, err, l)

	result := types.Substrate{UUID: id}

	if err = db.
		QueryRowContext(ctx, psqls["substrate"]["select"], id).
		Scan(
			&result.Name,
			&result.Type,
			&result.Vendor.UUID,
			&result.Vendor.Name,
			&result.Vendor.Website); err == nil {

		err = db.GetAllIngredients(ctx, &result, "SelectSubstrate")
	}

	return result, err
}

func (db *Conn) InsertSubstrate(ctx context.Context, s types.Substrate, cid types.CID) (types.Substrate, error) {
	var err error

	s.UUID = types.UUID(db.generateUUID().String())

	deferred, start, l := initAccessFuncs("InsertSubstrate", db.logger, s.UUID, cid)
	defer deferred(start, err, l)

	result, err := db.ExecContext(ctx, psqls["substrate"]["insert"], s.UUID, s.Name, s.Type, s.Vendor.UUID)

	if err != nil {
		if isPrimaryKeyViolation(err) {
			return db.InsertSubstrate(ctx, s, cid) // FIXME: infinite loop?
		}
		return s, err
	} else if rows, err := result.RowsAffected(); err != nil {
		return s, err
	} else if rows != 1 { // most likely cause is a bad vendor.uuid
		return s, fmt.Errorf("substrate was not added")
	}

	return s, err
}

func (db *Conn) UpdateSubstrate(ctx context.Context, id types.UUID, s types.Substrate, cid types.CID) error {
	var err error

	deferred, start, l := initAccessFuncs("UpdateSubstrate", db.logger, id, cid)
	defer deferred(start, err, l)

	// result, err := db.ExecContext(ctx, psqls["substrate"]["update"], s.Name, s.Type, s.Vendor.UUID, id)
	result, err := db.ExecContext(ctx, psqls["substrate"]["update"], s.Name, id)
	if err != nil {
		return err
	} else if rows, err := result.RowsAffected(); err != nil {
		return err
	} else if rows != 1 {
		return fmt.Errorf("substrate was not updated: '%s'", id)
	}
	return nil
}

func (db *Conn) DeleteSubstrate(ctx context.Context, id types.UUID, cid types.CID) error {
	// FIXME: delete all substrateingredients first
	return db.deleteByUUID(ctx, id, cid, "DeleteSubstrate", "substrate", db.logger)
}
