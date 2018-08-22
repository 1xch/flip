package flip

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"syscall"
	"unsafe"
)

func failOnly(f *FlagSet, format string, a ...interface{}) error {
	err := fmt.Errorf(format, a...)
	fmt.Fprintln(f.Out(), err)
	return err
}

func failFmt(f *FlagSet, format string, a ...interface{}) error {
	err := fmt.Errorf(format, a...)
	fmt.Fprintln(f.Out(), err)
	f.Usage(f.Out())
	return err
}

type color struct {
	params []Attribute
}

// Should work in most terminals.
// See github.com/mattn/go-colorable for tweaking tips by os.
func Color(value ...Attribute) func(io.Writer, ...interface{}) {
	c := &color{params: make([]Attribute, 0)}
	c.Add(value...)
	return c.Fprint
}

func (c *color) Add(value ...Attribute) *color {
	c.params = append(c.params, value...)
	return c
}

func (c *color) Fprint(w io.Writer, a ...interface{}) {
	c.wrap(w, a...)
}

func (c *color) Fprintf(w io.Writer, f string, a ...interface{}) {
	c.wrap(w, fmt.Sprintf(f, a...))
}

func (c *color) sequence() string {
	format := make([]string, len(c.params))
	for i, v := range c.params {
		format[i] = strconv.Itoa(int(v))
	}

	return strings.Join(format, ";")
}

func (c *color) wrap(w io.Writer, a ...interface{}) {
	if c.noColor() {
		fmt.Fprint(w, a...)
	}

	c.format(w)
	fmt.Fprint(w, a...)
	c.unformat(w)
}

func (c *color) format(w io.Writer) {
	fmt.Fprintf(w, "%s[%sm", escape, c.sequence())
}

func (c *color) unformat(w io.Writer) {
	fmt.Fprintf(w, "%s[%dm", escape, Reset)
}

var NoColor = !IsTerminal(os.Stdout.Fd())

const ioctlReadTermios = syscall.TCGETS

// IsTerminal return true if the file descriptor is terminal.
// see github.com/mattn/go-isatty
// You WILL want to change this if you are using an os other than a Linux variant.
func IsTerminal(fd uintptr) bool {
	var termios syscall.Termios
	_, _, err := syscall.Syscall6(syscall.SYS_IOCTL, fd, ioctlReadTermios, uintptr(unsafe.Pointer(&termios)), 0, 0, 0)
	return err == 0
}

func (c *color) noColor() bool {
	return NoColor
}

const escape = "\x1b"

type Attribute int

const (
	Reset Attribute = iota
	Bold
	Faint
	Italic
	Underline
	BlinkSlow
	BlinkRapid
	ReverseVideo
	Concealed
	CrossedOut
)

const (
	FgBlack Attribute = iota + 30
	FgRed
	FgGreen
	FgYellow
	FgBlue
	FgMagenta
	FgCyan
	FgWhite
)

const (
	FgHiBlack Attribute = iota + 90
	FgHiRed
	FgHiGreen
	FgHiYellow
	FgHiBlue
	FgHiMagenta
	FgHiCyan
	FgHiWhite
)

const (
	BgBlack Attribute = iota + 40
	BgRed
	BgGreen
	BgYellow
	BgBlue
	BgMagenta
	BgCyan
	BgWhite
)

const (
	BgHiBlack Attribute = iota + 100
	BgHiRed
	BgHiGreen
	BgHiYellow
	BgHiBlue
	BgHiMagenta
	BgHiCyan
	BgHiWhite
)

var (
	black   = Color(FgHiBlack)
	red     = Color(FgHiRed)
	green   = Color(FgHiGreen)
	yellow  = Color(FgHiYellow)
	blue    = Color(FgHiBlue)
	magenta = Color(FgHiMagenta)
	cyan    = Color(FgHiCyan)
	white   = Color(FgHiWhite)
)

//paramBool

//paramString

//func paramInt(r *regexp.Regexp, s, k string) int {
//	match := r.FindStringSubmatch(s)
//	for i, name := range r.SubexpNames() {
//		if i > 0 && i <= len(match) {
//			if k == name {
//				v, err := strconv.Atoi(match[i])
//				if err == nil {
//					return v
//				}
//			}
//		}
//	}
//	return 0
//}

//func paramFloat64(r *regexp.Regexp, s, k string) float64 {
//	match := r.FindStringSubmatch(s)
//	for i, name := range r.SubexpNames() {
//		if i > 0 && i <= len(match) {
//			if k == name {
//				v, err := strconv.ParseFloat(match[i], 64)
//				if err == nil {
//					return v
//				}
//			}
//		}
//	}
//	return 0.0
//}
