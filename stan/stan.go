package main

import (
	"flag"
	"fmt"
	"github.com/evolbioinf/clio"
	"github.com/evolbioinf/fasta"
	"github.com/evolbioinf/nwk"
	"github.com/evolbioinf/stan/util"
	"log"
	"math"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"
)

type Region struct {
	s, e int
}
type RegionSlice []Region
type Mutation struct {
	background, marker int
}
type Haplotypes struct {
	hap            [][]byte
	pos            []int
	ms, bs, mn, nn int
	n2m            map[int]Mutation
	l2h            map[int]int
	r              *rand.Rand
}

func (r RegionSlice) Len() int {
	return len(r)
}
func (r RegionSlice) Less(i, j int) bool {
	return r[i].s < r[j].s
}
func (r RegionSlice) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}
func checkDir(dir string, overwrite bool) {
	_, err := os.Stat(dir)
	if err == nil {
		if overwrite {
			err = os.RemoveAll(dir)
			util.Check(err)
		} else {
			m := fmt.Sprintf("directory %s already exists", dir)
			fmt.Fprintf(os.Stderr, "%s\n", m)
			os.Exit(1)
		}
	}
}
func genPartCoal(np, n int, b string,
	ran *rand.Rand) *nwk.Node {
	var root *nwk.Node
	tree := make([]*nwk.Node, 2*np-1)
	for i := 0; i < 2*np-1; i++ {
		tree[i] = nwk.NewNode()
	}
	for i := 0; i < np; i++ {
		l := b + strconv.Itoa(i+1)
		tree[i].Label = l
	}
	root = tree[2*np-2]
	for i := 0; i < np; i++ {
		tree[i].HasLength = true
	}
	t := 0.0
	f := float64(np) / float64(n)
	for i := np; i > 1; i-- {
		lambda := f * float64(np*(np-1)/2)
		t += ran.ExpFloat64() / lambda
		j := 2*np - i
		tree[j].Length = t
		tree[j].HasLength = true
	}
	for i := np; i > 1; i-- {
		p := tree[2*np-i]
		r := ran.Intn(i)
		c := tree[r]
		addChild(p, c)
		tree[r] = tree[i-1]
		r = ran.Intn(i - 1)
		c = tree[r]
		addChild(p, c)
		tree[r] = p
	}
	return root
}
func addChild(p, c *nwk.Node) {
	c.Parent = p
	if p.Child == nil {
		p.Child = c
	} else {
		v := p.Child
		for v.Sib != nil {
			v = v.Sib
		}
		v.Sib = c
	}
}
func branchLength(v *nwk.Node) {
	var l float64
	if v == nil {
		return
	}
	branchLength(v.Child)
	if v.Parent != nil {
		l = v.Parent.Length - v.Length
	}
	v.Length = l
	branchLength(v.Sib)
}
func mutate(v *nwk.Node, mm, bm float64, r *rand.Rand,
	n2m map[int]Mutation) {
	if v == nil {
		return
	}
	mn := calcMut(mm, v.Length, r)
	bn := calcMut(bm, v.Length, r)
	u := Mutation{marker: mn, background: bn}
	n2m[v.Id] = u
	mutate(v.Child, mm, bm, r, n2m)
	mutate(v.Sib, mm, bm, r, n2m)
}
func calcMut(t, l float64, r *rand.Rand) int {
	lambda := t * l / 2.0
	x := math.Exp(-lambda)
	p := 1.0
	c := 0
	for p > x {
		p *= r.Float64()
		c++
	}
	return c - 1
}
func genHaps(v *nwk.Node, haps *Haplotypes, m2p,
	b2p map[int]int) {
	if v == nil {
		return
	}
	m := haps.n2m[v.Id].marker
	n := haps.mn + haps.nn
	for i := 0; i < m; i++ {
		p := haps.r.Intn(haps.ms)
		haps.pos = append(haps.pos, m2p[p])
		ss := make([]byte, n)
		if v.Child != nil {
			recSeg(v.Child, ss, haps.l2h)
		} else {
			ss[haps.l2h[v.Id]] = 1
		}
		haps.hap = append(haps.hap, ss)
	}
	m = haps.n2m[v.Id].background
	for i := 0; i < m; i++ {
		p := haps.r.Intn(haps.bs)
		haps.pos = append(haps.pos, b2p[p])
		ss := make([]byte, n)
		if v.Child != nil {
			recSeg(v.Child, ss, haps.l2h)
		} else {
			ss[haps.l2h[v.Id]] = 1
		}
		haps.hap = append(haps.hap, ss)
	}
	genHaps(v.Child, haps, m2p, b2p)
	genHaps(v.Sib, haps, m2p, b2p)
}
func recSeg(v *nwk.Node, ss []byte, l2h map[int]int) {
	if v == nil {
		return
	}
	if v.Child == nil {
		ss[l2h[v.Id]] = 1
	}
	recSeg(v.Child, ss, l2h)
	recSeg(v.Sib, ss, l2h)
}
func printSeqs(seqs []*fasta.Sequence, dir, name string) {
	err := os.Mkdir(dir, 0750)
	util.Check(err)
	for i, seq := range seqs {
		p := dir + "/" + name + strconv.Itoa(i+1) +
			".fasta"
		f := util.Create(p)
		fmt.Fprintf(f, "%s\n", seq)
		f.Close()
	}
}
func main() {
	util.Name("stan")
	u := "stan [option]..."
	p := "Simulate targets and neighbors with marker regions."
	e := "stan -t 7 -n 13 -r 1501-2000,3501-4000"
	clio.Usage(u, p, e)
	optT := flag.Int("t", 10, "targets")
	optN := flag.Int("n", 10, "neighbors")
	optTT := flag.String("T", "targets", "target directory")
	optNN := flag.String("N", "neighbors", "neighbor directory")
	optL := flag.Int("l", 10000, "sequence length")
	optR := flag.String("r", "4501-5500", "marker regions")
	optM := flag.Float64("m", 0.01, "background mutation rate, "+
		"theta per nucleotide")
	optMM := flag.Float64("M", -0.1, "marker region mutation "+
		"rate in neighbors, theta per nucleotide; "+
		"delete if negative")
	optC := flag.Bool("c", false, "print coalescent tree")
	optA := flag.Bool("a", false, "print haplotypes")
	optO := flag.Bool("o", false, "overwrite existing directories")
	optS := flag.Int("s", 0, "seed for random number generator")
	optV := flag.Bool("v", false, "version")
	flag.Parse()
	if *optV {
		util.Version()
	}
	checkDir(*optTT, *optO)
	checkDir(*optNN, *optO)
	seed := int64(*optS)
	if seed == 0 {
		seed = time.Now().UnixNano()
	}
	ran := rand.New(rand.NewSource(seed))
	strs := strings.Split(*optR, ",")
	var regions []Region
	for _, str := range strs {
		coords := strings.Split(str, "-")
		s, err := strconv.Atoi(coords[0])
		util.Check(err)
		e, err := strconv.Atoi(coords[1])
		if e > *optL {
			log.Fatalf("Marker region %d-%d extends beyond "+
				"the end of the sequence.\n", s, e)
		}
		reg := Region{s: s - 1, e: e - 1}
		regions = append(regions, reg)
	}
	sort.Sort(RegionSlice(regions))
	n := *optT + *optN
	tc := genPartCoal(*optT, n, "t", ran)
	nc := genPartCoal(*optN, n, "n", ran)
	root := nwk.NewNode()
	tc.Parent = root
	nc.Parent = root
	tc.Sib = nc
	root.Child = tc
	root.Length = tc.Length
	if root.Length < nc.Length {
		root.Length = nc.Length
	}
	root.Length += ran.ExpFloat64()

	branchLength(root)
	ms := 0
	for _, region := range regions {
		ms += region.e - region.s + 1
	}
	bs := *optL - ms
	mm := float64(ms) * *optM
	bm := float64(bs) * *optM
	node2mut := make(map[int]Mutation)
	tc.Sib = nil
	mutate(tc, mm, bm, ran, node2mut)
	tc.Sib = nc
	mm = float64(ms) * *optMM
	mutate(nc, mm, bm, ran, node2mut)
	if *optC {
		fmt.Printf("%s\n", root)
	}
	leaf2hap := make(map[int]int)
	for i := 1; i <= *optT; i++ {
		leaf2hap[i] = i - 1
	}
	start := 2 * *optT
	end := start + *optN
	for i, j := start, *optT; i <= end; i, j = i+1, j+1 {
		leaf2hap[i] = j
	}
	haps := new(Haplotypes)
	haps.ms = ms
	haps.bs = bs
	haps.mn = *optT
	haps.nn = *optN
	haps.l2h = leaf2hap
	haps.n2m = node2mut
	haps.r = ran
	m2p := make(map[int]int)
	isMarker := make(map[int]bool)
	i := 0
	for _, region := range regions {
		for j := region.s; j <= region.e; j++ {
			m2p[i] = j
			i++
			isMarker[j] = true
		}
	}
	b2p := make(map[int]int)
	i = 0
	for j := 0; j < *optL; j++ {
		if !isMarker[j] {
			b2p[i] = j
			i++
		}
	}
	genHaps(root, haps, m2p, b2p)
	if *optA {
		fmt.Printf("Positions:")
		for _, p := range haps.pos {
			fmt.Printf(" %d", p+1)
		}
		fmt.Printf("\n")
		w := tabwriter.NewWriter(os.Stdout, 0, 1, 1, ' ', 0)
		for j := 0; j < haps.mn; j++ {
			fmt.Fprintf(w, "t%d\t", j+1)
			for i := 0; i < len(haps.hap); i++ {
				c := strconv.Itoa(int(haps.hap[i][j]))
				fmt.Fprintf(w, "%s", c)
			}
			fmt.Fprintf(w, "\t\n")
		}
		for j := haps.mn; j < haps.mn+haps.nn; j++ {
			fmt.Fprintf(w, "n%d\t", j-haps.mn+1)
			for i := 0; i < len(haps.hap); i++ {
				c := strconv.Itoa(int(haps.hap[i][j]))
				fmt.Fprintf(w, "%s", c)
			}
			fmt.Fprintf(w, "\t\n")
		}
		w.Flush()
	}
	dic := []byte{'A', 'C', 'G', 'T'}
	anc := make([]byte, 0)
	for i := 0; i < *optL; i++ {
		r := ran.Intn(4)
		anc = append(anc, dic[r])
	}
	al := make([][]byte, 0)
	for i := 0; i < *optL; i++ {
		s := make([]byte, 0)
		for j := 0; j < n; j++ {
			s = append(s, anc[i])
		}
		al = append(al, s)
	}
	for i, ss := range haps.hap {
		p := haps.pos[i]
		c1 := anc[p]
		r := ran.Intn(4)
		c2 := dic[r]
		for c2 == c1 {
			r = ran.Intn(4)
			c2 = dic[r]
		}
		for j, s := range ss {
			if s == 1 {
				al[p][j] = c2
			}
		}
	}
	if *optMM < 0 {
		start := 0
		row := 0
		cols := len(al[0])
		for _, region := range regions {
			end := region.s
			for i := start; i < end; i++ {
				for j := *optT; j < cols; j++ {
					al[row][j] = al[i][j]
				}
				row++
			}
			start = region.e + 1
		}
		end := len(al)
		for i := start; i < end; i++ {
			for j := *optT; j < cols; j++ {
				al[row][j] = al[i][j]
			}
			row++
		}
		if row < len(al) {
			for j := *optT; j < cols; j++ {
				al[row][j] = 0
			}
		}
	}
	targets := make([]*fasta.Sequence, 0)
	neighbors := make([]*fasta.Sequence, 0)
	for i := 0; i < *optT; i++ {
		h := "t" + strconv.Itoa(i+1)
		d := make([]byte, 0)
		for j := 0; j < len(al); j++ {
			d = append(d, al[j][i])
		}
		seq := fasta.NewSequence(h, d)
		targets = append(targets, seq)
	}
	n = *optT + *optN
	for i := *optT; i < n; i++ {
		h := "n" + strconv.Itoa(i-*optT+1)
		d := make([]byte, 0)
		for j := 0; j < len(al); j++ {
			c := al[j][i]
			if c == 0 {
				break
			}
			d = append(d, c)
		}
		seq := fasta.NewSequence(h, d)
		neighbors = append(neighbors, seq)
	}
	printSeqs(targets, *optTT, "t")
	printSeqs(neighbors, *optNN, "n")
}
