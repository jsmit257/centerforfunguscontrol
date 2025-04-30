package data

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"time"

	"github.com/jsmit257/huautla/types"
)

func (db *Conn) SelectGenerationIndex(ctx context.Context, cid types.CID) ([]types.Generation, error) {
	var err error

	deferred, l := initAccessFuncs("SelectGenerationIndex", db.logger, "nil", cid)
	defer deferred(&err, l)

	var rows *sql.Rows

	result := make([]types.Generation, 0, 100)

	rows, err = db.query.QueryContext(ctx, psqls["generation"]["ndx"])
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var row *types.Generation

	var lcID *types.UUID

	type (
		source struct {
			uuid *types.UUID
			typ  *string
		}
		vendor struct {
			uuid          *types.UUID
			name, website *string
		}
		strain struct {
			uuid *types.UUID
			name,
			species *string
			ctime *time.Time
		}
	)
	for rows.Next() {
		temp := types.Generation{}
		so := source{}
		v := vendor{}
		st := strain{}

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
			&so.uuid,
			&so.typ,
			&lcID,
			&st.uuid,
			&st.name,
			&st.species,
			&st.ctime,
			&v.uuid,
			&v.name,
			&v.website,
			&temp.MTime,
			&temp.CTime,
			&temp.DTime,
		); err != nil {
			return result, err
		}

		if so.uuid != nil {
			source := types.Source{
				UUID: *so.uuid,
				Type: *so.typ,
				Strain: types.Strain{
					UUID:    *st.uuid,
					Name:    *st.name,
					Species: *st.species,
					CTime:   *st.ctime,
					Vendor: types.Vendor{
						UUID:    *v.uuid,
						Name:    *v.name,
						Website: *v.website,
					},
				},
			}

			if lcID != nil {
				source.Lifecycle = &types.Lifecycle{UUID: *lcID}
			}

			temp.Sources = []types.Source{source}
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
	var result []types.Generation

	deferred, l := initAccessFuncs("SelectGeneration", db.logger, id, cid)
	defer deferred(&err, l)

	param, err := types.NewReportAttrs(map[string][]string{"generation-id": {string(id)}})
	if err != nil {
		return types.Generation{}, err
	}

	result, err = db.selectGenerations(ctx, param, cid)
	if err != nil {
		return types.Generation{}, err
	} else if len(result) == 1 {
		return result[0], nil
	} else {
		err = sql.ErrNoRows
	}

	return types.Generation{}, err
}

func (db *Conn) selectGenerations(ctx context.Context, p types.ReportAttrs, cid types.CID) ([]types.Generation, error) {
	var err error

	deferred, l := initAccessFuncs("selectGenerations", db.logger, "nil", cid)
	defer deferred(&err, l)

	result := make([]types.Generation, 0, 100)

	if !p.Contains("generation-id", "strain-id", "plating-id", "liquid-id", "eventtype-id") {
		err = fmt.Errorf("request doesn't contain at least 1 required field")
		return result, err
	}

	rows, err := db.query.QueryContext(ctx, psqls["generation"]["select"],
		p.Get("generation-id"),
		p.Get("strain-id"),
		p.Get("plating-id"),
		p.Get("liquid-id"),
		p.Get("eventtype-id"))
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		row := types.Generation{}

		if err = rows.Scan(
			&row.UUID,
			&row.PlatingSubstrate.UUID,
			&row.PlatingSubstrate.Name,
			&row.PlatingSubstrate.Type,
			&row.PlatingSubstrate.Vendor.UUID,
			&row.PlatingSubstrate.Vendor.Name,
			&row.PlatingSubstrate.Vendor.Website,
			&row.LiquidSubstrate.UUID,
			&row.LiquidSubstrate.Name,
			&row.LiquidSubstrate.Type,
			&row.LiquidSubstrate.Vendor.UUID,
			&row.LiquidSubstrate.Vendor.Name,
			&row.LiquidSubstrate.Vendor.Website,
			&row.MTime,
			&row.CTime,
			&row.DTime,
		); err != nil {
			break
		} else if err = db.GetGenerationEvents(ctx, &row, cid); err != nil {
			break
		} else if err = db.GetSources(ctx, &row, cid); err != nil {
			break
		}

		result = append(result, row)
	}

	return result, err
}

