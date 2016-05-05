package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

//Memory - type for showing memory usage
type Memory struct {
	total   int
	used    int
	free    int
	cache   int
	buffers int
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
				log.Printf("MemTotal parsing error: %v", err)
			}
		case strings.HasPrefix(str, "MemFree"):
			m.free, err = strconv.Atoi(strings.Trim(str, "MemFree: kB\n"))
			if err != nil {
				log.Printf("MemFree parsing error: %v", err)
			}
		case strings.HasPrefix(str, "Cached"):
			m.cache, err = strconv.Atoi(strings.Trim(str, "Cached: kB\n"))
			if err != nil {
				log.Printf("Cached parsing error: %v", err)
			}
		case strings.HasPrefix(str, "Buffers"):
			m.buffers, err = strconv.Atoi(strings.Trim(str, "Buffers: kB\n"))
			if err != nil {
				log.Printf("Buffers parsing error: %v", err)
			}
		}
	}
	m.used = m.total - m.cache - m.free - m.buffers
	return m
}

func main() {
	m := Memory{}
	m.parse()
	// fmt.Printf("#[fg=blue,bg=black,bright] %v/%vMB #[default]", m.used/1024, m.total/1024)
	fmt.Printf("%v/%vMB\n", m.used/1024, m.total/1024)
}
