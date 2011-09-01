/*
 * testplan.go 
 *
 * History:
 *  0.1   May11 MR Initial version, limited testing
 */

package atf

import ("fmt"
        "os"
        "json")

/*
 * TestPlan -
 */
type TestPlan struct {
    Name string
    Description string
    Setup *Action
    Cleanup *Action
    Configs []Configuration
}

/*
 * TestPlan.String - function that creates the TestPlan struct
 */
func (tp *TestPlan) String() string {
    s := fmt.Sprintf("TestPlan: %q\n", tp.Name)
    return s
}

/*
 * TestPlan.Xml - function that creates the TestPlan struct
 */
func (tp *TestPlan) Xml() string {
    xml := fmt.Sprintf("<TestPlan name=%q\n", tp.Name)
    if tp.Setup != nil {
        xml += tp.Setup.Xml()
    } else {
        xml += "<Setup />\n"
    }
    if tp.Configs != nil {

        for _, cfg := range(tp.Configs) {
            xml += cfg.Xml()
        }
        //xml += tp.Configs.Xml()
    } else {
        xml += "<Configuration />\n"
    }
    if tp.Cleanup != nil {
        xml += tp.Cleanup.Xml()
    } else {
        xml += "<Cleanup />\n"
    }
    xml += fmt.Sprintln("</TestPlan>")
    return xml
}

/*
 * TestPlan.Json - function that 
 */
func (tp *TestPlan) Json() (string, os.Error) {
    b, err := json.Marshal(tp)
    if err != nil { return "", err }
    return string(b[:]), err
}

/*
 * CreateTestPlan - function that creates the TestPlan struct
 */
const defTpListCap = 10
func CreateTestPlan(name string,
                   descr string,
                   setup *Action,
                   cleanup *Action) *TestPlan {
    cfgs := make([]Configuration, 0, defTpListCap)
    return &TestPlan{name, descr, setup, cleanup, cfgs}
}

/*
 * TestPlan.findEmpty - function that creates the TestPlan struct
 */
func (tp *TestPlan) findEmpty() int {
    for ix, cfg := range(tp.Configs) {
        if &cfg != nil && cfg.Name == "" { return ix }
    }
    return -1
}

/*
 * TestPlan.AppendConfig - function that 
 */
func (tp *TestPlan) AppendConfig(cfg *Configuration) []Configuration {
    if cfg.Name != "" {
        l := len(tp.Configs)
        c := cap(tp.Configs)
        if l+1 > c {
            newlst := make([]Configuration, 0, 2*c)
            copy(newlst, tp.Configs)
            tp.Configs = newlst
        }
        tp.Configs = tp.Configs[0:l+1]
        ix := tp.findEmpty()
        if ix != -1 {
            tp.Configs[ix] = *cfg
        }
    }
    return tp.Configs
}

/*
 * TestPlan.ExtendConfigList - function that 
 */
func (tp *TestPlan) ExtendConfigList(cfgl []Configuration) []Configuration {
    l := len(tp.Configs)
    if l+len(cfgl) > cap(tp.Configs) {
        newlst := make([]Configuration, 0, cap(tp.Configs)+len(cfgl))
        copy(newlst, tp.Configs)
        tp.Configs = newlst
    }
    tp.Configs = tp.Configs[0:l+len(cfgl)]
    empty := tp.findEmpty()
    if empty != -1 {
        for ix, cfg := range(cfgl) {
            tp.Configs[empty+ix] = cfg
        }
    }
    return tp.Configs
}

