package data

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jsmit257/huautla/types"
)

func (db *Conn) KnownAttributeNames(ctx context.Context, cid types.CID) ([]string, error) {
	var err error
	var result = []string{}
	var s string

	deferred, start, l := initAccessFuncs("KnownAttributeNames", db.logger, "nil", cid)
	defer deferred(start, err, l)

	var rows *sql.Rows

	rows, err = db.query.QueryContext(ctx, psqls["strainattribute"]["get-unique-names"])
	if err != nil {
		return result, err
	}

	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&s); err != nil {
			return []string{}, err
		}
		result = append(result, s)
	}

	return result, err
}

func (db *Conn) GetAllAttributes(ctx context.Context, s *types.Strain, cid types.CID) error {
	var err error

	deferred, start, l := initAccessFuncs("GetAllAttributes", db.logger, "nil", cid)
	defer deferred(start, err, l)

	var rows *sql.Rows

	s.Attributes = make([]types.StrainAttribute, 0, 100)

	rows, err = db.query.QueryContext(ctx, psqls["strainattribute"]["all"], s.UUID)
	if err != nil {
		return err
	}

	defer rows.Close()

	for rows.Next() {
		row := types.StrainAttribute{}
		if err = rows.Scan(
			&row.UUID,
			&row.Name,
			&row.Value); err != nil {

			return err
		}
		s.Attributes = append(s.Attributes, row)
	}

	return err
}

func (db *Conn) AddAttribute(ctx context.Context, s *types.Strain, n, v string, cid types.CID) error {
	var err error

	id := types.UUID(db.generateUUID().String())

	deferred, start, l := initAccessFuncs("AddAttribute", db.logger, s.UUID, cid)
	defer deferred(start, err, l)

	result, err := db.ExecContext(ctx, psqls["strainattribute"]["add"], id, n, v, s.UUID)

	if err != nil {
		if isPrimaryKeyViolation(err) {
			return db.AddAttribute(ctx, s, n, v, cid) // FIXME: infinite loop?
		}
		return err
	} else if rows, err := result.RowsAffected(); err != nil {
		return err
	} else if rows != 1 { // most likely cause is a bad vendor.uuid
		return fmt.Errorf("attribute was not added")
	}

	s.Attributes = append(s.Attributes, types.StrainAttribute{UUID: id, Name: n, Value: v})

	return err
}

func (db *Conn) ChangeAttribute(ctx context.Context, s *types.Strain, n, v string, cid types.CID) error {
	var err error

	db.logger.Errorf("fuck u bitch: '%s', '%s'", n, v)

	deferred, start, l := initAccessFuncs("ChangeAttribute", db.logger, s.UUID, cid)
	defer deferred(start, err, l)

	result, err := db.ExecContext(ctx, psqls["strainattribute"]["change"], v, n, s.UUID)

	if err != nil {
		return err
	} else if rows, err := result.RowsAffected(); err != nil {
		return err
	} else if rows != 1 { // most likely cause is a bad vendor.uuid
		return fmt.Errorf("attribute was not changed")
	}

	i, j := 0, len(s.Attributes)
	for i < j && s.Attributes[i].Name != n {
		db.logger.Errorf("fuck u bitch %d: '%s', '%s'", i, s.Attributes[i].Name, n)
		i++
	}
	db.logger.Errorf("fuck u bitch %d: '%s', '%s'", i, s.Attributes[i].Name, n)
	s.Attributes[i].Value = v

	return nil
}

func (db *Conn) RemoveAttribute(ctx context.Context, s *types.Strain, id types.UUID, cid types.CID) error {
	var err error

	deferred, start, l := initAccessFuncs("RemoveAttribute", db.logger, s.UUID, cid)
	defer deferred(start, err, l)

	result, err := db.ExecContext(ctx, psqls["strainattribute"]["remove"], id)

	if err != nil {
		return err
	} else if rows, err := result.RowsAffected(); err != nil {
		return err
	} else if rows != 1 { // most likely cause is a bad vendor.uuid
		return fmt.Errorf("attribute was not removed")
	}

	i, j := 0, len(s.Attributes)
	for i < j && s.Attributes[i].UUID != id {
		i++
	}
	s.Attributes = append(s.Attributes[:i], s.Attributes[i+1:]...)

	return nil
}
