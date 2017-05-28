/* gbstrings, a quick command-line tool for finding GB2312 strings in
   firmware images. */

package main

import (
	//Standard libraries.
	"flag"
	"fmt"
	"io/ioutil"
	"strings"

	//go get github.com/djimenez/iconv-go
	iconv "github.com/djimenez/iconv-go"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

//Returns the maximum length of a GB2312 string at an address.
func gblen(file []byte, index int) int {
	i := 0
	for i = index; i+1 < len(file); i++ {
		//If neither one byte nor two appended bytes make a
		//valid string, we've hit the end of the string.
		if !fastvalidgb(file[index:i+1]) && !fastvalidgb(file[index:i+2]) {
			return i - index
		}
	}

	//We're done by default at the end of the file.
	return i - index
}

//Searches a file for the next potential GB2312 string.
func findnextstring(file []byte, index int) (int, []byte, int) {
	i := 0

	for i = index; i < len(file); i++ {
		len := gblen(file, i)
		//fmt.Printf("Length is %d\n",len);
		if len >= *minlength && validgb(file[i:i+len]) {
			return i, file[i : i+len], i + len
		}
	}

	//Returns the string address, the string, and the next index
	//to search.
	return -1, []byte{}, index + 1
}

//Is the byte string a valid GB2312 candidate?  This only performs a
//brief check, and libiconv should be the final judge.
func fastvalidgb(in []byte) bool {
	for i := 0; i < len(in); i++ {
		if in[i] == 0 {
			//No null bytes.
			return false
		} else if in[i] > 0 && in[i] < 0x20 {
			//No control bytes.  Might exclude format strings.
			return false
		} else if in[i] > 0x7f && in[i] < 0xA1 {
			//In extended range, but not GB2312.
			return false
		}
	}

	//Return true in that this is a valid candidate.  Libiconv
	//should be checked in validgb() to be sure.
	return true
}

//Can iconv convert the string to Unicode?
func validgb(in []byte) bool {
	out := make([]byte, len(in)*4)

	//First we try to reject the string for obvious reasons.
	if !fastvalidgb(in) {
		return false
	}

	//If all else passes, we try to perform the conversion.  This
	//is slow.
	_, _, err := iconv.Convert(in, out, "gb2312", "utf-8")
	return (err == nil)
}

//Converts a GB2312 string to Unicode.
func fromgb(in []byte) string {
	out := make([]byte, len(in)*4)

	_, _, err := iconv.Convert(in, out, "gb2312", "utf-8")
	check(err)
	return strings.Trim(string(out), " \n\r\x00")
}

//Handles an input file, printing all strings.
func handlefile(filename string) {
	dat, err := ioutil.ReadFile(filename)
	check(err)

	for i := 0; i != -1 && i < len(dat); {
		at, str, j := findnextstring(dat, i)
		if len(str) > *minlength {
			fmt.Printf("0x%08x: %s\n",
				at, fromgb(str))
			i = j
		} else {
			i++
		}
	}
}

//Tests the code by converting known strings.
func test() {
	//fmt.Printf("Based at 0x%08x.\n", *baseadr)

	hello := []byte{
		0xc4, 0xe3, //你
		0xba, 0xc3, //好
		0xa3, 0xac, //Chinese comma "，"
		//Travis
		0x54, 0x72, 0x61, 0x76, 0x69, 0x73, 0x2e,
		//Newline
		0x0a}
	fmt.Printf("%s\n", fromgb(hello))

	//FM收音机  (FM Radio)
	fm := []byte{
		0x46, 0x4d, //FM
		0xca, 0xd5, //收
		0xd2, 0xf4, //音
		0xbb, 0xfa, //机
		//Newline
		0x0a}
	fmt.Printf("%s\n", fromgb(fm))

	//快速组队 (Quick Team)
	b := []byte{
		0xbf, 0xec, 0xcb, 0xd9, 0xd7, 0xe9, 0xb6, 0xd3,
		//Newline
		0x0a}
	fmt.Printf("%s\n", fromgb(b))

}

var testflag = flag.Bool("test", false, "Tests the app.")
var baseadr = flag.Int("b", 0, "Base address of target file.")
var input = flag.String("i", "", "Input file.")
var minlength = flag.Int("n", 8, "Minimum string length in bytes.")

func main() {
	flag.Parse()

	if *testflag {
		test()
	} else if strings.Compare(*input, "") != 0 {
		handlefile(*input)
	} else {
		flag.Usage()
	}
}
