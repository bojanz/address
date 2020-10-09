// Copyright (c) 2020 Bojan Zivanovic and contributors
// SPDX-License-Identifier: MIT

package address_test

import (
	"reflect"
	"testing"

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

func TestGetFormat(t *testing.T) {
	// Existing format.
	got := address.GetFormat("RS")
	want := address.Format{
		Layout:            "%1\n%2\n%P %L",
		Required:          []address.Field{address.FieldLine1, address.FieldLocality},
		PostalCodePattern: "\\d{5,6}",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}

	// Generic format.
	generic := address.GetFormat("ZZ")
	want = address.Format{
		Layout:   "%1\n%2\n%L",
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
