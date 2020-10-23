// Copyright (c) 2020 Bojan Zivanovic and contributors
// SPDX-License-Identifier: MIT

package address

import (
	"regexp"
	"sort"
)

// Address represents an address.
type Address struct {
	Line1 string
	Line2 string
	Line3 string
	// Sublocality is the neighborhood/suburb/district.
	Sublocality string
	// Locality is the city/village/post town.
	Locality string
	// Region is the state/province/prefecture.
	// An ISO code is used when available.
	Region string
	// PostalCode is the postal/zip/pin code.
	PostalCode string
	// CountryCode is the two-letter code as defined by CLDR.
	CountryCode string
}

// IsEmpty returns whether a is empty.
func (a Address) IsEmpty() bool {
	// An address must at least have a country code.
	return a.CountryCode == ""
}

// Format represents an address format.
type Format struct {
	Locale            Locale
	Layout            string
	LocalLayout       string
	Required          []Field
	SublocalityType   SublocalityType
	LocalityType      LocalityType
	RegionType        RegionType
	PostalCodeType    PostalCodeType
	PostalCodePattern string
	ShowRegionID      bool
	Regions           map[string]string
	LocalRegions      map[string]string
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
	if region == "" || len(f.Regions) == 0 {
		return true
	}
	_, ok := f.Regions[region]
	return ok
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
func (f Format) SelectRegions(locale Locale) map[string]string {
	if len(f.LocalRegions) > 0 && f.useLocalData(locale) {
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
