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
	"math"
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
	for {
		str, err := fileRead.ReadString('\n')
		if err != nil {
			break
		}
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
	str, _ := fileRead.ReadString('\n')
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
	if *interval < 1 {
		*interval = 1
	}
	time.Sleep(time.Duration(*interval) * time.Second)
	cpuAfter := CPU{}
	cpuAfter.parse()
	cpuDiff := CPU{
		user	:		cpuAfter.user - cpuBefore.user,
		sys 	:		cpuAfter.sys - cpuBefore.sys,
		nice	: 		cpuAfter.nice - cpuBefore.nice,
		idle	:		cpuAfter.idle - cpuBefore.idle,
	}
	cpu.utl = float64(cpuDiff.user + cpuDiff.sys + cpuDiff.nice) /
				float64(cpuDiff.user + cpuDiff.sys + cpuDiff.nice +
						cpuDiff.idle)*100.0
	cpu.utl = round(cpu.utl, 1)
}

func powerline(value float64, background, foreground *string) string {
	line := []byte("▁▂▃▄▅▆▇█")
	tick := len(line) / 3
	tickPosition := (tick * int(value)) / 100
	colorLine := "#[fg=green]"
	switch {
	case value > 75:
		colorLine = "#[fg=red]"
	case value > 50:
		colorLine = "#[fg=yellow]"
	}
	return fmt.Sprintf("%s%s#[fg=%s,bg=%s]", colorLine,
						line[tickPosition*3:tickPosition*3+3],
						*foreground, *background)
}

func round(value float64, offset int) (float64) {
	shift := math.Pow(10, float64(offset))
	return math.Floor((value * shift)+.5) / shift;
}

var interval = flag.Int("i", 2,	"Interval in seconds for calculating CPU " +
							"utilization")
var background = flag.String("b", "black", "Background color for cpu " +
							"and mem status bar")
var foreground = flag.String("f", "white", "Foreground color for cpu " +
							"and mem status bar")
func main() {
	flag.Parse()
	mem, cpu := Memory{}, CPU{}
	cpu.utilization(interval)
	mem.parse()
	cpuLine := powerline(cpu.utl, background, foreground)
	memLine := powerline(float64(mem.used) / float64(mem.total) * 100,
						background,foreground)
	fmt.Printf("#[fg=%s,bg=%s] %v/%vMB %v %.1f%% %v%s", *foreground, *background,
				mem.used/1024, mem.total/1024, string(memLine), cpu.utl,
				string(cpuLine), "#[default]")
}
