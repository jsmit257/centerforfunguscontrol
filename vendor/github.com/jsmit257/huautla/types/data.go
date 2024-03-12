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
		Species    string    `json:"species,omitempty"`
		Name       string    `json:"name"`
		CTime      time.Time `json:"create_date"`
		Vendor     `json:"vendor"`
		Attributes []StrainAttribute `json:"attributes,omitempty"`
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
		Location       string    `json:"location"`
		StrainCost     float32   `json:"strain_cost,omitempty"`
		GrainCost      float32   `json:"grain_cost,omitempty"`
		BulkCost       float32   `json:"bulk_cost,omitempty"`
		Yield          float32   `json:"yield,omitempty"`
		Count          int16     `json:"count,omitempty"`
		Gross          float32   `json:"gross,omitempty"`
		MTime          time.Time `json:"modified_date,omitempty"`
		CTime          time.Time `json:"create_date"`
		Strain         `json:"strain,omitempty"`
		GrainSubstrate Substrate `json:"grain_substrate,omitempty"`
		BulkSubstrate  Substrate `json:"bulk_substrate,omitempty"`
		Events         []Event   `json:"events,omitempty"`
	}

	Event struct {
		UUID        `json:"id"`
		Temperature float32   `json:"temperature"`
		Humidity    int8      `json:"humidity,omitempty"`
		MTime       time.Time `json:"modified_date"`
		CTime       time.Time `json:"create_date"`
		EventType   EventType `json:"event_type"`
	}
)
