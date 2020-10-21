// Copyright (c) 2020 Bojan Zivanovic and contributors
// SPDX-License-Identifier: MIT

package address_test

import (
	"testing"

	"github.com/bojanz/address"
)

func TestNewLocale(t *testing.T) {
	tests := []struct {
		id   string
		want address.Locale
	}{
		{"", address.Locale{}},
		{"de", address.Locale{Language: "de"}},
		{"de-CH", address.Locale{Language: "de", Territory: "CH"}},
		{"es-419", address.Locale{Language: "es", Territory: "419"}},
		{"sr-Cyrl", address.Locale{Language: "sr", Script: "Cyrl"}},
		{"sr-Latn-RS", address.Locale{Language: "sr", Script: "Latn", Territory: "RS"}},
		{"yue-Hans", address.Locale{Language: "yue", Script: "Hans"}},
		// ID with the wrong case, ordering, delimeter.
		{"SR_rs_LATN", address.Locale{Language: "sr", Script: "Latn", Territory: "RS"}},
		// ID with a variant. Variants are unsupported and ignored.
		{"ca-ES-VALENCIA", address.Locale{Language: "ca", Territory: "ES"}},
	}
	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			got := address.NewLocale(tt.id)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLocale_String(t *testing.T) {
	tests := []struct {
		locale address.Locale
		want   string
	}{
		{address.Locale{}, ""},
		{address.Locale{Language: "de"}, "de"},
		{address.Locale{Language: "de", Territory: "CH"}, "de-CH"},
		{address.Locale{Language: "sr", Script: "Cyrl"}, "sr-Cyrl"},
		{address.Locale{Language: "sr", Script: "Latn", Territory: "RS"}, "sr-Latn-RS"},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			id := tt.locale.String()
			if id != tt.want {
				t.Errorf("got %v, want %v", id, tt.want)
			}
		})
	}
}

func TestLocale_IsEmpty(t *testing.T) {
	tests := []struct {
		locale address.Locale
		want   bool
	}{
		{address.Locale{}, true},
		{address.Locale{Language: "de"}, false},
		{address.Locale{Language: "de", Territory: "CH"}, false},
		{address.Locale{Language: "sr", Script: "Cyrl"}, false},
		{address.Locale{Language: "sr", Script: "Latn", Territory: "RS"}, false},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			empty := tt.locale.IsEmpty()
			if empty != tt.want {
				t.Errorf("got %v, want %v", empty, tt.want)
			}
		})
	}
}
