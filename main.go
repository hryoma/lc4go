package main

import (
	"fmt"
	"github.com/chzyer/readline"
	"github.com/hryoma/lc4go/emulator"
	"github.com/spf13/cobra"
	"strings"
)

var breakpointCmd = &cobra.Command{
	Use:     "breakpoint",
	Short:   "Set a breakpoint in the code",
	Aliases: []string{"b"},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("breakpoint / b")
		emulator.Breakpoint()
	},
}

var continueCmd = &cobra.Command{
	Use:     "continue",
	Short:   "Continue running the instructions until termination",
	Aliases: []string{"c"},
	Run: func(cmd *cobra.Command, args []string) {
		emulator.Continue()
	},
}

var loadCmd = &cobra.Command{
	Use:     "load",
	Short:   "Load a file",
	Aliases: []string{"l"},
	Run: func(cmd *cobra.Command, args []string) {
		fileName, err := cmd.Flags().GetString("obj")
		if err != nil {
			fmt.Println(err)
		} else {
			emulator.Load(fileName)
		}
	},
}

var nextCmd = &cobra.Command{
	Use:     "next",
	Short:   "Run until the program counter reaches PC + 1",
	Aliases: []string{"n"},
	Run: func(cmd *cobra.Command, args []string) {
		emulator.Next()
	},
}

var printCmd = &cobra.Command{
	Use:     "print",
	Short:   "Print register values, PSR bits, code lines, or content in memory",
	Aliases: []string{"p"},
	Run: func(cmd *cobra.Command, args []string) {
		emulator.Print()
	},
}

var printCodeCmd = &cobra.Command{
	Use:     "code",
	Short:   "Print code lines",
	Aliases: []string{"c"},
	Run: func(cmd *cobra.Command, args []string) {
		emulator.PrintCode()
	},
}

var printMemCmd = &cobra.Command{
	Use:     "mem",
	Short:   "Print content in memory",
	Aliases: []string{"m"},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("Invalid number of arguments provided")
			return
		}

		emulator.PrintMem(args[0])
	},
}

var printPsrCmd = &cobra.Command{
	Use:     "psr",
	Short:   "Print NZP and privilege bits",
	Aliases: []string{"n"},
	Run: func(cmd *cobra.Command, args []string) {
		emulator.PrintPsr()
	},
}

var printRegCmd = &cobra.Command{
	Use:     "reg",
	Short:   "Print register values",
	Aliases: []string{"r"},
	Run: func(cmd *cobra.Command, args []string) {
		emulator.PrintReg()
	},
}

var runCmd = &cobra.Command{
	Use:     "run",
	Short:   "Run the file from the beginning",
	Aliases: []string{"r"},
	Run: func(cmd *cobra.Command, args []string) {
		emulator.Run()
	},
}

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset all values to the initial state",
	Run: func(cmd *cobra.Command, args []string) {
		emulator.Reset()
	},
}

var stepCmd = &cobra.Command{
	Use:     "step",
	Short:   "Execute one instruction",
	Aliases: []string{"s"},
	Run: func(cmd *cobra.Command, args []string) {
		emulator.Step()
	},
}

var rootCmd = &cobra.Command{}

func init() {
	// initialize state
	emulator.Reset()

	// register commands
	rootCmd.AddCommand(breakpointCmd)
	rootCmd.AddCommand(continueCmd)
	rootCmd.AddCommand(loadCmd)
	loadCmd.Flags().StringP("obj", "b", "", "Input object file path")
	rootCmd.AddCommand(nextCmd)
	rootCmd.AddCommand(printCmd)
	printCmd.AddCommand(printCodeCmd)
	printCmd.AddCommand(printMemCmd)
	printCmd.AddCommand(printPsrCmd)
	printCmd.AddCommand(printRegCmd)
	rootCmd.AddCommand(resetCmd)
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(stepCmd)
}

func main() {
	fmt.Println("LC4 ISA Emulator")

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
