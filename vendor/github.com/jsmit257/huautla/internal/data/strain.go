package data

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"time"

	"github.com/jsmit257/huautla/types"
)

func (db *Conn) SelectAllStrains(ctx context.Context, cid types.CID) ([]types.Strain, error) {
	var err error
	deferred, l := initAccessFuncs("SelectAllStrains", db.logger, "nil", cid)
	defer deferred(&err, l)

	rows, err := db.query.QueryContext(ctx, psqls["strain"]["select-all"])
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var generationID *types.UUID
	result := make([]types.Strain, 0, 100)
	for rows.Next() {
		row := types.Strain{}

		if err = rows.Scan(
			&row.UUID,
			&row.Species,
			&row.Name,
			&row.CTime,
			&row.DTime,
			&row.Vendor.UUID,
			&row.Vendor.Name,
			&row.Vendor.Website,
			&generationID,
		); err != nil {
			break
		}

		if generationID != nil {
			row.Generation = &types.Generation{UUID: *generationID}
		}

		result = append(result, row)
	}

	return result, err
}

func (db *Conn) SelectStrain(ctx context.Context, id types.UUID, cid types.CID) (types.Strain, error) {
	var err error
	deferred, l := initAccessFuncs("SelectStrain", db.logger, id, cid)
	defer deferred(&err, l)

	p, _ := types.NewReportAttrs(url.Values{"strain-id": []string{string(id)}})

	strs, err := db.selectStrains(ctx, p, cid)
	if err != nil {
		return types.Strain{}, err
	} else if len(strs) == 1 {
		return strs[0], nil // bury the happy path in the middle
	} else {
		err = sql.ErrNoRows
	}

	return types.Strain{}, err
}

func (db *Conn) selectStrains(ctx context.Context, p types.ReportAttrs, cid types.CID) ([]types.Strain, error) {
	var err error
	deferred, l := initAccessFuncs("selectStrains", db.logger, types.UUID(fmt.Sprintf("%v", p)), cid)
	defer deferred(&err, l)

	rows, err := db.query.QueryContext(ctx, psqls["strain"]["select"],
		p.Get("strain-id"),
		p.Get("vendor-id"))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var generationID *types.UUID
	result := make([]types.Strain, 0, 100)
	for rows.Next() {
		row := types.Strain{}

		if err = rows.Scan(
			&row.UUID,
			&row.Species,
			&row.Name,
			&row.CTime,
			&row.DTime,
			&row.Vendor.UUID,
			&row.Vendor.Name,
			&row.Vendor.Website,
			&generationID,
		); err != nil {
			break
		}

		if generationID != nil {
			row.Generation = &types.Generation{UUID: *generationID}
		}

		err = db.GetAllAttributes(ctx, &row, cid)

		result = append(result, row)
	}

	return result, err
}

func (db *Conn) InsertStrain(ctx context.Context, s types.Strain, cid types.CID) (types.Strain, error) {
	var err error

	s.UUID = types.UUID(db.generateUUID().String())
	s.CTime = time.Now().UTC()

	deferred, l := initAccessFuncs("InsertStrain", db.logger, s.UUID, cid)
	defer deferred(&err, l)

	var rows int64
	result, err := db.ExecContext(ctx, psqls["strain"]["insert"], s.UUID, s.Species, s.Name, s.CTime, s.Vendor.UUID)
	if err != nil {
		if isPrimaryKeyViolation(err) {
			return db.InsertStrain(ctx, s, cid) // FIXME: infinite loop?
		}
		return s, err
	} else if rows, err = result.RowsAffected(); err != nil {
		return s, err
	} else if rows != 1 { // most likely cause is a bad vendor.uuid
		err = fmt.Errorf("strain was not added")
	}

	return s, err
}

func (db *Conn) UpdateStrain(ctx context.Context, id types.UUID, s types.Strain, cid types.CID) error {
	var err error
	deferred, l := initAccessFuncs("UpdateStrain", db.logger, id, cid)
	defer deferred(&err, l)

	var rows int64
	result, err := db.ExecContext(ctx, psqls["strain"]["update"], s.Species, s.Name, s.Vendor.UUID, id)
	if err != nil {
		return err
	} else if rows, err = result.RowsAffected(); err != nil {
		return err
	} else if rows != 1 {
		return fmt.Errorf("strain was not updated: '%s'", id)
	}

	return nil
}

