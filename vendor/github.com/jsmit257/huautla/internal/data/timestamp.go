package data

import (
	"context"
	"fmt"
	"time"

	"github.com/jsmit257/huautla/types"
)

func (db *Conn) updateMTime(ctx context.Context, table string, modified time.Time, id types.UUID, _ types.CID) (time.Time, error) {
	var rows int64

	if result, err := db.ExecContext(
		ctx,
		fmt.Sprintf(psqls["timestamp"]["touch"], table),
		modified,
		id,
	); err != nil {
		return modified, err
	} else if rows, err = result.RowsAffected(); err != nil {
		return modified, err
	} else if rows != 1 {
		return modified, fmt.Errorf("mtime was not updated")
	}

	return modified, nil
}

func (db *Conn) UpdateTimestamps(ctx context.Context, table string, id types.UUID, data types.Timestamp) error {
	var rows int64

	if updt, err := data.UpdateString(); err != nil {
		return err
	} else if result, err := db.ExecContext(
		ctx,
		fmt.Sprintf(psqls["timestamp"]["update"], "uuids" /*table*/, updt),
		id,
	); err != nil {
		return err
	} else if rows, err = result.RowsAffected(); err != nil {
		return err
	} else if rows != 1 {
		return fmt.Errorf("timestamps were not updated")
	}

	return nil
}

func (db *Conn) Undelete(ctx context.Context, table string, id types.UUID) error {
	var rows int64

	if result, err := db.ExecContext(
		ctx,
		fmt.Sprintf(psqls["timestamp"]["undelete"], "uuids" /*table*/),
		id,
	); err != nil {
		return err
	} else if rows, err = result.RowsAffected(); err != nil {
		return err
	} else if rows != 1 {
		return fmt.Errorf("record could not be undeleted")
	}

	return nil
}
