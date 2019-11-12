package flip

import (
	"bytes"
	"context"
	"fmt"
	"io"
)

// An interface for providing instruction i.e. writes usage strings.
type Instructer interface {
	SwapInstructer(Instructer)
	Instruction(context.Context)
	SubsetInstruction(c ...Command) func(context.Context)
	Writer
}

type iswapper struct {
	Instructer
}

func (s *iswapper) SwapInstructer(i Instructer) {
	s.Instructer = i
}

type instructer struct {
	titleFmtString string
	output         io.Writer
	ifn            Cleanup
}

func newInstructer(tag string, cm Commander, o io.Writer) *iswapper {
	i := &instructer{"%s [OPTIONS...] {COMMAND} ...\n\n", o, nil}
	i.ifn = defaultInstruction(tag, cm, i)
	return &iswapper{i}
}

func (i *instructer) SwapInstructer(Instructer) {}

// Given a context.Context writes the what the Instructer is configured to write.
func (i *instructer) Instruction(c context.Context) {
	i.ifn(c)
}

// Returns a function to write instructions for a subset of provided Commands.
func (i *instructer) SubsetInstruction(cs ...Command) func(context.Context) {
	return func(c context.Context) {
		out := i.Out()
		b := new(bytes.Buffer)
		for _, cmd := range cs {
			cmd.Use(b)
		}
		fmt.Fprint(out, b)
	}
}

func titleString(titleFmtString, name string, b io.Writer) {
	title := Color(Bold, FgHiWhite)
	title(b, fmt.Sprintf(titleFmtString, name))
}

func defaultInstruction(tag string, cm Commander, i *instructer) Cleanup {
	return func(c context.Context) {
		out := i.Out()
		b := new(bytes.Buffer)
		titleString(i.titleFmtString, tag, b)

		gs := cm.Groups()
		gs.SortGroupsBy("")
		for _, g := range gs.Has {
			g.Use(b)
		}

		fmt.Fprint(out, b)
	}
}

// The io.Writer configured to the Instructor.
func (i *instructer) Out() io.Writer {
	return i.output
}

// Set the provided io.Writer to the Instructor.
func (i *instructer) SetOut(w io.Writer) {
	i.output = w
}
