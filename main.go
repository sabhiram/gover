package main

////////////////////////////////////////////////////////////////////////////////

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"text/template"
)

////////////////////////////////////////////////////////////////////////////////

const (
	tmpl = `package {{ .PackageName }}

// WARNING: Auto generated version file. Do not edit this file by hand.
// WARNING: go get github.com/sabhiram/gover to manage this file.
// Version: {{ .Major }}.{{ .Minor }}.{{ .Patch }}

const (
    Major = {{ .Major }}
    Minor = {{ .Minor }}
    Patch = {{ .Patch }}

    Version = "{{ .Major }}.{{ .Minor }}.{{ .Patch }}"
)
`
	versionKey     = "// Version: "
	defaultVersion = "0.0.1"

	usageStr = `gover <cmd> [options]

Where "cmd" is one of:

    init [<version>]        Create a "version_gen.go" file with the specified
                            version tag.  If the version is not specified it 
                            defaults to "0.0.1".

    increment [<opt>]       Increment the <opt> section of the version where 
                            "opt" can be one of: "patch", "minor", or "major". 
                            If unspecified, "opt" defaults to "patch".  Once 
                            incremented, all parts of the version of less 
                            significance are reset.

    version                 Print the current version found in the managed 
                            "version_gen.go" file. 

If "cmd" is unspecified, the current version in "version_gen.go" is reported.
`
)

////////////////////////////////////////////////////////////////////////////////

type version struct {
	Major uint64
	Minor uint64
	Patch uint64
}

func (v *version) unmarshal(s string) error {
	var err error

	items := strings.Split(s, ".")
	if len(items) < 3 {
		return fmt.Errorf("%s is an invalid version", s)
	}

	v.Major, err = strconv.ParseUint(items[0], 10, 64)
	if err != nil {
		return err
	}

	v.Minor, err = strconv.ParseUint(items[1], 10, 64)
	if err != nil {
		return err
	}

	v.Patch, err = strconv.ParseUint(items[2], 10, 64)
	return err
}

func (v *version) incrMajor() {
	v.Major += 1
	v.Minor = 0
	v.Patch = 0
}

func (v *version) incrMinor() {
	v.Minor += 1
	v.Patch = 0
}

func (v *version) incrPatch() {
	v.Patch += 1
}

func (v *version) update() error {
	t, err := template.New("gover").Parse(tmpl)
	if err != nil {
		return err
	}

	f, err := os.Create("version_gen.go")
	if err != nil {
		return err
	}
	defer f.Close()

	return t.Execute(f, &context{
		version:     v,
		PackageName: "main",
	})
}

func fromFile() (*version, error) {
	bs, err := ioutil.ReadFile("version_gen.go")
	if err != nil {
		return nil, err
	}

	fc := string(bs)
	if idx := strings.Index(fc, versionKey); idx > 0 {
		ver := fc[idx+len(versionKey):]
		if idx = strings.Index(ver, "\n"); idx > 0 {
			ver = ver[:idx]
		}

		v := &version{}
		return v, v.unmarshal(ver)
	}

	return nil, errors.New("version not found in gen file")
}

////////////////////////////////////////////////////////////////////////////////

type context struct {
	*version
	PackageName string
}

////////////////////////////////////////////////////////////////////////////////

func initFn(s string) error {
	var v version
	if err := v.unmarshal(s); err != nil {
		return err
	}

	return v.update()
}

func versFn() error {
	v, err := fromFile()
	if err != nil {
		return err
	}

	fmt.Printf("%d.%d.%d\n", v.Major, v.Minor, v.Patch)
	return nil
}

func incrFn(inc string) error {
	v, err := fromFile()
	if err != nil {
		return err
	}

	switch inc {
	case "major":
		fmt.Printf("Major incremented.\n")
		v.incrMajor()
	case "minor":
		fmt.Printf("Minor incremented.\n")
		v.incrMinor()
	case "patch":
		fmt.Printf("Patch incremented.\n")
		v.incrPatch()
	default:
		return fmt.Errorf("cannot increment %s", inc)
	}

	return v.update()
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	var (
		err           = error(nil)
		reportVersion = false
		printUsage    = false
		cmd           = "version"
	)

	if len(os.Args) >= 2 {
		cmd = strings.ToLower(os.Args[1])
	}

	switch cmd {
	case "init":
		vers := defaultVersion
		if len(os.Args) >= 3 {
			vers = os.Args[2]
		}
		err = initFn(vers)

	case "version", "vers":
		if _, err := os.Stat("version_gen.go"); os.IsNotExist(err) {
			err = fmt.Errorf("version file missing, did you mean to run `gover init`?")
			printUsage = true
		} else {
			reportVersion = true
		}

	case "increment", "inc":
		o := "patch"
		if len(os.Args) >= 3 {
			o = strings.ToLower(os.Args[2])
		}
		err = incrFn(o)
		if err == nil {
			reportVersion = true
		}

	default:
		err = fmt.Errorf("invalid command specified")
		printUsage = true
	}

	if err == nil && reportVersion == true {
		err = versFn()
	}

	if err != nil {
		fmt.Printf("Fatal error: %s\n", err.Error())
		defer os.Exit(1)
	}

	if printUsage {
		fmt.Printf("%s\n", usageStr)
	}
}

////////////////////////////////////////////////////////////////////////////////
