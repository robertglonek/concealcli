package main

import (
	"bufio"
	"concealcli/aerospike"
	"concealcli/conceal"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"
)

func main() {
	// parse args
	genkey := flag.Bool("genkey", false, "set to generate a new key")
	key := flag.String("key", "key.dat", "specify name of file to read key from")
	input := flag.String("in", "", "name of input file; empty = stdin-pipe")
	output := flag.String("out", "", "name of output file; empty = stdout")
	mapout := flag.String("map", "", "name of file to output the mapping to; empty = stderr")
	nsleep := flag.Duration("sleep", 0, "sleep between lines; e.g. 5us - 5 microseconds")
	flag.Parse()

	var k []byte
	var err error

	// handle key
	if *genkey {
		k, err = conceal.GenKey()
		if err != nil {
			log.Fatal(err)
		}
		err = ioutil.WriteFile(*key, k, 0644)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		k, err = ioutil.ReadFile(*key)
		if err != nil {
			log.Fatal(err)
		}
	}

	// open input/output
	in := os.Stdin
	out := os.Stdout
	mout := os.Stderr
	if *input != "" {
		in, err = os.Open(*input)
		if err != nil {
			log.Fatal(err)
		}
		defer in.Close()
	} else {
		fi, _ := os.Stdin.Stat()
		if (fi.Mode() & os.ModeCharDevice) != 0 {
			log.Fatalf("USAGE ERROR: either use switches to define input/output files, or pipe aerospike log to this program. For more information, run %s -help", os.Args[0])
		}
	}
	if *output != "" {
		out, err = os.Create(*output)
		if err != nil {
			log.Fatal(err)
		}
		defer out.Close()
	}
	if *mapout != "" {
		mout, err = os.Create(*mapout)
		if err != nil {
			log.Fatal(err)
		}
		defer mout.Close()
	}
	_, err = mout.WriteString("source\treplacement\n")
	if err != nil {
		log.Fatal(err)
	}
	nmap := &saveMap{
		f:  mout,
		wg: new(sync.WaitGroup),
	}

	// new aerospike handler
	a, err := aerospike.Init(k, nmap.run)
	if err != nil {
		log.Fatal(err)
	}

	// r/w loop
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		if err = scanner.Err(); err != nil {
			log.Fatal(err)
		}
		newline, err := a.LogLine(scanner.Text())
		if err != nil {
			log.Fatal(err)
		}
		_, err = out.WriteString(newline + "\n")
		if err != nil {
			log.Fatal(err)
		}
		if *nsleep != 0 {
			time.Sleep(*nsleep)
		}
	}

	// wait for map writes to finish
	nmap.wg.Wait()
}

type saveMap struct {
	f  *os.File
	wg *sync.WaitGroup
}

func (s *saveMap) run(source, dest string) error {
	s.wg.Add(1)
	defer s.wg.Done()
	_, err := s.f.WriteString(source + "\t" + dest + "\n")
	return err
}
