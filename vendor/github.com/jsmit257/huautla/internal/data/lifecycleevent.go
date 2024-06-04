package data

import (
	"context"
	"time"

	"github.com/jsmit257/huautla/types"
)

func (db *Conn) GetLifecycleEvents(ctx context.Context, lc *types.Lifecycle, cid types.CID) error {
	var err error

	deferred, start, l := initAccessFuncs("GetLifecycleEvents", db.logger, lc.UUID, cid)
	defer deferred(start, err, l)

	lc.Events, err = db.selectEventsList(ctx, psqls["event"]["all-by-observable"], lc.UUID, cid)

	return err
}

func (db *Conn) AddLifecycleEvent(ctx context.Context, lc *types.Lifecycle, e types.Event, cid types.CID) error {
	var err error

	deferred, start, l := initAccessFuncs("AddEvent", db.logger, lc.UUID, cid)
	defer deferred(start, err, l)

	if lc.Events, err = db.addEvent(ctx, lc.UUID, lc.Events, &e, cid); err != nil {
		return err
	} else if _, err = db.updateMTime(ctx, "lifecycles", lc.Events[0].MTime, lc.UUID, cid); err != nil {
		return err
	}
	return err
}

func (db *Conn) ChangeLifecycleEvent(ctx context.Context, lc *types.Lifecycle, e types.Event, cid types.CID) (types.Event, error) {
	var err error

	deferred, start, l := initAccessFuncs("ChangeEvent", db.logger, lc.UUID, cid)
	defer deferred(start, err, l)

	if lc.Events, err = db.changeEvent(ctx, lc.Events, &e, cid); err != nil {
		return e, err
	} else if _, err = db.updateMTime(ctx, "lifecycles", lc.Events[0].MTime, lc.UUID, cid); err != nil {
		return e, err
	}

	return e, err
}

func (db *Conn) RemoveLifecycleEvent(ctx context.Context, lc *types.Lifecycle, id types.UUID, cid types.CID) error {
	var err error

	deferred, start, l := initAccessFuncs("RemoveEvent", db.logger, lc.UUID, cid)
	defer deferred(start, err, l)

	if lc.Events, err = db.removeEvent(ctx, lc.Events, id, cid); err != nil {
		return err
	} else if _, err = db.updateMTime(ctx, "lifecycles", time.Now().UTC(), lc.UUID, cid); err != nil {
		return err
	}

	return nil
}
