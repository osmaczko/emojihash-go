package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"
	"strings"
)

func ToBigBase(value *big.Int, base uint64) (res [](uint64)) {
	toBigBaseImpl(value, base, &res)
	return res
}

func toBigBaseImpl(value *big.Int, base uint64, res *[](uint64)) {
	bigBase := new(big.Int).SetUint64(base)
	quotient := new(big.Int).Div(value, bigBase)
	if quotient.Cmp(new(big.Int).SetUint64(0)) != 0 {
		toBigBaseImpl(quotient, base, res)
	}

	*res = append(*res, new(big.Int).Mod(value, bigBase).Uint64())
}

func ToEmojiHash(value *big.Int, hashLen int, alphabet *[]string) (hash []string, err error) {
	valueBitLen := value.BitLen()
	alphabetLen := new(big.Int).SetInt64(int64(len(*alphabet)))

	indexes := ToBigBase(value, alphabetLen.Uint64())
	if hashLen == 0 {
		hashLen = len(indexes)
	} else if hashLen > len(indexes) {
		prependLen := hashLen - len(indexes)
		for i := 0; i < prependLen; i++ {
			indexes = append([](uint64){0}, indexes...)
		}
	}

	// alphabetLen^hashLen
	possibleCombinations := new(big.Int).Exp(alphabetLen, new(big.Int).SetInt64(int64(hashLen)), nil)

	// 2^valueBitLen
	requiredCombinations := new(big.Int).Exp(new(big.Int).SetInt64(2), new(big.Int).SetInt64(int64(valueBitLen)), nil)

	if possibleCombinations.Cmp(requiredCombinations) == -1 {
		return nil, errors.New("alphabet or hash length is too short to encode given value")
	}

	for _, v := range indexes {
		hash = append(hash, (*alphabet)[v])
	}

	return hash, nil
}

func loadAlphabet(maxLen int) ([]string, error) {
	file, err := os.Open("bare_emojis.txt")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	alphabet := make([]string, 0, 8)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		alphabet = append(alphabet, strings.Replace(scanner.Text(), "\n", "", -1))
	}

	if maxLen > 0 && len(alphabet) > maxLen {
		alphabet = alphabet[:maxLen]
	}

	return alphabet, nil
}

func main() {
	value, _ := new(big.Int).SetString("0x86138b210f21d41c757ae8a5d2a4cb29c1350f73", 0)
	hashLen := 0 // auto
	alphabetLen := 2757

	if len(os.Args[1:]) > 0 {
		readValue, ok := new(big.Int).SetString(os.Args[1], 0)
		if !ok {
			log.Fatal("invalid value")
		}
		value = readValue
	}

	if len(os.Args[1:]) > 1 {
		readValue, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatal(err)
		}
		hashLen = readValue
	}

	if len(os.Args[1:]) > 2 {
		readValue, err := strconv.Atoi(os.Args[3])
		if err != nil {
			log.Fatal(err)
		}
		alphabetLen = readValue
	}

	emojis, err := loadAlphabet(alphabetLen)
	if err != nil {
		log.Fatal(err)
	}

	hash, err := ToEmojiHash(value, hashLen, &emojis)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(hash)
}
