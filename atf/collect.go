/*
 * collect.go - implementation of the collector module
 *
 * Collector is a module that collects the configuration (from configuration
 * file) and builds the type hierarchy (that is: scripts) to be executed. 
 * The configuration can be encoded as JSON or XML or plain text (that one is 
 * not implemented yet and frankly I'm not sure that is actually needed; so it 
 * might be omitted in the end...)
 *
 * History:
 *  0.1.0   Apr10   MR  The initial version
 *  0.1.1   Jul11   MR  JSON works, XML is out for good (too complex to handle)
 *  0.2     Mar12   MR  XML is back and it works (with xml.Unmarshal()!), too;
 *                      had to change XML schema and add an <Action> tag
 *                      into <TestStep>
 */

package atf

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"path"
	//  "fmt"
	"bitbucket.org/miranr/goatf/atf/utils"
)

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
 * JsonCollector - defines the JSON collector
 */
type JsonCollector struct {
	Path string
}

func (c *JsonCollector) Collect() (ts *TestSet) {
	ts = new(TestSet)
	text, err := utils.ReadTextFile(c.Path)
	if err != nil && err != io.EOF {
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
	// read the XMl file
	text, err := utils.ReadTextFile(c.Path)
	if err != nil && err != io.EOF {
		return
	}
	// let's parse the XML ; 
	err = xml.Unmarshal([]uint8(text), ts)
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
	ts = new(TestSet)
	// FIXME: no implementation yet, returning empty pointer
	return
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
