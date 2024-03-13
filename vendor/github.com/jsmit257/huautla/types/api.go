package types

import "context"

type (
	DB interface {
		Eventer
		EventTyper
		Ingredienter
		Lifecycler
		Stager
		StrainAttributer
		Strainer
		SubstrateIngredienter
		Substrater
		Vendorer
	}

	Eventer interface {
		SelectByEventType(ctx context.Context, et EventType, cid CID) ([]Event, error)
		SelectEvent(ctx context.Context, id UUID, cid CID) (Event, error)
		GetLifecycleEvents(ctx context.Context, lc *Lifecycle, cid CID) error
		AddEvent(ctx context.Context, lc *Lifecycle, e Event, cid CID) error
		ChangeEvent(ctx context.Context, lc *Lifecycle, e Event, cid CID) (Event, error)
		RemoveEvent(ctx context.Context, lc *Lifecycle, id UUID, cid CID) error
	}

	EventTyper interface {
		SelectAllEventTypes(ctx context.Context, cid CID) ([]EventType, error)
		SelectEventType(ctx context.Context, id UUID, cid CID) (EventType, error)
		InsertEventType(ctx context.Context, e EventType, cid CID) (EventType, error)
		UpdateEventType(ctx context.Context, id UUID, e EventType, cid CID) error
		DeleteEventType(ctx context.Context, id UUID, cid CID) error
	}

	Ingredienter interface {
		SelectAllIngredients(ctx context.Context, cid CID) ([]Ingredient, error)
		SelectIngredient(ctx context.Context, id UUID, cid CID) (Ingredient, error)
		InsertIngredient(ctx context.Context, i Ingredient, cid CID) (Ingredient, error)
		UpdateIngredient(ctx context.Context, id UUID, i Ingredient, cid CID) error
		DeleteIngredient(ctx context.Context, id UUID, cid CID) error
	}

	Lifecycler interface {
		SelectLifecycleIndex(ctx context.Context, cid CID) ([]Lifecycle, error)
		SelectLifecycle(ctx context.Context, id UUID, cid CID) (Lifecycle, error)
		InsertLifecycle(ctx context.Context, lc Lifecycle, cid CID) (Lifecycle, error)
		UpdateLifecycle(ctx context.Context, lc Lifecycle, cid CID) (Lifecycle, error)
		DeleteLifecycle(ctx context.Context, id UUID, cid CID) error
	}

	Stager interface {
		SelectAllStages(ctx context.Context, cid CID) ([]Stage, error)
		SelectStage(ctx context.Context, id UUID, cid CID) (Stage, error)
		InsertStage(ctx context.Context, s Stage, cid CID) (Stage, error)
		UpdateStage(ctx context.Context, id UUID, s Stage, cid CID) error
		DeleteStage(ctx context.Context, id UUID, cid CID) error
	}

	StrainAttributer interface {
		KnownAttributeNames(ctx context.Context, cid CID) ([]string, error)
		GetAllAttributes(ctx context.Context, s *Strain, cid CID) error
		AddAttribute(ctx context.Context, s *Strain, a StrainAttribute, cid CID) (StrainAttribute, error)
		ChangeAttribute(ctx context.Context, s *Strain, a StrainAttribute, cid CID) error
		RemoveAttribute(ctx context.Context, s *Strain, id UUID, cid CID) error
	}

	Strainer interface {
		SelectAllStrains(ctx context.Context, cid CID) ([]Strain, error)
		SelectStrain(ctx context.Context, id UUID, cid CID) (Strain, error)
		InsertStrain(ctx context.Context, s Strain, cid CID) (Strain, error)
		UpdateStrain(ctx context.Context, id UUID, s Strain, cid CID) error
		DeleteStrain(ctx context.Context, id UUID, cid CID) error
	}

	SubstrateIngredienter interface {
		GetAllIngredients(ctx context.Context, s *Substrate, cid CID) error
		AddIngredient(ctx context.Context, s *Substrate, i Ingredient, cid CID) error
		ChangeIngredient(ctx context.Context, s *Substrate, oldI, newI Ingredient, cid CID) error
		RemoveIngredient(ctx context.Context, s *Substrate, i Ingredient, cid CID) error
	}

	Substrater interface {
		SelectAllSubstrates(ctx context.Context, cid CID) ([]Substrate, error)
		SelectSubstrate(ctx context.Context, id UUID, cid CID) (Substrate, error)
		InsertSubstrate(ctx context.Context, s Substrate, cid CID) (Substrate, error)
		UpdateSubstrate(ctx context.Context, id UUID, s Substrate, cid CID) error
		DeleteSubstrate(ctx context.Context, id UUID, cid CID) error
	}

	// vendors aren't meant to be a comprehensive list of attributes, really just
	// a name that makes a unique constraint with Strain.Name in the strain table;
	// storing them is just for convenience, and to avoid spelling errors
	Vendorer interface {
		// this doesn't support pagination, because our list of vendors will only
		// range into the hundreds *at most*
		SelectAllVendors(ctx context.Context, cid CID) ([]Vendor, error)
		// mostly useful for filling in complex objects like strains or substrates
		SelectVendor(ctx context.Context, id UUID, cid CID) (Vendor, error)
		// assigns a generated uuid, ctime and mtime to `v` and inserts it into vendors;
		// the modified `v` struct is returned
		InsertVendor(ctx context.Context, v Vendor, cid CID) (Vendor, error)
		// assigns a generated mtime to the `v` arg, then updates every vendor attribute
		// except uuid based on the `id` arg; overwrites the uuid in `v` and returns `v`
		UpdateVendor(ctx context.Context, id UUID, v Vendor, cid CID) error
		// this should really never be used; none of the foreign keys have cascade
		// delete, so any referenced vendors will cause an error here; just don't use it
		DeleteVendor(ctx context.Context, id UUID, cid CID) error
	}
)
