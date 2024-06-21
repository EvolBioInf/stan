package main

import (
	"flag"
	"fmt"
	"github.com/evolbioinf/clio"
	"github.com/evolbioinf/fasta"
	"github.com/evolbioinf/stan/util"
	"io"
	"math"
	"math/rand"
	"time"
)

func scan(r io.Reader, args ...interface{}) {
	optN := args[0].(float64)
	optL := args[1].(float64)
	optD := args[2].(float64)
	ran := args[3].(*rand.Rand)
	sc := fasta.NewScanner(r)
	for sc.ScanSequence() {
		seq := sc.Sequence()
		seqLen := len(seq.Data())
		deleted := make([]bool, seqLen)
		nDel := rpoi(optN, ran)
		for i := 0; i < nDel; i++ {
			delLen := int(math.Round(ran.NormFloat64()*
				optD + float64(optL)))
			mid := ran.Intn(seqLen)
			start := mid - delLen/2
			end := start + delLen - 1
			if start < 0 {
				start = 0
			}
			if end >= seqLen {
				end = seqLen - 1
			}
			for i := start; i <= end; i++ {
				deleted[i] = true
			}
		}
		d1 := seq.Data()
		d2 := make([]byte, 0, len(d1))
		for i, c := range d1 {
			if !deleted[i] {
				d2 = append(d2, c)
			}
		}
		seq = fasta.NewSequence(seq.Header(), d2)
		fmt.Println(seq)
	}
}
func rpoi(m float64, ran *rand.Rand) int {
	x := math.Exp(-m)
	pr := 1.0
	N := 0
	for pr > x {
		pr *= ran.Float64()
		N++
	}
	return N - 1
}
func main() {
	util.Name("rad")
	u := "rad [option]... [foo.fasta]..."
	p := "Randomly delete regions from sequences."
	e := "rad seq.fasta"
	clio.Usage(u, p, e)
	optV := flag.Bool("v", false, "version")
	optN := flag.Float64("n", 3.0, "mean number of deletions")
	optL := flag.Float64("l", 200.0, "mean length of deletion")
	optD := flag.Float64("d", 200.0, "standard deviation of "+
		"length of deletion")
	optS := flag.Int64("s", 0, "seed for random number "+
		"generator (default internal)")
	flag.Parse()
	if *optV {
		util.Version()
	}
	seed := *optS
	if seed == 0 {
		seed = time.Now().UnixNano()
	}
	source := rand.NewSource(seed)
	ran := rand.New(source)
	files := flag.Args()
	clio.ParseFiles(files, scan, *optN, *optL, *optD, ran)
}
