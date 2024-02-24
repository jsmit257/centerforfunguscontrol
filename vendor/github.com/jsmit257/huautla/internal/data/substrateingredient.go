package data

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jsmit257/huautla/types"
)

func (db *Conn) GetAllIngredients(ctx context.Context, s *types.Substrate, cid types.CID) error {
	var err error

	deferred, start, l := initAccessFuncs("GetAllIngredients", db.logger, "nil", cid)
	defer deferred(start, err, l)

	var rows *sql.Rows

	s.Ingredients = make([]types.Ingredient, 0, 100)

	rows, err = db.query.QueryContext(ctx, psqls["substrate-ingredient"]["all"], s.UUID)
	if err != nil {
		return err
	}

	for rows.Next() {
		row := types.Ingredient{}
		rows.Scan(
			&row.UUID,
			&row.Name)
		s.Ingredients = append(s.Ingredients, row)
	}

	return err
}

func (db *Conn) AddIngredient(ctx context.Context, s *types.Substrate, i types.Ingredient, cid types.CID) error {
	var err error
	var result sql.Result
	var rows int64

	deferred, start, l := initAccessFuncs("AddIngredient", db.logger, "nil", cid)
	defer deferred(start, err, l)

	result, err = db.query.ExecContext(ctx, psqls["substrate-ingredient"]["add"], db.generateUUID(), s.UUID, i.UUID)
	if err != nil {
		if isPrimaryKeyViolation(err) {
			return db.AddIngredient(ctx, s, i, cid) // FIXME: infinite loop?
		}
		return err
	} else if rows, err = result.RowsAffected(); err != nil {
		return err
	} else if rows != 1 {
		err = fmt.Errorf("substrateingredient was not added")
	} else {
		s.Ingredients = append(s.Ingredients, i)
	}

	return err
}

func (db *Conn) ChangeIngredient(ctx context.Context, s *types.Substrate, oldI, newI types.Ingredient, cid types.CID) error {
	var err error
	var result sql.Result

	deferred, start, l := initAccessFuncs("ChangeIngredient", db.logger, "nil", cid)
	defer deferred(start, err, l)

	result, err = db.query.ExecContext(ctx, psqls["substrate-ingredient"]["change"], newI.UUID, s.UUID, oldI.UUID)
	if err != nil {
		return err
	} else if rows, err := result.RowsAffected(); err != nil {
		return err
	} else if rows != 1 { // most likely cause is a bad vendor.uuid
		return fmt.Errorf("substrateingredient was not changed")
	}

	i, j := 0, len(s.Ingredients)
	for i < j && s.Ingredients[i] != oldI {
		i++
	}
	s.Ingredients[i] = newI

	return err
}

func (db *Conn) RemoveIngredient(ctx context.Context, s *types.Substrate, i types.Ingredient, cid types.CID) error {
	var err error
	var result sql.Result

	deferred, start, l := initAccessFuncs("RemoveIngredient", db.logger, s.UUID, cid)
	defer deferred(start, err, l)

	result, err = db.query.ExecContext(ctx, psqls["substrate-ingredient"]["remove"], s.UUID, i.UUID)
	if err != nil {
		return err
	} else if rows, err := result.RowsAffected(); err != nil {
		return err
	} else if rows != 1 { // most likely cause is a bad vendor.uuid
		return fmt.Errorf("substrateingredient was not removed")
	}

	ndx, j := 0, len(s.Ingredients)
	for ndx < j && s.Ingredients[ndx] != i {
		ndx++
	}
	s.Ingredients = append(s.Ingredients[:ndx], s.Ingredients[ndx+1:]...)

	return err
}
