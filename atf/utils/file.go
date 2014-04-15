/*
 * file.go -  misc utility functions for working with files 
 *
 * History:
 *  0.1.0   Jul11   MR  The initial version
 */

package utils

import (
	"bufio"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

/*
    Read a file with given path and return the contents as a string.
 */
func LoadFile(path string) (text string, err error) {
	text = ""
	// open the file as read-only
	file, err := os.Open(path)
	if err != nil { return }
	defer file.Close() // always close the file

	// read the file line by line
	read := bufio.NewReader(file)
	str, err := read.ReadString('\n')
	text += str
	for err != io.EOF {
		str, err = read.ReadString('\n')
		text += str
	}
	return
}

/*
    Read a text file and return the contents as a string. 
 */
func ReadTextFile(filename string) (string, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil { return "", err }
	return string(data), err
}

/*
    Read a text file and return contents as a list of lines.
 */
func ReadLines(filename string) (lines []string, err error) {
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

/*
    Write a text file with given path and contents.
 */
func WriteTextFile(path string, contents string) (err error) {
	f, err := os.Create(path)
	if err != nil {
		return
	}
	defer f.Close()

	_, err = f.Write([]byte(contents))
	if err != nil {
		return
	}
	return
}

/*
    Copy file from source to destination.
 */
func CopyFile(dst, src string) (int64, error) {
	sf, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer sf.Close()

	df, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer df.Close()

	return io.Copy(df, sf)
}
