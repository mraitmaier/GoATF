/*
 * testcase.go  - implementation of the TestCase type
 *
 * This type represents the test case and is the central data struct of the
 * complete application. TestCase is built from separate test steps (that are
 * self-evaluated: pass/fail) , including setup and cleanup actions, and the 
 * TestCase itself uses the evaluation algorithm to self evaluate (pass/fail) 
 * itself according to expected result.
 *
 * History:
 *  0.1   Apr10 MR Initial version, limited testing
 *  0.2   Mar12 MR heavy refactoring: changed the Execute() method to work with
 *                 registered closure; xml.Unmarshal() parsing definitions; 
 *  0.3   Mar12 MR case evaluation fixed
 */

package atf

import (
	"encoding/json"
	"fmt"
)

/*
 * TestCase -
 */
type TestCase struct {

	// a name of the test case; in XML, this is an attribute
	Name string `xml:"name,attr"`

	// a test case setup action
	Setup *Action `xml:"Setup"`

	// a test case cleanup action
	Cleanup *Action `xml:"Cleanup"`

	// expected result for this test case: either pass or expected fail;
	// in XML, this is an attribute
	Expected TestResult `xml:"expected,attr"`

	// actual result for this test case after execution;
	// in XML, this is an attribute
	Status TestResult `xml:"status,attr"`

	// a list of test steps; in XML, this is a sequence of <TestStep> tags
	Steps []TestStep `xml:"TestStep"`

	// a detailed description of the test case
	Description string
}

/*
 * TestCase.String 
 */
func (tc *TestCase) String() string {
	s := fmt.Sprintf("Test Case: %q\n\tstatus: %s \n", tc.Name, tc.Status)
	s += fmt.Sprintf("\tDescription: %q\n", tc.Description)
	s += fmt.Sprintf("\tExpected: %s \n", tc.Expected)
	if tc.Setup != nil {
		s += fmt.Sprintf("\tSetup: %s", tc.Setup.String())
	} else {
		s += fmt.Sprintf("\tSetup: none")
	}
	if tc.Cleanup != nil {
		s += fmt.Sprintf("\tCleanup: %s\n", tc.Cleanup.String())
	} else {
		s += fmt.Sprintf("\tCleanup: none\n")
	}
	if tc.Steps != nil {
		for _, step := range tc.Steps {
			s += fmt.Sprintf("%s\n", step.String())
		}
	} else {
		s += fmt.Sprintln("\tActions: empty\n")
	}
	return s
}

/*
 * TestCase.Xml -
 */
func (tc *TestCase) Xml() string {
	s := "<TestCase>\n"
	s += fmt.Sprintf("<Description>%s</Description>\n", tc.Description)
	if tc.Setup != nil {
		s += fmt.Sprintf("<Setup>%s</Setup>\n", tc.Setup.Xml())
	} else {
		s += "<Setup />\n"
	}
	if tc.Steps != nil {
		for _, step := range tc.Steps {
			s += step.Xml()
		}
	} else {
		s += "<Step />\n"
	}
	if tc.Cleanup != nil {
		s += fmt.Sprintf("<Cleanup>%s</Cleanup>\n", tc.Cleanup.Xml())
	} else {
		s += "<Cleanup />\n"
	}
	s += "</TestCase>\n"
	return s
}

/*
 * Testcase.Json -
 */
func (tc *TestCase) Json() (string, error) {
	b, err := json.Marshal(tc)
	if err != nil {
		return "", err
	}
	return string(b[:]), err
}

/*
 * Testcase.Html -
 */
func (tc *TestCase) Html() (string, error) {
	// TODO
	return "", nil
}

/*
 * TestCase.findEmpty -
 */
func (tc *TestCase) findEmpty() int {
	for ix, step := range tc.Steps {
		if step.Name == "" {
			return ix
		}
	}
	return -1
}

/*
 * AppendStep -
 */
func (tc *TestCase) AppendStep(ts *TestStep) []TestStep {
	if ts.Name != "" {
		// first we check length of the steps' list
		// if needed, we double the capacity
		l := len(tc.Steps)
		c := cap(tc.Steps)
		if l+1 > c {
			newlst := make([]TestStep, 2*c)
			copy(newlst, tc.Steps)
			tc.Steps = newlst
		} // if l+1 > c
		tc.Steps = tc.Steps[0 : l+1]
		ix := tc.findEmpty()
		// we do the insetion only if valid index was found
		if ix != -1 {
			tc.Steps[ix] = *ts
		}
	} else {
		// WHAT TODO
	}
	return tc.Steps
}

