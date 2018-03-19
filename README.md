# gover

Version golang applications.

## Install

```
go get github.com/sabhiram/gover
```

## Usage

In a directory with a `main` package:
```
gover init [version]
```

This creates a `version_gen.go` file in the current directory, and sets the starting version to `0.0.1` unless specified in the `init` step above.

To increment the version:
```
gover increment [patch=default|minor|major]
```

To read the current version from the command line:
```
gover [version]
```

Your `main` package will have access to:
```
const (
    Major = 0
    Minor = 0
    Patch = 1

    Version = "0.0.1"
)
```

## Workflow

When trying to tag releases, we can use `gover` to construct the git release tag.  For example:

```
gover increment major
go build .
git -tag `gover version`
git push ... -tags
```
