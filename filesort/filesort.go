package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
)

type sorter struct {
	source *os.File
	parts  []*os.File
	ints   sort.IntSlice
}

func (s *sorter) parting() {
	s.ints = make(sort.IntSlice, 0, 1e1)

	for n := range s.read(s.source) {
		s.ints = append(s.ints, n)

		if len(s.ints) >= cap(s.ints)-1 {
			s.savePart()
		}
	}

	s.savePart()
}
func (s *sorter) savePart() {
	if len(s.ints) == 0 {
		return
	}

	name := fmt.Sprintf("%s_part%d", s.source.Name(), len(s.parts))
	part, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		panic(fmt.Errorf("error creating %s: %w", name, err))
	}

	s.parts = append(s.parts, part)
	s.ints.Sort()

	for _, n := range s.ints {
		fmt.Fprintln(part, n)
	}
	part.Sync()
	part.Seek(0, 0)

	s.ints = s.ints[:0]
}
func (s *sorter) merge() {
	output, err := os.OpenFile(fmt.Sprintf("%s_sorted", s.source.Name()), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
	if err == nil {
		defer output.Close()
	} else {
		panic(fmt.Errorf("cannot create sorted output: %w", err))
	}

	caching := make(map[chan int]*int)
	for _, fd := range s.parts {
		caching[s.read(fd)] = nil
	}

	for {
		var min chan int

		for intc := range caching {
			if caching[intc] == nil {
				intv, ok := <-intc
				if ok {
					caching[intc] = &intv
				} else {
					delete(caching, intc)
					continue
				}
			}

			if min == nil || *caching[min] > *caching[intc] {
				min = intc
			}
		}

		if len(caching) == 0 {
			break
		} else {
			fmt.Fprintln(output, *caching[min])
			caching[min] = nil
		}
	}
}
func (s *sorter) read(fd *os.File) (output chan int) {
	output = make(chan int)

	go func() {
		defer fd.Close()

		sc := bufio.NewScanner(fd)
		for sc.Scan() {
			text := sc.Text()
			if len(text) > 0 {
				n, err := strconv.Atoi(text)
				if err != nil {
					panic(fmt.Errorf("text: %s, err: %s", text, err))
				}
				output <- n
			}
		}

		close(output)
	}()

	return
}
func (s *sorter) sort() {
	s.parting()
	s.merge()
}
