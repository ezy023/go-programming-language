// count the number of bits different between two SHA256 hashes
package main

import (
	"crypto/sha256"
	"fmt"
)

func numTable() {
	for i := 0; i < 16; i++ {
		fmt.Printf("%[1]d\t %02[1]x\t %08[1]b\n", i)
	}
}

func bitsRep(i uint8) {
	fmt.Printf("%10.08b\n", i)
	fmt.Printf("%10.08b\n", i>>1)
	j := i >> 1
	fmt.Printf("%10.08b\t ^(i>>1)\n", ^j)
	fmt.Printf("%10.08b\t i&^(i>>1)\n", i&^(i>>1))
}

func bitCount(b byte) uint {
	var count uint
	for i := 0; i < 8; i++ {
		if (b >> i & 1) == 1 {
			count++
		}
	}
	return count
}

func bitComp(a, b [32]byte) uint {
	var count uint
	for i := 0; i < len(a); i++ {
		for j := 0; j < 8; j++ {
			if ((a[i] >> j) & 1) != ((b[i] >> j) & 1) {
				count++
			}
		}
	}
	return count
}

func main() {
	var i, j, f byte = 0x5a, 0xdb, 0xfe
	fmt.Printf("%[1]d\t %02[1]x\t %08[1]b\n", i)
	fmt.Printf("%[1]d\t %02[1]x\t %08[1]b\n", j)
	fmt.Printf("%[1]d\t %02[1]x\t %08[1]b\n", f)
	fmt.Printf("%[1]d\t %02[1]x\t %08[1]b\n", ^f)
	fmt.Printf("-----------\n")
	// bitsRep(uint8(i))

	a := sha256.Sum256([]byte("Hello World"))
	b := sha256.Sum256([]byte("hello World"))
	fmt.Printf("Diff Bits: %d\n", bitComp(a, b))
}
