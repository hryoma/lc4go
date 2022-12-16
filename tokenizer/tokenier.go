package tokenizer

import (
	"fmt"
	// "log"
	"os"
	// "github.com/hryoma/lc4go/machine"
)

func TokenizeObj(fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		// log.Fatal(err)
		return
	}
	defer file.Close()

	buf := make([]byte, 2)
	for {
		_, err := file.Read(buf)
		if err != nil {
			// log.Fatal(err)
			fmt.Println("")
			return
		}

		// read each word
		for _, b := range buf {
			fmt.Printf("%x", b)
		}
		fmt.Printf(" ")
	}
}

