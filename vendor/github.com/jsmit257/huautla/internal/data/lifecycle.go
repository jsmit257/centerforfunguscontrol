package data

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jsmit257/huautla/types"
)

func (db *Conn) SelectLifecycleIndex(ctx context.Context, cid types.CID) ([]types.Lifecycle, error) {
	var err error

	deferred, start, l := initAccessFuncs("SelectLifecycleIndex", db.logger, "nil", cid)
	defer deferred(start, err, l)

	var rows *sql.Rows

	result := make([]types.Lifecycle, 0, 1000)

	rows, err = db.query.QueryContext(ctx, psqls["lifecycle"]["index"])
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var eID, etID, stID *types.UUID
		var etName, etSev, stName *string
		var mtime, ctime *time.Time
		var temp *float32
		var hum *int8

		row := types.Lifecycle{}

		if err = rows.Scan(
			&row.UUID,
			&row.Location,
			&row.MTime,
			&row.CTime,
			&row.Strain.UUID,
			&row.Strain.Species,
			&row.Strain.Name,
			&row.Strain.CTime,
			&row.Strain.Vendor.UUID,
			&row.Strain.Vendor.Name,
			&row.Strain.Vendor.Website,
			&eID,
			&temp,
			&hum,
			&mtime,
			&ctime,
			&etID,
			&etName,
			&etSev,
			&stID,
			&stName,
		); err != nil {
			break
		}

		if eID != nil {
			row.Events = []types.Event{{
				UUID:        *eID,
				Temperature: *temp,
				Humidity:    *hum,
				MTime:       *mtime,
				CTime:       *ctime,
				EventType: types.EventType{
					UUID:     *etID,
					Name:     *etName,
					Severity: *etSev,
					Stage: types.Stage{
						UUID: *stID,
						Name: *stName,
					},
				},
			}}
		}

		if curr := len(result) - 1; curr < 0 || result[curr].UUID != row.UUID {
			result = append(result, row)
		} else {
			result[curr].Events = append(result[curr].Events, row.Events...)
		}
	}

	return result, err
}

func (db *Conn) SelectLifecycle(ctx context.Context, id types.UUID, cid types.CID) (types.Lifecycle, error) {
	var err error
	var result []types.Lifecycle

	deferred, start, l := initAccessFuncs("SelectLifecycle", db.logger, id, cid)
	defer deferred(start, err, l)

	p, _ := types.NewReportAttrs(map[string][]string{"lifecycle-id": {string(id)}})

	result, err = db.SelectLifecyclesByAttrs(ctx, p, cid)
	if err != nil {
		return types.Lifecycle{}, err
	} else if l := len(result); l == 1 {
		return result[0], nil
	} else if l == 0 {
		err = sql.ErrNoRows
	} else {
		err = fmt.Errorf("too many rows returned for SelectLifecycle")
	}

	return types.Lifecycle{}, err
}

func (db *Conn) SelectLifecyclesByAttrs(ctx context.Context, p types.ReportAttrs, cid types.CID) ([]types.Lifecycle, error) {
	var err error
	var rows *sql.Rows

	deferred, start, l := initAccessFuncs("SelectLifecyclesByAttr", db.logger, "nil", cid)
	defer deferred(start, err, l)

	result := make([]types.Lifecycle, 0, 1000)

	var generationID *types.UUID

	if !p.Contains("lifecycle-id", "strain-id", "grain-id", "bulk-id") {
		err = fmt.Errorf("request doesn't contain at least 1 required field")
		return result, err
	}

	rows, err = db.QueryContext(ctx, psqls["lifecycle"]["select"],
		p.Get("lifecycle-id"),
		p.Get("strain-id"),
		p.Get("grain-id"),
		p.Get("bulk-id"))
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		row := types.Lifecycle{}

		if err = rows.Scan(

			&row.UUID,
			&row.Location,
			&row.StrainCost,
			&row.GrainCost,
			&row.BulkCost,
			&row.Yield,
			&row.Count,
			&row.Gross,
			&row.MTime,
			&row.CTime,
			&row.Strain.UUID,
			&row.Strain.Species,
			&row.Strain.Name,
			&generationID,
			&row.Strain.CTime,
			&row.Strain.Vendor.UUID,
			&row.Strain.Vendor.Name,
			&row.Strain.Vendor.Website,
			&row.GrainSubstrate.UUID,
			&row.GrainSubstrate.Name,
			&row.GrainSubstrate.Type,
			&row.GrainSubstrate.Vendor.UUID,
			&row.GrainSubstrate.Vendor.Name,
			&row.GrainSubstrate.Vendor.Website,
			&row.BulkSubstrate.UUID,
			&row.BulkSubstrate.Name,
			&row.BulkSubstrate.Type,
			&row.BulkSubstrate.Vendor.UUID,
			&row.BulkSubstrate.Vendor.Name,
			&row.BulkSubstrate.Vendor.Website,
		); err != nil {
			break
		}

		if generationID != nil {
			row.Strain.Generation = &types.Generation{UUID: *generationID}
		}

		if err = db.GetAllAttributes(ctx, &row.Strain, cid); err != nil {
			break
		} else if err = db.GetAllIngredients(ctx, &row.GrainSubstrate, cid); err != nil {
			break
		} else if err = db.GetAllIngredients(ctx, &row.BulkSubstrate, cid); err != nil {
			break
		} else if err = db.GetLifecycleEvents(ctx, &row, cid); err != nil {
			break
		}

		result = append(result, row)
	}

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
		lc.Location,
		lc.StrainCost,
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

func (db *Conn) UpdateLifecycle(ctx context.Context, lc types.Lifecycle, cid types.CID) (types.Lifecycle, error) {
	var err error
	var result sql.Result
	var rows int64

	lc.MTime = time.Now().UTC()

	deferred, start, l := initAccessFuncs("UpdateLifecycle", db.logger, lc.UUID, cid)
	defer deferred(start, err, l)

	if result, err = db.ExecContext(ctx, psqls["lifecycle"]["update"],
		lc.Location,
		lc.StrainCost,
		lc.GrainCost,
		lc.BulkCost,
		lc.Yield,
		lc.Count,
		lc.Gross,
		lc.MTime,
		lc.Strain.UUID,
		lc.GrainSubstrate.UUID,
		lc.BulkSubstrate.UUID,
		lc.UUID,
	); err != nil {
		return lc, err
	} else if rows, err = result.RowsAffected(); err != nil {
		return lc, err
	} else if rows != 1 {
		err = fmt.Errorf("lifecycle was not updated")
	}

	return lc, err
}

func (db *Conn) UpdateLifecycleMTime(ctx context.Context, lc *types.Lifecycle, modified time.Time, cid types.CID) (*types.Lifecycle, error) {
	var err error

	deferred, start, l := initAccessFuncs("UpdateLifecycleMTime", db.logger, lc.UUID, cid)
	defer deferred(start, err, l)

	lc.MTime, err = db.updateMTime(ctx, "lifecycles", modified, lc.UUID, cid)

	return lc, err
}

func (db *Conn) DeleteLifecycle(ctx context.Context, id types.UUID, cid types.CID) error {
	return db.deleteByUUID(ctx, id, cid, "DeleteLifecycle", "lifecycle", db.logger)
}
