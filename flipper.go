package flip

import (
	"context"
	"fmt"
	"io"
	"os"
	"sort"
)

// Flipper is the flag line processor interface.
type Flipper interface {
	Adder
	Instructer
	Commander
	Executer
	Cleaner
}

type flipper struct {
	Commander
	Instructer
	*executer //Executer
	*cleaner  //Cleaner
}

// Return a new package default Flipper corresponding to the provided string name.
func New(name string) *flipper {
	return newFlipper(
		func(f *flipper) { f.cleaner = newCleaner() },
		func(f *flipper) { f.Commander = newCommander(f) },
		func(f *flipper) { f.Instructer = newInstructer(name, f.Commander, os.Stdout) },
		func(f *flipper) { f.executer = newExecuter(f.Commander, f.RunCleanup) },
		func(f *flipper) {
			var ifn Cleanup
			ifn = f.Instruction
			f.SetCleanup(ExitUsageError, ifn)
		},
		func(f *flipper) { f.SetGroup("", 0) },
	)
}

type config func(*flipper)

func newFlipper(fns ...config) *flipper {
	f := &flipper{}
	for _, fn := range fns {
		fn(f)
	}
	return f
}

// Adds a builtin command by string name and string argument.
// Currently, commands added by this method are:
// - help (takes no other arguments)
// - version (followed by package, tag, version, and hash information strings, in that order)
func (f *flipper) AddBuiltIn(nc string, args ...string) *flipper {
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

// An interface for grouping and managing commands for a Flip instance.
type Commander interface {
	Grouper
	GetCommand(...string) []Command
	SetCommand(...Command) Flipper
}

type commander struct {
	f      *flipper
	groups *Groups
}

func newCommander(f *flipper) *commander {
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

// TODO: remove/control for duplicates
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

// Returns a new Command provided group, tag, use strings, priority integer,
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
	Has []*Group
}

func newGroups() *Groups {
	return &Groups{make([]*Group, 0)}
}

func (g Groups) SortGroupsBy(tag string) {
	var sfn func(int, int) bool
	switch tag {
	default:
		sfn = func(i, j int) bool { return g.Has[i].Priority < g.Has[j].Priority }
	}
	sort.SliceStable(g.Has, sfn)
}

//
type Group struct {
	Name     string
	Priority int
	Commands []Command
}

// Returns a new group provided the string name, priority integer, and any
// number of Command.
func NewGroup(name string, priority int, cs ...Command) *Group {
	return &Group{name, priority, cs}
}

// Set the groups sorting parameter. "alpha" indicating alphabetic sorting
// is the only currently available outside of the default sort by priority.
func (g *Group) SortCommandsBy(s string) {
	var sfn func(int, int) bool
	switch s {
	case "alpha":
		sfn = func(i, j int) bool { return g.Commands[i].Tag() < g.Commands[j].Tag() }
	default:
		sfn = func(i, j int) bool { return g.Commands[i].Priority() < g.Commands[j].Priority() }
	}
	sort.SliceStable(g.Commands, sfn)
}

// Writes the entire group usage to the provided io.Writer.
func (g *Group) Use(o io.Writer) {
	g.SortCommandsBy("default")
	for _, cmd := range g.Commands {
		cmd.Use(o)
	}
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

type queueCmd struct {
	*Group
	Command
	escapes bool
}

type isCommandFunc func(string) *queueCmd

func isCommand(cm Commander) isCommandFunc {
	return func(s string) *queueCmd {
		gs := cm.Groups()
		for _, g := range gs.Has {
			for _, cmd := range g.Commands {
				if s == cmd.Tag() {
					return &queueCmd{g, cmd, cmd.Escapes()}
				}
			}
		}
		return nil
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
	*queueCmd
	start, stop int
	v           []string
}

type pops []*pop

func (p pops) sort() {
	sort.SliceStable(p, func(i, j int) bool {
		return p[i].Group.Priority < p[j].Group.Priority
	})
	sort.SliceStable(p, func(i, j int) bool {
		if p[i].Group.Name == p[j].Group.Name {
			return p[i].Command.Priority() < p[j].Command.Priority()
		}
		return false
	})
}

func queue(fn isCommandFunc, arguments []string) pops {
	var ps pops

	for i, v := range arguments {
		if qc := fn(v); qc != nil {
			a := &pop{qc, i, 0, nil}
			ps = append(ps, a)
			if a.escapes {
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

	ps.sort()
	// trim out commands with lesser precedence than last escaping

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
			cmd := p.Command
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
