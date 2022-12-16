package machine

import (
	"testing"
)

func TestGetSignExtNPos(t *testing.T) {
	var nBits uint16 = 9
	//                    |   |   |   |   |
	var testData uint16 = 0b1010101010101010
	var expected uint16 = 0b0000000010101010
	var actual uint16 = getSignExtN(testData, nBits)
	if actual != expected {
		t.Error("Sign extension failed for positive int. Expected", expected, "but got", actual)
	}
}

func TestGetSignExtNNeg(t *testing.T) {
	var nBits uint16 = 9
	//                    |   |   |   |   |
	var testData uint16 = 0b0101010101010101
	var expected uint16 = 0b1111111101010101
	var actual uint16 = getSignExtN(testData, nBits)
	if actual != expected {
		t.Error("Sign extension failed for negative int. Expected", expected, "but got", actual)
	}
}

