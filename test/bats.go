/*
 *
 */
package atf

import ( "fmt"
         "os"
         "path"
)

func RunBats() {
    fmt.Println("#### start ####")
    testExecutor()
    testCollector()
    testStructure()
    testTime()
    fmt.Println("#### end ####")
}

func testTime() {
    n := Now()
    fmt.Printf("time: %s\n", n)
    fmt.Printf("file version of time: %s\n", FileConv(n))
    n = NowFile()
    fmt.Printf("another file version of time: %s\n", n)
}

func testCollector() {
    fmt.Println("#### collector test ####")
    testfile := "cfg/example_ts.xml"
    fmt.Println("Base: ", path.Base(testfile) )
    fmt.Println("Ext: ", path.Ext(testfile))
    fmt.Println("Clean path: ", path.Clean(testfile))
    fmt.Println("Is absolute path: ", path.IsAbs(testfile))
    fmt.Println(">>>>>>>> JSON >>>")
    var ts *TestSet
    var err os.Error
    ts, err = CollectTestSet("cfg/go_example_pretty.json")
    if ts == nil { fmt.Println("Empty Test set!") }
    if err != nil { panic(err) }
    fmt.Println(ts.String())
    fmt.Println(">>>>>>>> ReadLines() test >>>")
    lines, err := ReadLines(testfile)
    if err != nil { panic(err) }
    for _, val := range lines { fmt.Println(val) }
    fmt.Println(">>>>>>>> ReadTextFile() test >>>")
    data, err := ReadTextFile(testfile)
    if err != nil { panic(err) }
    fmt.Printf("$%s$\n", data)
}

func testExecutor()  {
    fmt.Println("#### EXEC ####")
    fmt.Println("excuting python script...")
    output, err := Execute("d:/test/test.py",
                            []string{})
    //output,_ := Execute("/home/miran/Code/go/hello/main", []string{""})
    fmt.Printf("### OUTPUT ###\n%s\n### END ###\n", output)
    //fmt.Println("Error code: " + err.String())
    fmt.Println("excuting perl script...")
    output,_ = Execute("d:/test/test.pl",
                           []string{})
    fmt.Printf("### OUTPUT ###\n%s\n### END ###\n", output)
    fmt.Println("excuting Tcl script...")
    output,_ = Execute("d:/test/test.tcl",
                           []string{})
    fmt.Printf("### OUTPUT ###\n%s\n### END ###\n", output)
    fmt.Println("excuting Expect script...")
    output, err = Execute("d:/test/test.exp", []string{})
    if err != nil { fmt.Println("expect script: ERROR", err) }
    fmt.Printf("### OUTPUT ###\n%s\n### END ###\n", output)
    fmt.Println("excuting native executable script...")
    output, err = Execute("d:/test/uname", []string{"-a"})
    if err != nil { fmt.Println("native script: ERROR", err) }
    fmt.Printf("### OUTPUT ###\n%s\n### END ###\n", output)
    fmt.Println("excuting native executable script...")
    output, err = Execute("d:/test/uname", []string{"--help"})
    if err != nil { fmt.Println("native script: ERROR", err) }
    fmt.Printf("### OUTPUT ###\n%s\n### END ###\n", output)
    fmt.Println("excuting native executable ...")
    output, err = Execute("d:/test/uname", []string{})
    if err != nil { fmt.Println("native script: ERROR", err) }
    fmt.Printf("### OUTPUT ###\n%s\n### END ###\n", output)
    fmt.Println("excuting java JAR...")
    output, err = Execute("d:/test/hello.jar", []string{})
    if err != nil { fmt.Println("java executable: ERROR", err) }
    fmt.Printf("### OUTPUT ###\n%s\n### END ###\n", output)
    fmt.Println("excuting ruby script...")
    output, err = Execute("d:/test/test.rb", []string{})
    if err != nil { fmt.Println("ruby script: ERROR", err) }
    fmt.Printf("### OUTPUT ###\n%s\n### END ###\n", output)
    fmt.Println("excuting groovy script...")
    output, err = Execute("d:/test/test.groovy", []string{})
    if err != nil { fmt.Println("groovy script: ERROR", err) }
    fmt.Printf("### OUTPUT ###\n%s\n### END ###\n", output)
 
    //
    fmt.Println("#### OS ####")
    fmt.Println("excuting python script...")
    os := os.Getenv("GOOS")
    fmt.Printf("GOOS=%q\n", os)

}

