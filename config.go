/*
 * config.go 
 *
 * History:
 *  0.1   Apr10 MR Initial version, limited testing
 */

package atf

import ("fmt"
        "json"
        "os")

/*
 * Configuration -
 */
type Configuration struct {
    Name string
    Sut *SysUnderTest
    Setup *Action
    Cleanup *Action
    Cases []TestCase
    Description string
}

func (cfg *Configuration) String() string {
    s := fmt.Sprintf("Configuration: %q\n", cfg.Name)
    s += fmt.Sprintf("  Description:\n'%s'\n", cfg.Description)
    s += fmt.Sprintf("%s\n", cfg.Sut.String())
    if cfg.Setup != nil {
        s += fmt.Sprintf("  Setup: %s", cfg.Setup.String())
    } else {
        s += fmt.Sprintln("  Setup: []")
    }
    if cfg.Cleanup != nil {
        s += fmt.Sprintf("  Cleanup: %s", cfg.Cleanup.String())
    } else {
        s += fmt.Sprintln("  Cleanup: []")
    }
    for _, v := range cfg.Cases { s += fmt.Sprintf("\n%s", v.String()) }
    return s
}

func (cfg *Configuration) Xml() string {
    xml := fmt.Sprintf("<Configuration name=%q>\n", cfg.Name )
    xml += fmt.Sprintf("<Description>%s</Description>\n", cfg.Description)
    xml += cfg.Sut.Xml()
    xml += cfg.Setup.Xml()
    if cfg.Cases != nil {
        for _, tc := range(cfg.Cases) {
            xml +=  tc.Xml()
        }
    } else {
        xml += "<Testcase />\n"
    }
    xml += cfg.Cleanup.Xml()
    xml += fmt.Sprintln("</Configuration>")
    return xml
}

/*
 * Configuration.Json -
 */
func (cfg *Configuration) Json() (s string, err os.Error) {
    var b []byte
    b, err = json.Marshal(cfg)
    s = ""
    if err == nil { s = string(b[:]) }
    return
}

/*
 * Configuration.findEmpty -
 */
func (cfg *Configuration) findEmpty() int {
    for ix, tcase := range cfg.Cases {
        if &tcase != nil && tcase.Name == "" { return ix }
    }
    return -1
}

/*
 * Configuration.AppendCase
 */
func (cfg *Configuration) AppendCase(tc *TestCase) []TestCase {
    // we append new test case only if it's valid: its name mustn't be empty
    if tc.Name != "" {
        l := len(cfg.Cases)
        c := cap(cfg.Cases)
        // we check the capacity of case list; if needed, double it 
        if l+1 > c {
            newlst := make([]TestCase, 0, 2*c)
            copy(newlst, cfg.Cases)
            cfg.Cases = newlst
        }
        // slice should be incremented
        cfg.Cases = cfg.Cases[0:l+1]
        // we need an index of first available slot
        // and do the insertion if we get the valid index
        ix := cfg.findEmpty()
        if ix != -1 {
            cfg.Cases[ix] = *tc
        }
    }
    return cfg.Cases
}

/*
 *
 */
func (cfg *Configuration) ExtendCaseList(tcl []TestCase) []TestCase{
    // first we check the capacity of the case list
    l := len(cfg.Cases)
    if l+len(tcl) > cap(cfg.Cases) {
        newlst := make([]TestCase, 0, cap(cfg.Cases)+len(tcl))
        copy(newlst, cfg.Cases)
        cfg.Cases = newlst
    }
    cfg.Cases = cfg.Cases[0:l+len(tcl)]
    empty := cfg.findEmpty()
    if empty != -1 {
        for ix, tc := range tcl {
            cfg.Cases[empty+ix] = tc
        }
    }
    return cfg.Cases
}

func (cfg *Configuration) cleanupAfterSetupFail() string {
    o := "Setup action has FAILED.\n"
    o += "Skipping the rest of the Configuration.\n"
    o += fmt.Sprintf("<<< Leaving configuration %q\n", cfg.Name)
    // mark all steps of all cases as skipped 
    for _, tc := range(cfg.Cases) {
        for _, step := range(tc.Steps) {
            step.Status = Skipped
        }
    }
    return o
}

func (cfg *Configuration) Execute() (output string) {
    output = fmt.Sprintf(">>> Entering Configuration %q\n", cfg.Name)
    // let's execute the setup action (if not empty)
    if cfg.Setup != nil {
        output += fmt.Sprintln("Executing setup script")
        output += cfg.Setup.Execute()
        // if setup action has failed, skip the rest of the configuration
        if !cfg.Setup.Success {
            output += cfg.cleanupAfterSetupFail()
            return output
        }
    } else {
        output += fmt.Sprintln("Setup action is not defined.")
    }
    // let's execute the test cases
    if cfg.Cases != nil {
        for _, tc := range(cfg.Cases) {
            output += tc.Execute()
        }
    }
    // let's execute the cleanup action (if not empty)
    if cfg.Cleanup != nil {
        output += fmt.Sprintln("Executing cleanup action")
        output += cfg.Cleanup.Execute()
    } else {
        output += fmt.Sprintln("Cleanup action is not defined.")
    }
    output += fmt.Sprintf("<<< Leaving Configuration %q\n", cfg.Name)
    return output
}

/*
 *
 */
const defCaseListCap = 10
func CreateConfiguration(name string,
                         sut *SysUnderTest,
                         setup *Action,
                         cleanup *Action,
                         descr string) *Configuration {
    cases := make([]TestCase, 0, defCaseListCap)
    return &Configuration{name, sut, setup, cleanup, cases, descr}
}


