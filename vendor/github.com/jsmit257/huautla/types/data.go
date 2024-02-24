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
		UUID `json:"-"`
		Name string `json:"name"`
	}

	Substrate struct {
		UUID        `json:"-"`
		Name        string        `json:"name"`
		Type        SubstrateType `json:"type"`
		Vendor      `json:"vendor"`
		Ingredients []Ingredient `json:"ingredients,omitempty"`
	}

	Ingredient struct {
		UUID `json:"-"`
		Name string `json:"name"`
	}

	Strain struct {
		UUID       `json:"-"`
		Name       string `json:"name"`
		Vendor     `json:"vendor"`
		Attributes []StrainAttribute `json:"attributes,omitempty"`
	}

	StrainAttribute struct {
		UUID  `json:"id"`
		Name  string `json:"name"`
		Value string `json:"value"`
	}

	Stage struct {
		UUID `json:"-"`
		Name string `json:"name"`
	}

	EventType struct {
		UUID     `json:"-"`
		Name     string `json:"name"`
		Severity string `json:"severity"`
		Stage    `json:"stage"`
	}

	Lifecycle struct {
		UUID           `json:"id"`
		Name           string    `json:"name"`
		Location       string    `json:"location"`
		GrainCost      float32   `json:"grain_cost"`
		BulkCost       float32   `json:"bulk_cost"`
		Yield          float32   `json:"yield"`
		Count          int16     `json:"count"`
		Gross          float32   `json:"gross"`
		MTime          time.Time `json:"modified_date"`
		CTime          time.Time `json:"create_date"`
		Strain         `json:"strain"`
		GrainSubstrate Substrate `json:"grain_substrate"`
		BulkSubstrate  Substrate `json:"bulk_substrate"`
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
