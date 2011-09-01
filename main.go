/*
 * main.go - the file to rule them all
 */
package main

import (
        "fmt"
        "flag"
        "os"
        "path"
        "path/filepath"
        "runtime"
        "atf"
)

const (
        numOfLoggers int = 3
      )
/*
 * Runner
 */
type Runner struct {
    tr *atf.TestReport // TestSet that's be run
    input string       // input configuration file (currently only JSON)
    workdir string
    logfile string
    syslog string
    report string
    cssfile string
    xml bool          // create XML report (beside HTML report)
    debug bool        // enable debug mode (for testing purposes only)
    logger *atf.Log   // a logger instance (
}

func NewRunner() *Runner {
    var r = new(Runner)
    r.logger = atf.NewLog(numOfLoggers)
    return r
}

/*
 * Runner.display - displays the contents of the Runner type
 * If complete flag is 'true', method will display the complete TestSet;
 * otherwise only name will be printed
 */
func (r *Runner) display(complete bool) {
    fmt.Printf("Input config file: %q\n", r.input)
    fmt.Printf("Working dir: %q\n", r.workdir)
    fmt.Printf("Log filename: %q\n", r.logfile)
    fmt.Printf("Syslog server IP: %q\n", r.syslog)
    fmt.Printf("Final report name: %q\n", r.report)
    fmt.Printf("(Optional) CCS file for HTML report: %q\n", r.cssfile)
    fmt.Printf("Debug node enabled? %t\n", r.debug)
    // display loggers
    fmt.Printf("Loggers (total # = %d):\n", r.logger.Len())
    if r.logger != nil {
        fmt.Println(r.logger.String())
    }
    // display test set
    if r.tr != nil {
        if complete {
            fmt.Println(r.tr.String())
        } else {
            fmt.Printf("TestSet: %q\n", r.tr.Name)
        }
    } else {
        fmt.Println("TestSet not defined yet.")
    }
}

func (r *Runner) setWorkDir(basedir string, tsName string) {
    var w = "results"
    if basedir == "" {
        if runtime.GOOS == "windows" {
            basedir = os.Getenv("USERPROFILE")
        } else {
            basedir = os.Getenv("HOME")
        }
        basedir = path.Join(basedir, w,
                                    fmt.Sprintf("%s_%s", tsName, atf.NowFile()))
    }
    r.workdir = filepath.ToSlash(basedir)
}

func (r *Runner) collect() os.Error {
    var ts *atf.TestSet = new(atf.TestSet)
    var err os.Error
    if r.input != "" {
        ts, err = atf.CollectTestSet(r.input)
        if err != nil { return err }
    } else {
        return os.NewError("There's no configuration file defined")
    }
    if ts == nil { return os.NewError("Test set is empty") }
    r.tr = atf.CreateTestReport(ts)
    return nil
}

// Let's define the default levels for different log handlers:
// all text goes only to file logger, console should take only the most
// important printous, while syslog handler should omit sending the execution
// outputs. 
const (
        defSyslogLevel atf.LogLevel = atf.NoticeLogLevel
        defFileLevel atf.LogLevel = atf.InfoLogLevel
        defStreamLevel atf.LogLevel = atf.NoticeLogLevel
      )

func (r *Runner) createLog() os.Error {
    logfile := ""
    // logfile input argument is NOT empty...
    if r.logfile != "" {
        // and represents absolute path, take it as it is.
        if path.IsAbs(r.logfile) {
            logfile = r.logfile
        } else {
        // if not absolute path, get working dir and join the path to filename
            logfile = path.Join(r.workdir, r.logfile)
        }
    } else {
        logfile = path.Join(r.workdir, "output.log")
    }
    r.logfile = logfile
    // now the real thing...
    format := "%s %s %s"
    err := r.createLoggers(format, r.debug)
    if err != nil { return err }
    // if logger is created, this message should print...
    r.logger.Warning("Log successfully created\n")
//    r.logger.Notice("Displaying Runner configuration:")
//    r.logger.Notice(r.display(false))
    return nil
}

