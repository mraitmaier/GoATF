/*
 * testset.go - implementation of the TestSet type
 *
 * The TestSet is an executable version of TestPlan: TestSet is executed, while
 * the TestPlan should serve only as a reference to a document that is stored
 * into database.
 *
 * NOTE: I'm not sure this is the right approach, but as a Go-learning step,
 * this is how it's currently done. It might even disappear in the future (or
 * the TestSet type).
 *
 * History:
 *  0.1   Apr10 MR Initial version, limited testing
 */

package atf

import (
	"encoding/json"
	"fmt"
)

/*****************************************************************************
 * TestSet - this is a subset of the TestPlan structure. They do differ in
 * two detailis: 
 *      1. TestSet implements the Execute() method, while TestPlan
 *      is not.
 *      2. Since TestSet is a subset of TestPlan, it holds a name of the
 *      originating test plan.
 * Otherwise they're completely the same. 
 */
type TestSet struct {

	// a test set name, of course; in XML, this is an attribute
	Name string `xml:"name,attr"`

	// a arbitrary long text description of the test set
	Description string

	// test set is a subset of test plan; we remember its name 
	TestPlan string

	// a system under test description
	*SysUnderTest

	// a setup action
	Setup *Action `xml:"Setup"`

	// a cleanup action
	Cleanup *Action `xml:"Cleanup"`

	// a list of test cases; in XML, this is a list of <TestCase> tags
	Cases []TestCase `xml:"TestCase"`
}

/*
 * TestSet.String
 */
func (ts *TestSet) String() string {
	s := fmt.Sprintf("TestSet: %q", ts.Name)
	s += fmt.Sprintf(" is owned by '%s' test plan.\n", ts.TestPlan)
	s += fmt.Sprintf("  Description:\n'%s'\n", ts.Description)
	if ts.Setup != nil {
		s += fmt.Sprintf("  Setup: %s", ts.Setup.String())
	} else {
		s += fmt.Sprintln("  Setup: []")
	}
	if ts.Cleanup != nil {
		s += fmt.Sprintf("  Cleanup: %s", ts.Cleanup.String())
	} else {
		s += fmt.Sprintln(" Cleanup: []")
	}
	for _, v := range ts.Cases {
		s += fmt.Sprintf("\n%s", v.String())
	}
	return s
}

/*
 * TestSet.Xml
 */
func (ts *TestSet) Xml() string {
	xml := fmt.Sprintf("<TestSet name=%q>\n", ts.Name)
	if ts.Setup != nil {
		xml += fmt.Sprintf("<Setup>\n%s</Setup>\n", ts.Setup.Xml())
	} else {
		xml += "<Setup />\n"
	}
	if ts.Cases != nil {
		for _, tc := range ts.Cases {
			xml += tc.Xml()
		}
		//xml += ts.Cases.Xml()
	} else {
		xml += "<TestCase />\n"
	}
	if ts.Cleanup != nil {
		xml += fmt.Sprintf("<Cleanup>\n%s</Cleanup>\n", ts.Cleanup.Xml())
	} else {
		xml += "<Cleanup />\n"
	}
	xml += fmt.Sprintln("</TestSet>")
	return xml
}

/*
 * TestSet.Json
 */
func (ts *TestSet) Json() (string, error) {
	b, err := json.Marshal(ts)
	if err != nil {
		return "", err
	}
	return string(b[:]), err
}

/*
 * TestSet.Html - HTML representation of the TestSet
 */
func (ts *TestSet) Html() (string, error) {
	// TODO
	return "", nil
}

/*
 * TestSet.findEmpty
 */
func (ts *TestSet) findEmpty() int {
	for ix, tc := range ts.Cases {
		if &tc != nil && tc.Name == "" {
			return ix
		}
	}
	return -1
}

/*
 * TestSet.AppendCase
 */
func (ts *TestSet) AppendCase(tc *TestCase) []TestCase {
	if tc.Name != "" {
		l := len(ts.Cases)
		c := cap(ts.Cases)
		if l+1 > c {
			newlst := make([]TestCase, 0, 2*c)
			copy(newlst, ts.Cases)
			ts.Cases = newlst
		}
		ts.Cases = ts.Cases[0 : l+1]
		ix := ts.findEmpty()
		if ix != -1 {
			ts.Cases[ix] = *tc
		}
	}
	return ts.Cases
}

/*
 * TestSet.ExtendCaseList
 */
func (ts *TestSet) ExtendCaseList(tcl []TestCase) []TestCase {
	l := len(ts.Cases)
	if l+len(tcl) > cap(ts.Cases) {
		newlst := make([]TestCase, 0, cap(ts.Cases)+len(tcl))
		copy(newlst, ts.Cases)
		ts.Cases = newlst
	}
	ts.Cases = ts.Cases[0 : l+len(tcl)]
	empty := ts.findEmpty()
	if empty != -1 {
		for ix, tc := range tcl {
			ts.Cases[empty+ix] = tc
		}
	}
	return ts.Cases
}

/*
 * TestSet.AppendConfig
 */
func (ts *TestSet) CleanupAfterTsetSetupFail() string {
	o := "Setup has FAILED\n"
	o += "Stopping the complete test set execution.\n"
	// mark all tcs & cases as skipped
	for _, tc := range ts.Cases {
		for _, step := range tc.Steps {
			step.Status.Set("NotTested")
		}
	}
	o += fmt.Sprintln("<<< Leaving test set %q", ts.Name)
	return o
}

/*
 * TestSet.Execute - executes the whole test set
 */
func (ts *TestSet) Execute(display *ExecDisplayFnCback) {
	output := ""
	// define function from function pointer
	_d := *display
	//
	_d("notice", fmt.Sprintf(">>> Entering Test Set %q\n", ts.Name))
	if ts.Setup != nil {
		output = ts.Setup.Execute()
		_d("notice", fmt.Sprintln("Executing setup script"))
		_d("info", FmtOutput(output))
		// if setup script has failed, there's no need to proceed...
		if ts.Setup.Status.Get() == "Fail" {
			_d("error", ts.CleanupAfterTsetSetupFail())
		}
	} else {
		_d("notice", fmt.Sprintln("Setup action is not defined."))
	}
	//
	if ts.Cases != nil {
		for _, tc := range ts.Cases {
			tc.Execute(display)
		}
	}
	//
	if ts.Cleanup != nil {
		_d("notice", fmt.Sprintln("Executing cleanup script"))
		_d("info", FmtOutput(ts.Cleanup.Execute()))
	} else {
		_d("notice", fmt.Sprintln("Cleanup action is not defined:"))
	}
	_d("notice", fmt.Sprintf("<<< Leaving test set %q\n", ts.Name))
}

/*
 * CreateTestSet - function that creates the TestSet struct
 */
const defCfgListCap = 10

func CreateTestSet(name string,
	descr string,
	sut *SysUnderTest,
	setup *Action,
	cleanup *Action) *TestSet {
	tcs := make([]TestCase, 0, defCfgListCap)
	return &TestSet{name, descr, "", sut, setup, cleanup, tcs}
}
