package atf

import (
	"encoding/json"
	"fmt"
	"strings"
)

/*****************************************************************************
 * actioner interface
 */
type Actioner interface {
	IsExecutable() bool /* is action executable? */
	IsManual() bool     /* is action Manual? */
}

/*****************************************************************************
 * Action - type representing a structure that holds two string
 * values: script and args. The former is the path to a script that will be
 * executed and the latter is the (optional) string holding the arguments to
 * the script. 
 */
type Action struct {

	/* script to be executed */
	Script string

	/* arguments to script (if needed) */
	Args string

	/* script execution success */
	Status TestResult `xml:"status,attr"`

	/* script execution output text */
	Output string

	/* description text, used mainly for manual actions */
	Description string

	/* is this action executable? */
	executable bool `xml:"executable,attr"`

	/* is this action manual? */
	manual bool `xml:"manual,attr"`
}

/*****************************************************************************
 * Action.String - a string represenation of the Action struct
 */
func (a *Action) String() string {
	if a.IsManual() {
		return fmt.Sprintf("Manual Action:\n%s", a.Description)
	} else {
		if a.IsExecutable() {
			s := fmt.Sprintf("%s %s\n", a.Script, a.Args)
			return s
		} // if isexecutable
	} // if ismanual
	return fmt.Sprint(a.Script, " ", a.Args)
}

/*****************************************************************************
 * Action.updateFlags - a private method that updates the two actions flags:
 *                      'executable' and 'manual'
 * This method is run every time the 'Execute()' method is invoked.
 */
func (a *Action) updateFlags() {
	// initialy, action is neither executable not manual
	a.executable = false
	a.manual = false
	// if the action script is defined, action is executable
	// we like executable actions, so we gave them precedence
	if a.Script != "" {
		a.executable = true
		a.manual = false
	} else {
		// otherwise, if only Description is defined, we have a manual action
		if a.Description != "" {
			a.executable = false
			a.manual = true
		}
	}
}

/*****************************************************************************
 * Action.Xml - is this action an executable (script, program) action?
 */
func (a *Action) IsExecutable() bool { return a.executable }

/*****************************************************************************
 * Action.Xml - is this action a manual action?
 */
func (a *Action) IsManual() bool { return a.manual }

/*****************************************************************************
 * Action.Xml - return an XML representation of the Action struct 
 */
func (a *Action) Xml() string {
	xml := ""
	if a.IsExecutable() {
		xml = fmt.Sprintf("<Script>%s</Script>\n", a.Script)
		xml += fmt.Sprintf("<Args>%s</Args>\n", a.Args)
		xml += fmt.Sprintf("<Output>%s</Output>\n", a.Output)
		xml += fmt.Sprintln("<Description />")
	} else {
		if a.IsManual() {
			xml = "<Script />\n<Args />\n<Output />\n"
			xml += fmt.Sprintf("<Description>%s</Description>\n", a.Description)
		} else {
			xml = "<Script />\n<Args />\n<Output />\n<Description />"
		}
	}
	return xml
}

/*****************************************************************************
 * Action.Json - return a JSON representation of the Action struct 
 */
func (a *Action) Json() (string, error) {
	b, err := json.Marshal(a) // marshal returns a []byte, not string!
	if err != nil {
		return "", err
	}
	return string(b[:]), err
}

/*****************************************************************************
 * Execute - execute the action
 *
 * This function executes the action. But. It is actually executed only if
 * 'executed' argument is set: consequently this means that a particular action
 * is an executable script or a program. If 'manual' flag is set, an action is
 * considered manual. If both arguments are reset, that action is considered 
 * an empty (do-nothing) action. 
 * If we deal with non-executable action, 'description' is simply copied to 
 * 'output' field. Also, 'success' has a meaning only if action is executed; 
 * if not, 'success' is always set.
 */
func (a *Action) Execute() (output string) {
	a.updateFlags() // let's update the actions flags...
	a.Status.Set("Pass") // we assume execution will be successful

	// We execute the action only if it's marked executable
	if a.IsExecutable() {
		var err error
		a.Output, err = Execute(a.Script, strings.Split(a.Args, " "))
		// if error has accured, script has failed
		if err != nil {
			a.Status.Set("Fail")
		}
	} else {
		// otherwise we just put description into output, success is already set
		a.Output = a.Description
	}
	return a.Output
}

/*****************************************************************************
 * CreateAction - create a normal scripted/automated action
 *
 * This is creation function for a executable action. The 'script' fields is
 * mandatory, the 'args' field can be empty string. Also, the 'executed' flag
 * must be set and the 'manual' flag reset. The 'success' flag is reset by
 * default. The 'description' field has no special meaning with executable 
 * actions.
 */
func CreateAction(script string, args string) *Action {
    return &Action{script, args, TestResult{"UnknownResult"}, 
        "", "", true, false}
}

/*****************************************************************************
 * CreateManualAction - create a manual action
 *
 * This is creation function for a manual action. The 'script' and 'args'
 * fields are left empty, only 'description' is needed.
 * The 'manual' flag is set and 'executable' flag is reset.
 * Since this action is not executable, 'success' is always set.
 */
func CreateManualAction(descr string) *Action {
	return &Action{"", "", TestResult{"Pass"}, "", descr, false, true}
}

/*****************************************************************************
 * CreateAction - create a normal scripted/automated action
 *
 * This is creation function for empty (do-nothing) action. All fields are set
 * apropriately: only flags are actually needed. The 'manual' and 'executable'
 * flags are reset, 'success' flag is set.
 */
func CreateEmptyAction() *Action {
	return &Action{"No action", "", TestResult{"Pass"}, "", 
        "No action", false, false}
}
