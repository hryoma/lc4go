package main

import (
	"fmt"
	"github.com/chzyer/readline"
	"github.com/hryoma/lc4go/emulator"
	"github.com/hryoma/lc4go/tokenizer"
	"github.com/spf13/cobra"
	"strings"
)

var initOptInput, initOptOutput string

func fInitLc4(cmd *cobra.Command, args []string) {
	fmt.Println("lc4")
	fmt.Println(args)
}

func fBreakpoint(cmd *cobra.Command, args []string) {
	fmt.Println("breakpoint / b")
	fmt.Println(args)
}

func fContinue(cmd *cobra.Command, args []string) {
	fmt.Println("continue / c")
	fmt.Println(args)
}

func fLoad(cmd *cobra.Command, args []string) {
	fmt.Println("load / l")
	fmt.Println(args)
}

func fNext(cmd *cobra.Command, args []string) {
	fmt.Println("next / n")
	fmt.Println(args)
}

func fPrint(cmd *cobra.Command, args []string) {
	fmt.Println("print / p")
	fmt.Println(args)
}

func fPrintCode(cmd *cobra.Command, args []string) {
	fmt.Println("print / p -c")
	fmt.Println(args)
}

func fPrintMem(cmd *cobra.Command, args []string) {
	fmt.Println("print / p -m")
	fmt.Println(args)
}

func fPrintNzp(cmd *cobra.Command, args []string) {
	fmt.Println("print / p -n")
	fmt.Println(args)
}

func fPrintReg(cmd *cobra.Command, args []string) {
	fmt.Println("print / p -r")
	fmt.Println(args)
}

func fRun(cmd *cobra.Command, args []string) {
	fmt.Println("run / r")
	fmt.Println(args)
}

func fReset(cmd *cobra.Command, args []string) {
	fmt.Println("reset")
	fmt.Println(args)
}

func fStep(cmd *cobra.Command, args []string) {
	fmt.Println("step / s")
	fmt.Println(args)
}

var breakpointCmd = &cobra.Command{
	Use:     "breakpoint",
	Short:   "Set a breakpoint in the code",
	Aliases: []string{"b"},
	Run:     fBreakpoint,
}

var continueCmd = &cobra.Command{
	Use:     "continue",
	Short:   "Continue running the instructions until termination",
	Aliases: []string{"c"},
	Run:     fContinue,
}

var loadCmd = &cobra.Command{
	Use:     "load",
	Short:   "Load a file",
	Aliases: []string{"l"},
	Run:     fLoad,
}

var nextCmd = &cobra.Command{
	Use:     "next",
	Short:   "Run until the program counter reaches PC + 1",
	Aliases: []string{"n"},
	Run:     fNext,
}

var printCmd = &cobra.Command{
	Use:     "print",
	Short:   "Print register values, NZP bits, code lines, or content in memory",
	Aliases: []string{"p"},
	Run:     fPrint,
}

var printCodeCmd = &cobra.Command{
	Use:     "code",
	Short:   "Print code lines",
	Aliases: []string{"c"},
	Run:     fPrintCode,
}

var printMemCmd = &cobra.Command{
	Use:     "mem",
	Short:   "Print content in memory",
	Aliases: []string{"m"},
	Run:     fPrintMem,
}

var printNzpCmd = &cobra.Command{
	Use:     "nzp",
	Short:   "Print NZP bits",
	Aliases: []string{"n"},
	Run:     fPrintNzp,
}

var printRegCmd = &cobra.Command{
	Use:     "reg",
	Short:   "Print register values",
	Aliases: []string{"r"},
	Run:     fPrintReg,
}

var runCmd = &cobra.Command{
	Use:     "run",
	Short:   "Run the file from the beginning",
	Aliases: []string{"r"},
	Run:     fRun,
}

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset all values to the initial state",
	Run:   fReset,
}
var stepCmd = &cobra.Command{
	Use:     "step",
	Short:   "Execute one instruction",
	Aliases: []string{"s"},
	Run:     fStep,
}

func main() {
	fmt.Println("LC4 ISA Emulator")

	// register commands
	var rootCmd = &cobra.Command{}
	rootCmd.AddCommand(breakpointCmd)
	rootCmd.AddCommand(continueCmd)
	rootCmd.AddCommand(loadCmd)
	rootCmd.AddCommand(nextCmd)
	rootCmd.AddCommand(printCmd)
	printCmd.AddCommand(printCodeCmd)
	printCmd.AddCommand(printMemCmd)
	printCmd.AddCommand(printNzpCmd)
	printCmd.AddCommand(printRegCmd)
	rootCmd.AddCommand(resetCmd)
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(stepCmd)

	// load input file
	tokenizer.Tokenize()

	// initialize shell
	shell, err := readline.NewEx(&readline.Config{
		Prompt:    "lc4> ",
		EOFPrompt: "exit",
	})
	if err != nil {
		panic(err)
	}
	defer shell.Close()

	// i/o loop
	for {
		line, err := shell.Readline()
		if err != nil {
			break
		}

		args := strings.Fields(line)
		if len(args) == 0 {
			continue
		}

		rootCmd.SetArgs(args)
		rootCmd.Execute()
		emulator.Emulate()
	}
}
