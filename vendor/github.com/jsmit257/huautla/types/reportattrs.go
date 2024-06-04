package types

import "fmt"

var valid = map[string]struct{}{
	"generation-id": {},
	"lifecycle-id":  {},
	"strain-id":     {},
	"plating-id":    {},
	"liquid-id":     {},
	"grain-id":      {},
	"bulk-id":       {},
}

func (ra ReportAttrs) Set(name, value string) error {
	if _, ok := valid[name]; !ok {
		return fmt.Errorf("unknown parameter: %s", name)
	}

	ra[name] = UUID(value)

	return nil
}

func (ra ReportAttrs) Get(name string) *UUID {
	temp, ok := ra[name]
	if !ok {
		return nil
	}
	result := UUID(temp)
	return &result
}