func (db *Conn) DeleteStrain(ctx context.Context, id types.UUID, cid types.CID) error {
	return db.deleteByUUID(ctx, id, cid, "DeleteStrain", "strain", db.logger)
}

func (db *Conn) GeneratedStrain(ctx context.Context, id types.UUID, cid types.CID) (types.Strain, error) {
	var err error
	deferred, l := initAccessFuncs("GeneratedStrains", db.logger, id, cid)
	defer deferred(&err, l)

	result := types.Strain{}

	return result, db.
		QueryRowContext(ctx, psqls["strain"]["generated-strain"], id).
		Scan(
			&result.UUID,
			&result.Species,
			&result.Name,
			&result.CTime,
			&result.DTime,
			&result.Vendor.UUID,
			&result.Vendor.Name,
			&result.Vendor.Website,
		)
}

func (db *Conn) UpdateGeneratedStrain(ctx context.Context, gid *types.UUID, sid types.UUID, cid types.CID) error {
	var err error
	deferred, l := initAccessFuncs("UpdateGeneratedStrain", db.logger, sid, cid)
	defer deferred(&err, l)

	var rows int64
	result, err := db.ExecContext(ctx, psqls["strain"]["update-gen-strain"], gid, sid)
	if err != nil {
		return err
	} else if rows, err = result.RowsAffected(); err != nil {
		return err
	} else if rows != 1 {
		return sql.ErrNoRows
	}

	return nil
}

func (s strain) children(db *Conn, ctx context.Context, cid types.CID, p *rpttree) error {
	var err error
	deferred, l := initAccessFuncs("strain::children", db.logger, s.UUID, cid)
	defer deferred(&err, l)

	param, _ := types.NewReportAttrs(url.Values{"strain-id": {string(s.UUID)}})

	gens, err := db.generationReport(ctx, param, cid, p)
	if err != nil {
		return err
	} else if len(gens) != 0 {
		p.data["generations"] = gens
	}

	lcs, err := db.lifecycleReport(ctx, param, cid, p)
	if err != nil {
		return err
	} else if len(lcs) != 0 {
		p.data["lifecycles"] = lcs
	}

	photos, err := db.photosReport(ctx, s.UUID, cid, p)
	if err != nil {
		return err
	} else if len(photos) != 0 {
		p.data["photos"] = photos
	}

	if s.Generation == nil {
		return nil
	} else if param, err = types.NewReportAttrs(url.Values{"generation-id": {string(s.Generation.UUID)}}); err != nil {
		return err
	} else if gens, err = db.generationReport(ctx, param, cid, p); err != nil {
		return err
	} else if len(gens) == 0 {
		err = fmt.Errorf("how does '%s' not identify a generation?", s.Generation.UUID)
	} else {
		p.data["generation"] = gens[0]
	}

	return err
}

func (db *Conn) StrainReport(ctx context.Context, id types.UUID, cid types.CID) (types.Entity, error) {
	var err error
	deferred, l := initAccessFuncs("StrainReport", db.logger, "nil", cid)
	defer deferred(&err, l)

	var result []types.Entity
	param, err := types.NewReportAttrs(url.Values{"strain-id": {string(id)}})
	if err != nil {
		return nil, err
	} else if result, err = db.strainReport(ctx, param, cid, nil); err != nil {
		return nil, err
	} else if len(result) == 1 {
		return result[0], nil
	} else {
		err = sql.ErrNoRows
	}

	return nil, err
}

func (db *Conn) strainReport(ctx context.Context, params types.ReportAttrs, cid types.CID, p *rpttree) ([]types.Entity, error) {
	var err error
	deferred, l := initAccessFuncs("strainReport", db.logger, "nil", cid)
	defer deferred(&err, l)

	var rpt rpt
	strs, err := db.selectStrains(ctx, params, cid)
	if err != nil {
		return nil, err
	}

	result := make([]types.Entity, 0, len(strs))
	for _, str := range strs {
		if rpt, err = db.newRpt(ctx, strain(str), cid, p); err != nil {
			return nil, err
		} else if rpt == nil {
			break
		} else {
			result = append(result, rpt.Data())
		}
	}

	return result, nil
}
