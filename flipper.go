package flip

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"sort"
)

// Flipper is the flag line processor interface.
type Flipper interface {
	Adder
	Commander
	Instructer
	Executer
	Cleaner
}

// A struct as the package default flag line processor, implementing Flipper.
type Flip struct {
	Commander
	Instructer
	Executer
	Cleaner
}

// Return a new package default Flip corresponding to the provided string name.
func New(name string) *Flip {
	return NewFlip(
		func(f *Flip) { f.Cleaner = newCleaner() },
		func(f *Flip) { f.Commander = newCommander(f) },
		func(f *Flip) { f.Instructer = newInstructer(name, f.Commander, os.Stdout) },
		func(f *Flip) { f.Executer = newExecuter(f.Commander, f.RunCleanup) },
		func(f *Flip) {
			var ifn Cleanup
			ifn = f.Instruction
			f.SetCleanup(ExitUsageError, ifn)
		},
		func(f *Flip) { f.SetGroup("", 0) },
	)
}

//
type FlipConfig func(*Flip)

//
func NewFlip(fns ...FlipConfig) *Flip {
	f := &Flip{}
	for _, fn := range fns {
		fn(f)
	}
	return f
}

// Adds a builtin command by string name and string argument.
// Currently, commands added by this method are:
// - help (takes no other arguments)
// - version (followed by package, tag, version, and hash information strings, in that order)
func (f *Flip) AddBuiltIn(nc string, args ...string) Flipper {
	switch nc {
	case "help":
		return f.addHelp()
	case "version":
		return f.addVersion(args...)
	}
	return f
}

// An interface for grouping commands.
type Grouper interface {
	Groups() *Groups
	GetGroup(string) *Group
	SetGroup(string, int, ...Command) Flipper
}

// An interface for grouping nad managing commands for a Flip instance.
type Commander interface {
	Grouper
	GetCommand(...string) []Command
	SetCommand(...Command) Flipper
}

type commander struct {
	f      Flipper
	groups *Groups
}

func newCommander(f Flipper) *commander {
	return &commander{f, newGroups()}
}

//
func (c *commander) Groups() *Groups {
	return c.groups
}

// Return a group corresponding to the provided name string, or nil is nothing found.
func (c *commander) GetGroup(name string) *Group {
	for _, g := range c.groups.Has {
		if name == g.Name {
			return g
		}
	}
	return nil
}

// Set a group with the string name and integer priority, containing the given commands.
func (c *commander) SetGroup(name string, priority int, cmds ...Command) Flipper {
	c.groups.Has = append(c.groups.Has, NewGroup(name, priority))
	for _, v := range cmds {
		v.SetGroup(name)
	}
	c.SetCommand(cmds...)
	return c.f
}

// Get commands corresponding to the provided string keys.
func (c *commander) GetCommand(ks ...string) []Command {
	var ret []Command
	for _, g := range c.groups.Has {
		for _, k := range ks {
			if k == g.Name {
				ret = append(ret, g.Commands...)
			}
		}
		for _, cmd := range g.Commands {
			for _, k := range ks {
				if k == cmd.Tag() {
					ret = append(ret, cmd)
				}
			}
		}
	}
	//remove dupes
	return ret
}

// Set the provided Commands, returning a Flip instance (useful for chaining).
func (c *commander) SetCommand(cmds ...Command) Flipper {
	for _, cmd := range cmds {
		g := c.GetGroup(cmd.Group())
		g.Commands = append(g.Commands, cmd)
	}
	return c.f
}

// A function taking context.Context, and a string slice, returns context.Context
// and an ExitStatus.
type CommandFunc func(context.Context, []string) (context.Context, ExitStatus)

// An interface for encapsulating a command.
type Command interface {
	Group() string
	SetGroup(string)
	Tag() string
	Priority() int
	Escapes() bool
	Use(io.Writer)
	Execute(context.Context, []string) (context.Context, ExitStatus)
	Flagger
}

type command struct {
	group, tag string
	use        string
	priority   int
	escapes    bool
	hasRun     bool
	cfn        CommandFunc
	*FlagSet
}