const ( step_descr = "This is step description."
        case_descr = "This is case description."
        cfg_descr  = "this is config description."
        set_descr  = "This is test set description."
        plan_descr = "This is test plan description."
      )

func testStructure() {
    fmt.Println()
    act1 := CreateAction("action1", "arg1")
    act2 := CreateAction("action2", "arg1 arg2")
    act3 := CreateAction("action3", "arg1 arg2 arg3")
    act4 := CreateAction("action4", "arg1 arg2 arg3 arg4")
    act5 := CreateAction("action5", "arg1 arg2 arg3 arg4 arg5")
    act6 := CreateAction("action6", "arg1 arg2 arg3 arg4")
    act7 := CreateAction("action7", "arg1 arg2 arg3")
    act8 := CreateAction("action8", "arg1 arg2")
    act9 := CreateAction("action9", "arg1")
    act0 := CreateAction("action0", "")
    blah, _ := act0.Json()
    fmt.Println(blah)
    fmt.Println("#### test structure test ####")
    step1 := CreateTestStep("step1", step_descr, Pass, Fail, act1)
    step2 := CreateTestStep("step2", step_descr, Pass, Fail, act2)
    step3 := CreateTestStep("step3", step_descr, Pass, Fail, act3)
    step4 := CreateTestStep("step4", step_descr, Pass, Fail, act4)
    step5 := CreateTestStep("step5", step_descr, Pass, Fail, act5)
    step6 := CreateTestStep("step6", step_descr, Pass, Fail, act6)
    step7 := CreateTestStep("step7", step_descr, Pass, Fail, act7)
    step8 := CreateTestStep("step8", step_descr, Pass, Fail, act8)
    step9 := CreateTestStep("step9", step_descr, Pass, Fail, act9)
    step0 := CreateTestStep("step0", step_descr, Pass, NotTested, act0)
    fmt.Println(">> displaying steps' data")
    fmt.Println(step1.String())
    fmt.Println(step1.Xml())
    fmt.Println(step1.Json())
    fmt.Println(step2.String())
    fmt.Println(step2.Xml())
    fmt.Println(step2.Json())
    fmt.Println(step3.String())
    fmt.Println(step3.Xml())
    fmt.Println(step3.Json())
    fmt.Println(step4.String())
    fmt.Println(step4.Xml())
    fmt.Println(step4.Json())
    fmt.Println(step5.String())
    fmt.Println(step5.Xml())
    fmt.Println(step5.Json())
    fmt.Println(step6.String())
    fmt.Println(step6.Xml())
    fmt.Println(step6.Json())
    fmt.Println(step7.String())
    fmt.Println(step7.Xml())
    fmt.Println(step7.Json())
    fmt.Println(step8.String())
    fmt.Println(step8.Xml())
    fmt.Println(step8.Json())
    fmt.Println(step9.String())
    fmt.Println(step9.String())
    fmt.Println(step9.Json())
    fmt.Println(step0.Xml())
    fmt.Println(step0.String())
    fmt.Println(step0.Json())
    // setups and cleanups
    fmt.Println(">> setups & cleanups")
    empty_setup := new(Action)
    empty_cleanup := new(Action)
    setup1 := act1
    setup2 := act2
    cleanup1 := act0
    cleanup2 := act9
    fmt.Println(">> test cases...")
    tcase1 := CreateTestCase("testcase1", setup1, cleanup1, 
            NotTested, NotTested, case_descr)
    tcase1.AppendStep(step1)
    tcase1.AppendStep(step2)
    tcase1.AppendStep(step3)
    fmt.Println(tcase1.String())
    fmt.Println(tcase1.Xml())
    fmt.Println(tcase1.Json())
    tcase2 := CreateTestCase("testcase2", setup2, cleanup2, 
            Pass, NotTested, case_descr)
    tcase2.AppendStep(step4)
    tcase2.AppendStep(step5)
    tcase2.AppendStep(step6)
    tcase2.AppendStep(step7)
    tcase2.AppendStep(step8)
    tcase2.AppendStep(step9)
    fmt.Println(tcase2.String())
    fmt.Println(tcase2.Xml())
    fmt.Println(tcase2.Json())
    tcase3 := CreateTestCase("testcase3", empty_setup, cleanup1, 
            XFail, NotTested, case_descr)
    tcase3.AppendStep(step0)
    fmt.Println(tcase3.String())
    fmt.Println(tcase3.Xml())
    fmt.Println(tcase3.Json())
    tcase4 := CreateTestCase("testcase4", setup1, empty_cleanup, 
            Pass, Skipped, case_descr)
    tcase4.AppendStep(step9)
    tcase4.AppendStep(step8)
    tcase4.AppendStep(step7)
    fmt.Println(tcase4.String())
    fmt.Println(tcase4.Xml())
    fmt.Println(tcase4.Json())
    tcase5 := CreateTestCase("testcase5", empty_setup, empty_cleanup, 
            XFail, NotTested, case_descr)
    tcase5.AppendStep(step6)
    tcase5.AppendStep(step5)
    tcase5.AppendStep(step4)
    tcase5.AppendStep(step3)
    fmt.Println(tcase5.String())
    fmt.Println(tcase5.Xml())
    fmt.Println(tcase5.Json())
    fmt.Println(">> SUTs...")
    sut1 := CreateSUT("SUT1", Hardware, "1.0.0", 
            "This is description.", "10.0.2.1")
    sut2 := CreateSUT("SUT2", Software, "blah", 
            "This is another description.", "")
    sut3 := CreateSUT("SUT3", Both, "???", 
            "This is yet another description.", "it depends")
    fmt.Println(">> configurations...")
    cfg1 := CreateConfiguration("cfg1", sut1, setup1, cleanup1, cfg_descr)
    cfg1.AppendCase(tcase1)
    fmt.Println(cfg1.String())
    fmt.Println(cfg1.Xml())
    fmt.Println(cfg1.Json())
    cfg2 := CreateConfiguration("cfg2", sut2, setup2, cleanup2, cfg_descr)
    cfg2.AppendCase(tcase2)
    cfg2.AppendCase(tcase3)
    fmt.Println(cfg2.String())
    fmt.Println(cfg2.Xml())
    fmt.Println(cfg2.Json())
    cfg3 := CreateConfiguration("cfg3", sut3, setup1, empty_cleanup,
            cfg_descr)
    cfg3.AppendCase(tcase1)
    cfg3.AppendCase(tcase2)
    cfg3.AppendCase(tcase3)
    cfg3.AppendCase(tcase4)
    cfg3.AppendCase(tcase5)
    fmt.Println(cfg3.String())
    fmt.Println(cfg3.Xml())
    fmt.Println(cfg3.Json())
    cfg4 := CreateConfiguration("cfg4", sut1, empty_setup, cleanup1,
            cfg_descr)
    cfg4.AppendCase(tcase5)
    fmt.Println(cfg4.String())
    fmt.Println(cfg4.Xml())
    fmt.Println(cfg4.Json())
    cfg5 := CreateConfiguration("cfg5", sut3, empty_setup, empty_cleanup,
            cfg_descr)
    cfg5.AppendCase(tcase5)
    cfg5.AppendCase(tcase4)
    cfg5.AppendCase(tcase3)
    cfg5.AppendCase(tcase2)
    cfg5.AppendCase(tcase1)
    fmt.Println(cfg5.String())
    fmt.Println(cfg5.Json())
    fmt.Println(cfg5.Xml())
    // TestSet
    fmt.Println(">> testset...")
    //
    ts := CreateTestSet("testset", set_descr, setup1, cleanup2)
    ts.AppendConfig(cfg1)
    ts.AppendConfig(cfg2)
    ts.AppendConfig(cfg3)
    ts.AppendConfig(cfg4)
    ts.AppendConfig(cfg5)
    fmt.Println(ts.String())
    fmt.Println(ts.Xml())
    json, _ := ts.Json()
    fmt.Println(json)
    // print JSON
}

func testActions() {
    fmt.Println("#### action ####")
    // Creating the first action
    act := CreateAction("action", "arg1 arg2")
    fmt.Printf("The script is %q\n", act.Script)
    fmt.Printf("The arguments are %q\n", act.Args)
    fmt.Println("Is action executable?", act.IsExecutable())
    fmt.Println("Is action manual?", act.IsManual())
    fmt.Println(act.Xml())
    // Creating the second action
    act2 := CreateAction("another action", "arg1")
    fmt.Printf("The script is %q\n", act2.Script)
    fmt.Printf("The arguments are %q\n", act2.Args)
    fmt.Println("Is action executable?", act.IsExecutable())
    fmt.Println("Is action manual?", act.IsManual())
    fmt.Println(act2.Xml())
    // Creating an empty action
    var empty_act Action // should be empty
    fmt.Printf("action: %q\n", empty_act.String())
    fmt.Println(empty_act.Xml())
}
