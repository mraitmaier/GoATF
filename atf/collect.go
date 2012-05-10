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
	"path"
	"xml"
	//   "strings"
	//  "fmt"
)

/* define white-space for local use */
//const white = " \n\t"

/*
 * ConfigType - an enum containing configuration types
 */
type ConfigType int

const (
	UnknownConfig ConfigType = iota
	JsonConfig
	TextConfig
	XmlConfig // currently unused
)

/****************************************************************************
 * Collector interface - defines the types that implement Collect() method
 */
type Collector interface {
	Collect() *TestSet
}

/****************************************************************************
 * JsonCollector - defines the JSON 
 */
type JsonCollector struct {
	Path string
}

func (c *JsonCollector) Collect() (ts *TestSet) {
	ts = new(TestSet)
	text, err := ReadTextFile(c.Path)
	if err != nil && err != os.EOF {
		return
	}
	err = json.Unmarshal([]uint8(text), ts)
	if err != nil {
		panic(err)
	}
	return
}

/****************************************************************************
 * XmlCollector - defines the  
 */
type XmlCollector struct {
	Path string
}

func (c *XmlCollector) Collect() (ts *TestSet) {
	ts = new(TestSet)
	//
	fin, err := os.Open(c.Path)
	if err != nil {
		panic(err)
	}
	defer fin.Close()
	// let's parse the XML ; 
	// FIXME: Unmarshal probably won't work, custom parser will have to be used
	err = xml.Unmarshal(fin, ts)
	if err != nil {
		panic(err)
	}
	return
}

/****************************************************************************
 * TextCollector - defines the 
 */
type TextCollector struct {
	Path string
}

func (c *TextCollector) Collect() (ts *TestSet) {
	//ts = new(TestSet)
	// FIXME: no implementation yet, returning empty pointer
	return nil
}

/*
 * getCfgType - a factory function that determines the type of the config file
 */
func getCfgType(pth string) (cfgtype ConfigType) {
	cfgtype = UnknownConfig
	ext := path.Ext(pth)
	switch ext {
	case ".json":
		cfgtype = JsonConfig
	case ".txt", ".cfg":
		cfgtype = TextConfig
	case ".xml":
		cfgtype = XmlConfig
	}
	return
}

/*****************************************************************************
 * CollectTestSet - function that creates the TestSet struct
 */
func CollectTestSet(path string) (ts *TestSet) {
	// let's create empty TestSet
	ts = new(TestSet)
	// we need one of the Collectors to get test set data
	var c Collector
	// determine the type of config file and unmarshal the data into TestSet 
	switch getCfgType(path) {
	case JsonConfig:
		c = new(JsonCollector)
		c.(*JsonCollector).Path = path
	case TextConfig:
		c = new(TextCollector)
		c.(*TextCollector).Path = path
	case XmlConfig:
		c = new(XmlCollector)
		c.(*XmlCollector).Path = path
	case UnknownConfig:
		return nil
	}
	// now collect the test set structure
	ts = c.Collect()
	//if ts == nil { fmt.Println("TestSet is empty.") }
	return
}

/******************************************************************************
 * resolveEndElem - take necesary when action when XML end element occurs
 */
/*
func resolveEndElem(elem *xml.EndElement, ts *TestSet) interface{} {
    var current_obj interface{} = ts
    e := elem.Name.Local
    switch e {
        case "TestSet":
        case "SystemUnderTest":
        case "TestCase":
        case "TestStep":
    }
    return current_obj
}
*/
/******************************************************************************
 * resolveStartElem - take necesary when action when XML start element occurs
 */
