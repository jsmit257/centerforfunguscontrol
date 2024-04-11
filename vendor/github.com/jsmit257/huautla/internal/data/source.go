package data

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jsmit257/huautla/types"
)

func (db *Conn) GetSources(ctx context.Context, g *types.Generation, cid types.CID) error {
	var err error

	deferred, start, l := initAccessFuncs("GetSources", db.logger, g.UUID, cid)
	defer deferred(start, err, l)

	var rows *sql.Rows

	if rows, err = db.query.QueryContext(ctx, psqls["source"]["get"], g.UUID); err != nil {
		return err
	}

	defer rows.Close()

	var lcID *types.UUID

	for rows.Next() {
		row := types.Source{}

		if err = rows.Scan(
			&row.UUID,
			&row.Type,
			&lcID,
			&row.Strain.UUID,
			&row.Strain.Name,
			&row.Strain.Species,
			&row.Strain.CTime,
			&row.Strain.Vendor.UUID,
			&row.Strain.Vendor.Name,
			&row.Strain.Vendor.Website,
		); err != nil {
			return err
		}

		if lcID != nil {
			row.Lifecycle = &types.Lifecycle{UUID: *lcID}
		}

		g.Sources = append(g.Sources, row)
	}

	return err
}

func (db *Conn) AddStrainSource(ctx context.Context, g *types.Generation, s types.Source, cid types.CID) error {
	var err error

	deferred, start, l := initAccessFuncs("AddStrainSource", db.logger, g.UUID, cid)
	defer deferred(start, err, l)

	s.UUID = types.UUID(db.generateUUID().String())
	s.CTime = time.Now().UTC()

	g.Sources, err = db.addSource(ctx, g.Sources, s, s.Strain.UUID, g.UUID, cid)

	return err
}

func (db *Conn) AddEventSource(ctx context.Context, g *types.Generation, e types.Event, cid types.CID) error {
	var err error

	deferred, start, l := initAccessFuncs("AddEventSource", db.logger, g.UUID, cid)
	defer deferred(start, err, l)

	s := types.Source{
		UUID: types.UUID(db.generateUUID().String()),
		Type: "Spore",
	}

	if s.Strain.UUID, err = db.getEventStrainID(ctx, e.UUID, cid); err != nil {
		return fmt.Errorf("couldn't get strain for AddEventSource (%#v)", e)
	} else if e.EventType, err = db.SelectEventType(ctx, e.EventType.UUID, cid); err != nil {
		return fmt.Errorf("couldn't get eventtype for AddEventSource")
	} else if e.EventType.Name == "Clone" {
		s.Type = "Clone"
	}

	g.Sources, err = db.addSource(ctx, g.Sources, s, e.UUID, g.UUID, cid)

	return err
}

func (db *Conn) getEventStrainID(ctx context.Context, id types.UUID, cid types.CID) (types.UUID, error) {
	var err error

	deferred, start, l := initAccessFuncs("getEventStrainID", db.logger, id, cid)
	defer deferred(start, err, l)

	var result types.UUID

	err = db.
		QueryRowContext(ctx, psqls["source"]["strain-from-event"], id).
		Scan(&result)

	return result, err
}

func (db *Conn) addSource(ctx context.Context, sources []types.Source, s types.Source, progenitor, generation types.UUID, cid types.CID) ([]types.Source, error) {
	var err error
	var result sql.Result

	if result, err = db.ExecContext(ctx, psqls["source"]["add"],
		s.UUID,
		s.Type,
		progenitor,
		generation,
	); err != nil {
		if isPrimaryKeyViolation(err) {
			return db.addSource(ctx, sources, s, progenitor, generation, cid)
		}
		return sources, err
	} else if rows, err := result.RowsAffected(); err != nil {
		return sources, err
	} else if rows != 1 { // most likely cause is a bad eventtype.uuid
		return sources, fmt.Errorf("source was not added")
	}

	if s.Strain, err = db.SelectStrain(ctx, s.Strain.UUID, cid); err != nil {
		return sources, fmt.Errorf("couldn't fetch strain")
	}

	return append([]types.Source{s}, sources...), err
}

func (db *Conn) ChangeSource(ctx context.Context, g *types.Generation, s types.Source, cid types.CID) error {
	var err error

	deferred, start, l := initAccessFuncs("ChangeSource", db.logger, g.UUID, cid)
	defer deferred(start, err, l)

	var result sql.Result

	result, err = db.ExecContext(ctx, psqls["source"]["change"], s.Type, s.UUID)
	if err != nil {
		return err
	} else if rows, err := result.RowsAffected(); err != nil {
		return err
	} else if rows != 1 { // most likely cause is a bad eventtype.uuid
		return fmt.Errorf("source was not changed")
	}

	i, j := 0, len(g.Sources)
	for i < j && g.Sources[i].UUID != s.UUID {
		i++
	}

	g.Sources = append(append([]types.Source{s}, g.Sources[:i]...), g.Sources[i+1:]...)

	return err
}

func (db *Conn) RemoveSource(ctx context.Context, g *types.Generation, id types.UUID, cid types.CID) error {
	var err error

	deferred, start, l := initAccessFuncs("RemoveSource", db.logger, g.UUID, cid)
	defer deferred(start, err, l)

	if err := db.deleteByUUID(ctx, id, cid, "RemoveSource", "source", l); err != nil {
		return err
	}

	i, j := 0, len(g.Sources)
	for i < j && g.Sources[i].UUID != id {
		i++
	}

	g.Sources = append(g.Sources[:i], g.Sources[i+1:]...)

	return err
}
