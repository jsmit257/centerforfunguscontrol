package types

import "context"

type (
	DB interface {
		LifecycleEventer
		EventTyper
		Generationer
		GenerationEventer
		Ingredienter
		Lifecycler
		Noter
		Observer
		Photoer
		Sourcer
		Stager
		StrainAttributer
		Strainer
		SubstrateIngredienter
		Substrater
		Vendorer
	}

	EventTyper interface {
		SelectAllEventTypes(ctx context.Context, cid CID) ([]EventType, error)
		SelectEventType(ctx context.Context, id UUID, cid CID) (EventType, error)
		InsertEventType(ctx context.Context, e EventType, cid CID) (EventType, error)
		UpdateEventType(ctx context.Context, id UUID, e EventType, cid CID) error
		DeleteEventType(ctx context.Context, id UUID, cid CID) error
	}

	Generationer interface {
		SelectGenerationIndex(context.Context, CID) ([]Generation, error)
		SelectGenerationsByAttrs(context.Context, ReportAttrs, CID) ([]Generation, error)
		SelectGeneration(context.Context, UUID, CID) (Generation, error)
		InsertGeneration(context.Context, Generation, CID) (Generation, error)
		UpdateGeneration(context.Context, Generation, CID) (Generation, error)
		DeleteGeneration(context.Context, UUID, CID) error
	}

	GenerationEventer interface {
		GetGenerationEvents(ctx context.Context, g *Generation, cid CID) error
		AddGenerationEvent(ctx context.Context, g *Generation, e Event, cid CID) error
		ChangeGenerationEvent(ctx context.Context, g *Generation, e Event, cid CID) (Event, error)
		RemoveGenerationEvent(ctx context.Context, g *Generation, id UUID, cid CID) error
	}

	Ingredienter interface {
		SelectAllIngredients(ctx context.Context, cid CID) ([]Ingredient, error)
		SelectIngredient(ctx context.Context, id UUID, cid CID) (Ingredient, error)
		InsertIngredient(ctx context.Context, i Ingredient, cid CID) (Ingredient, error)
		UpdateIngredient(ctx context.Context, id UUID, i Ingredient, cid CID) error
		DeleteIngredient(ctx context.Context, id UUID, cid CID) error
	}

	LifecycleEventer interface {
		GetLifecycleEvents(ctx context.Context, lc *Lifecycle, cid CID) error
		AddLifecycleEvent(ctx context.Context, lc *Lifecycle, e Event, cid CID) error
		ChangeLifecycleEvent(ctx context.Context, lc *Lifecycle, e Event, cid CID) (Event, error)
		RemoveLifecycleEvent(ctx context.Context, lc *Lifecycle, id UUID, cid CID) error
	}

	Lifecycler interface {
		SelectLifecycleIndex(ctx context.Context, cid CID) ([]Lifecycle, error)
		SelectLifecyclesByAttrs(ctx context.Context, p ReportAttrs, cid CID) ([]Lifecycle, error)
		SelectLifecycle(ctx context.Context, id UUID, cid CID) (Lifecycle, error)
		InsertLifecycle(ctx context.Context, lc Lifecycle, cid CID) (Lifecycle, error)
		UpdateLifecycle(ctx context.Context, lc Lifecycle, cid CID) (Lifecycle, error)
		DeleteLifecycle(ctx context.Context, id UUID, cid CID) error
	}

	Noter interface {
		GetNotes(ctx context.Context, id UUID, cid CID) ([]Note, error)
		AddNote(ctx context.Context, id UUID, notes []Note, n Note, cid CID) ([]Note, error)
		ChangeNote(ctx context.Context, notes []Note, n Note, cid CID) ([]Note, error)
		RemoveNote(ctx context.Context, notes []Note, id UUID, cid CID) ([]Note, error)
	}

	Observer interface {
		SelectByEventType(ctx context.Context, et EventType, cid CID) ([]Event, error)
		SelectEvent(ctx context.Context, id UUID, cid CID) (Event, error)
	}

	Photoer interface {
		GetPhotos(ctx context.Context, id UUID, cid CID) ([]Photo, error)
		AddPhoto(ctx context.Context, id UUID, photos []Photo, p Photo, cid CID) ([]Photo, error)
		ChangePhoto(ctx context.Context, photos []Photo, p Photo, cid CID) ([]Photo, error)
		RemovePhoto(ctx context.Context, photos []Photo, id UUID, cid CID) ([]Photo, error)
	}

	Sourcer interface {
		// GetSources(ctx context.Context, g *Generation, cid CID) error
		AddStrainSource(ctx context.Context, g *Generation, s Source, cid CID) error
		AddEventSource(ctx context.Context, g *Generation, e Event, cid CID) error
		ChangeSource(ctx context.Context, g *Generation, s Source, cid CID) error
		RemoveSource(ctx context.Context, g *Generation, id UUID, cid CID) error
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
		GeneratedStrain(ctx context.Context, id UUID, cid CID) (Strain, error)
		UpdateGeneratedStrain(ctx context.Context, gid *UUID, sid UUID, cid CID) error
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
