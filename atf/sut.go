/*
 * sut.go - file defining SysUnderTest struct and its methods
 * 
 * SUT is just descriptive structure that keeps some information about the
 * TestSet currently executed (used in configuration and in reports), it 
 * doesn't have any influence on execution.
 */

package atf

import (
	"encoding/json"
	"fmt"
)

/*
 * SysUnderTest - structure representing System Under Test (SUT)
 */
type SysUnderTest struct {

	Name        string `xml:"Name"`

	Systype     string `xml:"Type"`

	Version     string `xml:"Version"`

	Description string `xml:"Description"`

	IPaddr      string `xml:"IPAddress"`
}

/*
 * SysUnderTest.CreateSUT - create an SUT structure
 */
func CreateSUT(name, systype, version, descr, ip string) *SysUnderTest {
	return &SysUnderTest{name, systype, version, descr, ip}
}

/*
 * SysUnderTest.String - string representation of the SUT structure
 */
func (s *SysUnderTest) String() string {
	txt := "SystemUnderTest:\n"
	txt += fmt.Sprintf("   Name: %s\n", s.Name)
	txt += fmt.Sprintf("   Type: %s\n", s.Systype)
	txt += fmt.Sprintf("   Version: %s\n", s.Version)
	txt += fmt.Sprintf("   IP address: %s\n", s.IPaddr)
	txt += fmt.Sprintf("   Description:\n%s", s.Description)
	return txt
}

/*
 * SysUnderTest.Xml - XML representation of the SUT structure
 */
func (s *SysUnderTest) Xml() string {
	xml := "<SystemUnderTest>\n"
	xml += fmt.Sprintf("    <Name>%s</Name>\n", s.Name)
	xml += fmt.Sprintf("    <Type>%s</Type>\n", s.Systype)
	xml += fmt.Sprintf("    <Version>%s</Version>\n", s.Version)
	xml += fmt.Sprintf("    <IP>%s</IP>\n", s.IPaddr)
	xml += fmt.Sprintf("    <Description>%s</Description>\n", s.Description)
	xml += "</SystemUnderTest>\n"
	return xml
}

/*
 * SysUnderTest.Json - JSON representation of the SUT structure
 */
func (s *SysUnderTest) Json() (string, error) {
	b, err := json.Marshal(s)
	if err != nil {
		return "", err
	}
	return string(b[:]), err
}
