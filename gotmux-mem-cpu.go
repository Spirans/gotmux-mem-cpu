package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
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
	user		float64
	nice		float64
	sys			float64
	idle		float64
}

func (m *Memory) parse() *Memory {
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
			if err != nil {
				log.Printf("MemTotal parsing error: %v\n", err)
			}
		case strings.HasPrefix(str, "MemFree"):
			m.free, err = strconv.Atoi(strings.Trim(str, "MemFree: kB\n"))
			if err != nil {
				log.Printf("MemFree parsing error: %v\n", err)
			}
		case strings.HasPrefix(str, "Cached"):
			m.cache, err = strconv.Atoi(strings.Trim(str, "Cached: kB\n"))
			if err != nil {
				log.Printf("Cached parsing error: %v\n", err)
			}
		case strings.HasPrefix(str, "Buffers"):
			m.buffers, err = strconv.Atoi(strings.Trim(str, "Buffers: kB\n"))
			if err != nil {
				log.Printf("Buffers parsing error: %v\n", err)
			}
		}
	}
	m.used = m.total - m.cache - m.free - m.buffers
	return m
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
	cpu.user, err = strconv.ParseFloat(strSlice[0], 64)
	if err != nil {
		log.Printf("User CPU parsing error: %v\n", err)
	}
	cpu.nice, err = strconv.ParseFloat(strSlice[1], 64)
	if err != nil {
		log.Printf("Nice CPU parsing error: %v\n", err)
	}
	cpu.sys, err = strconv.ParseFloat(strSlice[2], 64)
	if err != nil {
		log.Printf("Sys CPU parsing error: %v\n", err)
	}
	cpu.idle, err = strconv.ParseFloat(strSlice[3], 64)
	if err != nil {
		log.Printf("IDLE CPU parsing error: %v\n", err)
	}
	return cpu
}

func (cpu *CPU) measureUsage() float64 {
	cpuBefore := CPU{}
	cpuBefore.parse()
	time.Sleep(time.Second)
	cpuAfter := CPU{}
	cpuAfter.parse()
	cpuDiff := CPU{}
	cpuDiff.user = cpuAfter.user - cpuBefore.user
	cpuDiff.sys = cpuAfter.sys - cpuBefore.sys
	cpuDiff.nice = cpuAfter.nice - cpuBefore.nice
	cpuDiff.idle = cpuAfter.idle - cpuBefore.idle
	avg := (cpuDiff.user + cpuDiff.sys + cpuDiff.nice) /
				 (cpuDiff.user + cpuDiff.sys + cpuDiff.nice + cpuDiff.idle) * 100.0
	return avg
}

func main() {
	mem, cpu := Memory{}, CPU{}
	mem.parse()
	fmt.Printf("%v/%vMB %.3v%%\n", mem.used/1024, mem.total/1024, cpu.measureUsage())
}
