package data

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jsmit257/huautla/types"
)

func (db *Conn) GetPhotos(ctx context.Context, id types.UUID, cid types.CID) ([]types.Photo, error) {
	var err error
	var rows *sql.Rows
	var result []types.Photo
	var p types.Photo

	deferred, start, l := initAccessFuncs("GetPhotos", db.logger, id, cid)
	defer deferred(start, err, l)

	rows, err = db.query.QueryContext(ctx, psqls["eventphoto"]["get"], id)
	if err != nil {
		return result, err
	}

	defer rows.Close()

	for rows.Next() {
		p = types.Photo{}
		var note_id *types.UUID
		var note *string
		var mtime, ctime *time.Time

		if err = rows.Scan(
			&p.UUID,
			&p.Filename,
			&p.CTime,
			&note_id,
			&note,
			&mtime,
			&ctime,
		); err != nil {
			return result, err
		} else if note_id != nil {
			p.Notes = []types.Note{{
				UUID:  *note_id,
				Note:  *note,
				MTime: *mtime,
				CTime: *ctime,
			}}
		}

		l := len(result)
		if l == 0 || result[l-1].UUID != p.UUID {
			result = append(result, p)
		} else {
			result[l-1].Notes = append(p.Notes, result[l-1].Notes...)
		}
	}

	return result, nil
}

func (db *Conn) AddPhoto(ctx context.Context, id types.UUID, photos []types.Photo, p types.Photo, cid types.CID) ([]types.Photo, error) {
	var err error
	var result sql.Result

	deferred, start, l := initAccessFuncs("AddPhoto", db.logger, id, cid)
	defer deferred(start, err, l)

	p.UUID = types.UUID(db.generateUUID().String())
	p.CTime = time.Now().UTC()

	if result, err = db.ExecContext(ctx, psqls["eventphoto"]["add"],
		p.UUID,
		p.Filename,
		id,
		p.CTime,
	); err != nil {
		if isPrimaryKeyViolation(err) {
			return db.AddPhoto(ctx, id, photos, p, cid)
		}
		return photos, err
	} else if rows, err := result.RowsAffected(); err != nil {
		return photos, err
	} else if rows != 1 {
		return photos, fmt.Errorf("photo was not added")
	}

	return append([]types.Photo{p}, photos...), err
}

func (db *Conn) ChangePhoto(ctx context.Context, photos []types.Photo, p types.Photo, cid types.CID) ([]types.Photo, error) {
	var err error
	var result sql.Result

	deferred, start, l := initAccessFuncs("ChangePhoto", db.logger, p.UUID, cid)
	defer deferred(start, err, l)

	if result, err = db.ExecContext(ctx, psqls["eventphoto"]["change"],
		p.Filename,
		p.UUID,
	); err != nil {
		return photos, err
	} else if rows, err := result.RowsAffected(); err != nil {
		return photos, err
	} else if rows != 1 {
		return photos, fmt.Errorf("photo was not changed")
	}

	i, j := 0, len(photos)
	for i < j && photos[i].UUID != p.UUID {
		i++
	}

	return append(append([]types.Photo{p}, photos[:i]...), photos[i+1:]...), nil
}

func (db *Conn) RemovePhoto(ctx context.Context, photos []types.Photo, id types.UUID, cid types.CID) ([]types.Photo, error) {
	var err error
	var result sql.Result

	deferred, start, l := initAccessFuncs("RemovePhoto", db.logger, id, cid)
	defer deferred(start, err, l)

	if result, err = db.ExecContext(ctx, psqls["eventphoto"]["remove"], id); err != nil {
		return photos, err
	} else if rows, err := result.RowsAffected(); err != nil {
		return photos, err
	} else if rows != 1 {
		return photos, fmt.Errorf("photo could not be removed")
	}

	i, j := 0, len(photos)
	for i < j && photos[i].UUID != id {
		i++
	}

	return append(photos[:i], photos[i+1:]...), nil
}
