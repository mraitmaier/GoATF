/*
 * testrpt.go 
 *
 * History:
 *  0.1   jun11 MR Initial version, limited testing
 */

package atf

import ("fmt"
        "json"
        "time"
        "os")


/*
 * TestReport - this is 
 */
type TestReport struct {
    *TestSet
    Started *time.Time
    Finished *time.Time
    Output string
}

/*
 * TestReport.String - function that 
 */
func (tr *TestReport) String() string {
    return fmt.Sprintf("TestReport %q\n\tstarted: %s\n\tfinished: %s\n",
            tr.TestSet.String(), tr.Started.String(), tr.Finished.String())
}

/*
 * TestReport.Xml - function that 
 */
func (tr *TestReport) Xml() string {
    xml := fmt.Sprintf("<TestReport>\n")
    xml += fmt.Sprintf("  <Started>%s</Started>\n", tr.Started.String())
    xml += fmt.Sprintf("  <Finished>%s</Finished>\n", tr.Finished.String())
    xml += fmt.Sprintf(tr.TestSet.Xml())
    xml += fmt.Sprintln("</TestReport>")
    return xml
}

/*
 * TestReport.Json - function that 
 */
func (tr *TestReport) Json() (string, os.Error) {
    b, err := json.Marshal(tr)
    if err != nil {return "", err }
    return string(b[:]), err
}

/*
 * TestReport.Json - function that 
 */
func (tr *TestReport) Html() (string, os.Error) {
    return "", nil
}

/*
 * CreateTestReport - function that creates the TestSet struct
 */
func CreateTestReport(ts *TestSet) *TestReport {
    var started *time.Time
    var finished *time.Time
    return &TestReport{ts, started, finished, ""}
}