/*
 * ExtendStepList -
 */
func (tc *TestCase) ExtendStepList(lst []TestStep) []TestStep {
	// first we check the capacity of the step list; if needed, double it
	l := len(tc.Steps)
	if l+len(lst) > cap(tc.Steps) {
		newlst := make([]TestStep, cap(tc.Steps)+len(lst))
		copy(newlst, lst)
		tc.Steps = newlst
	}
	tc.Steps = tc.Steps[0 : l+len(lst)]
	empty := tc.findEmpty()
	if empty != -1 {
		for ix, ts := range tc.Steps {
			tc.Steps[empty+ix] = ts
		}
	}
	return tc.Steps
}

func (tc *TestCase) cleanupAfterCaseSetupFail() string {
	output := "Setup action has FAILED.\n"
	output += "Skipping the rest of the case...\n"
	output += fmt.Sprintf("<<< Leaving TestCase %q\n", tc.Name)
	tc.Status.Set("Fail")
	// set all steps' status to NotTested
	for _, step := range tc.Steps {
		step.Status.Set("NotTested")
	}
	return output
}

func (tc *TestCase) Execute(display *ExecDisplayFnCback) {
	// we turn function ptr back to function
	_d := *display
	// and start with execution...
	_d("notice", fmt.Sprintf(">>> Entering TestCase %q\n", tc.Name))
	// let's execute setup action (if not empty)
	if tc.Setup != nil {
		_d("notice", fmt.Sprintln("Executing case setup action"))
		_d("info", FmtOutput(tc.Setup.Execute()))
		// if setup action has failed, skip the rest of the case
		if tc.Setup.Status.Get() == "Fail" {
			_d("error", tc.cleanupAfterCaseSetupFail())
		}
	} else {
		_d("notice", fmt.Sprintln("Setup action is not defined.\n"))
	}
	// now we execute the steps...
	if tc.Steps != nil {
		for _, step := range tc.Steps {
			step.Execute(display)
		}
	}
	// let's execute cleanup action (if not empty)
	if tc.Cleanup != nil {
		_d("notice", fmt.Sprintln("Executing case cleanup action"))
		_d("info", FmtOutput(tc.Setup.Execute()))
	} else {
		_d("notice", fmt.Sprintln("Cleanup action is not defined.\n"))
	}
	// now we evaluate the complete test case
	tc.evaluate()
	_d("notice", fmt.Sprintf("Test case evaluated to %q\n",
		tc.Status))
	_d("notice", fmt.Sprintf("<<< Leaving TestCase %q\n", tc.Name))
}

/*
 * TestCase.evaluate - private method evaluating the test case status
 */
func (tc *TestCase) evaluate() {
	tc.Status = TestResult{"Pass"} // initial values is Pass
	// if setup or cleanup have not Pass-ed, complete test case fails also 
	//if tc.Setup.Status != Pass || tc.Cleanup.Status != Pass {
	if tc.Setup.Status.Get() == "Fail" || tc.Cleanup.Status.Get() == "Fail" {
		tc.Status.Set("Fail")
		fmt.Println("DEBUG: setup or cleanup is not Pass") // DEBUG
		return
	}
	// otherwise compare steps' expected and final results
	for _, step := range tc.Steps {
		switch tc.Expected.Get() {
		case "Pass":
			if step.Status.Get() != "Pass" {
				tc.Status.Set("Fail")
				break
			}
		case "XFail":
			if step.Status.Get() != "Fail" {
				tc.Status.Set("Fail")
				break
			}
		default:
			// by definition, only PASS & XFAIL are allowed as expected results 
			tc.Status.Set("NotTested")
		} /* switch */
	} /* for */
}

/*
 * CreateTestCase -
 */
const defStepListCap = 10 /* default step list capacity */
func CreateTestCase(name string, setup *Action, cleanup *Action,
	expected TestResult, status TestResult, descr string) *TestCase {
	steps := make([]TestStep, 0, defStepListCap)
	return &TestCase{name, setup, cleanup, expected, status, steps, descr}
}
