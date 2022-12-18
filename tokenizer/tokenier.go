package tokenizer

import (
	"fmt"
	// "log"
	"os"
	"github.com/hryoma/lc4go/machine"
)

func readChar(file *os.File) (char byte, err error) {
	buf := make([]byte, 1)
	_, err = file.Read(buf)
	if err != nil {
		fmt.Println("Could not read word")
		return
	}

	// get the word
	char = buf[0]
	return 
}

func readWord(file *os.File) (word uint16, err error) {
	buf := make([]byte, 2)
	_, err = file.Read(buf)
	if err != nil {
		fmt.Println("Could not read word")
		return
	}

	// get the word
	byte1 := buf[0]
	byte2 := buf[1]

	word = uint16((byte1 << 8) + byte2)

	return
}

func parseCodeBlock(file *os.File) {
	// address
	addr, err := readWord(file)
	if err != nil {
		return
	}

	// number
	num, err := readWord(file)
	if err != nil {
		return
	}

	// read words
	for i := uint16(0); i < num; i++ {
		word, err := readWord(file)
		if err != nil {
			return
		}

		machine.Lc4.Mem[addr + i] = word;
	}
}

func parseDataBlock(file *os.File) {
	// address
	addr, err := readWord(file)
	if err != nil {
		return
	}

	// number
	num, err := readWord(file)
	if err != nil {
		return
	}

	// read words
	for i := uint16(0); i < num; i++ {
		word, err := readWord(file)
		if err != nil {
			return
		}

		machine.Lc4.Mem[addr + i] = word;
	}
}

func parseSymbol(file *os.File) {
	// address
	_, err := readWord(file)
	if err != nil {
		return
	}

	// number
	num, err := readWord(file)
	if err != nil {
		return
	}

	// read chars
	for i := uint16(0); i < num; i++ {
		// read the next char, but don't do anything with it
		_, err := readChar(file)
		if err != nil {
			return
		}
	}
}

func parseFileName(file *os.File) {
	// number
	num, err := readWord(file)
	if err != nil {
		return
	}

	// read chars
	for i := uint16(0); i < num; i++ {
		// read the next char, but don't do anything with it
		_, err := readChar(file)
		if err != nil {
			return
		}
	}
}

func parseLineNumber(file *os.File) {
	// address
	_, err := readWord(file)
	if err != nil {
		return
	}

	// line
	_, err = readWord(file)
	if err != nil {
		return
	}

	// file index
	_, err = readWord(file)
	if err != nil {
		return
	}

	// TODO: implement behavior for reading file name
}

func TokenizeObj(fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		// log.Fatal(err)
		return
	}
	defer file.Close()

	for {
		word, err := readWord(file)
		if err != nil {
			// log.Fatal(err)
			fmt.Println("")
			return
		}

		switch word {
		case 0xCADE:
			parseCodeBlock(file)
		case 0xDADA:
			parseDataBlock(file)
		case 0xC3B7:
			parseSymbol(file)
		case 0xF17E:
			parseFileName(file)
		case 0x715E:
			parseLineNumber(file)
		}
	}
}

