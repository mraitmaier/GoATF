/*
 * testset.go 
 *
 * History:
 *  0.1   Apr10 MR Initial version, limited testing
 */

package atf

import ("fmt"
        "json"
        "os")

/*
 * TestSet - this is just an alias for TestPlan structure. They do differ in
 * one single detail: TestSet implements the Execute() method, while TestPlan
 * is not. Otherwise they're completely the same. 
 * In general, TestSet is a subset of TestPlan. 
 */
type TestSet struct {
    Name string
    Description string
    TestPlan string
    Setup *Action
    Cleanup *Action
    Configs []Configuration
}

/*
 * TestSet.String
 */
func (ts *TestSet) String() string {
    s := fmt.Sprintf("TestSet: %q", ts.Name)
    s += fmt.Sprintf(" is owned by '%s' test plan.\n", ts.TestPlan)
    s += fmt.Sprintf("  Description:\n'%s'\n", ts.Description)
    if ts.Setup != nil { 
        s += fmt.Sprintf("  Setup: %s", ts.Setup.String())
    } else {
        s += fmt.Sprintln("  Setup: []")
    }
    if ts.Cleanup != nil { 
        s += fmt.Sprintf("  Cleanup: %s", ts.Cleanup.String())
    } else {
        s += fmt.Sprintln(" Cleanup: []")
    }
    for _, v := range ts.Configs { s += fmt.Sprintf("\n%s", v.String()) }
    return s
}

/*
 * TestSet.Xml
 */
func (ts *TestSet) Xml() string {
    xml := fmt.Sprintf("<TestSet name=%q>\n", ts.Name)
    if ts.Setup != nil {
        xml += fmt.Sprintf("<Setup>\n%s</Setup>\n", ts.Setup.Xml())
    } else {
        xml += "<Setup />\n"
    }
    if ts.Configs != nil {
        for _, cfg := range(ts.Configs) {
            xml += cfg.Xml()
        }
        //xml += ts.Configs.Xml()
    } else {
        xml += "<Configuration />\n"
    }
    if ts.Cleanup != nil {
        xml += fmt.Sprintf("<Cleanup>\n%s</Cleanup>\n", ts.Cleanup.Xml())
    } else {
        xml += "<Cleanup />\n"
    }
    xml += fmt.Sprintln("</TestSet>")
    return xml
}

/*
 * TestSet.Json
 */
func (ts *TestSet) Json() (string, os.Error) {
    b, err := json.Marshal(ts)
    if err != nil {return "", err }
    return string(b[:]), err
}

/*
 * TestSet.Html - HTML representation of the TestSet
 */
func (ts *TestSet) Html() (string, os.Error) {
    // TODO
    return "", nil
}

/*
 * TestSet.findEmpty
 */
func (ts *TestSet) findEmpty() int {
    for ix, cfg := range(ts.Configs) {
        if &cfg != nil && cfg.Name == "" { return ix }
    }
    return -1
}

/*
 * TestSet.AppendConfig
 */
func (ts *TestSet) AppendConfig(cfg *Configuration) []Configuration {
    if cfg.Name != "" {
        l := len(ts.Configs)
        c := cap(ts.Configs)
        if l+1 > c {
            newlst := make([]Configuration, 0, 2*c)
            copy(newlst, ts.Configs)
            ts.Configs = newlst
        }
        ts.Configs = ts.Configs[0:l+1]
        ix := ts.findEmpty()
        if ix != -1 {
            ts.Configs[ix] = *cfg
        }
    }
    return ts.Configs
}

/*
 * TestSet.AppendConfig
 */
func (ts *TestSet) ExtendConfigList(cfgl []Configuration) []Configuration {
    l := len(ts.Configs)
    if l+len(cfgl) > cap(ts.Configs) {
        newlst := make([]Configuration, 0, cap(ts.Configs)+len(cfgl))
        copy(newlst, ts.Configs)
        ts.Configs = newlst
    }
    ts.Configs = ts.Configs[0:l+len(cfgl)]
    empty := ts.findEmpty()
    if empty != -1 {
        for ix, cfg := range(cfgl) {
            ts.Configs[empty+ix] = cfg
        }
    }
    return ts.Configs
}

/*
 * TestSet.AppendConfig
 */
func (ts *TestSet) CleanupAfterTsetSetupFail() string {
    o := "Setup has FAILED\n"
    o += "Stopping the complete test set execution.\n"
    // mark all cfgs & cases as skipped
    for _, cfg := range(ts.Configs) {
        for _, tc := range(cfg.Cases) {
            for _, step := range(tc.Steps) {
                step.Status = Skipped
            }
        }
    }
    o += fmt.Sprintln("<<< Leaving test set %q", ts.Name)
    return o
}

/*
 * TestSet.Execute
 */
func (ts *TestSet) Execute() (output string) {
    output = fmt.Sprintf(">>> Entering Test Set %q\n", ts.Name)
    if ts.Setup != nil {
        output += fmt.Sprintln("Executing setup script")
        output += ts.Setup.Execute()
        // if setup script has failed, there's no need to proceed...
        if !ts.Setup.Success {
            output += ts.CleanupAfterTsetSetupFail()
            return output
        }
    } else {
        output += fmt.Sprintln("Setup action is not defined.")
    }
    //
    if ts.Configs != nil {
        for _, cfg := range(ts.Configs) {
            output += cfg.Execute()
        }
    }
    //
    if ts.Cleanup != nil {
        output += fmt.Sprintln("Executing cleanup script")
        output += ts.Cleanup.Execute()
    } else {
        output += fmt.Sprintln("Cleanup action is not defined:")
    }
    output += fmt.Sprintf("<<< Leaving test set %q\n", ts.Name)
    return output
}


/*
 * CreateTestSet - function that creates the TestSet struct
 */
const defCfgListCap = 10
func CreateTestSet(name string,
                   descr string,
                   setup *Action,
                   cleanup *Action,) *TestSet {
    cfgs := make([]Configuration, 0, defCfgListCap)
    return &TestSet{name, descr, "", setup, cleanup, cfgs}
}


