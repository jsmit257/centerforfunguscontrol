package data

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"time"

	"github.com/jsmit257/huautla/types"
)

func (db *Conn) SelectLifecycleIndex(ctx context.Context, cid types.CID) ([]types.Lifecycle, error) {
	var err error
	deferred, l := initAccessFuncs("SelectLifecycleIndex", db.logger, "nil", cid)
	defer deferred(&err, l)

	result := make([]types.Lifecycle, 0, 1000)

	rows, err := db.query.QueryContext(ctx, psqls["lifecycle"]["index"])
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
	deferred, l := initAccessFuncs("SelectLifecycle", db.logger, id, cid)
	defer deferred(&err, l)

	p, _ := types.NewReportAttrs(map[string][]string{"lifecycle-id": {string(id)}})

	result, err := db.selectLifecycles(ctx, p, cid)
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

func (db *Conn) selectLifecycles(ctx context.Context, p types.ReportAttrs, cid types.CID) ([]types.Lifecycle, error) {
	var err error
	deferred, l := initAccessFuncs("selectLifecycles", db.logger, "nil", cid)
	defer deferred(&err, l)

	var generationID *types.UUID

	result := make([]types.Lifecycle, 0, 1000)

	if !p.Contains("lifecycle-id", "strain-id", "grain-id", "bulk-id", "eventtype-id") {
		err = fmt.Errorf("request doesn't contain at least 1 required field")
		return result, err
	}

	rows, err := db.QueryContext(ctx, psqls["lifecycle"]["select"],
		p.Get("lifecycle-id"),
		p.Get("strain-id"),
		p.Get("grain-id"),
		p.Get("bulk-id"),
		p.Get("eventtype-id"))
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
			&row.Strain.DTime,
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

		if err = db.GetLifecycleEvents(ctx, &row, cid); err != nil {
			break
		}

		result = append(result, row)
	}

	return result, err
}

func (db *Conn) InsertLifecycle(ctx context.Context, lc types.Lifecycle, cid types.CID) (types.Lifecycle, error) {
	var err error
	deferred, l := initAccessFuncs("InsertLifecycle", db.logger, lc.UUID, cid)
	defer deferred(&err, l)

	var result sql.Result
	var rows int64

	lc.UUID = types.UUID(db.generateUUID().String())
	lc.MTime = time.Now().UTC()
	lc.CTime = lc.MTime

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
			return db.InsertLifecycle(ctx, lc, cid)
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
	deferred, l := initAccessFuncs("UpdateLifecycle", db.logger, lc.UUID, cid)
	defer deferred(&err, l)

	var result sql.Result
	var rows int64

	lc.MTime = time.Now().UTC()

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
		err = fmt.Errorf("one of strain, grain or bulk is not the right type")
	}

	return lc, err
}

func (db *Conn) UpdateLifecycleMTime(ctx context.Context, lc *types.Lifecycle, modified time.Time, cid types.CID) (*types.Lifecycle, error) {
	var err error
	deferred, l := initAccessFuncs("UpdateLifecycleMTime", db.logger, lc.UUID, cid)
	defer deferred(&err, l)

	lc.MTime, err = db.updateMTime(ctx, "lifecycles", modified, lc.UUID, cid)

	return lc, err
}

func (db *Conn) DeleteLifecycle(ctx context.Context, id types.UUID, cid types.CID) error {
	return db.deleteByUUID(ctx, id, cid, "DeleteLifecycle", "lifecycle", db.logger)
}

func (lc lifecycle) children(db *Conn, ctx context.Context, cid types.CID, p *rpttree) error {
	var err error

	deferred, l := initAccessFuncs("lifecycle::children", db.logger, lc.UUID, cid)
	defer deferred(&err, l)

	notes, err := db.notesReport(ctx, lc.UUID, cid, p)
	if err != nil {
		return err
	} else if len(notes) != 0 {
		p.data["notes"] = notes
	}

	photos, err := db.photosReport(ctx, lc.Strain.UUID, cid, p)
	if err != nil {
		return err
	} else if len(photos) != 0 {
		p.data["strain"].(map[string]interface{})["photos"] = photos
	}

	return nil
}

func (db *Conn) LifecycleReport(ctx context.Context, id types.UUID, cid types.CID) (types.Entity, error) {
	var err error

	deferred, l := initAccessFuncs("LifecycleReport", db.logger, id, cid)
	defer deferred(&err, l)

	var result []types.Entity

	param, err := types.NewReportAttrs(url.Values{"lifecycle-id": {string(id)}})
	if err != nil {
		return nil, err
	} else if result, err = db.lifecycleReport(ctx, param, cid, nil); err != nil {
		return nil, err
	} else if len(result) == 0 {
		err = sql.ErrNoRows
		return nil, err
	}

	return result[0], nil
}

func (db *Conn) lifecycleReport(ctx context.Context, params types.ReportAttrs, cid types.CID, p *rpttree) ([]types.Entity, error) {
	var err error
	var rpt rpt

	deferred, l := initAccessFuncs("lifecycleReport", db.logger, "nil", cid)
	defer deferred(&err, l)

	lcs, err := db.selectLifecycles(ctx, params, cid)
	if err != nil {
		return nil, err
	}

	result := make([]types.Entity, 0, len(lcs))
	for _, lc := range lcs {
		if err = db.GetAllAttributes(ctx, &lc.Strain, cid); err != nil {
			return nil, err
		} else if err = db.GetAllIngredients(ctx, &lc.GrainSubstrate, cid); err != nil {
			return nil, err
		} else if err = db.GetAllIngredients(ctx, &lc.BulkSubstrate, cid); err != nil {
			return nil, err
		} else if err = db.notesAndPhotos(ctx, lc.Events, lc.UUID, cid); err != nil {
			return nil, err
		} else if rpt, err = db.newRpt(ctx, lifecycle(lc), cid, p); err != nil {
			return nil, err
		} else if rpt == nil {
			break
		} else {
			result = append(result, rpt.Data())
		}
	}

	return result, nil
}
