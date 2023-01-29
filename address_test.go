// Copyright (c) 2020 Bojan Zivanovic and contributors
// SPDX-License-Identifier: MIT

package address_test

import (
	"encoding/json"
	"reflect"
	"testing"
	"unicode/utf8"

	"github.com/bojanz/address"
)

func TestAddress_IsEmpty(t *testing.T) {
	a := address.Address{}
	if !a.IsEmpty() {
		t.Errorf("expected address to be empty.")
	}

	a = address.Address{Locality: "Belgrade"}
	if !a.IsEmpty() {
		t.Errorf("expected address to be empty.")
	}

	a = address.Address{CountryCode: "RS"}
	if a.IsEmpty() {
		t.Errorf("expected address to not be empty.")
	}
}

func TestFormat_IsRequired(t *testing.T) {
	format := address.GetFormat("RS")
	got := format.IsRequired(address.FieldLine1)
	if !got {
		t.Errorf("expected FieldLine1 to be required.")
	}

	got = format.IsRequired(address.FieldLine2)
	if got {
		t.Errorf("expected FieldLine2 to not be required.")
	}
}

func TestFormat_CheckRequired(t *testing.T) {
	tests := []struct {
		field address.Field
		value string
		want  bool
	}{
		// Required and empty.
		{address.FieldLine1, "", false},
		// Optional and empty.
		{address.FieldLine2, "", true},
		// Required and non-empty.
		{address.FieldLocality, "Belgrade", true},
		// Optional and non-empty.
		{address.FieldPostalCode, "11000", true},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			format := address.GetFormat("RS")
			got := format.CheckRequired(tt.field, tt.value)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormat_CheckRegion(t *testing.T) {
	tests := []struct {
		countryCode string
		region      string
		want        bool
	}{
		// Empty value.
		{"US", "", true},
		// Valid ISO code.
		{"US", "CA", true},
		// Invalid ISO code.
		{"US", "California", false},
		// Country with no predefined regions.
		{"RS", "Vojvodina", true},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			format := address.GetFormat(tt.countryCode)
			got := format.CheckRegion(tt.region)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormat_CheckPostalCode(t *testing.T) {
	tests := []struct {
		countryCode string
		postalCode  string
		want        bool
	}{
		// Empty value.
		{"FR", "", true},
		// Valid postal code.
		{"FR", "75002", true},
		// Invalid postal code.
		{"FR", "INVALID", false},
		// Country with no predefined pattern.
		{"AG", "AG123", true},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			format := address.GetFormat(tt.countryCode)
			got := format.CheckPostalCode(tt.postalCode)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormat_SelectLayout(t *testing.T) {
	tests := []struct {
		countryCode string
		locale      string
		wantLocal   bool
	}{
		// China ("zh").
		{"CN", "en", false},
		{"CN", "ja", false},
		{"CN", "zh-Latn", false},
		{"CN", "zh", true},
		{"CN", "zh-Hant", true},
		// Hong Kong ("zh-Hant").
		{"HK", "en", false},
		{"HK", "ja", false},
		{"HK", "zh-Latn", false},
		{"HK", "zh", true},
		{"HK", "zh-Hant", true},
		// Serbia (no local layout defined).
		{"RS", "en", false},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			locale := address.NewLocale(tt.locale)
			format := address.GetFormat(tt.countryCode)
			got := format.SelectLayout(locale)
			want := format.Layout
			if tt.wantLocal {
				want = format.LocalLayout
			}

			if got != want {
				t.Errorf("got %v, want %v", got, want)
			}
		})
	}
}

func TestFormat_SelectRegions(t *testing.T) {
	tests := []struct {
		countryCode string
		locale      string
		wantLocal   bool
	}{
		// China ("zh").
		{"CN", "en", false},
		{"CN", "ja", false},
		{"CN", "zh-Latn", false},
		{"CN", "zh", true},
		{"CN", "zh-Hant", true},
		// Hong Kong ("zh-Hant").
		{"HK", "en", false},
		{"HK", "ja", false},
		{"HK", "zh-Latn", false},
		{"HK", "zh", true},
		{"HK", "zh-Hant", true},
		// Serbia (no local regions defined).
		{"RS", "en", false},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			format := address.GetFormat(tt.countryCode)
			locale := address.NewLocale(tt.locale)
			got := format.SelectRegions(locale)
			want := format.Regions
			if tt.wantLocal {
				want = format.LocalRegions
			}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("got %v, want %v", got, want)
			}
		})
	}
}

func TestRegionMap(t *testing.T) {
	r := address.NewRegionMap()
	if !reflect.DeepEqual(r, address.RegionMap{}) {
		t.Errorf(`got %v want %v`, r, address.RegionMap{})
	}
	if r.HasKey("06") {
		t.Error("got true, want false")
	}
	if r.Len() != 0 {
		t.Errorf(`got %v want 0`, r.Len())
	}

	r = address.NewRegionMap(
		"15", "Artemisa", "09", "Camagüey", "08", "Ciego de Ávila",
		"06", "Cienfuegos", "12", "Granma", "14", "Guantánamo",
	)

	region, ok := r.Get("06")
	if region != "Cienfuegos" || !ok {
		t.Errorf("got %v, %v want Cienfuegos, true", region, ok)
	}
	region, ok = r.Get("INVALID")
	if region != "" || ok {
		t.Errorf(`got %v, %v want "", false`, region, ok)
	}
	if !r.HasKey("06") {
		t.Error("got false, want true")
	}
	if r.HasKey("INVALID") {
		t.Error("got true, want false")
	}

	wantKeys := []string{
		"15", "09", "08", "06", "12", "14",
	}
	if !reflect.DeepEqual(r.Keys(), wantKeys) {
		t.Errorf(`got %v want %v`, r.Keys(), wantKeys)
	}
	if r.Len() != 6 {
		t.Errorf(`got %v want 6`, r.Len())
	}

	wantBytes := []byte(`{"15":"Artemisa","09":"Camagüey","08":"Ciego de Ávila","06":"Cienfuegos","12":"Granma","14":"Guantánamo"}`)
	gotBytes, err := json.Marshal(r)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	if !reflect.DeepEqual(gotBytes, wantBytes) {
		t.Errorf(`got %v want %v`, string(gotBytes), string(wantBytes))
	}
}

func TestCheckCountryCode(t *testing.T) {
	tests := []struct {
		countryCode string
		want        bool
	}{
		// Empty value.
		{"", true},
		// Valid country code.
		{"FR", true},
		// Invalid country code.
		{"ABCD", false},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := address.CheckCountryCode(tt.countryCode)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetCountryCodes(t *testing.T) {
	countryCodes := address.GetCountryCodes()
	for _, countryCode := range []string{"FR", "RS", "US"} {
		if !contains(countryCodes, countryCode) {
			t.Errorf("no %v country code found.", countryCode)
		}
	}
}

func TestGetCountryNames(t *testing.T) {
	got := address.GetCountryNames()
	want := map[string]string{
		"FR": "France",
		"RS": "Serbia",
	}
	for wantCode, wantName := range want {
		gotName, ok := got[wantCode]
		if !ok {
			t.Errorf("no %v country code found.", wantCode)
		}
		if gotName != wantName {
			t.Errorf("got %v, want %v", gotName, wantName)
		}
	}
}

func TestGetFormat(t *testing.T) {
	// Existing format.
	got := address.GetFormat("RS")
	want := address.Format{
		Layout:            "%1\n%2\n%3\n%P %L",
		Required:          []address.Field{address.FieldLine1, address.FieldLocality},
		PostalCodePattern: "\\d{5,6}",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}

	// Generic format.
	generic := address.GetFormat("ZZ")
	want = address.Format{
		Layout:   "%1\n%2\n%3\n%L",
		Required: []address.Field{address.FieldLine1, address.FieldLocality},
	}
	if !reflect.DeepEqual(generic, want) {
		t.Errorf("got %v, want %v", generic, want)
	}

	// Non-existent format.
	got = address.GetFormat("IC")
	if !reflect.DeepEqual(got, generic) {
		t.Errorf("got %v, want %v", got, generic)
	}
}

func TestGetFormats(t *testing.T) {
	formats := address.GetFormats()
	for _, countryCode := range []string{"ZZ", "RS"} {
		_, ok := formats[countryCode]
		if !ok {
			t.Errorf("no %v address format found.", countryCode)
		}
	}
}

func TestGetFormats_ValidRegionData(t *testing.T) {
	// Confirm that all regions contain valid utf8.
	// Avoids having the check at runtime, in RegionMap.MarshalJSON.
	formats := address.GetFormats()
	for countryCode, format := range formats {
		if format.Regions.Len() > 0 {
			keys := format.Regions.Keys()
			for _, key := range keys {
				value, _ := format.Regions.Get(key)
				if !utf8.ValidString(key) {
					t.Errorf("invalid key %v in %v regions", key, countryCode)
				}
				if !utf8.ValidString(value) {
					t.Errorf("invalid value %v for key %v in %v regions", value, key, countryCode)
				}
			}
		}
		if format.LocalRegions.Len() > 0 {
			keys := format.LocalRegions.Keys()
			for _, key := range keys {
				value, _ := format.LocalRegions.Get(key)
				if !utf8.ValidString(key) {
					t.Errorf("invalid key %v in %v local regions", key, countryCode)
				}
				if !utf8.ValidString(value) {
					t.Errorf("invalid value %v for key %v in %v regions", value, key, countryCode)
				}
			}
		}
	}
}

func contains(a []string, x string) bool {
	for _, v := range a {
		if v == x {
			return true
		}
	}
	return false
}
