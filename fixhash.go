// Fixhash finds x where hash(x) matches /regex/.
package main

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"hash"
	"os"
	"regexp"

	"github.com/lucasjones/reggen"
)

var algorithm = flag.String("f", "sha256", "Hash function [md5|sha256]")
var maxRepeat = flag.Int("r", 10, "Maximum [*,+] repeats")
var inPattern = flag.String("i", "[a-z]+", "Input pattern")
var outPattern = flag.String("o", "", "Output pattern")

func main() {
	flag.Parse()

	var h hash.Hash
	if *algorithm == "sha256" {
		h = sha256.New()
	} else if *algorithm == "md5" {
		h = md5.New()
	} else {
		fmt.Fprintf(os.Stderr, "Unsuported algorithm: %s\n", *algorithm)
		os.Exit(1)
	}

	if *maxRepeat < 1 {
		fmt.Fprintln(os.Stderr, "r must be positive.")
		os.Exit(1)
	}

	g, err := reggen.NewGenerator(*inPattern)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Bad input pattern.")
		os.Exit(1)
	}

	re, err := regexp.Compile(*outPattern)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Bad output pattern.")
		os.Exit(1)
	}

	var in, out string
	for {
		in = g.Generate(*maxRepeat)

		h.Reset()
		h.Write([]byte(in))
		out = hex.EncodeToString(h.Sum(nil))

		if re.MatchString(out) {
			break
		}
	}
	fmt.Printf("%s\t%s\n", out, in)
}
