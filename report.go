/*
 * report.go
 *
 * History:
 *  0.1.0   Jul10   MR  The initial version
 */

package atf

import (
	"os"
	"path"
	"path/filepath"
	//    "fmt"
)

type Reporter interface {
	Create(tr *TestReport) (string, os.Error)
}

/*
 * Report - a report structure to rule them all...
 *
 * It wraps all types of reports that ATF is aware of and defines the
 * operations on all of those reports.
 */
type Report struct {
	reports map[string]string
}

/****************************************************************************
 * CreateReport - create empty report structure 
 */
func CreateReport() *Report {
	var rpt = make(map[string]string)
	return &Report{rpt}
}

/****************************************************************************
 * Report.AddHtml - add a reference to HTML report 
 */
func (r *Report) AddHtml() { r.reports["html"] = "" }

/****************************************************************************
 * Report.AddXml - add a reference to XML report 
 */
func (r *Report) AddXml() { r.reports["xml"] = "" }

/****************************************************************************
 * Report.AddText - add a reference to text report 
 */
func (r *Report) AddJson() { r.reports["json"] = "" }

/****************************************************************************
 * Report.AddJson - add a reference to JSON report 
 */
func (r *Report) AddText() { r.reports["txt"] = "" }

/****************************************************************************
 * Report.create - private method that creates the report with given type 
 */
func (r *Report) create(tr *TestReport, typ string) (rpt string, err os.Error) {
	switch typ {
	case "html":
		rpt, err = tr.Html()
	case "xml":
		rpt = tr.Xml()
	case "txt": // TODO: TextReport not implemented yet
	case "json":
		rpt, err = tr.Json()
	default:
		rpt = "Unknown report type"
		err = os.EINVAL
	}
	return
}

/****************************************************************************
 * Report.Create - create all the defined reports and write them
 */
func (r *Report) Create(tr *TestReport, pth string) (err os.Error) {
	// if path is empty, create the default path
	if pth == "" {
		pth = "."
	}
	// iterate through existing report (types), create them and write them as
	// "report.<type>" into given path
	for i, contents := range r.reports {
		contents, err = r.create(tr, i)
		if err != nil {
			return err
		}
		filename := filepath.ToSlash(path.Join(pth, "report."+i))
		err = WriteTextFile(filename, contents)
		if err != nil {
			return err
		}
	}
	return
}
