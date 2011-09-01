/*
 * sut.go - file defining SysUnderTest struct and its methods
 */

package atf

import ("fmt"
        "json"
        "strings"
        "os")

/*
 * SUTType - enumeration defining System Under Test types
 */
type SUTType int
const (
    UnknownSUTType SUTType = iota
    Hardware
    Software
    Both
    )

/*
 * SUTType.String - return string representation of the SUTType
 */
func (t SUTType) String() string {
    var s string
    switch t {
        case UnknownSUTType: s = "Unknown System Under Test"
        case Hardware: s = "Hardware"
        case Software: s = "Software"
        case Both: s = "Hardware & Software"
    }
    return s
}

/*
 * SUTType.Value - return SUTType value from given string value
 */
func (t SUTType) Value(val string) SUTType {
    switch strings.ToLower(val) {
        case "hardware", "hw": t = Hardware
        case "software", "sw": t = Software
        case "both", "hardware & software", "hw & sw": t = Both
        default: t = UnknownSUTType
    }
    return t
}

/*
 * SysUnderTest - structure representing System Under Test (SUT)
 */
type SysUnderTest struct {
    Name string
    Systype SUTType
    Version string
    Description string
    IPaddr string
}

/*
 * SysUnderTest.CreateSUT - create an SUT structure
 */
func CreateSUT(name string, systype SUTType, version string,
        descr string, ip string) *SysUnderTest {
    return &SysUnderTest{name, systype, version, descr, ip}
}

/*
 * SysUnderTest.String - string representation of the SUT structure
 */
func (s *SysUnderTest) String() string {
    txt := fmt.Sprintf("SystemUnderTest: %s\n", s.Name )
    txt += fmt.Sprintf("   Type: %s\n", s.Systype.String())
    txt += fmt.Sprintf("   Version: %s\n", s.Version)
    txt += fmt.Sprintf("   IP address: %s\n", s.IPaddr)
    txt += fmt.Sprintf("   Description:\n%s", s.Description)
    return txt
}

/*
 * SysUnderTest.Xml - XML representation of the SUT structure
 */
func (s *SysUnderTest) Xml() string {
    xml := fmt.Sprintf("<SystemUnderTest name=\"%s\">\n", s.Name)
    xml += fmt.Sprintf("    <Type>%s</Type>\n", s.Systype.String())
    xml += fmt.Sprintf("    <Version>%s</Version>\n", s.Version)
    xml += fmt.Sprintf("    <IP>%s</IP\n", s.IPaddr)
    xml += fmt.Sprintf("    <Description>%s</Description>\n", s.Description)
    xml += "</SystemUnderTest>\n"
    return xml
}

/*
 * SysUnderTest.Json - JSON representation of the SUT structure
 */
func (s *SysUnderTest) Json() (string, os.Error) { 
    b, err := json.Marshal(s)
    if err != nil { return "", err }
    return string(b[:]), err
}
