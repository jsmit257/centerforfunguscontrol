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

	deferred, start, l := initAccessFuncs("GetPhotos", db.logger, id, cid)
	defer deferred(start, err, l)

	rows, err = db.query.QueryContext(ctx, psqls["photo"]["get"], id)
	if err != nil {
		return result, err
	}

	defer rows.Close()

	for rows.Next() {
		var p types.Photo
		var noteid *types.UUID
		var notetext *string
		var notemtime, notectime *time.Time

		if err = rows.Scan(
			&p.UUID,
			&p.Filename,
			&p.MTime,
			&p.CTime,
			&noteid,
			&notetext,
			&notemtime,
			&notectime,
		); err != nil {
			return result, err
		}

		if noteid != nil {
			p.Notes = []types.Note{{
				UUID:  *noteid,
				Note:  *notetext,
				MTime: *notemtime,
				CTime: *notectime,
			}}
		}

		if curr := len(result) - 1; curr == -1 || result[curr].UUID != p.UUID {
			result = append(result, p)
		} else {
			result[curr].Notes = append(result[curr].Notes, p.Notes...)
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
	p.MTime = p.CTime

	if result, err = db.ExecContext(ctx, psqls["photo"]["add"],
		p.UUID,
		p.Filename,
		id,
		p.MTime,
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

	p.MTime = time.Now().UTC()

	if result, err = db.ExecContext(ctx, psqls["photo"]["change"],
		p.Filename,
		p.MTime,
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

	if result, err = db.ExecContext(ctx, psqls["photo"]["remove"], id); err != nil {
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

func (db *Conn) photosReport(ctx context.Context, id types.UUID, cid types.CID, p *rpttree) ([]types.Entity, error) {
	photos, err := db.GetPhotos(ctx, id, cid)
	if err != nil {
		return nil, err
	} else if len(photos) == 0 {
		return nil, nil
	}

	result := make([]types.Entity, len(photos))
	for i, photo := range photos {
		rpt, err := db.newRpt(ctx, photo, cid, p)
		if err != nil {
			return nil, err
		}
		result[i] = rpt.Data()
	}

	return result, nil
}
