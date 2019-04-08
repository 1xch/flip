# flip

A flag line processor for the Go programming language


### Example

Examples speak a lot, so jump into some code.

    main.go:

        ```
        package main

        import (
            "context"
            "log"
            "os"
            "path"

            "github.com/Laughs-In-Flowers/flip"
        )

        var (
            F              flip.Flipper
            value          string = ""
            versionPackage string = path.Base(os.Args[0])
            versionTag     string = "Example"
            versionHash    string = "Ex#1"
            versionDate    string = "Today"
            output         int    = 0
        )

        func TopCommand() flip.Command {
            var val string
            fs := flip.NewFlagSet("t", flip.ContinueOnError)
            fs.StringVar(&val, "value", val, "A flag string value")

            return flip.NewCommand(
                "",
                "./example",
                "Top level options use.",
                1,
                false,
                func(c context.Context, a []string) (context.Context, flip.ExitStatus) {
                    value = value + " " + val
                    return c, flip.ExitNo
                },
                fs,
            )
        }

        func RunCommand1() flip.Command {
            var val string
            fs := flip.NewFlagSet("r1", flip.ContinueOnError)
            fs.StringVar(&val, "value", val, "A flag string value")

            return flip.NewCommand(
                "",
                "run1",
                "run1 command",
                1,
                false,
                func(c context.Context, a []string) (context.Context, flip.ExitStatus) {
                    value = value + " " + val
                    log.Printf("%v", value)
                    return c, flip.ExitSuccess
                },
                fs,
            )
        }

        func RunCommand2() flip.Command {
            var val string
            fs := flip.NewFlagSet("r2", flip.ContinueOnError)
            fs.StringVar(&val, "value", val, "A flag string value")

            return flip.NewCommand(
                "",
                "run2",
                "run2 command",
                2,
                false,
                func(c context.Context, a []string) (context.Context, flip.ExitStatus) {
                    value = value + " " + val
                    log.Printf("%v", value)
                    return c, flip.ExitSuccess
                },
                fs,
            )
        }

        func RunCommand3() flip.Command {
            var val string
            fs := flip.NewFlagSet("r3", flip.ContinueOnError)
            fs.StringVar(&val, "value", val, "A flag string value")

            return flip.NewCommand(
                "",
                "run3",
                "run3 command",
                2,
                false,
                func(c context.Context, a []string) (context.Context, flip.ExitStatus) {
                    value = value + " " + val
                    return c, flip.ExitNo
                },
                fs,
            )
        }

        func init() {
            F = flip.New("example")
            F.AddBuiltIn("version", versionPackage, versionTag, versionHash, versionDate).
                AddBuiltIn("help").
                SetGroup("gtop", -1, TopCommand()).
                SetGroup("grun", 1, RunCommand1(), RunCommand2(), RunCommand3())
        }

        func main() {
            os.Exit(F.Execute(context.Background(), os.Args))
        }               
        ```

    compile & run commands:

    1. ./example --value "Y" run1 -value "Z"

    2. ./example -value "X" run3 -value "Y" run2 -value "Z"

    2. ./example version

    3. ./example version -tag
        
    4. ./example help

    5. ./example help grun

    6. ./example help run2
