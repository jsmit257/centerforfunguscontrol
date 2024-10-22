package types

import (
	"fmt"
	"net/url"
)

var validReportAttrs = map[string]struct{}{
	"generation-id": {},
	"lifecycle-id":  {},
	"strain-id":     {},
	"substrate-id":  {},
	"plating-id":    {},
	"liquid-id":     {},
	"grain-id":      {},
	"bulk-id":       {},
	"eventtype-id":  {},
	"vendor-id":     {},
}

type reportAttrs map[string]UUID

func NewReportAttrs(m url.Values) (ReportAttrs, error) {
	result := reportAttrs{}
	return result, result.apply(m)
}

func (ra reportAttrs) Set(name, value string) error {
	if value == "" { // treat it like null
		return fmt.Errorf("empty value for key: %s", name)
	} else if _, ok := validReportAttrs[name]; !ok {
		return fmt.Errorf("unknown parameter: %s", name)
	}
	ra[name] = UUID(value)
	return nil
}

func (ra reportAttrs) Get(name string) *UUID {
	temp, ok := ra[name]
	if !ok {
		return nil
	}
	result := UUID(temp)
	return &result
}

/* returns true if any name in names is a key in this list */
func (ra reportAttrs) Contains(names ...string) bool {
	for _, k := range names {
		if _, ok := ra[k]; ok {
			return true
		}
	}
	return false
}

func (ra reportAttrs) apply(m url.Values) (err error) {
	errs := []string{}
	for k, v := range m {
		if err := ra.Set(k, v[0]); err != nil {
			errs = append(errs, k)
		}
	}
	if len(errs) > 0 {
		err = fmt.Errorf("failed to find param values in the following fields: %v", errs)
	}
	return err
}
