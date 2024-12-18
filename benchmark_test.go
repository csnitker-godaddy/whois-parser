package whoisparser

import "testing"

const benchmarkInput = `Domain Name: GIT.COM
Registry Domain ID: 72657455_DOMAIN_COM-VRSN
Registrar WHOIS Server: whois.uniregistrar.net
Registrar URL: http://uniregistry.com
Updated Date: 2019-05-17-T23:02:50Z
Creation Date: 2001-06-14-T10:32:43Z
Registrar Registration Expiration Date: 2020-06-14-T10:32:43Z
Registrar: UNIREGISTRAR CORP
Registrar IANA ID: 1659
Registrar Abuse Contact Email: abuse@uniregistry.com
Registrar Abuse Contact Phone: +1.4426008800
Domain Status: clientTransferProhibited http://www.icann.org/epp#clientTransferProhibited
Registry Registrant ID:
Registrant Name: PRIVACYDOTLINK CUSTOMER 1078347
Registrant Organization: 
Registrant Street: PO BOX 30485   
Registrant City: SEVEN MILE BEACH
Registrant State/Province: GRAND CAYMAN
Registrant Postal Code: KY1-1202
Registrant Country: KY
Registrant Phone: +1.3457495465
Registrant Phone Ext: 
Registrant Fax: 
Registrant Fax Ext: 
Registrant Email: 1078347@PRIVACY-LINK.COM
Registry Admin ID:
Admin Name: PRIVACYDOTLINK CUSTOMER 1078347
Admin Organization: 
Admin Street: PO BOX 30485   
Admin City: SEVEN MILE BEACH
Admin State/Province: GRAND CAYMAN
Admin Postal Code: KY1-1202
Admin Country: KY
Admin Phone: +1.3457495465
Admin Phone Ext: 
Admin Fax: 
Admin Fax Ext: 
Admin Email: 1078347@PRIVACY-LINK.COM
Registry Tech ID:
Tech Name: PRIVACYDOTLINK CUSTOMER 1078347
Tech Organization: 
Tech Street: PO BOX 30485   
Tech City: SEVEN MILE BEACH
Tech State/Province: GRAND CAYMAN
Tech Postal Code: KY1-1202
Tech Country: KY
Tech Phone: +1.3457495465
Tech Phone Ext: 
Tech Fax: 
Tech Fax Ext: 
Tech Email: 1078347@PRIVACY-LINK.COM
Name Server: ns2.venture.com
Name Server: ns1.venture.com
DNSSEC: unsigned
URL of the ICANN WHOIS Data Problem Reporting System: http://wdprs.internic.net/

>>> Last update of WHOIS database: 2019-09-30T15:00:39.872Z <<<

For more information on Whois status codes, please visit
https://www.icann.org/resources/pages/epp-status-codes-2014-06-16-en

TERMS OF USE: You  are  not  authorized  to  access or query our Whois
database through the use of high volume automated processes. Access to
the Whois database is  provided  solely to obtain information about or
related to a domain name  registration record, and no warranty is made
as to its accuracy or  fitness for any particular purpose..  You agree
that you may use this  Data only for lawful purposes and that under no
circumstances will you  use  this  data to allow, enable, or otherwise
support the transmission of  mass  unsolicited, commercial advertising
or solicitations via  e-mail,  telephone,  or facsimile.  Compilation,
repackaging,  dissemination  or  other  use  of this Data is expressly
prohibited   without  the  prior   written  consent  of  Uniregistrar,
Uniregistry Corp. or Uniregistry Ltd. (CA).  Uniregistrar reserves the
right to restrict your  access to  the  Whois  database  in  its  sole
discretion to ensure operational stability and police abuse.
`

func BenchmarkParse(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, err := Parse(benchmarkInput)
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkParseDateString(b *testing.B) {
	dates := []struct {
		date string
	}{
		{"09-Mar-2023"},
		{"31-Jul-2022"},
		{"2022-12-12T11:01:02Z"},
		{"2022-12-03"},
		{"2022. 12. 01."},
		{"2022-12-12 11:40:12"},
		{"2022.12.12 11:40:12"},
		{"28/06/2022 23:59:59"},
		{"24.10.2022"},
		{"2022-06-29 14:08:21+03"},
		{"31.8.2025 00:00:00"},
		{"01-10-2025"},
		{"20-Apr-2023 03:28:40"},
		{"2022-12-08 14:00:00 CLST"},
		{"December  2 2022"},
		{"Mon Jan  2 2006"},
		{"02/28/2025"},
		{"2001/03/22"},
		{"April 10 2023"},
		{"2025-Dec-11"},
		{"2025-Dec-11."},
		{"2024-06-05 00:00:00 (UTC+8)"},
		{"20221101 00:10:24"},
		{"Before 2006"},
		{"2001-06-14-T10:32:43Z"},
		{"02-01-2006 15:04:05 -07:00"},
	}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for _, tt := range dates {
				_, err := parseDateString(tt.date)
				if err != nil {
					b.Errorf("parseDateString(%s) error: %v", tt.date, err)
				}
			}
		}
	})
}
