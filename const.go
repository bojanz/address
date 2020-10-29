// Copyright (c) 2020 Bojan Zivanovic and contributors
// SPDX-License-Identifier: MIT

package address

import "fmt"

// Field represents an address field.
type Field string

const (
	FieldLine1       Field = "1"
	FieldLine2       Field = "2"
	FieldLine3       Field = "3"
	FieldSublocality Field = "S"
	FieldLocality    Field = "L"
	FieldRegion      Field = "R"
	FieldPostalCode  Field = "P"
)

// SublocalityType represents the sublocality type.
type SublocalityType uint8

const (
	SublocalityTypeSuburb SublocalityType = iota
	SublocalityTypeDistrict
	SublocalityTypeNeighborhood
	SublocalityTypeVillageTownship
	SublocalityTypeTownland
)

var sublocalityTypeNames = [...]string{"suburb", "district", "neighborhood", "village_township", "townland"}

// String returns the string representation of s.
func (s SublocalityType) String() string {
	if int(s) >= len(sublocalityTypeNames) {
		return ""
	}
	return sublocalityTypeNames[s]
}

// MarshalText implements the encoding.TextMarshaler interface.
func (s SublocalityType) MarshalText() ([]byte, error) {
	return []byte(s.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (s *SublocalityType) UnmarshalText(b []byte) error {
	aux := string(b)
	for i, name := range sublocalityTypeNames {
		if name == aux {
			*s = SublocalityType(i)
			return nil
		}
	}
	return fmt.Errorf("invalid sublocality type %q", aux)
}

// LocalityType represents the locality type.
type LocalityType uint8

const (
	LocalityTypeCity LocalityType = iota
	LocalityTypeDistrict
	LocalityTypePostTown
	LocalityTypeSuburb
)

var localityTypeNames = [...]string{"city", "district", "post_town", "suburb"}

// String returns the string representation of l.
func (l LocalityType) String() string {
	if int(l) >= len(localityTypeNames) {
		return ""
	}
	return localityTypeNames[l]
}

// MarshalText implements the encoding.TextMarshaler interface.
func (l LocalityType) MarshalText() ([]byte, error) {
	return []byte(l.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (l *LocalityType) UnmarshalText(b []byte) error {
	aux := string(b)
	for i, name := range localityTypeNames {
		if name == aux {
			*l = LocalityType(i)
			return nil
		}
	}
	return fmt.Errorf("invalid locality type %q", aux)
}

// RegionType represents the region type.
type RegionType uint8

const (
	RegionTypeProvince RegionType = iota
	RegionTypeArea
	RegionTypeCanton
	RegionTypeCounty
	RegionTypeDepartment
	RegionTypeDistrict
	RegionTypeDoSi
	RegionTypeEmirate
	RegionTypeIsland
	RegionTypeOblast
	RegionTypeParish
	RegionTypePrefecture
	RegionTypeState
)

var regionTypeNames = [...]string{
	"province", "area", "canton", "department", "distict", "do_si",
	"emirate", "island", "oblast", "parish", "prefecture", "state",
}

// String returns the string representation of r.
func (r RegionType) String() string {
	if int(r) >= len(regionTypeNames) {
		return ""
	}
	return regionTypeNames[r]
}

// MarshalText implements the encoding.TextMarshaler interface.
func (r RegionType) MarshalText() ([]byte, error) {
	return []byte(r.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (r *RegionType) UnmarshalText(b []byte) error {
	aux := string(b)
	for i, name := range regionTypeNames {
		if name == aux {
			*r = RegionType(i)
			return nil
		}
	}
	return fmt.Errorf("invalid region type %q", aux)
}

// PostalCodeType represents the postal code type.
type PostalCodeType uint8

const (
	PostalCodeTypePostal PostalCodeType = iota
	PostalCodeTypeEir
	PostalCodeTypePin
	PostalCodeTypeZip
)

var postalCodeTypeNames = [...]string{"postal", "eir", "pin", "zip"}

// String returns the string representation of p.
func (p PostalCodeType) String() string {
	if int(p) >= len(postalCodeTypeNames) {
		return ""
	}
	return postalCodeTypeNames[p]
}

// MarshalText implements the encoding.TextMarshaler interface.
func (p PostalCodeType) MarshalText() ([]byte, error) {
	return []byte(p.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (p *PostalCodeType) UnmarshalText(b []byte) error {
	aux := string(b)
	for i, name := range postalCodeTypeNames {
		if name == aux {
			*p = PostalCodeType(i)
			return nil
		}
	}
	return fmt.Errorf("invalid postal code type %q", aux)
}
