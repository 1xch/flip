package flip

import (
	"bytes"
	"context"
	"fmt"
	"strings"
)

type Help struct {
	f        Flip
	Full     bool
	Commands string
}

func NewHelp(f Flip) *Help {
	return &Help{f, true, ""}
}

func helpFlag(h *Help) *FlagSet {
	fs := NewFlagSet("help", ContinueOnError)
	fs.BoolVar(&h.Full, "full", true, "Print all help information.")
	fs.StringVar(&h.Commands, "commands", "", "Print help information for a subset of comma delimited commands or command groups")
	return fs
}

func (h *Help) Command() Command {
	return NewCommand(
		"",
		"help",
		`Print help information on demand.`,
		1,
		true,
		func(c context.Context, a []string) (context.Context, ExitStatus) {
			switch {
			case h.Commands != "":
				h.Full = false
			case len(a) > 0:
				h.Full = false
				h.Commands = strings.Join(a, ",")
			}
			switch {
			case h.Full:
				h.f.Instruction(c)
			case !h.Full:
				var cs []Command
				spl := strings.Split(h.Commands, ",")
				for _, v := range spl {
					gc := h.f.GetCommand(v)
					if len(gc) > 0 {
						cs = append(cs, gc...)
					}
				}
				h.f.SubsetInstruction(cs...)(c)
			}
			h.reset()
			return c, ExitSuccess
		},
		helpFlag(h),
	)
}

func (f *flip) addHelp() Flip {
	h := NewHelp(f)
	f.SetGroup("help", 1000, h.Command())
	return f
}

func (h *Help) reset() {
	h.Full = true
	h.Commands = ""
}

type Version struct {
	f                                                  Flip
	Package, Tag, Hash, Date                           string
	PrintPackage, PrintTag, PrintHash, PrintDate, Full bool
}

func NewVersion(f Flip, pkg, tag, hash, date string) *Version {
	return &Version{
		f, pkg, tag, hash, date, false, false, false, false, true,
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
	v.reset()
	return b.String()
}

func (v *Version) reset() {
	v.Full = true
	v.PrintPackage, v.PrintTag, v.PrintHash, v.PrintDate = false, false, false, false
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

func (v *Version) Command() Command {
	return NewCommand(
		"",
		"version",
		`Prints the package version and exits.`,
		1,
		true,
		func(c context.Context, a []string) (context.Context, ExitStatus) {
			switch {
			case v.PrintPackage, v.PrintTag, v.PrintHash, v.PrintDate:
				v.Full = false
			}
			b := new(bytes.Buffer)
			b.WriteString(v.String())
			o := v.f.Out()
			fmt.Fprint(o, b)
			return c, ExitSuccess
		},
		versionFlag(v),
	)
}

func (f *flip) addVersion(args ...string) Flip {
	p, t, h, d := "not provided", "not provided", "not provided", "not provided"
	in := len(args) - 1
	for i := -1; i <= 3; i++ {
		if i <= in {
			switch i {
			case 0:
				p = args[0]
			case 1:
				t = args[1]
			case 2:
				h = args[2]
			case 3:
				d = args[3]
			}
		}
	}
	v := NewVersion(f, p, t, h, d)
	f.SetGroup("version", 1000, v.Command())
	return f
}
