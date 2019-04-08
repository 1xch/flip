package flip

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

// A type representing one command line flag
type Flag struct {
	Name     string // name as it appears on command line
	Message  string // help message
	Value    Value  // value as set
	DefValue string // default value (as text); for usage message
}

// A package level interface for abstracting flag values
type Value interface {
	String() string
	Set(string) error
	Get() interface{}
}

type (
	setFn func(string) error
	getFn func() interface{}
)

type containValue struct {
	to   string
	sfn  setFn
	gfn  getFn
	kind string
}

// Value interface Set function for internal type containValue
func (v *containValue) Set(n string) error {
	return v.sfn(n)
}

// Value interface Get function for internal type containValue
func (v *containValue) Get() interface{} {
	return v.gfn()
}

// Value interface String function for internal type containValue
func (v *containValue) String() string {
	return fmt.Sprintf("%v", v.Get())
}

// internal boolFlag interface IsBoolFlag function for internal type containValue
func (v *containValue) IsBoolFlag() bool {
	if v.kind == "bool" {
		return true
	}
	return false
}

type boolValue bool

func newBoolValue(val bool, p *bool) *boolValue {
	*p = val
	return (*boolValue)(p)
}

// Value interface Set function for internal type boolValue
func (b *boolValue) Set(s string) error {
	v, err := strconv.ParseBool(s)
	*b = boolValue(v)
	return err
}

// Value interface Get function for internal type boolValue
func (b *boolValue) Get() interface{} { return bool(*b) }

// Value interface String function for internal type boolValue
func (b *boolValue) String() string { return fmt.Sprintf("%v", *b) }

// internal boolFlag interface IsBoolFlag function for internal type boolValue
func (b *boolValue) IsBoolFlag() bool { return true }

type boolFlag interface {
	Value
	IsBoolFlag() bool
}

//
type BoolContain interface {
	SetBool(string, bool)
	ToBool(string) bool
}

func boolContainValue(key string, value bool, v BoolContain) *containValue {
	v.SetBool(key, value)
	return &containValue{
		key,
		func(n string) error {
			s, err := strconv.ParseBool(n)
			v.SetBool(key, s)
			return err
		},
		func() interface{} {
			return v.ToBool(key)
		},
		"bool",
	}
}

type intValue int

func newIntValue(val int, p *int) *intValue {
	*p = val
	return (*intValue)(p)
}

// Value interface Set function for internal type intValue
func (i *intValue) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, 64)
	*i = intValue(v)
	return err
}

// Value interface Get function for internal type intValue
func (i *intValue) Get() interface{} { return int(*i) }

// Value interface String function for internal type intValue
func (i *intValue) String() string { return fmt.Sprintf("%v", *i) }

//
type IntContain interface {
	SetInt(string, int)
	ToInt(string) int
}

func intContainValue(key string, value int, v IntContain) *containValue {
	v.SetInt(key, value)
	return &containValue{
		key,
		func(n string) error {
			s, err := strconv.ParseInt(n, 0, 64)
			v.SetInt(key, int(s))
			return err
		},
		func() interface{} {
			return v.ToInt(key)
		},
		"int",
	}
}

//
type int64Value int64

func newInt64Value(val int64, p *int64) *int64Value {
	*p = val
	return (*int64Value)(p)
}

// Value interface Set function for internal type int64Value
func (i *int64Value) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, 64)
	*i = int64Value(v)
	return err
}

// Value interface Get function for internal type int64Value
func (i *int64Value) Get() interface{} { return int64(*i) }

// Value interface String function for internal type intValue
func (i *int64Value) String() string { return fmt.Sprintf("%v", *i) }

//
type Int64Contain interface {
	SetInt64(string, int64)
	ToInt64(string) int64
}

func int64ContainValue(key string, value int64, v Int64Contain) *containValue {
	v.SetInt64(key, value)
	return &containValue{
		key,
		func(n string) error {
			s, err := strconv.ParseInt(n, 0, 64)
			v.SetInt64(key, s)
			return err
		},
		func() interface{} {
			return v.ToInt64(key)
		},
		"int64",
	}
}

type uintValue uint

func newUintValue(val uint, p *uint) *uintValue {
	*p = val
	return (*uintValue)(p)
}