/*
func resolveStartElem(elem *xml.StartElement, obj interface{}) (interface{}, string) {
    var current_obj interface{} = obj
    e := elem.Name.Local
    switch e {
    case "TestSet":
        // extract the test set name (it is an attribute)
        if len(elem.Attr) > 0 { obj.(*TestSet).Name = elem.Attr[0].Value }
        fmt.Println("DEBUG: TestSet=" + obj.(*TestSet).Name)
    case "TestPlan":
    case "SystemUnderTest":
        // first create new SUT and then extract the name attribute
        obj.(*TestSet).SysUnderTest = new(SysUnderTest)
        if len(elem.Attr) > 0 {
            obj.(*TestSet).SysUnderTest.Name = elem.Attr[0].Value
        }
        current_obj = obj.(*TestSet).SysUnderTest
    case "TestCase":
        tc := new(TestCase)
        obj.(*TestSet).Cases = append(obj.(*TestSet).Cases, *tc)
        // extract the test case name (it is an attribute)
        if len(elem.Attr) > 0 { tc.Name = elem.Attr[0].Value }
        current_obj = tc
        fmt.Printf("DEBUG: TC=%q\n", tc.Name)
    case "TestStep":
        step := new(TestStep)
        if len(elem.Attr) > 0 {
            for _, attr := range elem.Attr {
                switch attr.Name.Local {
                case "name":     step.Name = attr.Value
                case "expected": step.Expected = ResolveResult(attr.Value)
                }
            }
        }
        step.Action = new(Action)
        obj.(*TestCase).Steps = append(obj.(*TestCase).Steps, *step)
        current_obj = step
    case "Setup":
    case "Cleanup":
    case "Script":
        a := new(Action)
        switch t := obj.(type) {
        case TestSet:
            t.Setup = a
        case TestCase:
            t.Setup = a
        }
    case "Args":
    case "Type":
    case "IP":
    case "Description":
    }
    return current_obj, e
}
*/

/******************************************************************************
 * resolveData - take necesary when action when XML (char)data occurs
 */
/*
func resolveData(tag string, data *xml.CharData, obj interface{}) {
    switch tag {
    case "TestSet":
    case "TestPlan":
        ts := obj.(*TestSet)
        tp := strings.Trim(string([]byte(data.Copy())), white)
        if tp != "" { ts.TestPlan = tp }
    case "SystemUnderTest":
    case "TestCase":
    case "TestStep":
    case "Setup":
    case "Cleanup":
    case "Script":
    case "Args":
    case "Type":
        t := strings.Trim(string([]byte(data.Copy())), white)
        if t != "" { obj.(*SysUnderTest).Systype = SUTTypeValue(t) }
    case "IP":
        // turn CharData into a string and trim the white-spaces
        ip := strings.Trim(string([]byte(data.Copy())), white)
        if ip != "" { obj.(*SysUnderTest).IPaddr = ip }
    case "Description":
        d := strings.Trim(string([]byte(data.Copy())), white)
        err := handleDesc(obj, d)
        if err != nil { panic(err) }
    }
    return
}
*/
/******************************************************************************
 * handleDesc - handle the Description XML tag's data
 *
 *  Description XML tag is used as a child to more XML tags and this function
 *  handles the these events.
 */
/*
func handleDesc(obj interface{}, s string) (err os.Error) {
    switch t := obj.(type) {
        case *SysUnderTest:
            if s != "" { obj.(*SysUnderTest).Description = s }
        case *TestSet:
            if s != "" { obj.(*TestSet).Description = s }
        case *TestCase:
            if s != "" { obj.(*TestCase).Description = s }
//        default:
//            err = os.EINVAL
    }
    return
}
*/
/******************************************************************************
 * collectXml - parses the XML file and creates the test set structure
 */
/*
func collectXml(path string) (ts *TestSet) {
    // open the XML file
    fin, err := os.Open(path)
    if err != nil { panic(err) }
    defer fin.Close()
    // let's parse the XML 
    var obj interface{}
    ts = new(TestSet)
    obj = ts
    p := xml.NewParser(fin)
    // iterate through XML and act according to XML token type...
    var tagname string
    for token, err := p.Token(); err == nil; token, err = p.Token() {
        // 
        switch t := token.(type) {
        case xml.StartElement:
                obj, tagname = resolveStartElem(&t, obj)
                / *
                switch objtype := obj.(type) {
                    case *TestCase: ts.Cases = append(ts.Cases, *objtype)
                }
                * /
        case xml.EndElement:
            obj = resolveEndElem(&t, ts)
        case xml.CharData:
            resolveData(tagname, &t, obj)
        }
    }
    return
}
*/
/*
func collectText(text string) (ts *TestSet, err os.Error) {
}
*/
