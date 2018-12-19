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

type registers [6]int

type vmstate struct {
	r     registers
	ip    int
	ipreg int
}

type instruction struct {
	opcode  string
	a, b, c int
}

func decodeInstruction(s string) (i instruction) {
	ipre := regexp.MustCompile("^#ip (\\d+)$")
	if ipre.MatchString(s) {
		match := ipre.FindStringSubmatch(s)
		i.opcode = "ip"
		i.a = atoi(match[1])
		return i
	}
	insre := regexp.MustCompile("^(\\S+) (\\d+) (\\d+) (\\d+)$")
	match := insre.FindStringSubmatch(s)
	i.opcode = match[1]
	i.a = atoi(match[2])
	i.b = atoi(match[3])
	i.c = atoi(match[4])
	return i
}

// registers = 0:a 1:b 2:c 3:ip 4:e 5:f

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

func main() {
	// read file
	file := "input"
	if len(os.Args) > 1 {
		file = os.Args[1]
	}
	lines := fileLines(file)

	// parse input file
	var initvm vmstate
	var program []instruction
	for _, l := range lines {
		ins := decodeInstruction(l)
		if ins.opcode == "ip" {
			initvm.ipreg = ins.a
		} else {
			program = append(program, ins)
		}
	}

	// part 1
	vm := initvm
	for {
		ins := program[vm.ip]
		vm.r[vm.ipreg] = vm.ip
		vm.r = handlers[ins.opcode](ins, vm.r)
		vm.ip = vm.r[vm.ipreg]
		vm.ip++
		if vm.ip < 0 || vm.ip >= len(program) {
			break
		}
	}
	fmt.Println(vm)

	// part 2
	vm = initvm
	vm.r[0] = 1
	for {
		ins := program[vm.ip]
		if vm.ip == 2 { // cheat once initial condition is set up
			total := 0
			for i := 1; i < vm.r[1]; i++ {
				if vm.r[1]%i == 0 {
					total += i
				}
			}
			vm.r[0] = total
			break
		}
		vm.r[vm.ipreg] = vm.ip
		vm.r = handlers[ins.opcode](ins, vm.r)
		vm.ip = vm.r[vm.ipreg]
		vm.ip++
		if vm.ip < 0 || vm.ip >= len(program) {
			break
		}
	}
	fmt.Println(vm)
}
