package data

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jsmit257/huautla/types"
)

func (db *Conn) SelectAllVendors(ctx context.Context, cid types.CID) ([]types.Vendor, error) {
	var err error

	deferred, start, l := initAccessFuncs("SelectAllVendors", db.logger, "nil", cid)
	defer deferred(start, err, l)

	var rows *sql.Rows

	result := make([]types.Vendor, 0, 100)

	rows, err = db.query.QueryContext(ctx, psqls["vendor"]["select-all"])
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		row := types.Vendor{}
		if err = rows.Scan(&row.UUID, &row.Name, &row.Website); err != nil {
			break
		}
		result = append(result, row)
	}

	return result, err
}

func (db *Conn) SelectVendor(ctx context.Context, id types.UUID, cid types.CID) (types.Vendor, error) {
	var err error

	deferred, start, l := initAccessFuncs("SelectVendor", db.logger, id, cid)
	defer deferred(start, err, l)

	result := types.Vendor{UUID: id}
	err = db.
		QueryRowContext(ctx, psqls["vendor"]["select"], id).
		Scan(&result.Name, &result.Website)

	return result, err
}

func (db *Conn) InsertVendor(ctx context.Context, v types.Vendor, cid types.CID) (types.Vendor, error) {
	var err error
	var result sql.Result
	var rows int64

	v.UUID = types.UUID(db.generateUUID().String())

	deferred, start, l := initAccessFuncs("InsertVendor", db.logger, v.UUID, cid)
	defer deferred(start, err, l)

	result, err = db.ExecContext(ctx, psqls["vendor"]["insert"], v.UUID, v.Name, v.Website)
	if err != nil {
		if isPrimaryKeyViolation(err) {
			return db.InsertVendor(ctx, v, cid) // FIXME: infinite loop?
		}
		return v, err
	} else if rows, err = result.RowsAffected(); err != nil {
		return v, err
	} else if rows != 1 {
		err = fmt.Errorf("vendor was not added")
	}

	return v, err
}

func (db *Conn) UpdateVendor(ctx context.Context, id types.UUID, v types.Vendor, cid types.CID) error {
	var err error

	deferred, start, l := initAccessFuncs("UpdateVendor", db.logger, id, cid)
	defer deferred(start, err, l)

	result, err := db.ExecContext(ctx, psqls["vendor"]["update"], v.Name, v.Website, id)
	if err != nil {
		return err
	} else if rows, err := result.RowsAffected(); err != nil {
		return err
	} else if rows != 1 {
		return fmt.Errorf("vendor was not updated: '%s'", id)
	}
	return nil
}

func (db *Conn) DeleteVendor(ctx context.Context, id types.UUID, cid types.CID) error {
	return db.deleteByUUID(ctx, id, cid, "DeleteVendor", "vendor", db.logger)
}
