package flip

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

type tflags struct {
	t1, b1, b2 bool
}

func flagsExpected(f *tflags, t *testing.T, fn func()) {}

func testFlagSet(label string, tf *tflags, b *bytes.Buffer) *FlagSet {
	fs := NewFlagSet(label, ContinueOnError)
	fs.BoolVar(&tf.t1, "t1", tf.t1, "boolean flag top")
	fs.BoolVar(&tf.b1, "b1", tf.b1, "boolean flag one")
	fs.BoolVar(&tf.b2, "b2", tf.b2, "boolean flag two")
	fs.SetOut(b)
	return fs
}

func cmdSet(tf *tflags, b *bytes.Buffer) [][]Command {
	tg0 := []Command{
		NewCommand(
			"", "testing", "top level package flags",
			1,
			false,
			func(c context.Context, s []string) (context.Context, ExitStatus) {
				return c, ExitNo
			},
			testFlagSet("testing", tf, b),
		),
	}

	tg1 := []Command{
		NewCommand(
			"one", "one-A", "command one-A",
			1,
			false,
			func(c context.Context, s []string) (context.Context, ExitStatus) {
				return c, ExitNo
			},
			testFlagSet("one-A", tf, b),
		),
		NewCommand(
			"one", "one-B", "command one-B",
			2,
			false,
			func(c context.Context, s []string) (context.Context, ExitStatus) {
				return c, ExitSuccess
			},
			testFlagSet("one-B", tf, b),
		),
		// nil cmdFunc command
	}

	tg2 := []Command{
		NewCommand(
			"two", "two-A", "command two-A",
			1,
			false,
			func(c context.Context, s []string) (context.Context, ExitStatus) {
				return c, ExitSuccess
			},
			testFlagSet("two-A", tf, b),
		),
		NewCommand(
			"two", "two-B", "command two-B",
			2,
			false,
			func(c context.Context, s []string) (context.Context, ExitStatus) {
				return c, ExitSuccess
			},
			testFlagSet("two-B", tf, b),
		),
		NewCommand(
			"two", "two-C", "command two-C",
			2,
			false,
			func(c context.Context, s []string) (context.Context, ExitStatus) {
				return c, ExitFailure
			},
			testFlagSet("two-C", tf, b),
		),
	}
	return [][]Command{tg1, tg2, tg0}
}

type flagsetTestFunc func(t *testing.T, fs *tflags)

func flagError(t *testing.T, tag string, have, expect bool) {
	if have != expect {
		t.Errorf("%s boolean glag error have %t, expected %t", tag, have, expect)
	}
}

var AnyExitMessage = "cleanup function on any exit code"

var flipExpect = []struct {
	expectExit   int
	fsFunc       flagsetTestFunc
	expectHelp   []string
	unexpectHelp []string
	cmd          []string
}{
	{
		-2,
		nil,
		[]string{"test [OPTIONS...] {COMMAND} ..."},
		nil,
		[]string{"testing"},
	},
	{
		-2,
		func(t *testing.T, fs *tflags) {
			flagError(t, "top", fs.t1, true)
		},
		[]string{"test [OPTIONS...] {COMMAND} ..."},
		nil,
		[]string{"testing", "-t1"},
	},
	{
		-2,
		nil,
		nil,
		nil,
		[]string{"testing", "one-A"},
	},
	{
		-2,
		func(t *testing.T, fs *tflags) {
			//flagError(t, "top", fs.t1, true)
		},
		[]string{"flag provided but not defined: -nonflag"},
		nil,
		[]string{"testing", "one-A", "-b1", "-nonflag"},
	},
	{
		0,
		func(t *testing.T, fs *tflags) {
			flagError(t, "one-B", fs.b1, false)
			flagError(t, "one-B", fs.b2, true)
		},
		nil,
		nil,
		[]string{"testing", "-t1", "one-B", "-b2"},
	},
	{
		0,
		func(t *testing.T, fs *tflags) {
			flagError(t, "one-A+one-B", fs.b1, true)
			flagError(t, "one-A+one-B", fs.b2, true)
		},
		nil,
		nil,
		[]string{"testing", "one-A", "-b1", "one-B", "-b2"},
	},
	{
		0,
		nil,
		nil,
		nil,
		[]string{"testing", "two-A"},
	},
	{
		0,
		nil,
		nil,
		nil,
		[]string{"testing", "two-B"},
	},
	{
		-1,
		nil,
		nil,
		nil,
		[]string{"testing", "two-C"},
	},
	{
		0,
		nil,
		[]string{"Print full version information. (default true)"},
		nil,
		[]string{"testing", "help"},
	},
	{
		0,
		nil,
		[]string{"one-A [<flags>]:"},
		nil,
		[]string{"testing", "help", "one-A"},
	},
	{
		0,
		nil,
		[]string{"one-A [<flags>]:", "one-B [<flags>]:"},
		nil,
		[]string{"testing", "help", "--commands", "one-A,one-B"},
	},
	{
		0,
		nil,
		[]string{"two-A [<flags>]:", "two-B [<flags>]:"},
		nil,
		[]string{"testing", "help", "two"},
	},
	{
		0,
		nil,
		[]string{"test package", "test tag", "test hash", "test date"},
		nil,
		[]string{"testing", "version"},
	},
	{
		0,
		nil,
		[]string{"test package", "test tag", "test hash", "test date"},
		nil,
		[]string{"testing", "version", "-full"},
	},
	{
		0,
		nil,
		[]string{"test package"},
		[]string{"test tag", "test hash", "test date"},
		[]string{"testing", "version", "-package"},
	},
	{
		0,
		nil,
		[]string{"test tag"},
		[]string{"test package", "test hash", "test date"},
		[]string{"testing", "version", "-tag"},
	},
	{
		0,
		nil,
		[]string{"test hash"},
		[]string{"test package", "test tag", "test date"},
		[]string{"testing", "version", "-hash"},
	},
	{
		0,
		nil,
		[]string{"test date", AnyExitMessage},
		[]string{"test package", "test tag", "test hash"},
		[]string{"testing", "version", "-date"},
	},
}

func TestFlip(t *testing.T) {
	for _, cmd := range flipExpect {
		fs := &tflags{}
		to := new(bytes.Buffer)
		sets := cmdSet(fs, to)
		f := New("test")
		f.SetOut(to)
		f.AddBuiltIn("help").
			AddBuiltIn("version", "test package", "test tag", "test hash", "test date").
			AddBuiltIn("NoExistingCommand").
			SetGroup("one", 1, sets[0]...).
			SetGroup("two", 2, sets[1]...).
			SetGroup("", -1, sets[2]...)
		f.SetCleanup(ExitAny, func(c context.Context) {
			to.WriteString(AnyExitMessage)
		})
		res := f.Execute(nil, cmd.cmd)
		if res != cmd.expectExit {
			t.Errorf("cmd %s expected %d, but received %d", cmd.cmd, cmd.expectExit, res)
		}
		if cmd.fsFunc != nil {
			cmd.fsFunc(t, fs)
		}
		help := to.String()
		if cmd.expectHelp != nil {
			for _, v := range cmd.expectHelp {
				if !strings.Contains(help, v) {
					t.Errorf("Expected help string did not contain %s:\n\n%s", v, help)
				}
			}
		}
		if cmd.unexpectHelp != nil {
			for _, v := range cmd.unexpectHelp {
				if strings.Contains(help, v) {
					t.Errorf("Expected help string contained %s but should not:\n\n%s", v, help)
				}
			}
		}
	}
}