// Value interface Set function for internal type uintValue
func (i *uintValue) Set(s string) error {
	v, err := strconv.ParseUint(s, 0, 64)
	*i = uintValue(v)
	return err
}

// Value interface Get function for internal type uintValue
func (i *uintValue) Get() interface{} { return uint(*i) }

// Value interface String function for internal type intValue
func (i *uintValue) String() string { return fmt.Sprintf("%v", *i) }

//
type UintContain interface {
	SetUint(string, uint)
	ToUint(string) uint
}

func uintContainValue(key string, value uint, v UintContain) *containValue {
	v.SetUint(key, value)
	return &containValue{
		key,
		func(n string) error {
			s, err := strconv.ParseUint(n, 0, 64)
			v.SetUint(key, uint(s))
			return err
		},
		func() interface{} {
			return v.ToUint(key)
		},
		"uint",
	}
}

type uint64Value uint64

func newUint64Value(val uint64, p *uint64) *uint64Value {
	*p = val
	return (*uint64Value)(p)
}

// Value interface Set function for internal type uint64Value
func (i *uint64Value) Set(s string) error {
	v, err := strconv.ParseUint(s, 0, 64)
	*i = uint64Value(v)
	return err
}

// Value interface Get function for internal type uint64Value
func (i *uint64Value) Get() interface{} { return uint64(*i) }

// Value interface String function for internal type intValue
func (i *uint64Value) String() string { return fmt.Sprintf("%v", *i) }

//
type Uint64Contain interface {
	SetUint64(string, uint64)
	ToUint64(string) uint64
}

func uint64ContainValue(key string, value uint64, v Uint64Contain) *containValue {
	v.SetUint64(key, value)
	return &containValue{
		key,
		func(n string) error {
			s, err := strconv.ParseUint(n, 0, 64)
			v.SetUint64(key, s)
			return err
		},
		func() interface{} {
			return v.ToUint64(key)
		},
		"uint64",
	}
}

type stringValue string

func newStringValue(val string, p *string) *stringValue {
	*p = val
	return (*stringValue)(p)
}

// Value interface Set function for internal type stringValue
func (s *stringValue) Set(val string) error {
	*s = stringValue(val)
	return nil
}

// Value interface Get function for internal type stringValue
func (s *stringValue) Get() interface{} { return string(*s) }

// Value interface String function for internal type stringValue
func (s *stringValue) String() string { return fmt.Sprintf("%s", *s) }

type StringContain interface {
	SetString(string, string)
	ToString(string) string
}

func stringContainValue(key, value string, v StringContain) *containValue {
	v.SetString(key, value)
	return &containValue{
		key,
		func(n string) error {
			v.SetString(key, n)
			return nil
		},
		func() interface{} {
			return v.ToString(key)
		},
		"string",
	}
}

type float64Value float64

func newFloat64Value(val float64, p *float64) *float64Value {
	*p = val
	return (*float64Value)(p)
}

// Value interface Set function for internal type float64Value
func (f *float64Value) Set(s string) error {
	v, err := strconv.ParseFloat(s, 64)
	*f = float64Value(v)
	return err
}

// Value interface Get function for internal type float64Value
func (f *float64Value) Get() interface{} { return float64(*f) }

// Value interface String function for internal type float64Value
func (f *float64Value) String() string { return fmt.Sprintf("%v", *f) }

type Float64Contain interface {
	SetFloat64(string, float64)
	ToFloat64(string) float64
}

func float64ContainValue(key string, value float64, v Float64Contain) *containValue {
	v.SetFloat64(key, value)
	return &containValue{
		key,
		func(n string) error {
			s, err := strconv.ParseFloat(n, 64)
			v.SetFloat64(key, s)
			return err
		},
		func() interface{} {
			return v.ToFloat64(key)
		},
		"float64",
	}
}

type durationValue time.Duration

func newDurationValue(val time.Duration, p *time.Duration) *durationValue {
	*p = val
	return (*durationValue)(p)
}

// Value interface Set function for internal type durationValue
func (d *durationValue) Set(s string) error {
	v, err := time.ParseDuration(s)
	*d = durationValue(v)
	return err
}

// Value interface Get function for internal type durationValue
func (d *durationValue) Get() interface{} { return time.Duration(*d) }

// Value interface String function for internal type durationValue
func (d *durationValue) String() string { return (*time.Duration)(d).String() }

