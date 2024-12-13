package main

//go:generate go test -run=TestDocHelp -update

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"module-builder/internal/module"
	"module-builder/internal/sort"
)

type command struct {
	flags   *flag.FlagSet
	run     func([]string)
	usage   string
	short   string
	long    string
	hasArgs bool
}

func (c command) name() string {
	name, _, _ := strings.Cut(c.usage, " ")
	return name
}

var (
	moduleFlags = flag.NewFlagSet("module", flag.ExitOnError)
	cleanFlags  = flag.NewFlagSet("clean", flag.ExitOnError)

	outputDir   string
	promoteType = module.PromoteNone

	commands = []*command{
		{
			usage:   "module args... [flags]",
			short:   "build archive(s) for module(s) and update the index.yaml",
			long:    ``, // TODO
			flags:   moduleFlags,
			run:     runModules,
			hasArgs: true,
		},
		{
			usage: "sort",
			short: "sorts index-dev.yaml and index.yaml",
			long:  ``, // TODO
			run:   runSort,
		},
	}
)

func init() {
	moduleFlags.StringVar(&outputDir, "output", "_artifacts", "output directory for archives")
	moduleFlags.Var(&promoteType, "promote", "promotion type for modules, disabled if empty")

	cleanFlags.StringVar(&outputDir, "output", "_artifacts", "output directory for archives")

	for _, cmd := range commands {
		name := cmd.name()
		if cmd.flags == nil {
			cmd.flags = flag.NewFlagSet(name, flag.ExitOnError)
		}
		cmd.flags.Usage = func() {
			help(name)
		}
	}
}

func output(msgs ...any) {
	fmt.Fprintln(flag.CommandLine.Output(), msgs...)
}

func usage() {
	printCommand := func(cmd *command) {
		output(fmt.Sprintf("\t%s\t%s", cmd.name(), cmd.short))
	}
	output("ModuleBuilder is a tool for managing host OS modules.")
	output()
	output("Usage:")
	output()
	output("\tmodule-builder <command> [arguments]")
	output()
	output("The commands are:")
	output()
	for _, cmd := range commands {
		printCommand(cmd)
	}
}

func failf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(1)
}

func findCommand(name string) *command {
	for _, cmd := range commands {
		if cmd.name() == name {
			return cmd
		}
	}
	return nil
}

func help(name string) {
	cmd := findCommand(name)
	if cmd == nil {
		failf("unknown command %q", name)
	}
	output(fmt.Sprintf("Usage: module-builder %s", cmd.usage))
	output()
	if cmd.long != "" {
		output(cmd.long)
	} else {
		output(fmt.Sprintf("ModuleBuilder %s is used to %s.", cmd.name(), cmd.short))
	}
	anyflags := false
	cmd.flags.VisitAll(func(*flag.Flag) {
		anyflags = true
	})
	if anyflags {
		output()
		output("Flags:")
		output()
		cmd.flags.PrintDefaults()
	}
}

func runModules(args []string) {
	if len(args) == 0 {
		fmt.Println("No modules set, nothing to do.")
		return
	}

	if err := module.Build(module.Config{
		Promote:   promoteType,
		Output:    outputDir,
		Dirs:      args,
		LogWriter: os.Stderr,
	}); err != nil {
		fmt.Printf("Build failed: %v\n", err)
		os.Exit(2)
	}

	fmt.Println("Build completed.")
}

func runSort(_ []string) {
	if err := sort.Index(sort.Config{
		LogWriter: os.Stderr,
	}); err != nil {
		fmt.Printf("Sorting failed: %v\n", err)
		os.Exit(2)
	}

	fmt.Println("Sorting completed.")
}

func main() {
	log.SetFlags(0)
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(2)
	}

	if args[0] == "help" {
		flag.CommandLine.SetOutput(os.Stdout)
		switch len(args) {
		case 1:
			flag.Usage()
		case 2:
			help(args[1])
		default:
			flag.Usage()
			failf("too many arguments to \"help\"")
		}
		os.Exit(0)
	}

	cmd := findCommand(args[0])
	if cmd == nil {
		flag.Usage()
		os.Exit(2)
	}

	_ = cmd.flags.Parse(args[1:]) // will exit on error
	args = cmd.flags.Args()
	if !cmd.hasArgs && len(args) > 0 {
		help(cmd.name())
		failf("command %s does not accept any arguments\n", cmd.name())
	}
	cmd.run(args)
}
