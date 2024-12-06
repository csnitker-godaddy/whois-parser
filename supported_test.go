package whoisparser

import (
	"encoding/json"
	"fmt"
	"github.com/likexian/gokit/xfile"
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
)

const allTLDDir = "testdata/alltlds"

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

func TestParseSupported(t *testing.T) {
	extensions := []string{}
	domains := map[string]map[string]bool{}
	parsedSuccessfullyFailed := map[string]map[string]bool{}
	markdownContent := "| TLD | Parsed Successfully | Domain Match | Created Date Valid | Updated Date Valid | Expiration Date Valid | Registrar ID Valid | Registrar Name Valid |\n"
	markdownContent += "|-----|--------------------|--------------|--------------------|--------------------|----------------------|-----------------|------------------|\n"

	dirs, err := xfile.ListDir(allTLDDir, xfile.TypeFile, -1)
	assert.Nil(t, err)

	for _, v := range dirs {
		if v.Name == "README.md" || v.Name == "SUPPORT.md" {
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
			whoisRaw, err := xfile.ReadText(allTLDDir + "/" + v.Name)
			assert.Nil(t, err)

			whoisInfo, err := Parse(whoisRaw)
			parsedSuccessfully := err == nil && whoisInfo.Domain != nil && whoisInfo.Domain.Punycode == domain

			domainMatch := false
			if whoisInfo.Domain != nil {
				domainMatch = whoisInfo.Domain.Punycode == domain
			}

			createdDateValid := whoisInfo.Domain != nil && whoisInfo.Domain.CreatedDate != "" && whoisInfo.Domain.CreatedDateInTime != nil
			updatedDateValid := whoisInfo.Domain != nil && whoisInfo.Domain.UpdatedDate != "" && whoisInfo.Domain.UpdatedDateInTime != nil
			expirationDateValid := whoisInfo.Domain != nil && whoisInfo.Domain.ExpirationDate != "" && whoisInfo.Domain.ExpirationDateInTime != nil

			registrarIDValid := whoisInfo.Registrar != nil && whoisInfo.Registrar.ID != ""
			registrarNameValid := whoisInfo.Registrar != nil && whoisInfo.Registrar.Name != ""

			if _, ok := domains[extension]; !ok {
				domains[extension] = map[string]bool{
					"ParsedSuccessfully":  false,
					"DomainMatch":         false,
					"CreatedDateValid":    false,
					"UpdatedDateValid":    false,
					"ExpirationDateValid": false,
					"RegistrarIDValid":    false,
					"RegistrarNameValid":  false,
				}
				extensions = append(extensions, extension)
			}

			domains[extension]["ParsedSuccessfully"] = domains[extension]["ParsedSuccessfully"] || parsedSuccessfully
			domains[extension]["DomainMatch"] = domains[extension]["DomainMatch"] || domainMatch
			domains[extension]["CreatedDateValid"] = domains[extension]["CreatedDateValid"] || createdDateValid
			domains[extension]["UpdatedDateValid"] = domains[extension]["UpdatedDateValid"] || updatedDateValid
			domains[extension]["ExpirationDateValid"] = domains[extension]["ExpirationDateValid"] || expirationDateValid
			domains[extension]["RegistrarIDValid"] = domains[extension]["RegistrarIDValid"] || registrarIDValid
			domains[extension]["RegistrarNameValid"] = domains[extension]["RegistrarNameValid"] || registrarNameValid

			if !parsedSuccessfully {
				if _, ok := parsedSuccessfullyFailed[extension]; !ok {
					parsedSuccessfullyFailed[extension] = map[string]bool{
						"ParsedSuccessfully":  false,
						"DomainMatch":         false,
						"CreatedDateValid":    false,
						"UpdatedDateValid":    false,
						"ExpirationDateValid": false,
						"RegistrarIDValid":    false,
						"RegistrarNameValid":  false,
					}
				}
				parsedSuccessfullyFailed[extension]["ParsedSuccessfully"] = parsedSuccessfullyFailed[extension]["ParsedSuccessfully"] || parsedSuccessfully
				parsedSuccessfullyFailed[extension]["DomainMatch"] = parsedSuccessfullyFailed[extension]["DomainMatch"] || domainMatch
				parsedSuccessfullyFailed[extension]["CreatedDateValid"] = parsedSuccessfullyFailed[extension]["CreatedDateValid"] || createdDateValid
				parsedSuccessfullyFailed[extension]["UpdatedDateValid"] = parsedSuccessfullyFailed[extension]["UpdatedDateValid"] || updatedDateValid
				parsedSuccessfullyFailed[extension]["ExpirationDateValid"] = parsedSuccessfullyFailed[extension]["ExpirationDateValid"] || expirationDateValid
				parsedSuccessfullyFailed[extension]["RegistrarIDValid"] = parsedSuccessfullyFailed[extension]["RegistrarIDValid"] || registrarIDValid
				parsedSuccessfullyFailed[extension]["RegistrarNameValid"] = parsedSuccessfullyFailed[extension]["RegistrarNameValid"] || registrarNameValid
			}

			if t.Failed() {
				t.Logf("whoisRaw: %s", whoisRaw)
				e := json.NewEncoder(os.Stdout)
				e.SetIndent("", "  ")
				e.Encode(whoisInfo)
			}
		})
	}

	for _, extension := range extensions {
		domainInfo := domains[extension]
		if domainInfo["ParsedSuccessfully"] {
			markdownContent += fmt.Sprintf("| .%s | %s | %s | %s | %s | %s | %s | %s |\n",
				extension,
				boolToCheckbox(domainInfo["ParsedSuccessfully"]),
				boolToCheckbox(domainInfo["DomainMatch"]),
				boolToCheckbox(domainInfo["CreatedDateValid"]),
				boolToCheckbox(domainInfo["UpdatedDateValid"]),
				boolToCheckbox(domainInfo["ExpirationDateValid"]),
				boolToCheckbox(domainInfo["RegistrarIDValid"]),
				boolToCheckbox(domainInfo["RegistrarNameValid"]))
		}
	}

	markdownContent += "\n\n### Failed TLDs\n"
	markdownContent += "| TLD | Parsed Successfully | Domain Match | Created Date Valid | Updated Date Valid | Expiration Date Valid | Registrar ID Valid | Registrar Name Valid |\n"
	markdownContent += "|-----|--------------------|--------------|--------------------|--------------------|----------------------|-----------------|------------------|\n"

	for extension, domainInfo := range parsedSuccessfullyFailed {
		markdownContent += fmt.Sprintf("| .%s | %s | %s | %s | %s | %s | %s | %s |\n",
			extension,
			boolToCheckbox(domainInfo["ParsedSuccessfully"]),
			boolToCheckbox(domainInfo["DomainMatch"]),
			boolToCheckbox(domainInfo["CreatedDateValid"]),
			boolToCheckbox(domainInfo["UpdatedDateValid"]),
			boolToCheckbox(domainInfo["ExpirationDateValid"]),
			boolToCheckbox(domainInfo["RegistrarIDValid"]),
			boolToCheckbox(domainInfo["RegistrarNameValid"]))
	}

	err = xfile.WriteText(allTLDDir+"/SUPPORT.md", strings.TrimSpace(markdownContent))
	assert.Nil(t, err)
}

// Helper function to convert boolean to checkbox representation
func boolToCheckbox(value bool) string {
	if value {
		return "[x]"
	}
	return "[ ]"
}

// Helper function to convert boolean to colored checkbox representation
func boolToColoredCheckbox(value bool) string {
	if value {
		return "[x]"
	}
	return "<span style=\"color:red;\">[ ]</span>"
}
