// Copyright (c) 2020 Bojan Zivanovic and contributors
// SPDX-License-Identifier: MIT

package address

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

// LocalityType represents the locality type.
type LocalityType uint8

const (
	LocalityTypeCity LocalityType = iota
	LocalityTypeDistrict
	LocalityTypePostTown
	LocalityTypeSuburb
)

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

// PostalCodeType represents the postal code type.
type PostalCodeType uint8

const (
	PostalCodeTypePostal PostalCodeType = iota
	PostalCodeTypeEir
	PostalCodeTypePin
	PostalCodeTypeZip
)
