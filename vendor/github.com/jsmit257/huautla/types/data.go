package types

import (
	"time"
)

type (
	UUID string

	CID string

	SubstrateType string

	Config struct {
		PGHost string
		PGUser string
		PGPass string
		PGPort uint
		PGSSL  string
	}

	Note struct {
		UUID  `json:"id,omitempty"`
		Note  string    `json:"note,omitempty"`
		MTime time.Time `json:"modified_at,omitempty"`
		CTime time.Time `json:"create_at,omitempty"`
	}

	Photo struct {
		UUID     `json:"id"`
		Filename string    `json:"image"`
		Notes    []Note    `json:"notes"`
		CTime    time.Time `json:"create_at"`
	}

	Vendor struct {
		UUID    `json:"id"`
		Name    string `json:"name"`
		Website string `json:"website,omitempty"`
	}

	Substrate struct {
		UUID        `json:"id"`
		Name        string        `json:"name"`
		Type        SubstrateType `json:"type"`
		Vendor      `json:"vendor"`
		Ingredients []Ingredient `json:"ingredients,omitempty"`
	}

	Ingredient struct {
		UUID `json:"id"`
		Name string `json:"name"`
	}

	Strain struct {
		UUID       `json:"id"`
		Species    string `json:"species,omitempty"`
		Name       string `json:"name"`
		Vendor     `json:"vendor"`
		Generation *Generation       `json:"generation,omitempty"`
		Attributes []StrainAttribute `json:"attributes,omitempty"`
		CTime      time.Time         `json:"create_date"`
	}

	StrainAttribute struct {
		UUID  `json:"id"`
		Name  string `json:"name"`
		Value string `json:"value"`
	}

	Stage struct {
		UUID `json:"id"`
		Name string `json:"name"`
	}

	EventType struct {
		UUID     `json:"id"`
		Name     string `json:"name"`
		Severity string `json:"severity"`
		Stage    `json:"stage"`
	}

	Lifecycle struct {
		UUID           `json:"id"`
		Location       string  `json:"location"`
		StrainCost     float32 `json:"strain_cost,omitempty"`
		GrainCost      float32 `json:"grain_cost,omitempty"`
		BulkCost       float32 `json:"bulk_cost,omitempty"`
		Yield          float32 `json:"yield,omitempty"`
		Count          int16   `json:"count,omitempty"`
		Gross          float32 `json:"gross,omitempty"`
		Strain         `json:"strain,omitempty"`
		GrainSubstrate Substrate `json:"grain_substrate,omitempty"`
		BulkSubstrate  Substrate `json:"bulk_substrate,omitempty"`
		Events         []Event   `json:"events,omitempty"`
		Notes          []Note    `json:"notes,omitempty"`
		MTime          time.Time `json:"modified_date,omitempty"`
		CTime          time.Time `json:"create_date"`
	}

	Event struct {
		UUID        `json:"id"`
		Temperature float32   `json:"temperature"`
		Humidity    int8      `json:"humidity,omitempty"`
		EventType   EventType `json:"event_type"`
		Photos      []Photo   `json:"photos,omitempty"`
		Notes       []Note    `json:"notes,omitempty"`
		MTime       time.Time `json:"modified_date"`
		CTime       time.Time `json:"create_date"`
	}

	Source struct {
		UUID      `json:"id"`
		Type      string     `json:"type"`
		Lifecycle *Lifecycle `json:"lifecycle,omitempty"`
		Strain    `json:"strain"`
	}

	Generation struct {
		UUID             `json:"id"`
		PlatingSubstrate Substrate `json:"plating_substrate"`
		LiquidSubstrate  Substrate `json:"liquid_substrate"`
		Sources          []Source  `json:"sources,omitempty"`
		Events           []Event   `json:"events,omitempty"`
		Notes            []Note    `json:"notes,omitempty"`
		MTime            time.Time `json:"modified_date"`
		CTime            time.Time `json:"create_date"`
	}
)
