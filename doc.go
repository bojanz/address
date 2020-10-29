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
// Address struct
//
// Represents an address as commonly handled by web applications and APIs.
//
//   type Address struct {
//      Line1 string
//      Line2 string
//      Line3 string
//      // Sublocality is the neighborhood/suburb/district.
//      Sublocality string
//      // Locality is the city/village/post town.
//      Locality string
//      // Region is the state/province/prefecture.
//      // An ISO code is used when available.
//      Region string
//      // PostalCode is the postal/zip/pin code.
//      PostalCode string
//      // CountryCode is the two-letter code as defined by CLDR.
//      CountryCode string
//   }
//
// Recipient fields such as FirstName/LastName are not included since they are usually
// present on the parent (Contact/Customer/User) struct. This allows the package to
// avoid being opinionated about name handling.
//
// There are three line fields, matching the HTML5 autocomplete spec and many shipping APIs.
// This leaves enough space for specifying an organization and department, house or hotel
// name, and other similar "care of" use cases. When mapping to an API that only has two address
// lines, Line3 can be appended to Line2, separated by a comma.
//
// Address formats
//
// The following information is available:
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