func durationContainValue(key, value string, v StringContain) *containValue {
	v.SetString(key, value)
	return &containValue{
		key,
		func(n string) error {
			v.SetString(key, n)
			return nil
		},
		func() interface{} {
			s := v.ToString(key)
			d, err := time.ParseDuration(s)
			if err != nil {
				return err
			}
			return d
		},
		"duration",
	}
}

//
type RgxBridgeFunc func(string, ...*regexp.Regexp) error

type regexValue struct {
	raw string
	rgx []*regexp.Regexp
	xfn RgxBridgeFunc
}

func newRegexValue(xfn RgxBridgeFunc, raw ...string) *regexValue {
	var rgx []*regexp.Regexp
	for _, r := range raw {
		rgx = append(rgx, regexp.MustCompile(r))
	}
	return &regexValue{strings.Join(raw, ","), rgx, xfn}
}

//
func (r *regexValue) Set(s string) error {
	return r.xfn(s, r.rgx...)
}

//
func (r *regexValue) Get() interface{} {
	return r.raw
}

//
func (r *regexValue) String() string { return r.raw }

//
type RgxContainBridgeFunc func(string, StringContain, ...*regexp.Regexp) error

func regexContainValue(key string, xfn RgxContainBridgeFunc, v StringContain, raw ...string) *containValue {
	v.SetString(key, strings.Join(raw, ","))
	var rgx []*regexp.Regexp
	for _, r := range raw {
		rgx = append(rgx, regexp.MustCompile(r))
	}
	return &containValue{
		key,
		func(n string) error {
			return xfn(n, v, rgx...)
		},
		func() interface{} {
			return v.ToString(key)
		},
		"regex",
	}
}

// An integer type representing method for handling errors.
type ErrorHandling int

const (
	ContinueOnError ErrorHandling = iota // continue on error
	ExitOnError                          // exit on error
	PanicOnError                         // panic on error
)

// An interface encapsulating behavior for sets of flags.
type Flagger interface {
	GetterSetter
	Parser
	Stater
	Visiter
	Writer
	Wuser
}

// A type for handling any number of flags and fulfilling the Flagger interface.
type FlagSet struct {
	name          string
	parsed        bool
	actual        map[string]*Flag
	formal        map[string]*Flag
	args          []string
	errorHandling ErrorHandling
	output        io.Writer
}

// Creates a new *FlagSet with the string name and the provided error handling.
func NewFlagSet(name string, errorHandling ErrorHandling) *FlagSet {
	return &FlagSet{
		name:          name,
		errorHandling: errorHandling,
	}
}

// An interface for gettign & setting flags.
type GetterSetter interface {
	Lookup(string) *Flag
	Set(string, string) error
	Var(Value, string, string)
}

// Return a *Flag by the provided name, or nil if nothing is found.
func (f *FlagSet) Lookup(name string) *Flag {
	if fl, ok := f.formal[name]; ok {
		return fl
	}
	return nil
}

// Sets a flag by string name and value, returning an error.
func (f *FlagSet) Set(name, value string) error {
	flag, ok := f.formal[name]
	if !ok {
		return fmt.Errorf("no such flag -%v", name)
	}
	err := flag.Value.Set(value)
	if err != nil {
		return err
	}
	if f.actual == nil {
		f.actual = make(map[string]*Flag)
	}
	f.actual[name] = flag
	return nil
}

// An interface handling flag parsing from a string slice, and returning details
// of parsing status.
type Parser interface {
	Parse([]string) error
	Parsed() bool
}

// *FlagSet function satisfying the Parser interface Parse function.
func (f *FlagSet) Parse(arguments []string) error {
	f.parsed = true
	f.args = arguments
	for {
		seen, err := f.parseOne()
		if seen {
			continue
		}
		if err == nil {
			break
		}
		switch f.errorHandling {
		case ContinueOnError:
			return err
		case ExitOnError:
			os.Exit(2)
		case PanicOnError:
			panic(err)
		}
	}
	return nil
}

