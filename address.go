// Copyright (c) 2020 Bojan Zivanovic and contributors
// SPDX-License-Identifier: MIT

// Package address handles address representation, validation and formatting.
package address

import (
	"bytes"
	"fmt"
	"regexp"
	"sort"
)

// Address represents an address.
type Address struct {
	Line1 string `json:"line1"`
	Line2 string `json:"line2"`
	Line3 string `json:"line3"`
	// Sublocality is the neighborhood/suburb/district.
	Sublocality string `json:"sublocality"`
	// Locality is the city/village/post town.
	Locality string `json:"locality"`
	// Region is the state/province/prefecture.
	// An ISO code is used when available.
	Region string `json:"region"`
	// PostalCode is the postal/zip/pin code.
	PostalCode string `json:"postal_code"`
	// CountryCode is the two-letter code as defined by CLDR.
	CountryCode string `json:"country"`
}

// IsEmpty returns whether a is empty.
func (a Address) IsEmpty() bool {
	// An address must at least have a country code.
	return a.CountryCode == ""
}

// Format represents an address format.
type Format struct {
	Locale            Locale          `json:"locale,omitempty"`
	Layout            string          `json:"layout,omitempty"`
	LocalLayout       string          `json:"local_layout,omitempty"`
	Required          []Field         `json:"required,omitempty"`
	SublocalityType   SublocalityType `json:"sublocality_type,omitempty"`
	LocalityType      LocalityType    `json:"locality_type,omitempty"`
	RegionType        RegionType      `json:"region_type,omitempty"`
	PostalCodeType    PostalCodeType  `json:"postal_code_type,omitempty"`
	PostalCodePattern string          `json:"postal_code_pattern,omitempty"`
	ShowRegionID      bool            `json:"show_region_id,omitempty"`
	Regions           RegionMap       `json:"regions,omitempty"`
	LocalRegions      RegionMap       `json:"local_regions,omitempty"`
}

// IsRequired returns whether the given field is required.
func (f Format) IsRequired(field Field) bool {
	for _, ff := range f.Required {
		if ff == field {
			return true
		}
	}
	return false
}

// CheckRequired checks whether a required field is valid (non-blank).
//
// Non-required fields are considered valid even if they're blank.
func (f Format) CheckRequired(field Field, value string) bool {
	required := f.IsRequired(field)
	return !required || (required && value != "")
}

// CheckRegion checks whether the given region is valid.
//
// An empty region is considered valid.
func (f Format) CheckRegion(region string) bool {
	if region == "" || f.Regions.Len() == 0 {
		return true
	}
	return f.Regions.HasKey(region)
}

// CheckPostalCode checks whether the given postal code is valid.
//
// An empty postal code is considered valid.
func (f Format) CheckPostalCode(postalCode string) bool {
	if postalCode == "" || f.PostalCodePattern == "" {
		return true
	}
	rx := regexp.MustCompile(f.PostalCodePattern)
	return rx.MatchString(postalCode)
}

// SelectLayout selects the correct layout for the given locale.
func (f Format) SelectLayout(locale Locale) string {
	if f.LocalLayout != "" && f.useLocalData(locale) {
		return f.LocalLayout
	}
	return f.Layout
}

// SelectRegions selects the correct regions for the given locale.
func (f Format) SelectRegions(locale Locale) RegionMap {
	if f.LocalRegions.Len() > 0 && f.useLocalData(locale) {
		return f.LocalRegions
	}
	return f.Regions
}

// useLocalData returns whether local data should be used for the given locale.
func (f Format) useLocalData(locale Locale) bool {
	if locale.Script == "Latn" {
		// Allow locales to opt out of local data. E.g: zh-Latn.
		return false
	}
	// Scripts are not compared, matching libaddressinput behavior. This means
	// that zh-Hant data will be shown to zh-Hans users, and vice-versa.
	return locale.Language == f.Locale.Language
}

// RegionMap represents a read-only ordered map of regions.
type RegionMap struct {
	keys   []string
	values map[string]string
}

// NewRegionMap creates a new region map from the given pairs.
func NewRegionMap(pairs ...string) RegionMap {
	if len(pairs) == 0 {
		return RegionMap{}
	}
	if len(pairs)%2 != 0 {
		panic(fmt.Errorf("uneven number of pairs given to NewRegionMap()"))
	}
	r := RegionMap{}
	r.keys = make([]string, 0, len(pairs)/2)
	r.values = make(map[string]string, len(pairs)/2)
	for i := 0; i < len(pairs)-1; i += 2 {
		r.keys = append(r.keys, pairs[i])
		r.values[pairs[i]] = pairs[i+1]
	}

	return r
}

// Get returns the value for the given key, or an empty string if none found.
func (r RegionMap) Get(key string) (string, bool) {
	v, ok := r.values[key]
	return v, ok
}

// HasKey returns whether the given key exists in the map.
func (r RegionMap) HasKey(key string) bool {
	_, ok := r.values[key]
	return ok
}

// Keys returns a list of keys.
func (r RegionMap) Keys() []string {
	return r.keys
}

// Len returns the number of keys in the map.
func (r RegionMap) Len() int {
	return len(r.keys)
}

func (r RegionMap) MarshalJSON() ([]byte, error) {
	if r.Len() == 0 {
		return []byte("{}"), nil
	}
	buf := &bytes.Buffer{}
	buf.Grow(r.Len() * 30)
	buf.WriteByte('{')
	for i, key := range r.keys {
		if i > 0 {
			buf.WriteByte(',')
		}
		// A fully generic MarshalJSON would call encoder.Encode() to ensure both key and value
		// are escaped and contain valid utf8. Since all regions are defined in the package, they
		// are assumed to be valid, avoiding thousands of allocs when marshalling the format list.
		buf.WriteString(`"` + key + `":"`)
		buf.WriteString(r.values[key])
		buf.WriteString(`"`)
	}
	buf.WriteByte('}')

	return buf.Bytes(), nil
}

// CheckCountryCode checks whether the given country code is valid.
//
// An empty country code is considered valid.
func CheckCountryCode(countryCode string) bool {
	if countryCode == "" {
		return true
	}
	_, ok := countries[countryCode]
	return ok
}

// GetCountryCodes returns all known country codes.
func GetCountryCodes() []string {
	countryCodes := make([]string, 0, len(countries))
	for countryCode := range countries {
		countryCodes = append(countryCodes, countryCode)
	}
	sort.Strings(countryCodes)
	return countryCodes
}

// GetCountryNames returns all known country names, keyed by country code.
func GetCountryNames() map[string]string {
	return countries
}

// GetFormats returns all known address formats, keyed by country code.
//
// The ZZ address format represents the generic fallback.
func GetFormats() map[string]Format {
	return formats
}

// GetFormat returns an address format for the given country code.
func GetFormat(countryCode string) Format {
	format, ok := formats[countryCode]
	if !ok {
		return formats["ZZ"]
	}
	return format
}
