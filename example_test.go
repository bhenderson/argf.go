package argf_test

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/bhenderson/argf"
)

func ExampleFrom() {
	r1 := strings.NewReader("first reader ")
	r2 := strings.NewReader("second reader ")
	r3 := strings.NewReader("third reader\n")
	reader, _ := argf.From(r1, r2, r3)

	if _, err := io.Copy(os.Stdout, reader); err != nil {
		log.Fatal(err)
	}
	// Output: first reader second reader third reader
}

func ExampleFromNames() {
	paths := []string{
		writeTestFile("1.txt", "first file\n"),
		writeTestFile("2.txt", "second file\n"),
		writeTestFile("3.txt", "third file\n"),
	}
	defer cleanTestFiles()

	reader, err := argf.FromNames(paths...)
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		// must be called before Read.
		name := argf.Name(reader)
		fmt.Printf("%s: %s\n", name, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	// Output: test_data/1.txt: first file
	// test_data/2.txt: second file
	// test_data/3.txt: third file
}

var testData = "test_data"

func writeTestFile(n, s string) string {
	err := os.Mkdir(testData, 0777)
	if err != nil && !os.IsExist(err) {
		panic(err)
	}
	filename := path.Join(testData, n)
	err = ioutil.WriteFile(filename, []byte(s), 0444)
	if err != nil {
		panic(err)
	}
	return filename
}

func cleanTestFiles() {
	err := os.RemoveAll(testData)
	if err != nil {
		panic(err)
	}
}