func (r *Runner) createLoggers(fmt string, debug bool) os.Error {
    // first, we define log levels (severity) 
    fLevel := defFileLevel   // this is level for file handler
    sLevel := defSyslogLevel // this is level for syslog & console handlers
    if debug {
        fLevel = atf.DebugLogLevel
        sLevel = atf.DebugLogLevel
    }
    // now create file logger
    f, err := atf.NewFileHandler(r.logfile, fmt, fLevel)
    if err != nil { return err }
    if f != nil { r.logger.Handlers = r.logger.AddHandler(f) }
    // and create console logger
    l := atf.NewStreamHandler(fmt, sLevel)
    if l != nil { r.logger.Handlers = r.logger.AddHandler(l) }
    // and finally create syslog logger if needed
    if r.syslog != "" {
        var s *atf.SyslogHandler
        s = atf.NewSyslogHandler(r.syslog, fmt, sLevel)
        if s != nil { r.logger.Handlers = r.logger.AddHandler(s) }
    }
    return err
}

func (r *Runner) Init() os.Error {
    // let's collect the configuration
    err := r.collect()
    if err != nil { return err }
    // check working dir value; if empty, redefine to default: '$HOME/results'
    r.setWorkDir(r.workdir, r.tr.TestSet.Name)
    // if this dir is not existent, create it
    err = os.MkdirAll(r.workdir, 0755)
    if err != nil { return err }
    // create log file
    err = r.createLog()
    return err
}

func (r *Runner) fmtOutput(o string) string {
    s := "Displaying output:\n################### OUTPUT ##################\n"
    s += o
    s += "################ OUTPUT END #################\n"
    return s
}

/*
 *
 */
func (r *Runner) runStep(step *atf.TestStep) {
    output := ""
    if step == nil {
        r.logger.Error("Empty test step.\n")
        return
    }
    r.logger.Notice(">>>>>>>>> Starting action\n")
    output = step.Execute()
    r.logger.Notice(fmt.Sprintf("Action status: %t\n", step.Success))
    r.logger.Info(r.fmtOutput(output))
    r.logger.Notice(">>>>>>>>> Action end.\n")
}

/*
 *
 */
func (r *Runner) runTestcase(tc *atf.TestCase) {
    if tc == nil {
        r.logger.Error("Empty test case\n")
        return
    }
    r.logger.Notice(fmt.Sprintf("### Starting test case: %q\n", tc.Name))
    r.runSetup(tc.Setup)
    for _, step := range tc.Steps { r.runStep(&step) }
    r.runCleanup(tc.Cleanup)
    r.logger.Notice(fmt.Sprintf("### Test case: %q end.\n", tc.Name))
}

/*
 *
 */
func (r *Runner) runConfig(cfg *atf.Configuration) {
    if cfg == nil {
        r.logger.Error("Empty configuration\n")
        return
    }
    r.logger.Notice(fmt.Sprintf("### Starting configuration: %q\n", cfg.Name))
    r.runSetup(cfg.Setup)
    for _, tcase := range cfg.Cases {
        r.runTestcase(&tcase)
    }
    r.runCleanup(cfg.Cleanup)
    r.logger.Notice(fmt.Sprintf("### configuration: %q end.\n", cfg.Name))
}

/*
 *
 */
func (r *Runner) runSetup(act *atf.Action) {
    var output string = ""
    // run test set setup action (if it exists)
    if act != nil {
        r.logger.Notice(">>>>>>>>> Starting setup action\n")
        output = r.tr.TestSet.Setup.Execute()
    }
    r.logger.Notice(fmt.Sprintf("Setup action status: %t\n",r.tr.Setup.Success))
    r.logger.Info(r.fmtOutput(output))
}

/*
 *
 */
func (r *Runner) runCleanup(act *atf.Action) {
    var output string = ""
    // run test set cleanup action (if it exists)
    if act != nil {
        r.logger.Notice(">>>>>>>>> Starting cleanup action\n")
        output = r.tr.TestSet.Cleanup.Execute()
    }
    r.logger.Notice(fmt.Sprintf("Cleanup action status: %t\n",
                r.tr.Setup.Success))
    r.logger.Info(r.fmtOutput(output))
}

