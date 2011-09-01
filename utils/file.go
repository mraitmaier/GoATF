/*
 * file.go -  misc utility functions for working with files 
 *
 * History:
 *  0.1.0   Jul11   MR  The initial version
 */

package atf

import (
	"os"
	"io/ioutil"
	"strings"
	"bufio"
)

/*
 * LoadFile - read a file with 'filename' and return the contents as a string
 */
func LoadFile(path string) (text string, err os.Error) {
	text = ""
	// open the file as read-only
	file, err := os.Open(path)
	if err != nil {
		return text, err
	}
	defer file.Close() // always close the file
	// read the file line by line
	read := bufio.NewReader(file)
	str, err := read.ReadString('\n')
	text += str
	for err != os.EOF {
		str, err = read.ReadString('\n')
		text += str
	}
	return text, err
}

/*
 * ReadTextFile - read a text file and return the contents as a string 
 *
 * If an error occurs during file read, we return an empty string (and 
 * an os.Error, of course).
 */
func ReadTextFile(filename string) (string, os.Error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(data), err
}

/*
 * ReadLines - read a text file and return a list of lines 
 *
 * If an error occurs during file read, we return only a list with single empty
 * string (and an os.Error, of course).
 */
func ReadLines(filename string) (lines []string, err os.Error) {
	// we read a file
	data, err := ioutil.ReadFile(filename)
	// if there's an error reading a file, we return a list with single empty
	// string and error
	if err != nil {
		return []string{""}, err
	}
	// now we convert the text into an array of lines
	lines = strings.Split(string(data), "\n")
	return
}
