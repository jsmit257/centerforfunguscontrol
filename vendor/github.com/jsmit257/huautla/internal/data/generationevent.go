package data

import (
	"context"
	"fmt"
	"time"

	"github.com/jsmit257/huautla/types"
)

func (db *Conn) GetGenerationEvents(ctx context.Context, g *types.Generation, cid types.CID) error {
	var err error

	deferred, start, l := initAccessFuncs("GetGenerationEvents", db.logger, g.UUID, cid)
	defer deferred(start, err, l)

	g.Events, err = db.selectEventsList(ctx, psqls["event"]["all-by-observable"], g.UUID, cid)

	return err
}

func (db *Conn) AddGenerationEvent(ctx context.Context, g *types.Generation, e types.Event, cid types.CID) error {
	var err error

	deferred, start, l := initAccessFuncs("AddGenerationEvent", db.logger, g.UUID, cid)
	defer deferred(start, err, l)

	if g.Events, err = db.addEvent(ctx, g.UUID, g.Events, &e, cid); err != nil {
		return err
	} else if _, err = db.UpdateGenerationMTime(ctx, g, e.MTime, cid); err != nil {
		return fmt.Errorf("couldn't update Generation.mtime")
	}
	return err
}

func (db *Conn) ChangeGenerationEvent(ctx context.Context, g *types.Generation, e types.Event, cid types.CID) (types.Event, error) {
	var err error

	deferred, start, l := initAccessFuncs("ChangeEvent", db.logger, g.UUID, cid)
	defer deferred(start, err, l)

	if g.Events, err = db.changeEvent(ctx, g.Events, &e, cid); err != nil {
		return e, err
	} else if _, err = db.UpdateGenerationMTime(ctx, g, g.Events[0].MTime, cid); err != nil {
		return e, err
	}

	return e, err
}

func (db *Conn) RemoveGenerationEvent(ctx context.Context, g *types.Generation, id types.UUID, cid types.CID) error {
	var err error

	deferred, start, l := initAccessFuncs("RemoveEvent", db.logger, g.UUID, cid)
	defer deferred(start, err, l)

	if g.Events, err = db.removeEvent(ctx, g.Events, id, cid); err != nil {
		return err
	} else if _, err = db.UpdateGenerationMTime(ctx, g, time.Now().UTC(), cid); err != nil {
		return err
	}

	return nil
}