/*
 *
 */
func (r *Runner) Run() {
    r.logger.Notice(fmt.Sprintf("# Starting Test set: %q\n", r.tr.TestSet.Name))
    // run the test set setup action
    r.runSetup(r.tr.TestSet.Setup)
    // now execute the configurations
    for _, cfg := range r.tr.TestSet.Configs {
        r.runConfig(&cfg)
    }
    // run test set cleanup action (if it exists)
    r.runCleanup(r.tr.TestSet.Cleanup)
    r.logger.Notice(fmt.Sprintf("# Test set: %q end.\n", r.tr.TestSet.Name))
    // This is the end of execution
}

/*
 *
 */
func (r *Runner) createHtmlHeader(name string) string {
    s := "<!DOCTYPE html>\n"
    s += "<html>\n<head>\n"
    s += fmt.Sprintf("<meta charset=%q>\n", "utf-8")
    s += fmt.Sprintf("<title>Report: %s</title>\n", name)
    s += "</head>\n"
    return s
}

/*
 *
 */
func (r *Runner) createXmlReport(filename string) os.Error {
    xml := fmt.Sprintf("<?xml version=%q encoding=%q?>", "1.0", "UTF-9")
    xml += r.tr.Xml()
    fout, err := os.OpenFile(filename, os.O_CREATE | os.O_WRONLY, 0755)
    if err != nil { return err }
    defer fout.Close()
    fmt.Fprint(fout, xml)
    return nil
}

/*
 *
 */
func (r *Runner) createHtmlReport(filename string) os.Error {
    // HTML report is always created
    html := r.createHtmlHeader(r.tr.TestSet.Name)
    html += "<body>\n"
    h, err := r.tr.Html()
    if err != nil { return err }
    html += "</body>\n</html>\n"
    html += h
    // the file itself
    fout, err := os.OpenFile(filename, os.O_CREATE | os.O_WRONLY, 0755)
    if err != nil { return err }
    defer fout.Close()
    fmt.Fprint(fout, html)
    return nil
}

/*
 *
 */
func (r *Runner) CreateReports() {
    // always create HTML report
    filename := filepath.ToSlash(path.Join(r.workdir, "report.html"))
    err := r.createHtmlReport(filename)
    if err != nil {
        r.logger.Error("XML report could not be created.\n")
        r.logger.Error(fmt.Sprintf("Reason: %s\n", err))
        return
    }
    r.logger.Notice(fmt.Sprintf("HTML report %q created.\n", filename))
    // create XML report, if needed
    if r.xml {
        filename = filepath.ToSlash(path.Join(r.workdir, "report.xml"))
        err := r.createXmlReport(filename)
        if err != nil {
            r.logger.Error("XML report could not be created.\n")
            r.logger.Error(fmt.Sprintf("Reason: %s\n", err))
            return
        }
        r.logger.Notice(fmt.Sprintf("XML report %q created.\n", filename))
    }
}

/************************************************
 * parseArgs - parge command-line arguments
 */
func parseArgs(r *Runner) {
    flag.StringVar(&r.input, "i", "", "Input configuration path")
    flag.StringVar(&r.workdir, "w", "", "Working directory path")
    flag.StringVar(&r.logfile, "l", "", "Logfile name")
    flag.StringVar(&r.syslog, "s", "", "Syslog server IP")
    flag.StringVar(&r.report, "r", "", "final report filename")
    flag.StringVar(&r.cssfile, "c", "", "custom CSS file for HTML report")
    flag.BoolVar(&r.xml, "X", false, "create XML report (beside HTML report)")
    flag.BoolVar(&r.debug, "d", false,
            "enable debug mode (for testing purposes)")
    //
    flag.Parse()
}

/*
 *
 */
func main() {
    //atf.RunBats() // for testing purposes : test/bats.go
    r := NewRunner()
    // parse CLI arguments
    parseArgs(r)
    // initialize new Runner 
    err := r.Init()
    if err != nil { panic(err) }
    r.display(false) // DEBUG
    // now, run the damn thing....
    r.Run()
    //
    r.CreateReports()
    // close the logger
    r.logger.Close()
}

