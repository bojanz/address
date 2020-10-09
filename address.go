// Copyright (c) 2020 Bojan Zivanovic and contributors
// SPDX-License-Identifier: MIT

package address

import "regexp"

// Address represents an address.
type Address struct {
	Line1 string
	Line2 string
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
	Layout            string
	Required          []Field
	SublocalityType   SublocalityType
	LocalityType      LocalityType
	RegionType        RegionType
	PostalCodeType    PostalCodeType
	PostalCodePattern string
	ShowRegionID      bool
	Regions           map[string]string
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
