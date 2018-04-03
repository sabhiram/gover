package main

////////////////////////////////////////////////////////////////////////////////

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/sabhiram/gover/version"
)

////////////////////////////////////////////////////////////////////////////////

const (
	defaultVersion = "0.0.1"
	usageStr       = `gover <cmd> [options]

Where "cmd" is one of:

    init [-version <x.y.z>] [-package <p>]

        Create a "version_gen.go" file with the specified version tag.  If the
        version is not specified it defaults to "0.0.1".

        If the "-package" is specified, the file is generated under the package
        "<p>".  If omitted, the package defaults to "main".

    increment [<opt>]

        Increment the <opt> section of the version where "opt" can be one of:
        "patch", "minor", or "major".

        If unspecified, "opt" defaults to "patch".  Once incremented, all parts
        of the version of less significance are reset.

    major

    	Same as "increment major".

    minor

    	Same as "increment minor".

    patch

    	Same as "increment patch".

    version

        Print the current version found in the managed "version_gen.go" file.

If "cmd" is unspecified, the current version in "version_gen.go" is reported.
`
)

////////////////////////////////////////////////////////////////////////////////

func main() {
	var (
		err           = error(nil)
		reportVersion = false
		printUsage    = false
		cmd           = "version"
		args          = []string{}
	)

	if len(os.Args) >= 2 {
		cmd = strings.ToLower(os.Args[1])
		args = os.Args[2:]
	}

	switch cmd {
	case "init":
		var pn, vs string
		fs := flag.NewFlagSet("init-option", flag.PanicOnError)
		fs.SetOutput(bytes.NewBuffer([]byte{}))
		fs.StringVar(&vs, "version", defaultVersion, "Override init version")
		fs.StringVar(&vs, "v", defaultVersion, "Override init version (short)")
		fs.StringVar(&pn, "package", "main", "Override package name")
		fs.StringVar(&pn, "p", "main", "Override package name (short)")
		if err = fs.Parse(args); err != nil {
			printUsage = true
			break
		}
		err = version.Init(vs, pn)

	case "version", "vers":
		if _, err := os.Stat("version_gen.go"); os.IsNotExist(err) {
			err = fmt.Errorf("version file missing, did you mean to run `gover init`?")
			printUsage = true
		} else {
			reportVersion = true
		}

	case "increment", "inc", "major", "minor", "patch":
		o := "patch"
		if len(args) >= 1 {
			o = strings.ToLower(args[0])
		}
		switch cmd {
		case "major", "minor", "patch":
			o = cmd
		}
		if err = version.Increment(o); err == nil {
			reportVersion = true
		}

	default:
		err = fmt.Errorf("invalid command specified")
		printUsage = true
	}

	if err == nil && reportVersion == true {
		var v *version.Version
		v, err = version.Load()
		if err == nil {
			fmt.Printf("%s\n", v.String())
		}
	}

	if err != nil {
		fmt.Printf("Fatal error: %s\n", err.Error())
		defer os.Exit(1)
	}

	if printUsage {
		fmt.Printf("%s\n", usageStr)
	}
}
