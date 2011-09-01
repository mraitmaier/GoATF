/*
 * report.go
 *
 * History:
 *  0.1.0   Jul10   MR  The initial version
 */

package atf

import (
//    "fmt"
)


type Reporter interface {
    Report(TestReport) string
}

/*
 * ReportType - an enum containing configuration types
 */
type ReportType int
const (
	UnknownReport ReportType = iota
    HtmlReport
    XmlReport
)

