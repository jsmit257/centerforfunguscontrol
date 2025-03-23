package data

import (
	"context"
	"fmt"
	"net/url"

	"github.com/jsmit257/huautla/types"
)

func (db *Conn) SelectAllVendors(ctx context.Context, cid types.CID) ([]types.Vendor, error) {
	var err error
	deferred, l := initAccessFuncs("SelectAllVendors", db.logger, "nil", cid)
	defer deferred(&err, l)

	rows, err := db.query.QueryContext(ctx, psqls["vendor"]["select-all"])
	if err != nil {
		return nil, err
	}

	result := make([]types.Vendor, 0, 100)
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
	deferred, l := initAccessFuncs("SelectVendor", db.logger, id, cid)
	defer deferred(&err, l)

	result := types.Vendor{UUID: id}
	err = db.
		QueryRowContext(ctx, psqls["vendor"]["select"], id).
		Scan(&result.UUID, &result.Name, &result.Website)

	return result, err
}

func (db *Conn) InsertVendor(ctx context.Context, v types.Vendor, cid types.CID) (types.Vendor, error) {
	var err error
	deferred, l := initAccessFuncs("InsertVendor", db.logger, v.UUID, cid)
	defer deferred(&err, l)

	v.UUID = types.UUID(db.generateUUID().String())

	var rows int64
	result, err := db.ExecContext(ctx, psqls["vendor"]["insert"], v.UUID, v.Name, v.Website)
	if err != nil {
		if isPrimaryKeyViolation(err) {
			l.WithField("id", v.UUID).WithError(err).Error("da fuck?")
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
	deferred, l := initAccessFuncs("UpdateVendor", db.logger, id, cid)
	defer deferred(&err, l)

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

func (v vendor) children(db *Conn, ctx context.Context, cid types.CID, p *rpttree) error {
	var err error
	deferred, l := initAccessFuncs("vendor::children", db.logger, types.UUID(v.UUID), cid)
	defer deferred(&err, l)

	param, _ := types.NewReportAttrs(url.Values{"vendor-id": {string(v.UUID)}})

	subs, err := db.substrateReport(ctx, param, cid, p)
	if err != nil {
		return err
	} else if len(subs) != 0 {
		p.data["substrates"] = subs
	}

	strs, err := db.strainReport(ctx, param, cid, p)
	if err != nil {
		return err
	} else if len(strs) != 0 {
		p.data["strains"] = strs
	}

	return nil
}

func (db *Conn) VendorReport(ctx context.Context, id types.UUID, cid types.CID) (types.Entity, error) {
	var err error
	deferred, l := initAccessFuncs("VendorReport", db.logger, id, cid)
	defer deferred(&err, l)

	var rpt rpt
	result, err := db.SelectVendor(ctx, id, cid)
	if err != nil {
		return nil, err
	} else if rpt, err = db.newRpt(ctx, vendor(result), cid, nil); err != nil {
		return nil, err
	} else if rpt == nil {
		return nil, nil
	}

	return rpt.Data(), nil
}
