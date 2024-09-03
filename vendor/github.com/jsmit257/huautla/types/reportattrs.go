package types

import (
	"fmt"
	"net/url"
)

var valid = map[string]struct{}{
	"generation-id": {},
	"lifecycle-id":  {},
	"strain-id":     {},
	"plating-id":    {},
	"liquid-id":     {},
	"grain-id":      {},
	"bulk-id":       {},
}

func NewReportAttrs(m url.Values) (ReportAttrs, error) {
	result := ReportAttrs{}
	return result, result.Map(m)
}

func (ra ReportAttrs) Set(name, value string) error {
	if value == "" { // treat it like null
		return fmt.Errorf("empty value for key: %s", name)
	} else if _, ok := valid[name]; !ok {
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

/* returns true if any name in names is a key in this list */
func (ra ReportAttrs) Contains(names ...string) bool {
	for _, k := range names {
		if _, ok := ra[k]; ok {
			return true
		}
	}
	return false
}

func (ra ReportAttrs) Map(m url.Values) (err error) {
	errs := []string{}
	for k, v := range m {
		if err := ra.Set(k, v[0]); err != nil {
			errs = append(errs, k)
		}
	}
	if len(errs) > 0 {
		err = fmt.Errorf("unrecognizable param values in the following fields: %v", errs)
	}
	return err
}
