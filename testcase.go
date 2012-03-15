/*
 * structure.go 
 *
 * History:
 *  0.1   Apr10 MR Initial version, limited testing
 */

package atf

import (
	"fmt"
	"json"
	"os"
)

/*
 * TestCase -
 */
type TestCase struct {
	Name        string `xml:"attr"`
	Setup       *Action
	Cleanup     *Action
	Expected    TestResult `xml:"attr"`
	Status      TestResult `xml:"attr"`
	Steps       []TestStep
	Description string
}

/*
 * TestCase.String 
 */
func (tc *TestCase) String() string {
	s := fmt.Sprintf("Test Case: %q\n\tstatus: %s \n", tc.Name,
		tc.Status.String())
	s += fmt.Sprintf("\texpected: %s \n", tc.Expected)
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
		s += fmt.Sprintln("\tactions: empty\n")
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
func (tc *TestCase) Json() (string, os.Error) {
	b, err := json.Marshal(tc)
	if err != nil {
		return "", err
	}
	return string(b[:]), err
}

/*
 * Testcase.Html -
 */
func (tc *TestCase) Html() (string, os.Error) {
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
	tc.Status = Fail
	// set all steps' status to skipped
	for _, step := range tc.Steps {
		step.Status = Skipped
	}
	return output
}

func (tc *TestCase) Execute() (output string) {
	output = fmt.Sprintf(">>> Entering TestCase %q\n", tc.Name)
	// let's execute setup action (if not empty)
	if tc.Setup != nil {
		output += fmt.Sprintln("Executing setup action")
		output += tc.Setup.Execute()
		// if setup action has failed, skip the rest of the case
		if !tc.Setup.Success {
			output += tc.cleanupAfterCaseSetupFail()
			return output
		}
	} else {
		output += fmt.Sprintln("Setup action is not defined.")
	}
	// now we execute the steps...
	if tc.Steps != nil {
		for _, step := range tc.Steps {
			output += step.Execute()
		}
	}
	// let's execute cleanup action (if not empty)
	if tc.Cleanup != nil {
		output += fmt.Sprintln("Executing cleanup action")
		output += tc.Setup.Execute()
	} else {
		output += fmt.Sprintln("Cleanup action is not defined.")
	}
	// now we evaluate the complete test case
	tc.evaluate()
	output += fmt.Sprintf("Test case evaluated to %q\n", tc.Status.String())
	output += fmt.Sprintf("<<< Leaving TestCase %q\n", tc.Name)
	return output
}

func (tc *TestCase) checkSkipped() bool {
	// we assume that all steps are "skipped" by default
	status := true
	// we iterate through step list to check if all steps are "skipped"
	for _, step := range tc.Steps {
		if step.Status != Skipped {
			status = false
			break
		}
	}
	return status
}

func (tc *TestCase) evaluate() {
	// first we check is steps were skipped: mark case as "skipped", too 
	if tc.checkSkipped() {
		tc.Status = Skipped
	} else {
		// otherwise compare expected and final results
		switch tc.Expected {
		case Pass:
			tc.Status = Pass
			for _, step := range tc.Steps {
				if step.Status != Pass {
					tc.Status = Fail
					break
				}
			}
		case XFail:
			tc.Status = Pass
			for _, step := range tc.Steps {
				if step.Status != Fail {
					tc.Status = Fail
					break
				}
			}
		default:
			/* by definition, only PASS & EXPECTED_FAIL are allowed as
			   expected results */
			tc.Status = NotTested
		} /* switch */
	}
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
