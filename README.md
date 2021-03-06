# gover

[![Build Status](https://travis-ci.org/sabhiram/gover.svg?branch=master)](https://travis-ci.org/sabhiram/gover)

Version golang applications.

## Idea

A command line tool to generate and manage a `version.go` file which keeps track of the current [semantic version](https://semver.org/) of the go application.

## Install

```
go get github.com/sabhiram/gover
```

## Usage

```
$ gover <cmd> [options]

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

    version

        Print the current version found in the managed "version_gen.go" file.

If "cmd" is unspecified, the current version in "version_gen.go" is reported.
```

## Workflow

In a directory with a `main` package:
```
gover init -version 0.0.1 [-package main]
```

This creates a `version_gen.go` file in the current directory, and sets the starting version to `0.0.1` unless specified in the `init` step above.

Your `main` package will have access to:
```
const (
    Major = 0
    Minor = 0
    Patch = 1

    Version = "0.0.1"
)
```

Increment the version's major:
```
gover increment major

#   - or -

gover major
```

Increment the version's minor:
```
gover increment minor

#   - or -

gover minor
```

Increment the version's patch:
```
gover increment patch

#   - or -

gover patch
```

Read the current version from the command line:
```
gover version

#   - or -

gover
```

## CI/CD Tagging

When trying to tag releases, we can use `gover` to construct the git release tag.  For example:

```
gover increment major
go build .
git tag `gover`
git push ... --tags
```
