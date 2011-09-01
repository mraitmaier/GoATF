
package atf

import ("testing")

func TestString(t *testing.T) {
    val1 := Pass
    val2 := Fail
    val3 := XFail
    val4 := NotTested
    val5 := NotAvailable
    val0 := UnknownResult
    //
    answer := val1.String()
    if answer != "pass" {
        t.Errorf("TestResult: Pass test failed")
    }
    //
    answer = val2.String()
    if answer != "fail" {
        t.Errorf("TestResult: Fail test failed")
    }
    //
    answer = val3.String()
    if answer != "expected fail" {
        t.Errorf("TestResult: XFail test failed")
    }
    //
    answer = val4.String()
    if answer != "not tested" {
        t.Errorf("TestResult: NotTested test failed")
    }
    //
    answer = val5.String()
    if answer != "not available" {
        t.Errorf("TestResult: NotAvailable test failed")
    }
    //
    answer = val0.String()
    if answer != "unknown test result" {
        t.Errorf("TestResult: UnknownResult test failed")
    }

}
