package main

import (
	"bufio"
	"fmt"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

func parseFile(path string) map[int]int {
	buf, err := os.Open(path)
	if err != nil {
		fmt.Println("File opening error", err)
		return nil
	}

	defer func() {
		if err = buf.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	scanner := bufio.NewScanner(buf)
	lineCount := 0
	ids := make(map[int]int)
	for scanner.Scan() {
		line := scanner.Text()
		if lineCount == 0 {
			lineCount++
			continue
		}
		for i, val := range strings.Split(line, ",") {
			if val == "x" {
				continue
			}
			id, err := strconv.Atoi(val)
			if err != nil {
				fmt.Println("arg value is not an integer", err)
				return ids
			}
			ids[i] = id
		}
	}
	return ids
}

func main() {
	argc := len(os.Args[1:])

	if argc < 1 {
		fmt.Println("not enough args")
		return
	}

	input := os.Args[1]

	fp, err := filepath.Abs(input)
	if err != nil {
		fmt.Println("Filepath invalid", input, err)
		return
	}

	ids := parseFile(fp)
	// ids = map[int]int{
	// 	0: 3,
	// 	3: 4,
	// 	4: 5,
	// }
	fmt.Println(ids)
	comMod := int64(-1)
	multiple := int64(0)

	var keys []int
	for k := range ids {
		keys = append(keys, k)
	}

	sort.Ints(keys)
	for _, key := range keys {
		id := int64(ids[key])
		offset := int64(key)
		if comMod == -1 {
			comMod = offset
			multiple = id
			continue
		}
		a := id
		b := multiple
		c := offset
		d := comMod
		// if b < a {
		// 	b = id
		// 	a = multiple
		// 	d = offset
		// 	c = comMod

		// }
		comMod = calcComMod(int64(a), int64(b), int64(c), int64(d))
		multiple = multiple * id

		fmt.Println(comMod, multiple)
		comMod = comMod % multiple
		for comMod < 0 {
			comMod += multiple
		}
		fmt.Println(comMod, multiple)
		if a == 13 {
			break
		}
	}
	comMod = comMod % multiple
	for comMod < 0 {
		comMod += multiple
	}
	fmt.Println(comMod, multiple)
	fmt.Println(multiple - comMod)
	for _, idx := range keys {
		if (comMod % int64(ids[idx])) != int64(idx)%int64(ids[idx]) {
			println("Invaid answer for ", ids[idx], idx)
		}
	}
}

func bezoutID(a int64, b int64) (int64, int64) {
	rem := b
	x := a // becomes gcd(a, b)
	s := int64(0)
	y := int64(1) // the coefficient of a
	t := int64(1)
	z := int64(0) // the coefficient of b
	for rem > 0 {
		quotient := x / rem
		x, rem = rem, x%rem
		y, s = s, y-quotient*s
		z, t = t, z-quotient*t
	}
	return y % (b / x), z % (-a / x) // modulus in this way so that y is positive and z is negative
}

func calcComMod(a int64, b int64, ra int64, rb int64) int64 {
	ma, mb := bezoutID(a, b)
	fmt.Println(a, ra, ma, b, rb, mb)
	comMod := big.NewInt(1)
	tmp := big.NewInt(1)
	comMod = comMod.Mul(big.NewInt(rb), big.NewInt(a))
	fmt.Println(rb, "*", a, "*", ma)
	comMod = comMod.Mul(comMod, big.NewInt(ma))
	tmp = tmp.Mul(big.NewInt(ra), big.NewInt(b))
	tmp = tmp.Mul(tmp, big.NewInt(mb))
	fmt.Println(ra, "*", b, "*", mb)
	fmt.Println(comMod, tmp)
	comMod = comMod.Add(comMod, tmp)
	fmt.Println(comMod)
	comMod = comMod.Mod(comMod, big.NewInt(a*b))
	for comMod.Cmp(big.NewInt(0)) < 0 {
		comMod = comMod.Add(comMod, big.NewInt(a*b))
	}
	// if 249309811 == b {
	// 	return 5790158400
	// }
	// if 7229984519 == b {
	// 	return 49169825792
	// }
	// if 4721179890907 == b {
	// 	return 80309445984256
	// }
	// if 193568375527187 == b {
	// 	return 1231453023109120
	// }
	return comMod.Int64()
}
