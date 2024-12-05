/*
 * Copyright 2014-2024 Li Kexian
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * Go module for domain whois information parsing
 * https://www.likexian.com/
 */

package whoisparser

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"sort"
	"strings"
	"testing"

	"github.com/likexian/gokit/xfile"
	"github.com/likexian/gokit/xjson"
	"golang.org/x/net/idna"
)

const (
	noterrorDir  = "testdata/noterror"
	notfoundDir  = "testdata/notfound"
	verifiedList = `
# WhoisParser

## Overview

It is supposed to be working with all domain extensions,

but verified extensions as below must works, because I have checked them one by one manually.

If there is any problem, please feel free to open a new issue.

## Verified Extensions

| extension | whois | output | verified |
| --------- | ----- | ------ | :------: |
`
)

func TestVersion(t *testing.T) {
	assert.Contains(t, Version(), ".")
	assert.Contains(t, Author(), "likexian")
	assert.Contains(t, License(), "Apache License")
}

func TestParseError(t *testing.T) {
	tests := map[error]string{
		ErrNotFoundDomain:    "No matching record.",
		ErrReservedDomain:    "Reserved Domain Name",
		ErrPremiumDomain:     "This platinum domain is available for purchase.",
		ErrBlockedDomain:     "This name subscribes to the Uni EPS+ product",
		ErrDomainDataInvalid: "connect to whois server failed: dial tcp 43: i/o timeout",
		ErrDomainLimitExceed: "WHOIS LIMIT EXCEEDED - SEE WWW.PIR.ORG/WHOIS FOR DETAILS",
	}

	for e, v := range tests {
		_, err := Parse(v)
		assert.Equal(t, err, e)
	}

	_, err := Parse(`Domain Name: likexian-no-money-registe.ai
	Domain Status: No Object Found`)
	assert.Equal(t, err, ErrNotFoundDomain)
}

func ContainsAny(s string, subs []string) bool {
	for _, sub := range subs {
		if strings.Contains(s, sub) {
			return true
		}
	}

	return false
}

func IsContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}

	return false
}