func (f *FlagSet) parseOne() (bool, error) {
	if len(f.args) == 0 {
		return false, nil
	}
	s := f.args[0]
	if len(s) == 0 || s[0] != '-' || len(s) == 1 {
		return false, nil
	}
	numMinuses := 1
	if s[1] == '-' {
		numMinuses++
		if len(s) == 2 { // "--" terminates the flags
			f.args = f.args[1:]
			return false, nil
		}
	}
	name := s[numMinuses:]
	if len(name) == 0 || name[0] == '-' || name[0] == '=' {
		return false, failFmt(f, "bad flag syntax: %s", s)
	}

	// it's a flag. does it have an argument?
	f.args = f.args[1:]
	hasValue := false
	value := ""
	for i := 1; i < len(name); i++ { // equals cannot be first
		if name[i] == '=' {
			value = name[i+1:]
			hasValue = true
			name = name[0:i]
			break
		}
	}

	m := f.formal
	var flag *Flag
	var exists bool
	flag, exists = m[name]
	if !exists {
		return false, failOnly(f, "flag provided but not defined: -%s\n", name)
	}

	if fv, ok := flag.Value.(boolFlag); ok && fv.IsBoolFlag() { // special case: doesn't need an arg
		if hasValue {
			if err := fv.Set(value); err != nil {
				return false, failFmt(f, "invalid boolean value %q for -%s: %v", value, name, err)
			}
		} else {
			if err := fv.Set("true"); err != nil {
				return false, failFmt(f, "invalid boolean flag %s: %v", name, err)
			}
		}
	} else {
		// It must have a value, which might be the next argument.
		if !hasValue && len(f.args) > 0 {
			// value is the next arg
			hasValue = true
			value, f.args = f.args[0], f.args[1:]
		}
		if !hasValue {
			return false, failFmt(f, "flag needs an argument: -%s", name)
		}
		if err := flag.Value.Set(value); err != nil {
			return false, failFmt(f, "invalid value %q for flag -%s: %v", value, name, err)
		}
	}
	if f.actual == nil {
		f.actual = make(map[string]*Flag)
	}
	f.actual[name] = flag
	return true, nil
}

//  *FlagSet function satisfying the Parser interface Parsed function.
func (f *FlagSet) Parsed() bool {
	return f.parsed
}

// A package level writer interface.
type Writer interface {
	Out() io.Writer
	SetOut(io.Writer)
}

// Returns an io.Writer from the *FlagSet.
func (f *FlagSet) Out() io.Writer {
	if f.output == nil {
		return os.Stderr
	}
	return f.output
}

// Sets the provided io.Writer to the *FlagSet.
func (f *FlagSet) SetOut(output io.Writer) {
	f.output = output
}

// An interface to manage sending a function to all flags.
type Visiter interface {
	Visit(func(*Flag))
	VisitAll(func(*Flag))
}

// A *FlagSet VisitAll function satisfying the Visiter interface.
func (f *FlagSet) VisitAll(fn func(*Flag)) {
	for _, flag := range sortFlags(f.formal) {
		fn(flag)
	}
}

// A *FlagSet Visit function satisfying the Visiter interface.
func (f *FlagSet) Visit(fn func(*Flag)) {
	for _, flag := range sortFlags(f.actual) {
		fn(flag)
	}
}

func sortFlags(flags map[string]*Flag) []*Flag {
	list := make(sort.StringSlice, len(flags))
	i := 0
	for _, f := range flags {
		list[i] = f.Name
		i++
	}
	list.Sort()
	result := make([]*Flag, len(list))
	for i, name := range list {
		result[i] = flags[name]
	}
	return result
}

// Sets a Value, name string & usage string as a *Flag to be used by the *FlagSet
// This will panic for duplicate and/or  previously defined Flags.
func (f *FlagSet) Var(value Value, name string, usage string) {
	// Remember the default value as a string; it won't change.
	flag := &Flag{name, usage, value, value.String()}
	_, alreadythere := f.formal[name]
	if alreadythere {
		msg := fmt.Sprintf("%s flag redefined: %s", f.name, name)
		fmt.Fprintln(f.Out(), msg)
		panic(msg) // Happens only if flags are declared with identical names
	}
	if f.formal == nil {
		f.formal = make(map[string]*Flag)
	}
	f.formal[name] = flag
}

//
func (f *FlagSet) BoolVar(p *bool, name string, value bool, usage string) {
	f.Var(newBoolValue(value, p), name, usage)
}

//
func (f *FlagSet) Bool(name string, value bool, usage string) *bool {
	p := new(bool)
	f.BoolVar(p, name, value, usage)
	return p
}

