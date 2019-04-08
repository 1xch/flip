package flip

import (
	"bytes"
	"context"
	"fmt"
	"strings"
)

// An interface that handles adding(package builtin commands) by string parameters.
type Adder interface {
	AddBuiltIn(string, ...string) Flipper
}

type help struct {
	f        Flipper
	full     bool
	commands string
}

func newHelp(f Flipper) *help {
	return &help{f, true, ""}
}

func helpFlag(h *help) *FlagSet {
	fs := NewFlagSet("help", ContinueOnError)
	fs.BoolVar(&h.full, "full", true, "Print all help information.")
	fs.StringVar(&h.commands, "commands", "", "Print help information for a subset of comma delimited commands or command groups")
	return fs
}

func (h *help) command() Command {
	return NewCommand(
		"",
		"help",
		`Print help information on demand.`,
		1,
		true,
		func(c context.Context, a []string) (context.Context, ExitStatus) {
			switch {
			case h.commands != "":
				h.full = false
			case len(a) > 0:
				h.full = false
				h.commands = strings.Join(a, ",")
			}
			switch {
			case h.full:
				h.f.Instruction(c)
			case !h.full:
				var cs []Command
				spl := strings.Split(h.commands, ",")
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

func (f *Flip) addHelp() Flipper {
	h := newHelp(f)
	f.SetGroup("help", 1000, h.command())
	return f
}

func (h *help) reset() {
	h.full = true
	h.commands = ""
}

type version struct {
	f                                                  Flipper
	vpackage, tag, hash, date                          string
	printPackage, printTag, printHash, printDate, full bool
}

func newVersion(f Flipper, pkg, tag, hash, date string) *version {
	return &version{
		f, pkg, tag, hash, date, false, false, false, false, true,
	}
}

var space = []byte(" ")

func (v *version) String() string {
	b := new(bytes.Buffer)
	if v.printPackage {
		b.WriteString(v.vpackage)
		b.Write(space)
	}
	if v.printTag {
		b.WriteString(v.tag)
		b.Write(space)
	}
	if v.printHash {
		b.WriteString(v.hash)
		b.Write(space)
	}
	if v.printDate {
		b.WriteString(v.date)
		b.Write(space)
	}
	if v.full {
		b.WriteString(v.fullString())
	}
	v.reset()
	b.WriteString("\n")
	return b.String()
}

func (v *version) reset() {
	v.full = true
	v.printPackage, v.printTag, v.printHash, v.printDate = false, false, false, false
}

func (v *version) fullString() string {
	return fmt.Sprintf("%s %s %s %s", v.vpackage, v.tag, v.hash, v.date)
}

func versionFlag(v *version) *FlagSet {
	fs := NewFlagSet("version", ContinueOnError)
	fs.BoolVar(&v.full, "full", true, "Print full version information.")
	fs.BoolVar(&v.printPackage, "package", false, "Print available package information.")
	fs.BoolVar(&v.printTag, "tag", false, "Print available tag information.")
	fs.BoolVar(&v.printHash, "hash", false, "Print available hash informtion.")
	fs.BoolVar(&v.printDate, "date", false, "Print available date information.")
	return fs
}

func (v *version) command() Command {
	return NewCommand(
		"",
		"version",
		`Prints the package version and exits.`,
		1,
		true,
		func(c context.Context, a []string) (context.Context, ExitStatus) {
			switch {
			case v.printPackage, v.printTag, v.printHash, v.printDate:
				v.full = false
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

func (f *Flip) addVersion(args ...string) Flipper {
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
	v := newVersion(f, p, t, h, d)
	f.SetGroup("version", 1000, v.command())
	return f
}
