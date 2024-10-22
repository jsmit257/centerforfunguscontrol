package data

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"

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

	deferred, start, l := initAccessFuncs("SelectSubstrate", db.logger, "nil", cid)
	defer deferred(start, err, l)

	param, err := types.NewReportAttrs(url.Values{"substrate-id": []string{string(id)}})
	if err != nil {
		return types.Substrate{}, err
	}

	subs, err := db.selectSubstrateByAttrs(ctx, param, cid)
	if err != nil {
		return types.Substrate{}, err
	} else if l := len(subs); l == 1 {
		return subs[0], nil
	} else {
		err = sql.ErrNoRows
	}

	return types.Substrate{}, err
}

func (db *Conn) selectSubstrateByAttrs(ctx context.Context, param types.ReportAttrs, cid types.CID) ([]types.Substrate, error) {
	var err error

	db.logger.WithField("param", param).Errorf("incoming!!!")

	deferred, start, l := initAccessFuncs("SelectSubstrateByAttrs", db.logger, "nil", cid)
	defer deferred(start, err, l)

	var result []types.Substrate

	if !param.Contains("substrate-id", "vendor-id") {
		err = fmt.Errorf("request doesn't contain at least 1 required field: %#v", param)
		return result, err
	}

	rows, err := db.QueryContext(ctx, psqls["substrate"]["select"],
		param.Get("substrate-id"),
		param.Get("vendor-id"))
	if err != nil {
		return result, err
	}

	for rows.Next() {
		row := types.Substrate{}
		if err = rows.Scan(
			&row.UUID,
			&row.Name,
			&row.Type,
			&row.Vendor.UUID,
			&row.Vendor.Name,
			&row.Vendor.Website,
		); err != nil {
			break
		} else if err = db.GetAllIngredients(ctx, &row, "SelectSubstrate"); err != nil {
			break
		}
		result = append(result, row)
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

	result, err := db.ExecContext(ctx, psqls["substrate"]["update"], s.Name, s.Type, s.Vendor.UUID, id)
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

type Substrate types.Substrate

func (s Substrate) children(db *Conn, ctx context.Context, cid types.CID, p *rpttree) error {
	var err error

	deferred, start, l := initAccessFuncs("Substrate::children", db.logger, s.UUID, cid)
	defer deferred(start, err, l)

	getter := db.lifecycleReport
	key := "lifecycles"
	switch s.Type {
	case types.PlatingType, types.LiquidType:
		getter = db.generationReport
		key = "generations"
	}

	var values []types.Entity
	param, err := types.NewReportAttrs(url.Values{fmt.Sprintf("%s-id", s.Type): []string{string(s.UUID)}})
	if err != nil {
		return err
	} else if values, err = getter(ctx, param, cid, p); err != nil {
		return err
	} else if len(values) != 0 {
		p.data[key] = values
	}

	return nil
}

func (db *Conn) SubstrateReport(ctx context.Context, id types.UUID, cid types.CID) (types.Entity, error) {
	var err error

	deferred, start, l := initAccessFuncs("SubstrateReport", db.logger, id, cid)
	defer deferred(start, err, l)

	var rpt rpt
	sub, err := db.SelectSubstrate(ctx, id, cid)
	if err != nil {
		return nil, err
	} else if rpt, err = db.newRpt(ctx, Substrate(sub), cid, nil); err != nil {
		return nil, err
	} else if rpt == nil {
		return nil, nil
	}

	return rpt.Data(), nil
}
