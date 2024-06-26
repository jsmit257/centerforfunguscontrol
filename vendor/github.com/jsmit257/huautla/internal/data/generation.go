package data

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jsmit257/huautla/types"
)

func (db *Conn) SelectGenerationIndex(ctx context.Context, cid types.CID) ([]types.Generation, error) {
	var err error

	deferred, start, l := initAccessFuncs("SelectGenerationIndex", db.logger, "nil", cid)
	defer deferred(start, err, l)

	var rows *sql.Rows

	result := make([]types.Generation, 0, 100)

	rows, err = db.query.QueryContext(ctx, psqls["generation"]["ndx"])
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var row *types.Generation

	var lcID *types.UUID

	for rows.Next() {
		temp := types.Generation{Sources: []types.Source{{}}}

		if err = rows.Scan(
			&temp.UUID,
			&temp.PlatingSubstrate.UUID,
			&temp.PlatingSubstrate.Name,
			&temp.PlatingSubstrate.Type,
			&temp.PlatingSubstrate.Vendor.UUID,
			&temp.PlatingSubstrate.Vendor.Name,
			&temp.PlatingSubstrate.Vendor.Website,
			&temp.LiquidSubstrate.UUID,
			&temp.LiquidSubstrate.Name,
			&temp.LiquidSubstrate.Type,
			&temp.LiquidSubstrate.Vendor.UUID,
			&temp.LiquidSubstrate.Vendor.Name,
			&temp.LiquidSubstrate.Vendor.Website,
			&temp.Sources[0].UUID,
			&temp.Sources[0].Type,
			&lcID,
			&temp.Sources[0].Strain.UUID,
			&temp.Sources[0].Strain.Name,
			&temp.Sources[0].Strain.Species,
			&temp.Sources[0].Strain.CTime,
			&temp.Sources[0].Strain.Vendor.UUID,
			&temp.Sources[0].Strain.Vendor.Name,
			&temp.Sources[0].Strain.Vendor.Website,
			&temp.MTime,
			&temp.CTime,
		); err != nil {
			return result, err
		}

		if lcID != nil {
			temp.Sources[0].Lifecycle = &types.Lifecycle{UUID: *lcID}
		}

		if row != nil {
			if row.UUID == temp.UUID {
				temp.Sources = append(temp.Sources, row.Sources...)
			} else {
				result = append(result, *row)
			}
		}
		row = &temp
	}

	if row != nil {
		result = append(result, *row)
	}

	return result, err
}

func (db *Conn) SelectGeneration(ctx context.Context, id types.UUID, cid types.CID) (types.Generation, error) {
	var err error

	deferred, start, l := initAccessFuncs("SelectGeneration", db.logger, id, cid)
	defer deferred(start, err, l)

	result := types.Generation{UUID: id}

	if err = db.
		QueryRowContext(ctx, psqls["generation"]["select"], id).
		Scan(
			&result.PlatingSubstrate.UUID,
			&result.PlatingSubstrate.Name,
			&result.PlatingSubstrate.Type,
			&result.PlatingSubstrate.Vendor.UUID,
			&result.PlatingSubstrate.Vendor.Name,
			&result.PlatingSubstrate.Vendor.Website,
			&result.LiquidSubstrate.UUID,
			&result.LiquidSubstrate.Name,
			&result.LiquidSubstrate.Type,
			&result.LiquidSubstrate.Vendor.UUID,
			&result.LiquidSubstrate.Vendor.Name,
			&result.LiquidSubstrate.Vendor.Website,
			&result.MTime,
			&result.CTime,
		); err != nil {
		return result, err
	}

	if err = db.GetAllIngredients(ctx, &result.PlatingSubstrate, cid); err != nil {
		return result, err
	} else if err = db.GetAllIngredients(ctx, &result.LiquidSubstrate, cid); err != nil {
		return result, err
	} else if err = db.GetGenerationEvents(ctx, &result, cid); err != nil {
		return result, err
	} else if err = db.GetSources(ctx, &result, cid); err != nil {
		return result, err
	}

	return result, err
}

func (db *Conn) InsertGeneration(ctx context.Context, g types.Generation, cid types.CID) (types.Generation, error) {
	var err error
	var result sql.Result
	var rows int64

	g.UUID = types.UUID(db.generateUUID().String())
	g.CTime = time.Now().UTC()

	deferred, start, l := initAccessFuncs("InsertGeneration", db.logger, g.UUID, cid)
	defer deferred(start, err, l)

	if result, err = db.ExecContext(ctx, psqls["generation"]["insert"],
		g.UUID,
		g.PlatingSubstrate.UUID,
		g.LiquidSubstrate.UUID,
		g.CTime,
	); err != nil {
		if isPrimaryKeyViolation(err) {
			return db.InsertGeneration(ctx, g, cid) // FIXME: infinite loop?
		}
		return g, err
	} else if rows, err = result.RowsAffected(); err != nil {
		return g, err
	} else if rows != 1 {
		return g, fmt.Errorf("generation was not added: %d", rows)
	}

	return db.SelectGeneration(ctx, g.UUID, cid)
}

func (db *Conn) UpdateGeneration(ctx context.Context, g types.Generation, cid types.CID) (types.Generation, error) {
	var err error
	var result sql.Result
	var rows int64

	deferred, start, l := initAccessFuncs("UpdateGeneration", db.logger, g.UUID, cid)
	defer deferred(start, err, l)

	if result, err = db.ExecContext(ctx, psqls["generation"]["update"],
		g.PlatingSubstrate.UUID,
		g.LiquidSubstrate.UUID,
		g.UUID,
	); err != nil {
		return g, err
	} else if rows, err = result.RowsAffected(); err != nil {
		return g, err
	} else if rows != 1 {
		err = fmt.Errorf("generation was not updated")
	}

	return g, err
}

func (db *Conn) UpdateGenerationMTime(ctx context.Context, g *types.Generation, modified time.Time, cid types.CID) (*types.Generation, error) {
	var err error

	deferred, start, l := initAccessFuncs("UpdateGenerationMTime", db.logger, g.UUID, cid)
	defer deferred(start, err, l)

	g.MTime, err = db.updateMTime(ctx, "generations", modified, g.UUID, cid)

	return g, err
}

func (db *Conn) DeleteGeneration(ctx context.Context, id types.UUID, cid types.CID) error {
	return db.deleteByUUID(ctx, id, cid, "DeleteGeneration", "generation", db.logger)
}