func TestParse(t *testing.T) {
	extensions := []string{}
	domains := map[string][]string{}

	dirs, err := xfile.ListDir(noterrorDir, xfile.TypeFile, -1)
	assert.Nil(t, err)

	// TODO Remove this part
	f, err := os.Create("toberemoved")
	assert.Nil(t, err)

	defer f.Close()
	// END TODO

	for _, v := range dirs {
		if v.Name == "README.md" {
			continue
		}

		domain := strings.Split(v.Name, "_")[1]
		extension := ""
		if strings.Contains(v.Name, ".") {
			extension = domain[strings.LastIndex(domain, ".")+1:]
		}

		if IsContains([]string{"pre", "json"}, extension) {
			continue
		}

		t.Run(v.Name, func(t *testing.T) {
			whoisRaw, err := xfile.ReadText(noterrorDir + "/" + v.Name)
			assert.Nil(t, err)

			whoisInfo, err := Parse(whoisRaw)
			if !assert.Nil(t, err, v.Name) {
				t.Errorf("error: %s", err)
				return
			}

			if whoisInfo.Domain == nil {
				t.Errorf("domain is nil")
				t.FailNow()
			}
			assert.Equal(t, whoisInfo.Domain.Punycode, domain)
			assert.Equal(t, whoisInfo.Domain.Extension, extension)

			// TODO Remove this part
			if whoisInfo.Domain.CreatedDate != "" && whoisInfo.Domain.CreatedDateInTime == nil {
				f.WriteString(fmt.Sprintf("(CreatedDate) %s: %s\n", v.Name, whoisRaw))
			}

			if whoisInfo.Domain.UpdatedDate != "" && whoisInfo.Domain.UpdatedDateInTime == nil {
				f.WriteString(fmt.Sprintf("(UpdatedDate) %s: %s\n", v.Name, whoisRaw))
			}

			if whoisInfo.Domain.ExpirationDate != "" && whoisInfo.Domain.ExpirationDateInTime == nil {
				f.WriteString(fmt.Sprintf("(ExpirationDate) %s: %s\n", v.Name, whoisRaw))
			}
			// END TODO

			if !IsContains([]string{"", "at", "aq", "br", "ch", "cl", "de", "edu", "eu", "fr", "gov", "hk",
				"hm", "int", "it", "jp", "kr", "kz", "mo", "nl", "nz", "pl", "pm", "re", "ro", "ru", "su", "tf", "ee",
				"tk", "travel", "tv", "tw", "uk", "wf", "yt", "ir", "fi", "rs", "dk", "by", "ua",
				"xn--mgba3a4f16a", "xn--p1ai", "sg", "se", "sk", "nu", "hu"}, extension) {
				assert.NotZero(t, whoisInfo.Domain.ID, "domain id is empty")
			}

			if !IsContains([]string{"at", "ch", "edu", "eu", "int", "kr", "mo", "tw", "ir", "pl", "tk", "by",
				"xn--mgba3a4f16a", "hu"}, extension) {
				assert.NotZero(t, whoisInfo.Domain.Status, "status is empty")
			}

			if ContainsAny(whoisRaw, []string{"signedDelegation", "Signed delegation", "Signed "}) {
				assert.True(t, whoisInfo.Domain.DNSSec, "dnssec is false")
			} else {
				assert.False(t, whoisInfo.Domain.DNSSec, "dnssec is true")
			}

			if !IsContains([]string{"aero", "ai", "at", "aq", "asia", "berlin", "biz", "br", "ch", "cn",
				"co", "cymru", "cl", "cx", "de", "edu", "eu", "fr", "gov", "hk", "hm", "in", "int", "it", "jp", "kr",
				"la", "london", "me", "mo", "museum", "name", "nl", "nz", "pm", "re", "ro", "ru", "sh", "sk",
				"kz", "su", "tel", "ee", "tf", "tk", "travel", "tw", "uk", "us", "wales", "wf", "xxx",
				"yt", "ir", "fi", "rs", "dk", "by", "ua", "sg", "st", "xn--mgba3a4f16a", "xn--fiqs8s", "xn--p1ai",
				"se", "nu", "hu"}, extension) {
				if whoisInfo.Registrar != nil && !IsContains([]string{"9999", "119"}, whoisInfo.Registrar.ID) {
					assert.NotZero(t, whoisInfo.Domain.WhoisServer, "whois server is empty: %s", whoisRaw)
				}
			}

			if !IsContains([]string{"gov", "name", "tw", "hu"}, extension) {
				assert.NotZero(t, whoisInfo.Domain.NameServers, "name servers is empty")
			}

			if !IsContains([]string{"aq", "ai", "at", "au", "de", "eu", "gov", "hm", "name", "nl", "nz", "ir", "tk",
				"xn--mgba3a4f16a"}, extension) &&
				!strings.Contains(domain, "ac.jp") &&
				!strings.Contains(domain, "co.jp") &&
				!strings.Contains(domain, "go.jp") &&
				!strings.Contains(domain, "ne.jp") {
				assert.NotZero(t, whoisInfo.Domain.CreatedDate, "created date is empty")
				assert.NotNil(t, whoisInfo.Domain.CreatedDateInTime, "created date in time is empty")
			}

			if whoisInfo.Domain.UpdatedDate != "" && !IsContains([]string{"aq", "ai", "at", "ch", "cn", "eu", "gov", "hk", "hm", "mo",
				"name", "nl", "ro", "ru", "su", "tk", "tw", "dk", "xn--fiqs8s", "xn--p1ai", "hu"}, extension) {
				assert.NotNil(t, whoisInfo.Domain.UpdatedDateInTime, "updated date in time is empty")
			}

			if !IsContains([]string{"", "ai", "at", "aq", "au", "br", "ch", "de", "eu", "gov", "ee",
				"hm", "int", "name", "nl", "nz", "tk", "kz", "hu", "uz"}, extension) &&
				!strings.Contains(domain, "ac.jp") &&
				!strings.Contains(domain, "co.jp") &&
				!strings.Contains(domain, "go.jp") &&
				!strings.Contains(domain, "ne.jp") {
				assert.NotZero(t, whoisInfo.Domain.ExpirationDate)
				assert.NotNil(t, whoisInfo.Domain.ExpirationDateInTime)
			}

			if !IsContains([]string{"", "ai", "at", "aq", "au", "br", "bf", "ca", "ch", "cn", "cl", "cx", "de",
				"edu", "eu", "fr", "gov", "gs", "hk", "hm", "int", "it", "jp", "kr", "kz", "la", "mo", "nl",
				"nz", "pl", "pm", "re", "ro", "ru", "su", "sk", "tf", "tk", "tw", "uk", "wf", "yt", "ir", "fi", "rs",
				"ee", "dk", "by", "ua", "xn--mgba3a4f16a", "xn--fiqs8s", "xn--p1ai", "sg", "se", "nu", "hu"}, extension) {
				if assert.NotZero(t, whoisInfo.Registrar, "registrar is nil") {
					assert.NotZero(t, whoisInfo.Registrar.ID, "registrar id is empty")
				}
			}

			if !IsContains([]string{"", "at", "aq", "br", "de",
				"edu", "gov", "hm", "int", "jp", "mo", "tk", "ir", "dk", "xn--mgba3a4f16a", "hu"}, extension) {
				if assert.NotZero(t, whoisInfo.Registrar, "registrar is nil") {
					assert.NotZero(t, whoisInfo.Registrar.Name, "registrar name is empty")
				}
			}

			//if !IsContains([]string{"", "aero", "ai", "at", "aq", "asia", "au", "br", "ch", "cn", "de",
			//	"edu", "gov", "hk", "hm", "int", "jp", "kr", "kz", "la", "london", "love", "mo",
			//	"museum", "name", "nl", "nz", "pl", "ru", "sk", "sg", "su", "tk", "top", "ir", "fi", "rs", "dk", "by", "ua",
			//	"xn--mgba3a4f16a", "xn--fiqs8s", "xn--p1ai", "se", "nu", "hu"}, extension) {
			//	assert.NotZero(t, whoisInfo.Registrar.ReferralURL)
			//}

			err = xjson.Dump(noterrorDir+"/"+v.Name+".json", whoisInfo)
			assert.Nil(t, err)

			extension, _ = idna.ToUnicode(extension)
			if !IsContains(extensions, extension) {
				extensions = append(extensions, extension)
			}

			if _, ok := domains[extension]; !ok {
				domains[extension] = []string{}
			}

			domains[extension] = append(domains[extension], domain)

			if t.Failed() {
				t.Logf("whoisRaw: %s", whoisRaw)
				e := json.NewEncoder(os.Stdout)
				e.SetIndent("", "  ")
				e.Encode(whoisInfo)
			}
		})
	}

	sort.Strings(extensions)
	verified := verifiedList

	for _, extension := range extensions {
		sort.Strings(domains[extension])
		for _, domain := range domains[extension] {
			unicodeDomain, _ := idna.ToUnicode(domain)
			asciiExtension, _ := idna.ToASCII(extension)
			if asciiExtension == "" {
				asciiExtension = domain
			}
			verified += fmt.Sprintf("| .%s | [%s](%s_%s) | [%s](%s_%s.json) | √ |\n",
				extension, unicodeDomain, asciiExtension, domain, unicodeDomain, asciiExtension, domain)
		}
	}

	err = xfile.WriteText(noterrorDir+"/README.md", strings.TrimSpace(verified))
	assert.Nil(t, err)
}

func TestAssearchDomain(t *testing.T) {
	tests := []struct {
		whois     string
		name      string
		extension string
	}{
		{"Domain: example.com\n", "example", "com"},
		{"Domain Name: example.com\n", "example", "com"},
		{"Domain_Name: example.com\n", "example", "com"},

		{"Domain: com\n", "com", ""},
		{"Domain Name: com\n", "com", ""},
		{"Domain_Name: com\n", "com", ""},

		{"Domain Name: 示例.中国\n", "示例", "中国"},
		{"Domain Name: 中国\n", "中国", ""},
	}

	for _, v := range tests {
		name, extension := searchDomain(v.whois)
		assert.Equal(t, name, v.name)
		assert.Equal(t, extension, v.extension)
	}
}
