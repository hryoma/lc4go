# lc4go

LC4 ISA emulator, written in Go. A CLI app that acts like a GDB for LC4.

LC4 is an assembly language. For more information about the instruction set, see [here](https://www.cis.upenn.edu/~cis5710/current/lc4.html).

This tool takes in `.obj` files, which are the bytecodes.


## Usage

### Launching the CLI

From the root directory of the project, the application can be launched with the following command:

```bash
go run main.go
```

### CLI Commands

**Loading an .obj File**

To load an .obj file, you can use the `load` command:

```bash
lc4> load -b <path/to/obj/file>
```

**Breakpoint**

You can set breakpoints anywhere in memory with the `breakpoint`/`b` command:

```bash
lc4> breakpoint 0x1234
```

**Printing**

You can print the states and values stored in the machine at any point. This includes:
- `p code`: print code at the current program counter
- `p mem <addr>`: print the value stored in memory at the provided address
- `p reg`: print all register values
- `p psr`: print thet NZP bits and privilege bit
- `p`: print all of the above at once


**Execution**

Basic execution commands are provided, with convenient aliases:
- `step`/`s`: execute one instruction
- `next`/`n`: run until PC = current PC + 1
- `continue`/`c`: run from current PC to the end
- `run`/`r`: run from from the beginning to the end

**Miscellaneous**

Additionally, there are a couple of other helper commands:
- `reset`: reset only states, but keep memory and breakpoints
- `clear`: reset all states and values


## Example

```bash
# go run main.go

lc4> load -b example/os.obj
lc4> load -b example/math.obj
lc4> b 4
lc4> r
lc4> p
lc4> clear
```

Note that the `os.obj` file should always be loaded in.
