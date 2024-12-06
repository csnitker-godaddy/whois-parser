package whoisparser

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/likexian/gokit/xfile"
	"github.com/likexian/gokit/xjson"
	"github.com/stretchr/testify/assert"
)

const allTLDDir = "testdata/alltlds"

func IsContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}

	return false
}

//nolint:lll
func TestParseSupported(t *testing.T) {
	extensions := []string{}
	domains := map[string]map[string]bool{}
	markdownContent := "| TLD | Parsed Successfully | Domain Match | Created Date Valid | Updated Date Valid | Expiration Date Valid | Registrar ID Valid | Registrar Name Valid | Whois Server Valid |\n"
	markdownContent += "|-----|--------------------|--------------|--------------------|--------------------|----------------------|-----------------|------------------|----------------|\n"
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
			assert.Nil(t, err, v.Name)
			err = xjson.Dump(allTLDDir+"/"+v.Name+".json", whoisInfo)
			assert.Nil(t, err)

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

			whoisServerValid := whoisInfo.Domain != nil && whoisInfo.Domain.WhoisServer != ""

			if _, ok := domains[extension]; !ok {
				domains[extension] = map[string]bool{
					"ParsedSuccessfully":  false,
					"DomainMatch":         false,
					"CreatedDateValid":    false,
					"UpdatedDateValid":    false,
					"ExpirationDateValid": false,
					"RegistrarIDValid":    false,
					"RegistrarNameValid":  false,
					"WhoisServerValid":    false,
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
			domains[extension]["WhoisServerValid"] = domains[extension]["WhoisServerValid"] || whoisServerValid

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
		markdownContent += fmt.Sprintf("| .%s | %s | %s | %s | %s | %s | %s | %s | %s |\n",
			extension,
			boolToColoredCheckbox(domainInfo["ParsedSuccessfully"]),
			boolToColoredCheckbox(domainInfo["DomainMatch"]),
			boolToColoredCheckbox(domainInfo["CreatedDateValid"]),
			boolToColoredCheckbox(domainInfo["UpdatedDateValid"]),
			boolToColoredCheckbox(domainInfo["ExpirationDateValid"]),
			boolToColoredCheckbox(domainInfo["RegistrarIDValid"]),
			boolToColoredCheckbox(domainInfo["RegistrarNameValid"]),
			boolToColoredCheckbox(domainInfo["WhoisServerValid"]))
	}

	err = xfile.WriteText(allTLDDir+"/SUPPORT.md", strings.TrimSpace(markdownContent))
	assert.Nil(t, err)
}

// Helper function to convert boolean to colored checkbox representation
func boolToColoredCheckbox(value bool) string {
	if value {
		return "[x]"
	}
	return "<span style=\"color:red;\">[ ]</span>"
}
