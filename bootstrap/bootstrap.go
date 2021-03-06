/*
The `bootstrap` program is used to generated language patterns and store
them as the Go source file `../data.go`. It should not be called by users,
unless you want to add or change patterns and include them into the
library. You then need to rebuild the library.

For generating additional language patterns and load them into a running
program, see `../textpat/`.
*/
package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"strings"
)

var header = `// THIS IS A GENERATED FILE
// DO NOT EDIT

package textcat

`

func main() {
	out, err := os.Create("../data.go")
	checkErr(err)
	fmt.Fprint(out, header)

	sep := "\n"
	fmt.Fprint(out, "var data = map[string]map[string]int{")
	for _, utf := range []bool{true, false} {
		suffix := ".raw"
		other := "utf8"
		if utf {
			suffix = ".utf8"
			other = "raw"
		}
		for _, filename := range os.Args[1:] {
			if strings.HasSuffix(path.Dir(filename), other) {
				continue
			}
			lang := path.Base(filename)
			lang = lang[:len(lang)-len(path.Ext(lang))]
			fmt.Fprintf(out, "%s\t%q: {", sep, lang+suffix)
			sep = ",\n"

			r, e := os.Open(filename)
			checkErr(e)
			b, e := ioutil.ReadAll(r)
			r.Close()
			checkErr(e)

			sep2 := "\n"
			for n, p := range getPatterns(string(b), utf) {
				fmt.Fprintf(out, "%s\t\t%q: %d", sep2, p.S, n)
				sep2 = ",\n"
			}
			fmt.Fprint(out, "}")
		}
	}
	fmt.Fprint(out, "}\n\n")

	out.Close()
}

func checkErr(err error) {
	if err != nil {
		_, filename, lineno, ok := runtime.Caller(1)
		if ok {
			fmt.Fprintf(os.Stderr, "%v:%v: %v\n", filename, lineno, err)
		}
		panic(err)
	}
}
