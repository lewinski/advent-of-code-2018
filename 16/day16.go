package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

func atoi(x string) int {
	i, _ := strconv.Atoi(x)
	return i
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func fileLines(path string) []string {
	f, err := os.Open(path)
	check(err)

	lines := make([]string, 0)
	s := bufio.NewScanner(f)
	for s.Scan() {
		lines = append(lines, s.Text())
	}

	return lines
}

type registers [4]int

func decodeRegisters(s string) (r registers) {
	re := regexp.MustCompile("^(?:Before:|After:)\\s+\\[(\\d+), (\\d+), (\\d+), (\\d+)\\]$")
	match := re.FindStringSubmatch(s)
	r[0] = atoi(match[1])
	r[1] = atoi(match[2])
	r[2] = atoi(match[3])
	r[3] = atoi(match[4])
	return
}

type instruction struct {
	opcode, a, b, c int
}

func decodeInstruction(s string) (i instruction) {
	re := regexp.MustCompile("^(\\d+) (\\d+) (\\d+) (\\d+)$")
	match := re.FindStringSubmatch(s)
	i.opcode = atoi(match[1])
	i.a = atoi(match[2])
	i.b = atoi(match[3])
	i.c = atoi(match[4])
	return i
}

type operation func(instruction, registers) registers

func addr(i instruction, r registers) registers {
	r[i.c] = r[i.a] + r[i.b]
	return r
}

func addi(i instruction, r registers) registers {
	r[i.c] = r[i.a] + i.b
	return r
}

func mulr(i instruction, r registers) registers {
	r[i.c] = r[i.a] * r[i.b]
	return r
}

func muli(i instruction, r registers) registers {
	r[i.c] = r[i.a] * i.b
	return r
}

func banr(i instruction, r registers) registers {
	r[i.c] = r[i.a] & r[i.b]
	return r
}

func bani(i instruction, r registers) registers {
	r[i.c] = r[i.a] & i.b
	return r
}

func borr(i instruction, r registers) registers {
	r[i.c] = r[i.a] | r[i.b]
	return r
}

func bori(i instruction, r registers) registers {
	r[i.c] = r[i.a] | i.b
	return r
}

func setr(i instruction, r registers) registers {
	r[i.c] = r[i.a]
	return r
}

func seti(i instruction, r registers) registers {
	r[i.c] = i.a
	return r
}

func gtir(i instruction, r registers) registers {
	if i.a > r[i.b] {
		r[i.c] = 1
	} else {
		r[i.c] = 0
	}
	return r
}

func gtri(i instruction, r registers) registers {
	if r[i.a] > i.b {
		r[i.c] = 1
	} else {
		r[i.c] = 0
	}
	return r
}

func gtrr(i instruction, r registers) registers {
	if r[i.a] > r[i.b] {
		r[i.c] = 1
	} else {
		r[i.c] = 0
	}
	return r
}

func eqir(i instruction, r registers) registers {
	if i.a == r[i.b] {
		r[i.c] = 1
	} else {
		r[i.c] = 0
	}
	return r
}

func eqri(i instruction, r registers) registers {
	if r[i.a] == i.b {
		r[i.c] = 1
	} else {
		r[i.c] = 0
	}
	return r
}

func eqrr(i instruction, r registers) registers {
	if r[i.a] == r[i.b] {
		r[i.c] = 1
	} else {
		r[i.c] = 0
	}
	return r
}

var handlers = map[string]operation{
	"addr": addr, "addi": addi,
	"mulr": mulr, "muli": muli,
	"banr": banr, "bani": bani,
	"borr": borr, "bori": bori,
	"setr": setr, "seti": seti,
	"gtir": gtir, "gtri": gtri, "gtrr": gtrr,
	"eqir": eqir, "eqri": eqri, "eqrr": eqrr,
}

var opHandlers = []operation{
	// 0  [                        banr:47                 eqir:47 eqri:47 eqrr:47 gtir:47 gtri:47 gtrr:47                                ]
	banr,
	// 1  [                                                eqir:50         eqrr:50         gtri:50                                        ]
	eqrr,
	// 2  [                                                                        gtir:59                                         setr:59]
	setr,
	// 3  [                                                eqir:55                                 gtrr:55                                ]
	eqir,
	// 4  [addi:47 addr:47 bani:47 banr:47 bori:47 borr:47                         gtir:47 gtri:47 gtrr:47 muli:47 mulr:47 seti:47 setr:47]
	bori,
	// 5  [                bani:45 banr:45                 eqir:45 eqri:45 eqrr:45 gtir:45 gtri:45 gtrr:45 muli:45 mulr:45 seti:45 setr:45]
	muli,
	// 6  [                bani:54 banr:54                 eqir:54 eqri:54 eqrr:54 gtir:54 gtri:54 gtrr:54                                ]
	bani,
	// 7  [        addr:46 bani:46 banr:46         borr:46         eqri:46         gtir:46 gtri:46 gtrr:46 muli:46 mulr:46 seti:46 setr:46]
	borr,
	// 8  [                                                eqir:56 eqri:56 eqrr:56 gtir:56 gtri:56 gtrr:56                                ]
	gtir,
	// 9  [                                                                                        gtrr:47                                ]
	gtrr,
	// 10 [addi:56                                                                                                 mulr:56 seti:56        ]
	addi,
	// 11 [                                                eqir:59                         gtri:59 gtrr:59                                ]
	gtri,
	// 12 [                                                eqir:56 eqri:56                 gtri:56 gtrr:56                                ]
	eqri,
	// 13 [addi:53 addr:53                                                                                 muli:53 mulr:53                ]
	addr,
	// 14 [                bani:49 banr:49                         eqri:49 eqrr:49 gtir:49 gtri:49                 mulr:49 seti:49        ]
	mulr,
	// 15 [                bani:59 banr:59                 eqir:59 eqri:59 eqrr:59 gtir:59 gtri:59 gtrr:59                 seti:59        ]
	seti,
}

func main() {
	// read file
	file := "input"
	if len(os.Args) > 1 {
		file = os.Args[1]
	}
	lines := fileLines(file)

	// part 1
	sampleCount := 0
	sampleEnd := len(lines)
	for i := 0; i < len(lines); i += 4 {
		if lines[i] == "" && lines[i+1] == "" {
			sampleEnd = i + 2
			break
		}

		before := decodeRegisters(lines[i])
		ins := decodeInstruction(lines[i+1])
		after := decodeRegisters(lines[i+2])

		count := 0
		valid := []string{}
		for name, handler := range handlers {
			test := handler(ins, before)
			if test == after {
				count++
				valid = append(valid, name)
			}
		}

		test := opHandlers[ins.opcode](ins, before)
		if test != after {
			fmt.Println("opcode", ins.opcode, "is invalid")
		}

		if count >= 3 {
			sampleCount++
		}
	}

	fmt.Println(sampleCount)

	// part 2
	reg := registers{}
	for i := sampleEnd; i < len(lines); i++ {
		ins := decodeInstruction(lines[i])
		reg = opHandlers[ins.opcode](ins, reg)
	}
	fmt.Println(reg)
}
