package data

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jsmit257/huautla/types"
)

type (
	attr       types.StrainAttribute
	event      types.Event
	eventtype  types.EventType
	generation types.Generation
	ingredient types.Ingredient
	lifecycle  types.Lifecycle
	note       types.Note
	photo      types.Photo
	strain     types.Strain
	substrate  types.Substrate
	vendor     types.Vendor

	rpt interface {
		Data() types.Entity
	}

	rpttree struct {
		id     string
		data   types.Entity
		parent *rpttree
	}

	rpttrees []*rpttree
)

func (db *Conn) newRpt(ctx context.Context, e any, cid types.CID, p *rpttree) (rpt, error) {
	result := &rpttree{
		id:     rptID(e),
		data:   make(types.Entity),
		parent: p,
	}

	if result.id == "" {
		return nil, fmt.Errorf("couldn't determine entity type: '%v' '%T'", e, e)
	} else if result.cycle(result.id) {
		return nil, nil //fmt.Errorf("detected cycle in the graph at: '%s', %v", result.id, result.parents().String())
	}

	js, err := json.Marshal(e)
	if err != nil {
		return nil, fmt.Errorf("marshal error: '%w' for : '%w'", err, result.err())
	} else if err = json.Unmarshal(js, &result.data); err != nil {
		return nil, fmt.Errorf("unmarshal error: '%w' for : '%w'", err, result.err())
	} else if T, ok := e.(interface {
		children(*Conn, context.Context, types.CID, *rpttree) error
	}); ok {
		err = T.children(db, ctx, cid, result)
	}

	return result, err
}

func (r *rpttree) Data() types.Entity {
	return r.data
}

func (r *rpttree) cycle(test string) bool {
	if r.parent == nil {
		return false
	} else if r.parent.id == test {
		return true
	}

	return r.parent.cycle(test)
}

func (r *rpttree) err() error {
	return fmt.Errorf("rpttree: '%v'", r.parents())
}

func (r *rpttree) parents() rpttrees {
	if r.parent == nil {
		return nil
	}
	return append(r.parent.parents(), r.parent)
}

func rptID(e any) string {
	switch T := e.(type) {
	case lifecycle:
		return fmt.Sprintf("lifecycle#%s", T.UUID)
	case generation:
		return fmt.Sprintf("generation#%s", T.UUID)
	case strain:
		return fmt.Sprintf("strain#%s", T.UUID)
	case types.StrainAttribute:
		return fmt.Sprintf("strainattribute#%s", T.UUID)
	case substrate:
		return fmt.Sprintf("substrate#%s", T.UUID)
	case types.Ingredient:
		return fmt.Sprintf("ingredient#%s", T.UUID)
	case types.Event:
		return fmt.Sprintf("event#%s", T.UUID)
	case eventtype:
		return fmt.Sprintf("eventtype#%s", T.UUID)
	case types.Stage:
		return fmt.Sprintf("stage#%s", T.UUID)
	case types.Photo:
		return fmt.Sprintf("photo#%s", T.UUID)
	case types.Note:
		return fmt.Sprintf("note#%s", T.UUID)
	case types.Vendor:
		return fmt.Sprintf("vendor#%s", T.UUID)
	case vendor:
		return fmt.Sprintf("vendor#%s", T.UUID)
	}
	return ""
}
