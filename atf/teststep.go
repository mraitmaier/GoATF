/*
 * teststep.go - implementation of the TestStep type 
 *
 * This data structure represents the single test case (executable) step (or
 * action). It is always expected for a step to pass, so the self-evaluation is
 * as simple as possible.
 *
 * History:
 *  0.1   Apr10 MR Initial version, limited testing
 */

package atf

import (
	"encoding/json"
	"fmt"
)

/*
 * TestStep
 */
type TestStep struct {

	/* name of the test step; in XML, this is an attribute */
	Name string

	/* expected status of the step; in XML, this is an attribute */
	Expected TestResult `xml:"expected,attr"`

	/* status of the step; in XML, this is an attribute */
	Status TestResult `xml:"status,attr"`

	/* every test step needs an action: either manual or executable */
	*Action `xml:"Action"`
}

/*
 * TestStep.String 
 */
func (ts *TestStep) String() string {
	var act string
	// let's check the action first...
	if ts.Action != nil {
		act = ts.Action.String()
	} else {
		act = "none"
	}
	return fmt.Sprintf(
		"TestStep: %q expected status: %q status: %q action: %q\n",
		ts.Name, ts.Expected, ts.Status, act)
}

/*
 * TestStep.Display
 */
func (ts *TestStep) Display() string {
	txt := "TestStep\n\n"
	txt += fmt.Sprintf("Name: %q\n", ts.Name)
	txt += fmt.Sprintf("Expected status: %q\n", ts.Expected)
	txt += fmt.Sprintf("Status: %q", ts.Status)
	if ts.Action != nil {
		txt += fmt.Sprintf("Action: %q\n", ts.Action.String())
	} else {
		txt += "Action: N/A\n"
	}
	return txt
}

/*
 * TestStep.Xml 
 */
func (ts *TestStep) Xml() string {
	s := "<TestStep />\n"
	if ts.Action != nil {
		s = fmt.Sprintf("<TestStep expected=%q status=%q>",
			ts.Expected, ts.Status)
        s += fmt.Sprintf("<Name>%s</Name>\n", ts.Name)
		s += fmt.Sprintf("%s</TestStep>\n", ts.Action.Xml())
	}
	return s
}

/*
 * TestStep.Json
 */
func (ts *TestStep) Json() (string, error) {
	b, err := json.Marshal(ts)
	if err != nil {
		return "", err
	}
	return string(b[:]), err
}

/*
 * TestStep.Html
 */
func (ts *TestStep) Html() (string, error) {
	// TODO
	return "", nil
}

/*
 * TestStep.Execute
 */
func (ts *TestStep) Execute(display *ExecDisplayFnCback) {
	// we turn the function ptr back to function 
	_d := *display
	// and start the execution
	_d("notice", fmt.Sprintf(">>> Entering test step %q\n", ts.Name))
	// we execute the action when it's not empty
	if ts.Action != nil {
		_d("info", FmtOutput(ts.Action.Execute()))
	} else {
		_d("error", fmt.Sprintln("Action is EMPTY?????"))
	}
	// let's evaluate expectations and final status of the step
	switch ts.Expected {
	case "Pass":
		if ts.Action.Status == "Pass" {
			ts.Status.Set("Pass")
		} else {
			ts.Status.Set("Fail")
		}
	case "XFail":
		if ts.Action.Status == "Pass" {
			ts.Status.Set("Fail")
		} else {
			ts.Status.Set("Pass")
		}
	default:
		//only Pass & XFail are allowed as expected status 
		ts.Status.Set("NotTested")
	}
	_d("info", fmt.Sprintf("Test step evaluated to %q\n", ts.Status))
	_d("notice", fmt.Sprintf("<<< Leaving test step %q\n", ts.Name))
}

/*
 * CreateTestStep
 */
func CreateTestStep(name string, descr string, expected TestResult,
	status TestResult, act *Action) *TestStep {
	return &TestStep{name, expected, status, act}
}
