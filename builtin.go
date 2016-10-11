package flip

import (
	"bytes"
	"context"
	"fmt"
)

type Version struct {
	Package, Tag, Hash, Date                           string
	PrintPackage, PrintTag, PrintHash, PrintDate, Full bool
}

func NewVersion(pkg, tag, hash, date string) *Version {
	return &Version{
		pkg, tag, hash, date, false, false, false, false, true,
	}
}

var space = []byte(" ")

func (v *Version) String() string {
	b := new(bytes.Buffer)
	if v.PrintPackage {
		b.WriteString(v.Package)
		b.Write(space)
	}
	if v.PrintTag {
		b.WriteString(v.Tag)
		b.Write(space)
	}
	if v.PrintHash {
		b.WriteString(v.Hash)
		b.Write(space)
	}
	if v.PrintDate {
		b.WriteString(v.Date)
		b.Write(space)
	}
	if v.Full {
		b.WriteString(v.full())
	}
	return b.String()
}

func (v *Version) full() string {
	return fmt.Sprintf("%s %s %s %s", v.Package, v.Tag, v.Hash, v.Date)
}

func versionFlag(v *Version) *FlagSet {
	fs := NewFlagSet("version", ContinueOnError)
	fs.BoolVar(&v.Full, "full", true, "Print full version information.")
	fs.BoolVar(&v.PrintPackage, "package", false, "Print available package information.")
	fs.BoolVar(&v.PrintTag, "tag", false, "Print available tag information.")
	fs.BoolVar(&v.PrintHash, "hash", false, "Print available hash informtion.")
	fs.BoolVar(&v.PrintDate, "date", false, "Print available date information.")
	return fs
}

const versionUse = `Prints the package version and exits.`

func (v *Version) Command() Command {
	return NewCommand(
		"",
		"version",
		versionUse,
		1,
		func(c context.Context, a []string) ExitStatus {
			switch {
			case v.PrintPackage, v.PrintTag, v.PrintHash, v.PrintDate:
				v.Full = false
			}
			fmt.Println(v.String())
			return ExitSuccess
		},
		versionFlag(v),
	)
}

func BaseWithVersion(p, t, h, d string) *Commander {
	v := NewVersion(p, t, h, d)
	Base.RegisterGroup("version", 1000, v.Command())
	return Base
}