//
func (f *FlagSet) BoolContainVar(d BoolContain, name, key string, value bool, usage string) {
	f.Var(boolContainValue(key, value, d), name, usage)
}

//
func (f *FlagSet) BoolContain(d BoolContain, name, key, usage string) BoolContain {
	f.BoolContainVar(d, name, key, false, usage)
	return d
}

//
func (f *FlagSet) IntVar(p *int, name string, value int, usage string) {
	f.Var(newIntValue(value, p), name, usage)
}

//
func (f *FlagSet) Int(name string, value int, usage string) *int {
	p := new(int)
	f.IntVar(p, name, value, usage)
	return p
}

//
func (f *FlagSet) IntContainVar(d IntContain, name, key string, value int, usage string) {
	f.Var(intContainValue(key, value, d), name, usage)
}

//
func (f *FlagSet) IntContain(d IntContain, name, key, usage string) IntContain {
	f.IntContainVar(d, name, key, 0, usage)
	return d
}

//
func (f *FlagSet) Int64Var(p *int64, name string, value int64, usage string) {
	f.Var(newInt64Value(value, p), name, usage)
}

//
func (f *FlagSet) Int64(name string, value int64, usage string) *int64 {
	p := new(int64)
	f.Int64Var(p, name, value, usage)
	return p
}

//
func (f *FlagSet) Int64ContainVar(d Int64Contain, name, key string, value int64, usage string) {
	f.Var(int64ContainValue(key, value, d), name, usage)
}

//
func (f *FlagSet) Int64Contain(d Int64Contain, name, key, usage string) Int64Contain {
	f.Int64ContainVar(d, name, key, 0, usage)
	return d
}

//
func (f *FlagSet) UintVar(p *uint, name string, value uint, usage string) {
	f.Var(newUintValue(value, p), name, usage)
}

//
func (f *FlagSet) Uint(name string, value uint, usage string) *uint {
	p := new(uint)
	f.UintVar(p, name, value, usage)
	return p
}

//
func (f *FlagSet) UintContainVar(d UintContain, name, key string, value uint, usage string) {
	f.Var(uintContainValue(key, value, d), name, usage)
}

//
func (f *FlagSet) UintContain(d UintContain, name, key, usage string) UintContain {
	f.UintContainVar(d, name, key, 0, usage)
	return d
}

//
func (f *FlagSet) Uint64Var(p *uint64, name string, value uint64, usage string) {
	f.Var(newUint64Value(value, p), name, usage)
}

//
func (f *FlagSet) Uint64(name string, value uint64, usage string) *uint64 {
	p := new(uint64)
	f.Uint64Var(p, name, value, usage)
	return p
}

//
func (f *FlagSet) Uint64ContainVar(d Uint64Contain, name, key string, value uint64, usage string) {
	f.Var(uint64ContainValue(key, value, d), name, usage)
}

//
func (f *FlagSet) Uint64Contain(d Uint64Contain, name, key, usage string) Uint64Contain {
	f.Uint64ContainVar(d, name, key, 0, usage)
	return d
}

//
func (f *FlagSet) StringVar(p *string, name string, value string, usage string) {
	f.Var(newStringValue(value, p), name, usage)
}

//
func (f *FlagSet) String(name string, value string, usage string) *string {
	p := new(string)
	f.StringVar(p, name, value, usage)
	return p
}

//
func (f *FlagSet) StringContainVar(d StringContain, name, key, value, usage string) {
	f.Var(stringContainValue(key, value, d), name, usage)
}

//
func (f *FlagSet) StringContain(d StringContain, name, key, usage string) StringContain {
	f.StringContainVar(d, name, key, "", usage)
	return d
}

//
func (f *FlagSet) Float64Var(p *float64, name string, value float64, usage string) {
	f.Var(newFloat64Value(value, p), name, usage)
}

//
func (f *FlagSet) Float64(name string, value float64, usage string) *float64 {
	p := new(float64)
	f.Float64Var(p, name, value, usage)
	return p
}

//
func (f *FlagSet) Float64ContainVar(d Float64Contain, name, key string, value float64, usage string) {
	f.Var(float64ContainValue(key, value, d), name, usage)
}

//
func (f *FlagSet) Float64Contain(d Float64Contain, name, key, usage string) Float64Contain {
	f.Float64ContainVar(d, name, key, 0, usage)
	return d
}

