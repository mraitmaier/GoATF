/*
 * collect.go
 *
 * History:
 *  0.1.0   Apr10   MR  The initial version
 *  0.1.1   Jul11   MR  JSON works, XML is out for good (too complex to handle)
 */

package atf

import (
	"os"
	"json"
//    "fmt"
    "path"
)

/*
 * ConfigType - an enum containing configuration types
 */
type ConfigType int
const (
	UnknownConfig ConfigType = iota
	JsonConfig
	TextConfig
)

/*
 * Collector interface
 *
 */
/*
type Collector interface {
    Collect(source string) *TestSet
}
*/

/*
 * getCfgType - a factory function that determines the type of the config file
 */
func getCfgType(pth string) (cfgtype ConfigType){
    cfgtype = UnknownConfig
    ext := path.Ext(pth)
    switch ext {
        case ".json"        : cfgtype = JsonConfig
        case ".txt", ".cfg" : cfgtype = TextConfig
    }
    return
}

func collectJson(text string) (ts *TestSet) {
    // use standard Unmarshal function to create data from JSON
    err := json.Unmarshal([]uint8(text), ts)
    if err != nil { panic(err) }
    return
}

/*
 * CollectTestSet - function that creates the TestSet struct
 */
func CollectTestSet (path string) (ts *TestSet, err os.Error) {
    // let's create empty TestSet
    ts = new(TestSet)
    // determine the type of config file and unmarshal the data into TestSet 
    switch getCfgType(path) {
        case JsonConfig:
            text, err := ReadTextFile(path)
            if err != nil && err != os.EOF { return nil, err }
            err = json.Unmarshal([]uint8(text), ts)
            if err != nil { panic (err) }
        case TextConfig:
             /* XXX this is not implemented: currently we return error */
             return nil, os.EINVAL
            //ts = collectText(text)
        case UnknownConfig: return nil, os.EINVAL
    }
    return
}

/*
func collectText(text string) (ts *TestSet, err os.Error) {
}
*/



