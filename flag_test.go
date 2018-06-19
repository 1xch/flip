package flip

import (
	"bytes"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/Laughs-In-Flowers/data"
)

var flagExpect []fexp

type fexp struct {
	mthd  string
	key   string
	vals  []interface{}
	exp   interface{}
	prs   [][]string
	usg   []string
	catch bool
}

var vec *data.Vector

func init() {
	vec = data.New("test_vector")

	flagExpect = []fexp{
		{
			"Bool", "b",
			[]interface{}{"b", false, "A boolean flag"},
			true,
			[][]string{{"-b"}, {"--b"}},
			[]string{"-b", "A boolean flag"},
			false,
		},
		{
			"BoolVector", "bv",
			[]interface{}{vec, "bv", "bvKey", "A vector backed boolean flag"},
			true,
			[][]string{{"-bv"}, {"--bv"}},
			[]string{"-bv", "A vector backed boolean flag"},
			false,
		},
		{
			"Int", "i",
			[]interface{}{"i", 0, "An integer flag"},
			500,
			[][]string{{"-i", "500"}},
			[]string{"-i int", "An integer flag"},
			false,
		},
		{
			"IntVector", "iv",
			[]interface{}{vec, "iv", "ivKey", "A vector backed integer flag"},
			500,
			[][]string{{"-iv", "500"}},
			[]string{"-iv int", "A vector backed integer flag"},
			false,
		},
		{
			"Int64", "i64",
			[]interface{}{"i64", int64(0), "An integer64 flag"},
			int64(500),
			[][]string{{"-i64", "500"}},
			[]string{"-i64 int", "An integer64 flag"},
			false,
		},
		{
			"Int64Vector", "iv64",
			[]interface{}{vec, "iv64", "iv64Key", "A vector backed int64 flag"},
			int64(500),
			[][]string{{"-iv64", "500"}},
			[]string{"-iv64 int", "A vector backed int64 flag"},
			false,
		},
		{
			"Uint", "u",
			[]interface{}{"u", uint(0), "A uint flag"},
			uint(500),
			[][]string{{"-u", "500"}},
			[]string{"-u uint", "A uint flag"},
			false,
		},
		{
			"UintVector", "uv",
			[]interface{}{vec, "uv", "uvKey", "A vector backed uint flag"},
			uint(500),
			[][]string{{"-uv", "500"}},
			[]string{"-uv uint", "A vector backed uint flag"},
			false,
		},
		{
			"Uint64", "u64",
			[]interface{}{"u64", uint64(0), "A uint64 flag"},
			uint64(500),
			[][]string{{"-u64", "500"}},
			[]string{"-u64 uint", "A uint64 flag"},
			false,
		},
		{
			"Uint64Vector", "u64v",
			[]interface{}{vec, "u64v", "u64vKey", "A vector backed uint64 flag"},
			uint64(500),
			[][]string{{"-u64v", "500"}},
			[]string{"-u64v uint", "A vector backed uint64 flag"},
			false,
		},
		{
			"String", "s",
			[]interface{}{"s", string(""), "A string flag"},
			"hello",
			[][]string{{"-s", "hello"}, {"--s", "hello"}},
			[]string{"-s string", "A string flag"},
			false,
		},
		{
			"StringVector", "sv",
			[]interface{}{vec, "sv", "svKey", "A vector backed string flag"},
			"hello",
			[][]string{{"-sv", "hello"}},
			[]string{"-sv string", "A vector backed string flag"},
			false,
		},
		{
			"Float64", "f64",
			[]interface{}{"f64", float64(0.0), "A float64 flag `FLOAT64`"},
			float64(500.0),
			[][]string{{"-f64", "500.0"}},
			[]string{"-f64 FLOAT64", "A float64 flag"},
			false,
		},
		{
			"Float64Vector", "f64v",
			[]interface{}{vec, "f64v", "f64vKey", "A vector backed float64 flag"},
			float64(500.0),
			[][]string{{"-f64v", "500.0"}, {"-f64v", "500"}},
			[]string{"-f64v float", "A vector backed float64 flag"},
			false,
		},
		{
			"Duration", "d",
			[]interface{}{"d", time.Second * 1, "A duration flag"},
			time.Second * 500,
			[][]string{{"-d", "500s"}},
			[]string{"-d duration", "A duration flag"},
			false,
		},
		{
			"Float64", "f64",
			[]interface{}{"f64", float64(0.0), "Catching float64 flag errors"},
			float64(500.0),
			[][]string{{"-f64"}},
			[]string{"-f64 float", "Catching float64 flag errors"},
			true,
		},
		{
			"Float64", "f64",
			[]interface{}{"f64", float64(0.0), "Catching float64 flag errors"},
			float64(500.0),
			[][]string{{"-f64", "red"}},
			[]string{"-f64 float", "Catching float64 flag errors"},
			true,
		},
		{
			"String", "s",
			[]interface{}{"s", string(""), "Catching string flag errors"},
			"",
			[][]string{{"---s", "hello"}},
			[]string{"-s string", "Catching string flag errors"},
			true,
		},
		{
			"Bool", "b",
			[]interface{}{"b", false, "Catching boolean flag errors"},
			true,
			[][]string{{"-------b"}},
			[]string{"-b", "Catching boolean flag errors"},
			true,
		},
		{
			"Bool", "b",
			[]interface{}{"b", false, "Catching boolean flag errors"},
			true,
			[][]string{{"-b=red"}},
			[]string{"-b", "Catching boolean flag errors"},
			true,
		},
	}
}

