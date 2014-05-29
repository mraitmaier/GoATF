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
	"encoding/xml"
	"fmt"
)

/*
 * TestStep
 */
type TestStep struct {

	/* name of the test step; in XML, this is an attribute */
	Name string         `xml:"name,attr"`

	/* expected status of the step; in XML, this is an attribute */
	Expected TestResult `xml:"expected,attr"`

	/* status of the step; in XML, this is an attribute */
	Status TestResult   `xml:"status,attr"`

	/* every test step needs an action: either manual or executable */
	Action *Action      `xml:"Action"`
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
		"TestStep: %q expected: %q status: %q action: %q\n",
		ts.Name, ts.Expected, ts.Status, act)
}

/*
 * TestStep.Display
 */
func (ts *TestStep) Display() string {
	txt := fmt.Sprintf("TestStep: %q\n", ts.Name)
	txt += fmt.Sprintf("Expected status: %q\n", ts.Expected)
	txt += fmt.Sprintf("Status: %q\n", ts.Status)
	if ts.Action != nil {
		txt += fmt.Sprintf("Action: %q\n", ts.Action.String())
	} else {
		txt += "Action: N/A\n"
	}
	return txt
}

func (ts *TestStep) Xml() (string, error) {

    output, err := xml.MarshalIndent(ts, "", "  ")
    if err != nil {
        return "", nil
    }

    return string(output), nil
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

//  Tidy up action: check flags; if both flags are false, just clear the action.
func (ts *TestStep) Normalize() {
    ts.Action.UpdateFlags()
    if !ts.Action.IsManual() && !ts.Action.IsExecutable() {
        ts.Action = nil
    } else {
        ts.Action.Result = "NotTested"
        // if expected status is empty for executable action, force "Pass"
        if ts.Action.IsExecutable() && ts.Expected == "" {
            ts.Expected = "Pass"
        }
    }
}

/*
 * TestStep.Execute
 */
func (ts *TestStep) Execute(display *ExecDisplayFnCback) {

	// we turn the function ptr back to function 
	disp := *display

	// and start the execution
	disp("info", fmt.Sprintf(">>> Entering test step %q\n", ts.Name))

	// we execute the action when it's not empty
	if ts.Action != nil {
	    disp("notice", fmt.Sprintf("Executing test step action: %q\n",
                ts.Action.String()))
		disp("info", FmtOutput(ts.Action.Execute()))
	} else {
		disp("error", fmt.Sprintln("Action is EMPTY?????"))
	}

	// let's evaluate expectations and final status of the step
	switch ts.Expected {
	case "Pass":
		if ts.Action.Result == "Pass" {
			ts.Status = "Pass"
		} else {
			ts.Status = "Fail"
		}
	case "XFail":
		if ts.Action.Result == "Pass" {
			ts.Status = "Fail"
		} else {
			ts.Status = "Pass"
		}
	default:
		//only Pass & XFail are allowed as expected status 
		ts.Status = "NotTested"
	}
	disp("notice", fmt.Sprintf("Test step evaluated to %q\n", ts.Status))
    disp("info", fmt.Sprintf("<<< Leaving test step %q\n", ts.Name))
}

/*
 * CreateTestStep
 */
func CreateTestStep(name string, descr string, expected TestResult,
	status TestResult, act *Action) *TestStep {
	return &TestStep{name, expected, status, act}
}