func (db *Conn) InsertGeneration(ctx context.Context, g types.Generation, cid types.CID) (types.Generation, error) {
	var err error
	var result sql.Result
	var rows int64

	g.UUID = types.UUID(db.generateUUID().String())
	g.CTime = time.Now().UTC()

	deferred, l := initAccessFuncs("InsertGeneration", db.logger, g.UUID, cid)
	defer deferred(&err, l)

	if result, err = db.ExecContext(ctx, psqls["generation"]["insert"],
		g.UUID,
		g.PlatingSubstrate.UUID,
		g.LiquidSubstrate.UUID,
		g.CTime,
	); err != nil {
		if isPrimaryKeyViolation(err) {
			return db.InsertGeneration(ctx, g, cid)
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

	deferred, l := initAccessFuncs("UpdateGeneration", db.logger, g.UUID, cid)
	defer deferred(&err, l)

	g.MTime = time.Now().UTC()

	if result, err = db.ExecContext(ctx, psqls["generation"]["update"],
		g.PlatingSubstrate.UUID,
		g.LiquidSubstrate.UUID,
		g.UUID,
		g.MTime,
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

	deferred, l := initAccessFuncs("UpdateGenerationMTime", db.logger, g.UUID, cid)
	defer deferred(&err, l)

	g.MTime, err = db.updateMTime(ctx, "generations", modified, g.UUID, cid)

	return g, err
}

func (db *Conn) DeleteGeneration(ctx context.Context, id types.UUID, cid types.CID) error {
	return db.deleteByUUID(ctx, id, cid, "DeleteGeneration", "generation", db.logger)
}

func (g generation) children(db *Conn, ctx context.Context, cid types.CID, p *rpttree) error {
	var err error

	deferred, l := initAccessFuncs("generation::children", db.logger, g.UUID, cid)
	defer deferred(&err, l)

	notes, err := db.notesReport(ctx, g.UUID, cid, p)
	if err != nil {
		return err
	} else if len(notes) != 0 {
		p.data["notes"] = notes
	}

	var rpt rpt
	progeny, err := db.GeneratedStrain(ctx, g.UUID, cid)
	if err != sql.ErrNoRows {
		if err != nil {
			return err
		} else if rpt, err = db.newRpt(ctx, strain(progeny), cid, p); err != nil {
			return err
		} else if rpt != nil {
			p.data["progeny"] = rpt.Data()
		}
	}

	return nil
}

func (db *Conn) GenerationReport(ctx context.Context, id types.UUID, cid types.CID) (types.Entity, error) {
	var err error

	deferred, l := initAccessFuncs("GenerationReport", db.logger, id, cid)
	defer deferred(&err, l)

	var result []types.Entity
	p, err := types.NewReportAttrs(url.Values{"generation-id": {string(id)}})
	if err != nil {
		return nil, err
	} else if result, err = db.generationReport(ctx, p, cid, nil); err != nil {
		return nil, err
	} else if len(result) == 0 {
		err = sql.ErrNoRows
		return nil, err
	}

	return result[0], nil
}

func (db *Conn) generationReport(ctx context.Context, params types.ReportAttrs, cid types.CID, p *rpttree) ([]types.Entity, error) {
	var err error

	deferred, l := initAccessFuncs("generationReport", db.logger, "nil", cid)
	defer deferred(&err, l)

	gens, err := db.selectGenerations(ctx, params, cid)
	if err != nil {
		return nil, err
	}

	var rpt rpt
	result := make([]types.Entity, 0, len(gens))
	for _, gen := range gens {
		if err = db.GetAllIngredients(ctx, &gen.PlatingSubstrate, cid); err != nil {
			return nil, err
		} else if err = db.GetAllIngredients(ctx, &gen.LiquidSubstrate, cid); err != nil {
			return nil, err
		} else if err = db.notesAndPhotos(ctx, gen.Events, gen.UUID, cid); err != nil {
			return nil, err
		} else if rpt, err = db.newRpt(ctx, generation(gen), cid, p); err != nil {
			return nil, err
		} else if rpt != nil {
			result = append(result, rpt.Data())
		}
	}
	return result, nil
}
