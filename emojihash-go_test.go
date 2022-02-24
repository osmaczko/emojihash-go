package main

import (
	"math"
	"math/big"
	"reflect"
	"testing"
)

func TestToBigBase(t *testing.T) {
	checker := func(value *big.Int, base uint64, expected *[](uint64)) {
		res := ToBigBase(value, base)
		if !reflect.DeepEqual(res, *expected) {
			t.Fatalf("invalid big base conversion %v != %v", res, *expected)
		}
	}

	lengthChecker := func(value *big.Int, base, expectedLength uint64) {
		res := ToBigBase(value, base)
		if len(res) != int(expectedLength) {
			t.Fatalf("invalid big base conversion %d != %d", len(res), expectedLength)
		}
	}

	checker(new(big.Int).SetUint64(15), 16, &[](uint64){15})
	checker(new(big.Int).SetUint64(495), 16, &[](uint64){1, 14, 15})
	checker(new(big.Int).SetUint64(495), 30, &[](uint64){16, 15})
	checker(new(big.Int).SetUint64(495), 1024, &[](uint64){495})
	checker(new(big.Int).SetUint64(2048), 1024, &[](uint64){2, 0})

	val, _ := new(big.Int).SetString("0xFFFFFFFFFFFFFF", 0)
	base := uint64(math.Pow(2, 7*4))
	checker(val, base, &[](uint64){base - 1, base - 1})

	val, _ = new(big.Int).SetString("0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF", 0)
	lengthChecker(val, 2757, 14)
	lengthChecker(val, 2756, 15)
}

func TestToEmojiHash(t *testing.T) {
	alphabet := [](string){"😇", "🤐", "🥵", "🙊", "🤌"}

	checker := func(value *big.Int, hashLen int, expected *[](string)) {
		res, err := ToEmojiHash(value, hashLen, &alphabet)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(res, *expected) {
			t.Fatalf("invalid emojihash conversion %v != %v", res, *expected)
		}
	}

	val := new(big.Int).SetUint64(777)
	checker(val, 5, &[](string){"🤐", "🤐", "🤐", "😇", "🥵"})
	checker(val, 10, &[](string){"😇", "😇", "😇", "😇", "😇", "🤐", "🤐", "🤐", "😇", "🥵"})

	// 20bytes of data described by 14 emojis requires at least 2757 length alphabet
	alphabet = make([](string), 2757)
	val, _ = new(big.Int).SetString("0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF", 0) // 20 bytes
	_, err := ToEmojiHash(val, 14, &alphabet)
	if err != nil {
		t.Fatal("mismatched emoji hash alphabet length requirement")
	}

	alphabet = make([](string), 2757-1)
	_, err = ToEmojiHash(val, 14, &alphabet)
	if err == nil {
		t.Fatal("mismatched emoji hash alphabet length requirement")
	}
}