// Returns a new Command provided group, tag, use strings, priority integer
// a boolean indicating escape (stop processing command for other commands
// after this command is found, passing the params to the current command instead
// of going to another command), A CommandFunc to process the command, and a
// corresponding FlagSet for the Command).
func NewCommand(group, tag, use string,
	priority int,
	escapes bool,
	cfn CommandFunc,
	fs *FlagSet) Command {
	return &command{group, tag, use, priority, escapes, false, cfn, fs}
}

// Set the Command group to the provided string.
func (c *command) SetGroup(k string) {
	c.group = k
}

// Returns the Command group as a string.
func (c *command) Group() string {
	return c.group
}

// Returns the Command tag(i.e. primary title) as a string
func (c *command) Tag() string {
	return c.tag
}

// Returns the Command priority in its group as an integer.
func (c *command) Priority() int {
	return c.priority
}

// Returns a boolean indicating if the Command escapes processing further commands.
func (c *command) Escapes() bool {
	return c.escapes
}

func (c *command) useHead(o io.Writer) {
	white(o, fmt.Sprintf("-----\n%s [<flags>]:\n", c.tag))
}

func (c *command) useString(o io.Writer) {
	white(o, fmt.Sprintf("\t%s\n\n", c.use))
}

// Writes the Command's entire usage to the provided io.Writer.
func (c *command) Use(o io.Writer) {
	c.useHead(o)
	c.useString(o)
	c.Usage(o)
	fmt.Fprint(o, "\n")
}

// Executes the Commands CommandFunc.
func (c *command) Execute(ctx context.Context, v []string) (context.Context, ExitStatus) {
	if c.cfn != nil {
		c.hasRun = true
		return c.cfn(ctx, v)
	}
	return ctx, ExitFailure
}

//
type Groups struct {
	SortBy string
	Has    []*Group
}

func newGroups() *Groups {
	return &Groups{"default", make([]*Group, 0)}
}

// groups Len function for sort.Sort
func (g Groups) Len() int { return len(g.Has) }

// groups Less function for sort.Sort
func (g Groups) Less(i, j int) bool {
	switch g.SortBy {
	default:
		return g.Has[i].Priority < g.Has[j].Priority
	}
	return false
}

// groups Swap function for sort.Sort
func (g Groups) Swap(i, j int) { g.Has[i], g.Has[j] = g.Has[j], g.Has[i] }

//
type Group struct {
	Name     string
	Priority int
	sortBy   string
	Commands []Command
}

// Returns a new group provided the string name, priority integer, and any
// number of Command.
func NewGroup(name string, priority int, cs ...Command) *Group {
	return &Group{name, priority, "", cs}
}

// group Len function for sort.Sort
func (g Group) Len() int { return len(g.Commands) }

// group Less function for sort.Sort
func (g Group) Less(i, j int) bool {
	switch g.sortBy {
	case "alpha":
		return g.Commands[i].Tag() < g.Commands[j].Tag()
	default:
		return g.Commands[i].Priority() < g.Commands[j].Priority()
	}
	return false
}

// group Swap function for sort.Sort
func (g Group) Swap(i, j int) {
	g.Commands[i], g.Commands[j] = g.Commands[j], g.Commands[i]
}

// Set the groups sorting parameter. "alpha" indicating alphabetic sorting
// is the only currently available outside of the default sort by priority.
func (g *Group) SortBy(s string) {
	g.sortBy = s
	sort.Sort(g)
}

// Writes the entire group usage to the provided io.Writer.
func (g *Group) Use(o io.Writer) {
	g.SortBy("default")
	for _, cmd := range g.Commands {
		cmd.Use(o)
	}
}

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

func titleString(titleFmtString, name string, b *bytes.Buffer) {
	title := Color(Bold, FgHiWhite)
	title(b, fmt.Sprintf(titleFmtString, name))
}

