/*
 * testrslt.go
 */

package atf

//import ("fmt")
/*
 *
 */
type TestResult int
const (
    UnknownResult TestResult = iota
    Pass
    Fail
    XFail
    Skipped
    NotTested
    NotAvailable
)

func ResolveResult(status string) TestResult {
    val := UnknownResult
    switch (status) {
        case "unknown":
            val = UnknownResult
        case "pass", "Pass", "PASS":
            val = Pass
        case "fail", "Fail", "FAIL":
            val = Fail
        case "xfail", "Xfail", "XFail", "XFAIL", "expected fail":
            val = XFail
        case "nottested", "not tested", "NotTested", "Not Tested", "NOT TESTED":
            val = NotTested
        case "skipped", "SKIPPED", "Skipped":
            val = Skipped
        case "n/a", "N/A", "NA", "notavailable", "not available",
                           "Not Available", "NOT AVAILABLE":
            val = NotAvailable
        default:
            val = UnknownResult
    }
    return val
}

func (tr TestResult) String() string {
    var s string
    switch tr {
        case UnknownResult: s = "unknown test result"
        case Pass: s = "pass"
        case Fail: s = "fail"
        case XFail: s = "expected fail"
        case NotTested: s = "not tested"
        case Skipped: s = "skipped"
        case NotAvailable: s = "not available"
    }
    return s
}

