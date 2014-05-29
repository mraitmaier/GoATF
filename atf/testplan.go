/*
 * testplan.go 
 *
 * History:
 *  1   May11 MR Initial version, limited testing
 *  2   May14 MR Updated version: XML handling simplified, added conversion to
 *               TestSet, appending test cases simplified
 */

package atf

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
    "bitbucket.org/miranr/goatf/atf/utils"
)

// Represents a test plan.
// Note that the TestPlan has a sibling in TestSet type: TestSet is an
// executable version of TestPlan.
type TestPlan struct {
	Name        string      `xml:"name,attr"`
	Description string      `xml:"Description"`
	Setup       *Action     `xml:"Setup"`
	Cleanup     *Action     `xml:"Cleanup"`
	Cases       []TestCase  `xml:"Cases>TestCase"`
}

//  Returns a plan text representation of the TestPlan instance.
func (tp *TestPlan) String() string {
	s := fmt.Sprintf("TestPlan: %q\n", tp.Name)
	return s
}

//  Returns a XML-encoded representation of the TestPlan instance.
func (tp *TestPlan) Xml() (string, error) {

    output, err := xml.MarshalIndent(tp, "  ", "    ")
    if err != nil {
        return "", err
    }

    return string(output), nil
}

//  Returns a JSON-encoded representation of the TestPlan instance.
func (tp *TestPlan) Json() (string, error) {
	b, err := json.Marshal(tp)
	if err != nil {
		return "", err
	}
	return string(b[:]), err
}

// Append one or more test cases to the list of test cases.
func (tp *TestPlan) Append (cases ...TestCase)  {
    tp.Cases = append(tp.Cases, cases...)
}

// Convert a TestPlan into a TestSet instance.
// Note that we force deep copy of data.
func (tp *TestPlan) ToTestSet() *TestSet {

    ts := new(TestSet)
    ts.Name = utils.CopyS(tp.Name) // TestSet name can (and should) be changed
    ts.Description = utils.CopyS(tp.Description)
    ts.TestPlan = utils.CopyS(tp.Name)
    *ts.Setup = *tp.Setup
    *ts.Cleanup = *tp.Cleanup
    ts.Sut = new(SysUnderTest) // return empty instance
    copy(ts.Cases, tp.Cases)

    return ts
}

// Creates a new TestPlan instance.
func CreateTestPlan(name, descr string,
	                setup, cleanup *Action) *TestPlan {
	tcs := make([]TestCase, 0)
	return &TestPlan{name, descr, setup, cleanup, tcs}
}

/*
 * TestPlan.findEmpty - function that creates the TestPlan struct
func (tp *TestPlan) findEmpty() int {
	for ix, tc := range tp.Cases {
		if &tc != nil && tc.Name == "" {
			return ix
		}
	}
	return -1
}
 */

/*
 * TestPlan.AppendCase - function that 
func (tp *TestPlan) AppendCase(tc *TestCase) []TestCase {
	if tc.Name != "" {
		l := len(tp.Cases)
		c := cap(tp.Cases)
		if l+1 > c {
			newlst := make([]TestCase, 0, 2*c)
			copy(newlst, tp.Cases)
			tp.Cases = newlst
		}
		tp.Cases = tp.Cases[0 : l+1]
		ix := tp.findEmpty()
		if ix != -1 {
			tp.Cases[ix] = *tc
		}
	}
	return tp.Cases
}
 */

/*
 * TestPlan.ExtendCaseList - function that extends tha list of test cases
func (tp *TestPlan) ExtendCaseList(tcl []TestCase) []TestCase {
	l := len(tp.Cases)
	if l+len(tcl) > cap(tp.Cases) {
		newlst := make([]TestCase, 0, cap(tp.Cases)+len(tcl))
		copy(newlst, tp.Cases)
		tp.Cases = newlst
	}
	tp.Cases = tp.Cases[0 : l+len(tcl)]
	empty := tp.findEmpty()
	if empty != -1 {
		for ix, tc := range tcl {
			tp.Cases[empty+ix] = tc
		}
	}
	return tp.Cases
}
*/
