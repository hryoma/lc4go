package main

import (
	"fmt"
	"github.com/chzyer/readline"
	"github.com/hryoma/lc4go/emulator"
	"github.com/hryoma/lc4go/machine"
	"github.com/spf13/cobra"
	"strings"
)

func fInitLc4(cmd *cobra.Command, args []string) {
	fmt.Println("lc4")
	emulator.InitLc4(&lc4)
}

func fBreakpoint(cmd *cobra.Command, args []string) {
	fmt.Println("breakpoint / b")
	emulator.Breakpoint(&lc4)
}

func fContinue(cmd *cobra.Command, args []string) {
	fmt.Println("continue / c")
	emulator.Continue(&lc4)
}

func fLoad(cmd *cobra.Command, args []string) {
	fmt.Println("load / l")
	emulator.Load(&lc4)
}

func fNext(cmd *cobra.Command, args []string) {
	fmt.Println("next / n")
	emulator.Next(&lc4)
}

func fPrint(cmd *cobra.Command, args []string) {
	fmt.Println("print / p")
	emulator.Print(&lc4)
}

func fPrintCode(cmd *cobra.Command, args []string) {
	fmt.Println("print / p -c")
	emulator.PrintCode(&lc4)
}

func fPrintMem(cmd *cobra.Command, args []string) {
	fmt.Println("print / p -m")
	emulator.PrintMem(&lc4)
}

func fPrintNzp(cmd *cobra.Command, args []string) {
	fmt.Println("print / p -n")
	emulator.PrintNzp(&lc4)
}

func fPrintReg(cmd *cobra.Command, args []string) {
	fmt.Println("print / p -r")
	emulator.PrintReg(&lc4)
}

func fRun(cmd *cobra.Command, args []string) {
	fmt.Println("run / r")
	emulator.Run(&lc4)
}

func fReset(cmd *cobra.Command, args []string) {
	fmt.Println("reset")
	emulator.Reset(&lc4)
}

func fStep(cmd *cobra.Command, args []string) {
	fmt.Println("step / s")
	emulator.Step(&lc4)
}

var initOptInput, initOptOutput string

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

var lc4 machine.Machine

func init() {
	fmt.Println("init")
	lc4.Pc = 20
}

func main() {
	fmt.Println("LC4 ISA Emulator")

	// initialize state
	emulator.InitLc4(&lc4)

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
	}
}
