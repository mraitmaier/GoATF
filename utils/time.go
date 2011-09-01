/*
 * time.go -  misc utility functions for working with  date/time
 *
 * History:
 *  0.1.0   Jul11   MR  The initial version
 */

package atf

import (
        "time"
        "strings"
       )

func Now() string {
    t := time.LocalTime()
    return t.Format("2006-01-02 15:04:05")
}

func NowFile() string {
    t := time.LocalTime()
    return t.Format("2006_01_02_15_04_05")
}
func FileConv(o string) (n string) {
    n = strings.Replace(o, " ", "_", -1)
    n = strings.Replace(n, ":", "_", -1)
    n = strings.Replace(n, "-", "_", -1)
    return
}

