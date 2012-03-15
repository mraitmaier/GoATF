/*
 * testplan.go 
 *
 * History:
 *  0.1   May11 MR Initial version, limited testing
 */

package atf

import (
	"fmt"
	"os"
	"json"
)

/*
 * TestPlan -
 */
type TestPlan struct {
	Name        string
	Description string
	Setup       *Action
	Cleanup     *Action
	Cases       []TestCase
}

/*
 * TestPlan.String - function that creates the TestPlan struct
 */
func (tp *TestPlan) String() string {
	s := fmt.Sprintf("TestPlan: %q\n", tp.Name)
	return s
}

/*
 * TestPlan.Xml - function that creates the TestPlan struct
 */
func (tp *TestPlan) Xml() string {
	xml := fmt.Sprintf("<TestPlan name=%q\n", tp.Name)
	if tp.Setup != nil {
		xml += tp.Setup.Xml()
	} else {
		xml += "<Setup />\n"
	}
	if tp.Cases != nil {

		for _, tc := range tp.Cases {
			xml += tc.Xml()
		}
		//xml += tp.Cases.Xml()
	} else {
		xml += "<TestCase />\n"
	}
	if tp.Cleanup != nil {
		xml += tp.Cleanup.Xml()
	} else {
		xml += "<Cleanup />\n"
	}
	xml += fmt.Sprintln("</TestPlan>")
	return xml
}

/*
 * TestPlan.Json - function that 
 */
func (tp *TestPlan) Json() (string, os.Error) {
	b, err := json.Marshal(tp)
	if err != nil {
		return "", err
	}
	return string(b[:]), err
}

/*
 * CreateTestPlan - function that creates the TestPlan struct
 */
const defTpListCap = 10

func CreateTestPlan(name string,
descr string,
setup *Action,
cleanup *Action) *TestPlan {
	tcs := make([]TestCase, 0, defTpListCap)
	return &TestPlan{name, descr, setup, cleanup, tcs}
}

/*
 * TestPlan.findEmpty - function that creates the TestPlan struct
 */
func (tp *TestPlan) findEmpty() int {
	for ix, tc := range tp.Cases {
		if &tc != nil && tc.Name == "" {
			return ix
		}
	}
	return -1
}

/*
 * TestPlan.AppendCase - function that 
 */
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

/*
 * TestPlan.ExtendCaseList - function that extends tha list of test cases
 */
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