func countVisit(t *testing.T, fs *FlagSet, v int) {
	c := 0
	fs.Visit(func(ff *Flag) {
		c = c + 1
	})
	if c < v || c > v {
		t.Errorf("visit counted less than or more than %d flags", v)
	}
}

func countVisitAll(t *testing.T, fs *FlagSet, v int) {
	c := 0
	fs.VisitAll(func(ff *Flag) {
		c = c + 1
	})
	if c < v || c > v {
		t.Errorf("visitall counted less than or more than %d flags", v)
	}
}

func testOneFlagExpectation(t *testing.T, x fexp) {
	for _, p := range x.prs {
		b := new(bytes.Buffer)
		if x.catch {
			defer catchShouldPanic(t, x.mthd, b)
		}
		fs := NewFlagSet("test_flagset", PanicOnError)
		fs.SetOut(b)
		fsMethod(t, false, fs, x.mthd, x.vals...)
		fsMethod(t, true, fs, x.mthd, x.vals...)
		fsMethod(t, false, fs, "String", []interface{}{"xs", string("do not visit"), "A string flag not visited"}...)
		fs.Parse(p)
		if prsd := fs.Parsed(); !prsd {
			t.Error("flagset not parsed, but should be")
		}
		f := fs.Lookup(x.key)
		v := f.Value.Get()
		if v != x.exp {
			t.Errorf("flag error: '%s' expected %v(%T), got %v(%T)", x.mthd, x.exp, x.exp, v, v)
		}
		fs.Usage(b)
		usg := b.String()
		for _, u := range x.usg {
			if !strings.Contains(usg, u) {
				t.Errorf("flag usage error: '%s' expected %s in %s", x.mthd, u, usg)
			}
		}
		countVisit(t, fs, 1)
		countVisitAll(t, fs, 2)
		if noFlag := fs.Lookup("noFlag"); noFlag != nil {
			t.Error("Lookup flag returned a flag where none should be")
		}
		fs.Set("xs", "now visited")
		countVisit(t, fs, 2)

		if fs.NFlag() != 2 {
			t.Error("FlagSet NFlag != 2")
		}
		if fs.Arg(0) != "" {
			t.Error("FlagSet Arg(0) != ''")
		}
		if fs.NArg() != 0 {
			t.Error("FlagSet NArg() != 0")
		}
		if len(fs.Args()) != 0 {
			t.Error("FlagSet len(Args()) != 0")
		}
	}
}

func TestFlagSet(t *testing.T) {
	for _, x := range flagExpect {
		testOneFlagExpectation(t, x)
	}
}

func catchShouldPanic(t *testing.T, m string, b *bytes.Buffer) {
	r := recover()
	if r == nil {
		t.Errorf("'%s:%s' should force exit in this situation, but did not\n\t\trecovery message: %s", m, b.String(), r)
	}
}

func fsMethod(t *testing.T, catch bool, fs *FlagSet, m string, in ...interface{}) {
	if catch {
		b := new(bytes.Buffer)
		b.WriteString("attempt to set existing flag")
		defer catchShouldPanic(t, m, b)
	}
	method := reflect.ValueOf(fs).MethodByName(m)
	params := make([]reflect.Value, method.Type().NumIn())
	for i := 0; i < method.Type().NumIn(); i++ {
		object := in[i]
		params[i] = reflect.ValueOf(object)
	}
	method.Call(params)
}
