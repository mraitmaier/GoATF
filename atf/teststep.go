/*
 * structure.go 
 *
 * History:
 *  0.1   Apr10 MR Initial version, limited testing
 */

package atf

import (
	"fmt"
	"os"
	"json"
)

/*
 * TestStep
 */
type TestStep struct {
	Name     string            `xml:"attr"` /* name of the test step */
	Expected TestResult        `xml:"attr"` /* expected status of the step */
	Status   TestResult        `xml:"attr"` /* status of the step */
	*Action  `xml:"TestStep>"` /* every test step needs an action: 
	   either manual or executable */
}
/*
 * TestStep.String 
 */
func (ts *TestStep) String() string {
	s := fmt.Sprintf("\tTestStep: %q\n\texpected status: %s\n\tstatus: %s\n",
		ts.Name, ts.Expected.String(), ts.Status.String())
	if ts.Action != nil {
		s += fmt.Sprintf("\tAction: %s", ts.Action.String())
	} else {
		s += "\tAction: none"
	}
	return s
}

/*
 * TestStep.Display
 */
func (ts *TestStep) Display() string {
	txt := "TestStep\n\n"
	txt += fmt.Sprintf("Name: %s\n", ts.Name)
	txt += fmt.Sprintf("Expected status: %s\n", ts.Expected.String())
	txt += fmt.Sprintf("Status: %s", ts.Status.String())
	if ts.Action != nil {
		txt += fmt.Sprintf("Action: %s\n", ts.Action.String())
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
		s = fmt.Sprintf("<TestStep name=%q expected=%q status=%q>",
			ts.Name, ts.Expected.String(), ts.Status.String())
		s += fmt.Sprintf("%s</TestStep>\n", ts.Action.Xml())
	}
	return s
}

/*
 * TestStep.Json
 */
func (ts *TestStep) Json() (string, os.Error) {
	b, err := json.Marshal(ts)
	if err != nil {
		return "", err
	}
	return string(b[:]), err
}

/*
 * TestStep.Html
 */
func (ts *TestStep) Html() (string, os.Error) {
	// TODO
	return "", nil
}

/*
 * TestStep.Execute
 */
func (ts *TestStep) Execute() (output string) {
	//    var rc TestResult = Fail // return code of the executed script/executable
	output = fmt.Sprintf(">>> Entering test step %q\n", ts.Name)
	// we execute the action when it's not empty
	if ts.Action != nil {
		output += ts.Action.Execute()
	} else {
		output += fmt.Sprintln("Action is EMPTY?????")
	}
	// let's evaluate expectations and final status of the step
	switch ts.Expected {
	case Pass:
		if ts.Action.Success {
			ts.Status = Pass
		} else {
			ts.Status = Fail
		}
	case XFail:
		if ts.Action.Success {
			ts.Status = Fail
		} else {
			ts.Status = Pass
		}
	default:
		ts.Status = NotTested
	}
	output += fmt.Sprintf("Test step evaluated to %q\n", ts.Status.String())
	output += fmt.Sprintf("<<< Leaving test step %q\n", ts.Name)
	return output
}

/*
 * CreateTestStep
 */
func CreateTestStep(name string, descr string, expected TestResult,
status TestResult, act *Action) *TestStep {
	return &TestStep{name, expected, status, act}
}
