package data

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jsmit257/huautla/types"
)

func (db *Conn) SelectLifecycle(ctx context.Context, id types.UUID, cid types.CID) (types.Lifecycle, error) {
	var err error

	deferred, start, l := initAccessFuncs("SelectLifecycle", db.logger, id, cid)
	defer deferred(start, err, l)

	result := types.Lifecycle{UUID: id}

	if err = db.
		QueryRowContext(ctx, psqls["lifecycle"]["select"], id).
		Scan(
			&result.Name,
			&result.Location,
			&result.GrainCost,
			&result.BulkCost,
			&result.Yield,
			&result.Count,
			&result.Gross,
			&result.MTime,
			&result.CTime,
			&result.Strain.UUID,
			&result.Strain.Name,
			&result.Strain.Vendor.UUID,
			&result.Strain.Vendor.Name,
			&result.GrainSubstrate.UUID,
			&result.GrainSubstrate.Name,
			&result.GrainSubstrate.Type,
			&result.GrainSubstrate.Vendor.UUID,
			&result.GrainSubstrate.Vendor.Name,
			&result.BulkSubstrate.UUID,
			&result.BulkSubstrate.Name,
			&result.BulkSubstrate.Type,
			&result.BulkSubstrate.Vendor.UUID,
			&result.BulkSubstrate.Vendor.Name); err != nil {

		return result, err
	}

	if err = db.GetAllAttributes(ctx, &result.Strain, cid); err != nil {
		return result, err
	} else if err = db.GetAllIngredients(ctx, &result.GrainSubstrate, cid); err != nil {
		return result, err
	} else if err = db.GetAllIngredients(ctx, &result.BulkSubstrate, cid); err != nil {
		return result, err
	}

	err = db.GetLifecycleEvents(ctx, &result, cid)

	return result, err
}

func (db *Conn) InsertLifecycle(ctx context.Context, lc types.Lifecycle, cid types.CID) (types.Lifecycle, error) {
	var err error
	var result sql.Result
	var rows int64

	lc.UUID = types.UUID(db.generateUUID().String())
	lc.MTime = time.Now().UTC()
	lc.CTime = lc.MTime

	deferred, start, l := initAccessFuncs("InsertLifecycle", db.logger, lc.UUID, cid)
	defer deferred(start, err, l)

	result, err = db.ExecContext(ctx, psqls["lifecycle"]["insert"],
		lc.UUID,
		lc.Name,
		lc.Location,
		lc.GrainCost,
		lc.BulkCost,
		lc.Yield,
		lc.Count,
		lc.Gross,
		lc.MTime,
		lc.CTime,
		lc.Strain.UUID,
		lc.GrainSubstrate.UUID,
		lc.BulkSubstrate.UUID)

	if err != nil {
		if isPrimaryKeyViolation(err) {
			return db.InsertLifecycle(ctx, lc, cid) // FIXME: infinite loop?
		}
		return lc, err
	} else if rows, err = result.RowsAffected(); err != nil {
		return lc, err
	} else if rows != 1 {
		return lc, fmt.Errorf("lifecycle was not added: %d", rows)
	}

	return db.SelectLifecycle(ctx, lc.UUID, cid)
}

func (db *Conn) UpdateLifecycle(ctx context.Context, lc types.Lifecycle, cid types.CID) error {
	var err error
	var result sql.Result
	var rows int64

	lc.MTime = time.Now().UTC()

	deferred, start, l := initAccessFuncs("UpdateLifecycle", db.logger, lc.UUID, cid)
	defer deferred(start, err, l)

	if result, err = db.ExecContext(ctx, psqls["lifecycle"]["update"],
		lc.Name,
		lc.Location,
		lc.GrainCost,
		lc.BulkCost,
		lc.Yield,
		lc.Count,
		lc.Gross,
		lc.MTime,
		lc.Strain.UUID,
		lc.GrainSubstrate.UUID,
		lc.BulkSubstrate.UUID,
		lc.UUID); err != nil {

		return err
	} else if rows, err = result.RowsAffected(); err != nil {
		return err
	} else if rows != 1 {
		err = fmt.Errorf("lifecycle was not updated")
	}

	return err
}

func (db *Conn) DeleteLifecycle(ctx context.Context, id types.UUID, cid types.CID) error {
	return db.deleteByUUID(ctx, id, cid, "DeleteLifecycle", "lifecycle", db.logger)
}