//
func (f *FlagSet) DurationVar(p *time.Duration, name string, value time.Duration, usage string) {
	f.Var(newDurationValue(value, p), name, usage)
}

//
func (f *FlagSet) Duration(name string, value time.Duration, usage string) *time.Duration {
	p := new(time.Duration)
	f.DurationVar(p, name, value, usage)
	return p
}

//
func (f *FlagSet) DurationContainVar(d StringContain, name, key, value, usage string) {
	f.Var(durationContainValue(key, value, d), name, usage)
}

//
func (f *FlagSet) DurationContain(d StringContain, name, key, usage string) StringContain {
	f.DurationContainVar(d, name, key, "0s", usage)
	return d
}

// A flag processing a regular expression
func (f *FlagSet) RegexVar(name, usage string, xfn RgxBridgeFunc, rawRegexps ...string) {
	f.Var(newRegexValue(xfn, rawRegexps...), name, usage)
}

// A contain backed flag processing a regular expression
func (f *FlagSet) RegexContainVar(d StringContain, name, key, usage string, xfn RgxContainBridgeFunc, rawRegexps ...string) StringContain {
	f.Var(regexContainValue(key, xfn, d, rawRegexps...), name, usage)
	return d
}

//
type Stater interface {
	NFlag() int
	NArg() int
	Arg(int) string
	Args() []string
}

//
func (f *FlagSet) NFlag() int { return len(f.actual) }

//
func (f *FlagSet) Arg(i int) string {
	if i < 0 || i >= len(f.args) {
		return ""
	}
	return f.args[i]
}

//
func (f *FlagSet) NArg() int { return len(f.args) }

//
func (f *FlagSet) Args() []string { return f.args }

// An interface calling usage strings of a flagset.
type Wuser interface {
	Usage(io.Writer)
}

//
func (f *FlagSet) Usage(o io.Writer) {
	f.VisitAll(func(flag *Flag) {
		s := fmt.Sprintf("\t-%s", flag.Name) // Two spaces before -; see next two comments.
		name, usage := UnquoteMessage(flag)
		if len(name) > 0 {
			s += " " + name
		}
		// Boolean flags of one ASCII letter are so common we
		// treat them specially, putting their usage on the same line.
		if len(s) <= 4 { // space, space, '-', 'x'.
			s += "\t"
		} else {
			// Four spaces before the tab triggers good alignment
			// for both 4- and 8-space tab stops.
			s += "\n    \t"
		}
		s += fmt.Sprintf("\t%s", usage)
		if !isZeroValue(flag.DefValue) {
			if _, ok := flag.Value.(*stringValue); ok {
				// put quotes on the value
				s += fmt.Sprintf(" (default %q)", flag.DefValue)
			} else {
				s += fmt.Sprintf(" (default %v)", flag.DefValue)
			}
		}
		white(o, s, "\n")
	})
}

//
func UnquoteMessage(flag *Flag) (name string, usage string) {
	// Look for a back-quoted name, but avoid the strings package.
	usage = flag.Message
	for i := 0; i < len(usage); i++ {
		if usage[i] == '`' {
			for j := i + 1; j < len(usage); j++ {
				if usage[j] == '`' {
					name = usage[i+1 : j]
					usage = usage[:i] + name + usage[j+1:]
					return name, usage
				}
			}
			break // Only one back quote; use type name.
		}
	}
	// No explicit name, so use type if we can find one.
	name = "value"
	switch flag.Value.(type) {
	case *boolValue:
		name = ""
	case *durationValue:
		name = "duration"
	case *float64Value:
		name = "float"
	case *intValue, *int64Value:
		name = "int"
	case *stringValue, *regexValue:
		name = "string"
	case *uintValue, *uint64Value:
		name = "uint"
	case *containValue:
		vv := flag.Value.(*containValue)
		switch vv.kind {
		case "bool":
			name = ""
		case "duration":
			name = "duration"
		case "float64":
			name = "float"
		case "int", "int64":
			name = "int"
		case "string", "regex":
			name = "string"
		case "uint", "uint64":
			name = "uint"
		}
	}
	return
}

func isZeroValue(value string) bool {
	switch value {
	case "false", "", "0", "0s":
		return true
	}
	return false
}
