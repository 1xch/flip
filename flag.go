package flip

import (
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/Laughs-In-Flowers/data"
)

//
type Flag struct {
	Name     string // name as it appears on command line
	Message  string // help message
	Value    Value  // value as set
	DefValue string // default value (as text); for usage message
}

//
type Value interface {
	String() string
	Set(string) error
	Get() interface{}
}

type (
	setFn func(string) error
	getFn func() interface{}
	//rgxToVector func(string, *regexp.Regexp, data.Vector) error
)

type vectorValue struct {
	to   string
	sfn  setFn
	gfn  getFn
	kind string
}

//
func (v *vectorValue) Set(n string) error {
	return v.sfn(n)
}

//
func (v *vectorValue) Get() interface{} {
	return v.gfn()
}

//
func (v *vectorValue) String() string {
	return fmt.Sprintf("%v", v.Get())
}

//
func (v *vectorValue) IsBoolFlag() bool {
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

//
func (b *boolValue) Set(s string) error {
	v, err := strconv.ParseBool(s)
	*b = boolValue(v)
	return err
}

//
func (b *boolValue) Get() interface{} { return bool(*b) }

//
func (b *boolValue) String() string { return fmt.Sprintf("%v", *b) }

//
func (b *boolValue) IsBoolFlag() bool { return true }

type boolFlag interface {
	Value
	IsBoolFlag() bool
}

func boolVectorValue(key string, value bool, v *data.Vector) *vectorValue {
	v.SetBool(key, value)
	return &vectorValue{
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

//
func (i *intValue) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, 64)
	*i = intValue(v)
	return err
}

//
func (i *intValue) Get() interface{} { return int(*i) }

//
func (i *intValue) String() string { return fmt.Sprintf("%v", *i) }

func intVectorValue(key string, value int, v *data.Vector) *vectorValue {
	v.SetInt(key, value)
	return &vectorValue{
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

//func intVectorRxValue(key, rx string, defaultI int, v *data.Vector) *vectorValue {
//	v.SetInt(key, defaultI)
//	rxp := regexp.MustCompile(rx)
//	fn := func(s string, r *regexp.Regexp, v *data.Vector) error {
//		fs := rxp.FindString(s)
//		i := paramInt(r, fs, key)
//		v.SetInt(key, i)
//		return nil
//	}
//	return &vectorValue{
//		key,
//		func(n string) error {
//			return fn(n, rxp, v)
//		},
//		func() interface{} {
//			return v.ToInt(key)
//		},
//	}
//}

//
type int64Value int64

func newInt64Value(val int64, p *int64) *int64Value {
	*p = val
	return (*int64Value)(p)
}

//
func (i *int64Value) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, 64)
	*i = int64Value(v)
	return err
}

//
func (i *int64Value) Get() interface{} { return int64(*i) }

//
func (i *int64Value) String() string { return fmt.Sprintf("%v", *i) }

func int64VectorValue(key string, value int64, v *data.Vector) *vectorValue {
	v.SetInt64(key, value)
	return &vectorValue{
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

//
func (i *uintValue) Set(s string) error {
	v, err := strconv.ParseUint(s, 0, 64)
	*i = uintValue(v)
	return err
}

//
func (i *uintValue) Get() interface{} { return uint(*i) }

//
func (i *uintValue) String() string { return fmt.Sprintf("%v", *i) }

func uintVectorValue(key string, value uint, v *data.Vector) *vectorValue {
	v.SetUint(key, value)
	return &vectorValue{
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

//
func (i *uint64Value) Set(s string) error {
	v, err := strconv.ParseUint(s, 0, 64)
	*i = uint64Value(v)
	return err
}

//
func (i *uint64Value) Get() interface{} { return uint64(*i) }

//
func (i *uint64Value) String() string { return fmt.Sprintf("%v", *i) }

func uint64VectorValue(key string, value uint64, v *data.Vector) *vectorValue {
	v.SetUint64(key, value)
	return &vectorValue{
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

//
func (s *stringValue) Set(val string) error {
	*s = stringValue(val)
	return nil
}

//
func (s *stringValue) Get() interface{} { return string(*s) }

//
func (s *stringValue) String() string { return fmt.Sprintf("%s", *s) }

func stringVectorValue(key, value string, v *data.Vector) *vectorValue {
	v.SetString(key, value)
	return &vectorValue{
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

//
func (f *float64Value) Set(s string) error {
	v, err := strconv.ParseFloat(s, 64)
	*f = float64Value(v)
	return err
}

//
func (f *float64Value) Get() interface{} { return float64(*f) }

//
func (f *float64Value) String() string { return fmt.Sprintf("%v", *f) }

func float64VectorValue(key string, value float64, v *data.Vector) *vectorValue {
	v.SetFloat64(key, value)
	return &vectorValue{
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

//
func (d *durationValue) Set(s string) error {
	v, err := time.ParseDuration(s)
	*d = durationValue(v)
	return err
}

//
func (d *durationValue) Get() interface{} { return time.Duration(*d) }

//
func (d *durationValue) String() string { return (*time.Duration)(d).String() }

//
type ErrorHandling int

const (
	ContinueOnError ErrorHandling = iota // continue on error
	ExitOnError                          // exit on error
	PanicOnError                         // panic on error
)

//
type Flagger interface {
	Parser
	Setter
	Stater
	Visiter
	Writer
	Wuser
}

//
type Setter interface {
	Lookup(string) *Flag
	Set(string, string) error
	Var(Value, string, string)
}

//
type Visiter interface {
	Visit(func(*Flag))
	VisitAll(func(*Flag))
}

//
type Writer interface {
	Out() io.Writer
	SetOut(io.Writer)
}

//
type Wuser interface {
	Usage(io.Writer)
}

//
type FlagSet struct {
	name          string
	parsed        bool
	actual        map[string]*Flag
	formal        map[string]*Flag
	args          []string
	errorHandling ErrorHandling
	output        io.Writer
}

//
func NewFlagSet(name string, errorHandling ErrorHandling) *FlagSet {
	return &FlagSet{
		name:          name,
		errorHandling: errorHandling,
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

//
func (f *FlagSet) Out() io.Writer {
	if f.output == nil {
		return os.Stderr
	}
	return f.output
}

//
func (f *FlagSet) SetOut(output io.Writer) {
	f.output = output
}

//
func (f *FlagSet) VisitAll(fn func(*Flag)) {
	for _, flag := range sortFlags(f.formal) {
		fn(flag)
	}
}

//
func (f *FlagSet) Visit(fn func(*Flag)) {
	for _, flag := range sortFlags(f.actual) {
		fn(flag)
	}
}

//
func (f *FlagSet) Lookup(name string) *Flag {
	if fl, ok := f.formal[name]; ok {
		return fl
	}
	return nil
}

//
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

func isZeroValue(value string) bool {
	switch value {
	case "false", "", "0", "1s":
		return true
	}
	return false
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
	case *stringValue:
		name = "string"
	case *uintValue, *uint64Value:
		name = "uint"
	case *vectorValue:
		vv := flag.Value.(*vectorValue)
		switch vv.kind {
		case "bool":
			name = ""
		case "float64":
			name = "float"
		case "int", "int64":
			name = "int"
		case "string":
			name = "string"
		case "uint", "uint64":
			name = "uint"
		}
	}
	return
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
func (f *FlagSet) BoolVectorVar(d *data.Vector, name, key string, value bool, usage string) {
	f.Var(boolVectorValue(key, value, d), name, usage)
}

//
func (f *FlagSet) BoolVector(d *data.Vector, name, key, usage string) *data.Vector {
	f.BoolVectorVar(d, name, key, false, usage)
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
func (f *FlagSet) IntVectorVar(d *data.Vector, name, key string, value int, usage string) {
	f.Var(intVectorValue(key, value, d), name, usage)
}

//
func (f *FlagSet) IntVector(d *data.Vector, name, key, usage string) *data.Vector {
	f.IntVectorVar(d, name, key, 0, usage)
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
func (f *FlagSet) Int64VectorVar(d *data.Vector, name, key string, value int64, usage string) {
	f.Var(int64VectorValue(key, value, d), name, usage)
}

//
func (f *FlagSet) Int64Vector(d *data.Vector, name, key, usage string) *data.Vector {
	f.Int64VectorVar(d, name, key, 0, usage)
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
func (f *FlagSet) UintVectorVar(d *data.Vector, name, key string, value uint, usage string) {
	f.Var(uintVectorValue(key, value, d), name, usage)
}

//
func (f *FlagSet) UintVector(d *data.Vector, name, key, usage string) *data.Vector {
	f.UintVectorVar(d, name, key, 0, usage)
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
func (f *FlagSet) Uint64VectorVar(d *data.Vector, name, key string, value uint64, usage string) {
	f.Var(uint64VectorValue(key, value, d), name, usage)
}

//
func (f *FlagSet) Uint64Vector(d *data.Vector, name, key, usage string) *data.Vector {
	f.Uint64VectorVar(d, name, key, 0, usage)
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
func (f *FlagSet) StringVectorVar(d *data.Vector, name, key, value, usage string) {
	f.Var(stringVectorValue(key, value, d), name, usage)
}

//
func (f *FlagSet) StringVector(d *data.Vector, name, key, usage string) *data.Vector {
	f.StringVectorVar(d, name, key, "", usage)
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
func (f *FlagSet) Float64VectorVar(d *data.Vector, name, key string, value float64, usage string) {
	f.Var(float64VectorValue(key, value, d), name, usage)
}

//
func (f *FlagSet) Float64Vector(d *data.Vector, name, key, usage string) *data.Vector {
	f.Float64VectorVar(d, name, key, 0, usage)
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

// duration vector

//
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

func (f *FlagSet) failOnly(format string, a ...interface{}) error {
	err := fmt.Errorf(format, a...)
	fmt.Fprintln(f.Out(), err)
	return err
}

func (f *FlagSet) failf(format string, a ...interface{}) error {
	err := fmt.Errorf(format, a...)
	fmt.Fprintln(f.Out(), err)
	f.Usage(f.Out())
	return err
}

//
type Parser interface {
	Parse([]string) error
	Parsed() bool
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
		return false, f.failf("bad flag syntax: %s", s)
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
		return false, f.failOnly("flag provided but not defined: -%s\n", name)
	}

	if fv, ok := flag.Value.(boolFlag); ok && fv.IsBoolFlag() { // special case: doesn't need an arg
		if hasValue {
			if err := fv.Set(value); err != nil {
				return false, f.failf("invalid boolean value %q for -%s: %v", value, name, err)
			}
		} else {
			if err := fv.Set("true"); err != nil {
				return false, f.failf("invalid boolean flag %s: %v", name, err)
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
			return false, f.failf("flag needs an argument: -%s", name)
		}
		if err := flag.Value.Set(value); err != nil {
			return false, f.failf("invalid value %q for flag -%s: %v", value, name, err)
		}
	}
	if f.actual == nil {
		f.actual = make(map[string]*Flag)
	}
	f.actual[name] = flag
	return true, nil
}

//
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

func (f *FlagSet) Parsed() bool {
	return f.parsed
}

var errHelp = errors.New("flag: help requested")
