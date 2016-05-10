package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	"flag"
)

//Memory - type for showing memory usage
type Memory struct {
	total   int
	used    int
	free    int
	cache   int
	buffers int
}

//CPU - type for showing CPU usage
type CPU struct {
	user		int
	nice		int
	sys			int
	idle		int
	utl			float64
}

func checkParsingError(e error, msg string) {
	if e != nil {
		log.Printf("%v parsing error: %v\n", msg, e)
	}
}

func (m *Memory) parse() {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	fileRead := bufio.NewReader(file)
	for i := 0; i < 5; i++ {
		str, _ := fileRead.ReadString(10)
		switch {
		case strings.HasPrefix(str, "MemTotal"):
			m.total, err = strconv.Atoi(strings.Trim(str, "MemTotal: kB\n"))
			checkParsingError(err, "MemTotal")
		case strings.HasPrefix(str, "MemFree"):
			m.free, err = strconv.Atoi(strings.Trim(str, "MemFree: kB\n"))
			checkParsingError(err, "MemFree")
		case strings.HasPrefix(str, "Cached"):
			m.cache, err = strconv.Atoi(strings.Trim(str, "Cached: kB\n"))
			checkParsingError(err, "Cached")
		case strings.HasPrefix(str, "Buffers"):
			m.buffers, err = strconv.Atoi(strings.Trim(str, "Buffers: kB\n"))
			checkParsingError(err, "Buffers")
		}
	}
	m.used = m.total - m.cache - m.free - m.buffers
}

func (cpu *CPU) parse() *CPU {
	file, err := os.Open("/proc/stat")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	fileRead := bufio.NewReader(file)
	str, _ := fileRead.ReadString(10)
	str = strings.Trim(str, "cpu ")
	strSlice := strings.Split(str, " ")
	cpu.user, err = strconv.Atoi(strSlice[0])
	checkParsingError(err, "UserCPU")
	cpu.nice, err = strconv.Atoi(strSlice[1])
	checkParsingError(err, "NiceCPU")
	cpu.sys, err = strconv.Atoi(strSlice[2])
	checkParsingError(err, "SysCPU")
	cpu.idle, err = strconv.Atoi(strSlice[3])
	checkParsingError(err, "IdleCPU")
	return cpu
}

func (cpu *CPU) utilization(interval *int) {
	cpuBefore := CPU{}
	cpuBefore.parse()
	time.Sleep(time.Duration(*interval) * time.Second)
	cpuAfter := CPU{}
	cpuAfter.parse()
	cpuDiff := CPU{
		user	:		cpuAfter.user - cpuBefore.user,
		sys		:		cpuAfter.sys - cpuBefore.sys,
		nice	: 	cpuAfter.nice - cpuBefore.nice,
		idle	:		cpuAfter.idle - cpuBefore.idle,
	}
	cpu.utl = float64(cpuDiff.user + cpuDiff.sys + cpuDiff.nice) /
				 float64(cpuDiff.user + cpuDiff.sys + cpuDiff.nice + cpuDiff.idle)*100.0
}

func main() {
	mem, cpu := Memory{}, CPU{}
	mem.parse()
	var interval = flag.Int("interval", 2,
													"Interval for calculate CPU utilization, 2sec by default")
	flag.Parse()
	cpu.utilization(interval)
	fmt.Printf("%v/%vMB %.3v%%\n", mem.used/1024, mem.total/1024, cpu.utl)
}
