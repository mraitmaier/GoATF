/*
 * testrslt.go - implementation of the TestResult type
 *
 * This type defines the valid test results (pass/fail/xfail...) and valid 
 * operations on them.
 */

package atf

import (
    "encoding/xml"
)

/*
 * validTestResult - a slice of valid test result (string) values
 */
var validTestResults = []string{"UnknownResult", "Pass", "Fail",
	"XFail", "NotTested"}

/*
 * ValidTestResult - a function that returns the slice of valid test result
 * values; we use a function to publicly make this slice a constant
 */
func ValidTestResults() []string { return validTestResults }

/*
 * IsValidTestResult - a function that checks the validity of the test result
 * value; returns true or false, of course
 */
func IsValidTestResult(val string) bool {
	status := false
	for _, v := range validTestResults {
		if v == val {
			status = true
			break
		}
	}
	return status
}

/*
 * TestResult - a struct hiding a string value for test result
type TestResult struct {
	result string // this data is private
}
 */

type TestResult string

/*
func (tr TestResult) UnmarshalXML(d *xml.Decoder, s xml.StartElement) error {

    var data string
    err := d.DecodeElement(&data, &s)
    if err != nil {
        return err
    }
    tr(data)
    return nil
}
*/

/*
 * TestResult.String - String method for TestResult is defined
func (tr *TestResult) String() string { return tr.result }
 */

/*
 * TestResult.Get - get a value of test result
func (tr *TestResult) Get() string { return tr.result }
 */

/*
 * TestResult.Set - set a value of test result
 */
func (tr TestResult) Set(val string) (err error) {
	if IsValidTestResult(val) {
		tr = TestResult(val)
	} else {
		err = ATFError_Invalid_Test_Result
	}
	return
}

func (tr *TestResult) Xml() (x string, err error) {

	x = ""
	b, err := xml.MarshalIndent(tr, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b[:]), nil
}

