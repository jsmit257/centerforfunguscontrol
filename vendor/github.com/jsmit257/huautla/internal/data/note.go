package data

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jsmit257/huautla/types"
)

func (db *Conn) GetNotes(ctx context.Context, id types.UUID, cid types.CID) ([]types.Note, error) {
	var err error
	var rows *sql.Rows
	var result []types.Note

	deferred, start, l := initAccessFuncs("GetNotes", db.logger, id, cid)
	defer deferred(start, err, l)

	rows, err = db.query.QueryContext(ctx, psqls["note"]["get"], id)
	if err != nil {
		return result, err
	}

	defer rows.Close()

	note := types.Note{}

	for rows.Next() {
		if err = rows.Scan(
			&note.UUID,
			&note.Note,
			&note.MTime,
			&note.CTime,
		); err != nil {
			return result, err
		}
		result = append(result, note)
	}

	return result, nil
}

func (db *Conn) AddNote(ctx context.Context, oID types.UUID, notes []types.Note, n types.Note, cid types.CID) ([]types.Note, error) {
	var err error
	var result sql.Result

	deferred, start, l := initAccessFuncs("AddNote", db.logger, oID, cid)
	defer deferred(start, err, l)

	n.UUID = types.UUID(db.generateUUID().String())
	n.MTime = time.Now().UTC()
	n.CTime = n.MTime

	if result, err = db.ExecContext(ctx, psqls["note"]["add"],
		n.UUID,
		n.Note,
		oID,
		n.CTime,
	); err != nil {
		if isPrimaryKeyViolation(err) {
			return db.AddNote(ctx, oID, notes, n, cid)
		}
		return notes, err
	} else if rows, err := result.RowsAffected(); err != nil {
		return notes, err
	} else if rows != 1 {
		return notes, fmt.Errorf("note was not added")
	}

	return append([]types.Note{n}, notes...), err
}

func (db *Conn) ChangeNote(ctx context.Context, notes []types.Note, n types.Note, cid types.CID) ([]types.Note, error) {
	var err error
	var result sql.Result

	deferred, start, l := initAccessFuncs("ChangeNote", db.logger, n.UUID, cid)
	defer deferred(start, err, l)

	n.MTime = time.Now().UTC()

	if result, err = db.ExecContext(ctx, psqls["note"]["change"],
		n.Note,
		n.MTime,
		n.UUID,
	); err != nil {
		return notes, err
	} else if rows, err := result.RowsAffected(); err != nil {
		return notes, err
	} else if rows != 1 {
		return notes, fmt.Errorf("note was not changed")
	}

	i, j := 0, len(notes)
	for i < j && notes[i].UUID != n.UUID {
		i++
	}

	return append(append([]types.Note{n}, notes[:i]...), notes[i+1:]...), nil
}

func (db *Conn) RemoveNote(ctx context.Context, notes []types.Note, id types.UUID, cid types.CID) ([]types.Note, error) {
	var err error
	var result sql.Result

	deferred, start, l := initAccessFuncs("RemoveNote", db.logger, id, cid)
	defer deferred(start, err, l)

	if result, err = db.ExecContext(ctx, psqls["note"]["remove"], id); err != nil {
		return notes, err
	} else if rows, err := result.RowsAffected(); err != nil {
		return notes, err
	} else if rows != 1 {
		return notes, fmt.Errorf("note could not be removed")
	}

	i, j := 0, len(notes)
	for i < j && notes[i].UUID != id {
		i++
	}

	return append(notes[:i], notes[i+1:]...), nil
}
