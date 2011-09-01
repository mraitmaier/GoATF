package atf

import (
//        "fmt"
        "testing"
        )


func testExeAction(t *testing.T) {
//    fmt.Println("Executing Action tests...")
    script := "action"
    args := "arg1 arg2"
    act := CreateAction(script, args)
    if act.Script != "action" { 
        t.Errorf("Script: %q != %q", script, act.Script) 
    }
    if act.Args != "arg1 arg2" { t.Errorf("Args: %q != %q", args, act.Args) }
    if !act.IsExecutable() {t.Errorf("%q must be executable", script) }
    if act.IsManual() {t.Errorf("%q must NOT be manual", script) }
/*
    fmt.Println(act.Xml())
    // Creating the second action
    act2 := CreateAction("another action", "arg1")
    fmt.Printf("The script is %q\n", act2.Script)
    fmt.Printf("The arguments are %q\n", act2.Args)
    fmt.Println("Is action executable?", act.IsExecutable())
    fmt.Println("Is action manual?", act.IsManual())
    fmt.Println(act2.Xml())
*/
}

func testEmptyAction(t *testing.T) {
    act := CreateEmptyAction()
    if act.IsExecutable() { t.Errorf("Empty cannot be executable") }
    if act.IsManual() { t.Errorf("Empty cannot be Manual") }
}


func testManualAction(t *testing.T) {
    desc := "This is a description of a manual action"
    act := CreateManualAction(desc)
    if !act.IsManual() { t.Errorf("Empty must be marked Manual") }
}
