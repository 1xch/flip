package flip

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path"
	"strings"
)

var Base Flip

func SetCommand(cmd Command) {
	Base.SetCommand(cmd)
}

func SetGroup(name string, priority int, cmds ...Command) {
	Base.SetGroup(name, priority, cmds...)
}

func init() {
	Base = New(path.Base(os.Args[0]))
}

type Help struct {
	f                    Flip
	Full, Subset, Single bool
	Commands             string
}

func NewHelp(f Flip) *Help {
	return &Help{f, true, false, false, ""}
}

func helpFlag(h *Help) *FlagSet {
	fs := NewFlagSet("help", ContinueOnError)
	fs.BoolVar(&h.Full, "full", true, "Print all help information.")
	fs.StringVar(&h.Commands, "commands", "", "Print help information for a subset of comma delimited commands")
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
				h.Subset = true
			case len(a) > 1:
				h.Commands = strings.Join(a, ",")
				h.Subset = true
			case len(a) == 1:
				h.Full = false
				h.Single = true
			}
			switch {
			case h.Full:
				h.f.Instruction(c)
			case h.Subset:
				var cs []Command
				spl := strings.Split(h.Commands, ",")
				for _, v := range spl {
					gc := h.f.GetCommand(v)
					if gc != nil {
						cs = append(cs, gc)
					}
				}
				h.f.SubsetInstruction(cs...)(c)
			case h.Single:
				gc := h.f.GetCommand(a[0])
				if gc != nil {
					h.f.SubsetInstruction(gc)(c)
				}
			}
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

func SetHelp() {
	Base.AddCommand("help")
}

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
			fmt.Println(v.String())
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
	v := NewVersion(p, t, h, d)
	f.SetGroup("version", 1000, v.Command())
	return f
}

func SetVersion(args ...string) {
	Base.AddCommand("version", args...)
}