func defaultInstruction(tag string, cm Commander, i *instructer) Cleanup {
	return func(c context.Context) {
		out := i.Out()
		b := new(bytes.Buffer)
		titleString(i.titleFmtString, tag, b)

		gs := cm.Groups()
		sort.Sort(gs)
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

// An interface for command execution.
type Executer interface {
	Execute(context.Context, []string) int
}

type executer struct {
	iscmdfn isCommandFunc
	cleanfn runCleanupFunc
}

func newExecuter(cm Commander, cu runCleanupFunc) *executer {
	return &executer{isCommand(cm), cu}
}

type isCommandFunc func(string) (Command, bool, bool)

func isCommand(cm Commander) isCommandFunc {
	return func(s string) (Command, bool, bool) {
		gs := cm.Groups()
		for _, g := range gs.Has {
			for _, cmd := range g.Commands {
				if s == cmd.Tag() {
					return cmd, true, cmd.Escapes()
				}
			}
		}
		return nil, false, false
	}
}

// An integer type useful for marking results of commands.
type ExitStatus int

const (
	ExitNo         ExitStatus = 999  // continue processing commands
	ExitSuccess    ExitStatus = 0    // return 0
	ExitFailure    ExitStatus = -1   // return -1
	ExitUsageError ExitStatus = -2   // return -2
	ExitAny        ExitStatus = -666 // status for cleaning function setup, never return
)

type pop struct {
	start, stop int
	c           Command
	v           []string
}

type pops []*pop

// internal pops type Len function for sort.Sort
func (p pops) Len() int { return len(p) }

//  internal pops type Less function for sort.Sort
func (p pops) Less(i, j int) bool { return p[i].c.Priority() < p[j].c.Priority() }

//  internal pops type Swap function for sort.Sort
func (p pops) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

func queue(fn isCommandFunc, arguments []string) pops {
	var ps pops

	for i, v := range arguments {
		if cmd, exists, escapes := fn(v); exists {
			a := &pop{i, 0, cmd, nil}
			ps = append(ps, a)
			if escapes {
				break
			}
		}
	}

	li := len(ps) - 1
	la := len(arguments)
	for i, v := range ps {
		if i+1 <= li {
			nx := ps[i+1]
			v.stop = nx.start
		} else {
			v.stop = la
		}
	}

	for _, p := range ps {
		p.v = arguments[p.start:p.stop]
	}

	sort.Sort(ps)

	return ps
}

func execute(ctx context.Context, cmd Command, arguments []string) (context.Context, ExitStatus) {
	err := cmd.Parse(arguments)
	if err != nil {
		return ctx, ExitUsageError
	}
	return cmd.Execute(ctx, arguments)
}

// The Execute function taking a context.Context, and slice of string arguments,
// returning an integer corresponding to an ExitStatus.
func (e *executer) Execute(ctx context.Context, arguments []string) int {
	var exit ExitStatus
	switch {
	case len(arguments) <= 1:
		goto INSTRUCTION
	default:
		q := queue(e.iscmdfn, arguments)
		for _, p := range q {
			cmd := p.c
			args := p.v[1:]
			ctx, exit = execute(ctx, cmd, args)
			switch exit {
			case ExitSuccess:
				return e.cleanfn(exit, ctx)
			case ExitFailure:
				return e.cleanfn(exit, ctx)
			case ExitUsageError:
				goto INSTRUCTION
			default:
				continue
			}
		}
	}

INSTRUCTION:
	return e.cleanfn(ExitUsageError, ctx)
}

// A cleanup function taking a context.Context only.
type Cleanup func(context.Context)

type runCleanupFunc func(ExitStatus, context.Context) int

// An interface for post-command actions.
type Cleaner interface {
	SetCleanup(ExitStatus, ...Cleanup)
	RunCleanup(ExitStatus, context.Context) int
}

type cleaner struct {
	cfns map[ExitStatus][]Cleanup
}

func newCleaner() *cleaner {
	return &cleaner{make(map[ExitStatus][]Cleanup)}
}

// Set the provided Cleanup functions to be run on the provided ExitStatus.
func (c *cleaner) SetCleanup(e ExitStatus, cfns ...Cleanup) {
	if c.cfns[e] == nil {
		c.cfns[e] = make([]Cleanup, 0)
	}
	c.cfns[e] = append(c.cfns[e], cfns...)
}

// Given an ExitStatus and a context.Context, runs any associated Cleanup functions.
func (c *cleaner) RunCleanup(e ExitStatus, ctx context.Context) int {
	if cfns, ok := c.cfns[e]; ok {
		for _, cfn := range cfns {
			cfn(ctx)
		}
	}
	if afns, ok := c.cfns[ExitAny]; ok {
		for _, afn := range afns {
			afn(ctx)
		}
	}

	return int(e)
}
