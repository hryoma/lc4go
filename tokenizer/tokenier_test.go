package tokenizer

import (
	"testing"
	"github.com/hryoma/lc4go/tokenizer"
	"github.com/hryoma/lc4go/machine"
)

var test_folder string = "/../test_objs/"

func TestTokenizeObjMultiplyObj(t *testing.T) {
	var fileName = test_folder + "multiply.obj"
	tokenizer.TokenizeObj(fileName)

	if machine.Lc4.Mem[0] != 0x9400 {
		t.Log("Data block not parsed correctly")
		t.Log("Expected:", 0x9400, "Actual:", machine.Lc4.Mem[0])
		t.Fail()
	}
	if machine.Lc4.Mem[1] != 0x2300 {
		t.Log("Data block not parsed correctly")
		t.Log("Expected:", 0x2300, "Actual:", machine.Lc4.Mem[1])
		t.Fail()
	}
	if machine.Lc4.Mem[2] != 0x0C03 {
		t.Log("Data block not parsed correctly")
		t.Log("Expected:", 0x0C03, "Actual:", machine.Lc4.Mem[2])
		t.Fail()
	}
	if machine.Lc4.Mem[3] != 0x1480 {
		t.Log("Data block not parsed correctly")
		t.Log("Expected:", 0x1480, "Actual:", machine.Lc4.Mem[3])
		t.Fail()
	}
	if machine.Lc4.Mem[4] != 0x127F {
		t.Log("Data block not parsed correctly")
		t.Log("Expected:", 0x127F, "Actual:", machine.Lc4.Mem[4])
		t.Fail()
	}
	if machine.Lc4.Mem[5] != 0x0FFB {
		t.Log("Data block not parsed correctly")
		t.Log("Expected:", 0x0FFB, "Actual:", machine.Lc4.Mem[5])
		t.Fail()
	}
	if machine.Lc4.Mem[6] != 0x0000{
		t.Log("Data block not parsed correctly")
		t.Log("Expected:", 0x0000, "Actual:", machine.Lc4.Mem[6])
		t.Fail()
	}
}

