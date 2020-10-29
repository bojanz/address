// Copyright (c) 2020 Bojan Zivanovic and contributors
// SPDX-License-Identifier: MIT

// +build ignore

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
	"sort"
	"strings"
	"text/template"
	"time"
)

const dataTemplate = `// Code generated by go generate; DO NOT EDIT.
//go:generate go run gen.go

package address

// CLDRVersion is the CLDR version from which the data is derived.
const CLDRVersion = "{{ .CLDRVersion }}"

var countries = map[string]string{
	{{ export .Countries }}
}
`

func main() {
	log.Println("Fetching data...")
	cldrVersion, err := fetchVersion()
	if err != nil {
		log.Fatal(err)
	}
	countries, err := fetchCountries()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Processing...")
	os.Remove("countries.go")
	f, err := os.Create("countries.go")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	funcMap := template.FuncMap{
		"export": export,
	}
	t, err := template.New("data").Funcs(funcMap).Parse(dataTemplate)
	if err != nil {
		log.Fatal(err)
	}
	t.Execute(f, struct {
		CLDRVersion string
		Countries   map[string]string
	}{
		CLDRVersion: cldrVersion,
		Countries:   countries,
	})

	log.Println("Done.")
}

// fetchVersion fetches the CLDR version from GitHub.
func fetchVersion() (string, error) {
	data, err := fetchURL("https://raw.githubusercontent.com/unicode-cldr/cldr-localenames-full/master/package.json")
	if err != nil {
		return "", fmt.Errorf("fetchVersion: %w", err)
	}
	aux := struct {
		Version string
	}{}
	if err := json.Unmarshal(data, &aux); err != nil {
		return "", fmt.Errorf("fetchVersion: %w", err)
	}

	return aux.Version, nil
}

// fetchCountries fetches the CLDR country names from GitHub.
//
// The JSON version of CLDR data is used because it is more convenient
// to parse. See https://github.com/unicode-cldr/cldr-json for details.
func fetchCountries() (map[string]string, error) {
	data, err := fetchURL("https://raw.githubusercontent.com/unicode-cldr/cldr-localenames-full/master/main/en/territories.json")
	if err != nil {
		return nil, fmt.Errorf("fetchCountries: %w", err)
	}
	aux := struct {
		Main struct {
			En struct {
				LocaleDisplayNames struct {
					Territories map[string]string
				}
			}
		}
		Version string
	}{}
	if err := json.Unmarshal(data, &aux); err != nil {
		return nil, fmt.Errorf("fetchCountries: %w", err)
	}
	countries := aux.Main.En.LocaleDisplayNames.Territories
	for countryCode := range countries {
		if len(countryCode) > 2 {
			delete(countries, countryCode)
		}
		if contains([]string{"EU", "EZ", "UN", "QO", "XA", "XB", "ZZ"}, countryCode) {
			delete(countries, countryCode)
		}
	}

	return countries, nil
}

func fetchURL(url string) ([]byte, error) {
	client := http.Client{Timeout: 15 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("fetchURL: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetchURL: Get %q: %v", url, resp.Status)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("fetchURL: Get %q: %w", url, err)
	}

	return data, nil
}

func contains(a []string, x string) bool {
	for _, v := range a {
		if v == x {
			return true
		}
	}
	return false
}

func export(i interface{}) string {
	v := reflect.ValueOf(i)
	switch v.Kind() {
	case reflect.Map:
		return exportMap(v)
	default:
		return fmt.Sprintf("%#v", i)
	}
}

func exportMap(v reflect.Value) string {
	var values []string
	flipped := make(map[string]string, v.Len())
	iter := v.MapRange()
	for iter.Next() {
		key := iter.Key().String()
		value := iter.Value().String()
		values = append(values, value)
		flipped[value] = key
	}
	sort.Slice(values, func(i, j int) bool {
		// Compare Å as A to avoid having it come after Z.
		v1 := strings.Replace(values[i], "Å", "A", 1)
		v2 := strings.Replace(values[j], "Å", "A", 1)
		return v1 < v2
	})

	b := strings.Builder{}
	for i, value := range values {
		key := flipped[value]
		fmt.Fprintf(&b, "%q: %#v,", key, value)
		if i+1 < len(values) {
			fmt.Fprintf(&b, "\n\t")
		}
	}

	return b.String()
}