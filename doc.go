// Copyright (c) 2020 Bojan Zivanovic and contributors
// SPDX-License-Identifier: MIT

// Package address handles address representation, validation and formatting.
//
// Features:
//   1. Address struct.
//   2. Address formats for ~200 countries.
//   3. Regions for ~40 countries, with local names where relevant (e.g: Okinawa / 沖縄県).
//   4. Country list, powered by CLDR v37.
//   5. HTML formatter.
//   6. HTTP handler for serving address formats and regions as JSON: only ~14kb gzipped!
//
// Address formats
//
// The following address format information is available for each country:
//   - Which fields are used, and in which order.
//   - Which fields are required.
//   - Labels for the sublocality, locality, region and postal code fields.
//   - Regular expression pattern for validating postal codes.
//   - Regions and how to display them in an address.
//
// Certain countries (e.g. China, Japan, Russia, Ukraine) have region names defined in both Latin
// and local scripts. The script is selected based on locale. For example, the "ru" locale will
// use Russian regions in Cyrilic,  while "ru-Latn" and other locales will use the Latin version.
//
// Helpers are provided for validating required fields, regions, postal codes.
//
// Format data was generated from Google's Address Data but isn't automatically regenerated, to
// allow the community to submit their own corrections directly to the package.
//
// Countries
//
// The country list is auto-generated from CLDR. Updating to the latest CLDR release is always
// one `go generate` away.
//
// Most software uses the CLDR country list instead of the ISO one because the CLDR country names
// match their colloquial usage more closely (e.g. "Russia" instead of "Russian Federation").
//
// To reduce the size of the included data, this package only includes country names in English.
// Translated country names can be fetched on the frontend via Intl.DisplayNames.
// Alternatively, one can plug in x/text/language/display by setting a custom CountryMapper on the formatter.
//
package address
