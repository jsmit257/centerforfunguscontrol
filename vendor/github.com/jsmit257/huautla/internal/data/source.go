package data

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jsmit257/huautla/types"
)

func (db *Conn) GetSources(ctx context.Context, g *types.Generation, cid types.CID) error {
	var err error
	deferred, l := initAccessFuncs("GetSources", db.logger, g.UUID, cid)
	defer deferred(&err, l)

	var rows *sql.Rows

	if rows, err = db.query.QueryContext(ctx, psqls["source"]["get"], g.UUID); err != nil {
		return err
	}

	defer rows.Close()

	var lcID *types.UUID
	var progenitor types.UUID

	for rows.Next() {
		row := types.Source{}

		if err = rows.Scan(
			&row.UUID,
			&row.Type,
			&progenitor,
			&lcID,
			&row.Strain.UUID,
			&row.Strain.Name,
			&row.Strain.Species,
			&row.Strain.CTime,
			&row.Strain.DTime,
			&row.Strain.Vendor.UUID,
			&row.Strain.Vendor.Name,
			&row.Strain.Vendor.Website,
		); err != nil {
			break
		}

		if lcID != nil {
			var lc types.Lifecycle
			if lc, err = db.SelectLifecycle(ctx, *lcID, cid); err != nil {
				break
			} else {
				row.Lifecycle = &lc
			}

			for _, e := range row.Lifecycle.Events {
				if e.UUID == progenitor {
					row.Lifecycle.Events = []types.Event{e}
					break
				}
			}
		}

		g.Sources = append(g.Sources, row)
	}

	return err
}

func (db *Conn) InsertSource(ctx context.Context, genid types.UUID, origin string, s types.Source, cid types.CID) (types.Source, error) {
	var err error
	deferred, l := initAccessFuncs("InsertSource", db.logger, genid, cid)
	defer deferred(&err, l)

	s.UUID = types.UUID(db.generateUUID().String())

	progenitor := s.Strain.UUID
	if origin == "event" {
		progenitor = s.Lifecycle.Events[0].UUID
	} else if origin != "strain" {
		return types.Source{}, fmt.Errorf("only origins of type 'strain' and 'event' are allowed: '%s'", origin)
	}

	var result sql.Result
	if result, err = db.ExecContext(ctx, psqls["source"]["add"],
		s.UUID,
		s.Type,
		progenitor,
		genid,
	); err != nil {
		if isPrimaryKeyViolation(err) {
			return db.InsertSource(ctx, genid, origin, s, cid)
		}
		return types.Source{}, err
	} else if rows, err := result.RowsAffected(); err != nil {
		return types.Source{}, err
	} else if rows != 1 {
		return types.Source{}, fmt.Errorf("source was not added")
	}

	return s, nil
}

func (db *Conn) UpdateSource(ctx context.Context, origin string, s types.Source, cid types.CID) error {
	var err error
	deferred, l := initAccessFuncs("UpdateSource", db.logger, nil, cid)
	defer deferred(&err, l)

	_ = s.Strain.UUID
	if origin == "event" {
		_ = s.Lifecycle.Events[0].UUID
	} else if origin != "strain" {
		return fmt.Errorf("only origins of type 'strain' and 'event' are allowed: '%s'", origin)
	}

	var result sql.Result

	result, err = db.ExecContext(ctx, psqls["source"]["change"], s.Type, s.UUID)
	if err != nil {
		return err
	} else if rows, err := result.RowsAffected(); err != nil {
		return err
	} else if rows != 1 { // most likely cause is a bad eventtype.uuid
		return fmt.Errorf("source was not changed")
	}

	return err
}

func (db *Conn) RemoveSource(ctx context.Context, g *types.Generation, id types.UUID, cid types.CID) error {
	var err error
	deferred, l := initAccessFuncs("RemoveSource", db.logger, g.UUID, cid)
	defer deferred(&err, l)

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
