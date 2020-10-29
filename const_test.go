// Copyright (c) 2020 Bojan Zivanovic and contributors
// SPDX-License-Identifier: MIT

package address_test

import (
	"testing"

	"github.com/bojanz/address"
)

func TestSublocalityType_String(t *testing.T) {
	got := address.SublocalityTypeTownland.String()
	want := "townland"
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	got = address.SublocalityType(20).String()
	if got != "" {
		t.Errorf("got %v, want an empty string", got)
	}
}

func TestSublocalityType_MarshalText(t *testing.T) {
	b, _ := address.SublocalityTypeTownland.MarshalText()
	got := string(b)
	want := "townland"
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}

	b, _ = address.SublocalityType(20).MarshalText()
	got = string(b)
	if got != "" {
		t.Errorf("got %v, want an empty string", got)
	}
}

func TestSublocalityType_UnmarshalText(t *testing.T) {
	var sublocalityType address.SublocalityType
	sublocalityType.UnmarshalText([]byte("townland"))
	if sublocalityType != address.SublocalityTypeTownland {
		t.Errorf("got %v, want %v", sublocalityType, address.SublocalityTypeTownland)
	}

	err := sublocalityType.UnmarshalText([]byte("abcd"))
	want := `invalid sublocality type "abcd"`
	if err == nil || err.Error() != want {
		t.Errorf("got %v, want %v", err, want)
	}
}

func TestLocalityType_String(t *testing.T) {
	got := address.LocalityTypePostTown.String()
	want := "post_town"
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	got = address.LocalityType(20).String()
	if got != "" {
		t.Errorf("got %v, want an empty string", got)
	}
}

func TestLocalityType_MarshalText(t *testing.T) {
	b, _ := address.LocalityTypePostTown.MarshalText()
	got := string(b)
	want := "post_town"
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}

	b, _ = address.LocalityType(20).MarshalText()
	got = string(b)
	if got != "" {
		t.Errorf("got %v, want an empty string", got)
	}
}

func TestLocalityType_UnmarshalText(t *testing.T) {
	var localityType address.LocalityType
	localityType.UnmarshalText([]byte("post_town"))
	if localityType != address.LocalityTypePostTown {
		t.Errorf("got %v, want %v", localityType, address.LocalityTypePostTown)
	}

	err := localityType.UnmarshalText([]byte("abcd"))
	want := `invalid locality type "abcd"`
	if err == nil || err.Error() != want {
		t.Errorf("got %v, want %v", err, want)
	}
}

func TestRegionType_String(t *testing.T) {
	got := address.RegionTypeProvince.String()
	want := "province"
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	got = address.RegionType(20).String()
	if got != "" {
		t.Errorf("got %v, want an empty string", got)
	}
}

func TestRegionType_MarshalText(t *testing.T) {
	b, _ := address.RegionTypeProvince.MarshalText()
	got := string(b)
	want := "province"
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}

	b, _ = address.RegionType(20).MarshalText()
	got = string(b)
	if got != "" {
		t.Errorf("got %v, want an empty string", got)
	}
}

func TestRegionType_UnmarshalText(t *testing.T) {
	var regionType address.RegionType
	regionType.UnmarshalText([]byte("province"))
	if regionType != address.RegionTypeProvince {
		t.Errorf("got %v, want %v", regionType, address.RegionTypeProvince)
	}

	err := regionType.UnmarshalText([]byte("abcd"))
	want := `invalid region type "abcd"`
	if err == nil || err.Error() != want {
		t.Errorf("got %v, want %v", err, want)
	}
}

func TestPostalCodeType_String(t *testing.T) {
	got := address.PostalCodeTypeZip.String()
	want := "zip"
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	got = address.PostalCodeType(20).String()
	if got != "" {
		t.Errorf("got %v, want an empty string", got)
	}
}

func TestPostalCodeType_MarshalText(t *testing.T) {
	b, _ := address.PostalCodeTypeZip.MarshalText()
	got := string(b)
	want := "zip"
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}

	b, _ = address.PostalCodeType(20).MarshalText()
	got = string(b)
	if got != "" {
		t.Errorf("got %v, want an empty string", got)
	}
}

func TestPostalCodeType_UnmarshalText(t *testing.T) {
	var postalCodeType address.PostalCodeType
	postalCodeType.UnmarshalText([]byte("zip"))
	if postalCodeType != address.PostalCodeTypeZip {
		t.Errorf("got %v, want %v", postalCodeType, address.PostalCodeTypeZip)
	}

	err := postalCodeType.UnmarshalText([]byte("abcd"))
	want := `invalid postal code type "abcd"`
	if err == nil || err.Error() != want {
		t.Errorf("got %v, want %v", err, want)
	}
}
