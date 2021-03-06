package constraints

import (
	"strconv"

	"github.com/vale1410/bule/glob"
	"github.com/vale1410/bule/sat"
	"github.com/vale1410/bule/sorters"
)

func (pb *Threshold) TranslateBySN() {
	pb.TransTyp = CSN
	pb.Normalize(LE, true)
	glob.A(pb.Typ == LE, "does not work on OPT or ==, but we have", pb.Typ)
	pb.SortDescending()
	sn := NewSortingNetwork(*pb)
	sn.CreateSorter()

	//glob.D("size of comparators", len(sn.Sorter.Comparators))

	//PrintThresholdTikZ("sn.tex", []SortingNetwork{sn})

	wh := 1
	var which [8]bool

	switch wh {
	case 1:
		which = [8]bool{false, false, false, true, true, true, false, false}
	case 2:
		which = [8]bool{false, false, false, true, true, true, false, true}
	case 3:
		which = [8]bool{false, true, true, true, true, true, true, false}
	case 4:
		which = [8]bool{false, true, true, true, true, true, true, true}
	}

	pred := sat.Pred("auxSN_" + strconv.Itoa(pb.Id))
	pb.Clauses.AddClauseSet(CreateEncoding(sn.LitIn, which, []sat.Literal{}, "BnB", pred, sn.Sorter))
}

//this construction is based on AtMost threshold constraints

type SortingNetwork struct {
	pb     Threshold
	Tare   int64
	Sorter sorters.Sorter
	Bags   [][]sat.Literal
	LitIn  []sat.Literal //Bags flattened, input to Sorter
	typ    sorters.SortingNetworkType
}

// TODO: update construction of sorting network!
func NewSortingNetwork(pb Threshold) (sn SortingNetwork) {
	// much more configuration in the future
	sn.pb = pb
	sn.typ = sorters.OddEven
	return
}

func (t *SortingNetwork) CreateSorter() {

	glob.A(!t.pb.Empty(), "No empty at this point.")

	t.CreateBags()

	layers := make([]sorters.Sorter, len(t.Bags))

	for i, bag := range t.Bags {

		layers[i] = sorters.CreateSortingNetwork(len(bag), -1, t.typ)

		t.LitIn = append(t.LitIn, bag...)
	}

	t.Sorter.In = make([]int, 0, len(t.LitIn))
	t.Sorter.Out = make([]int, 0, len(t.LitIn))

	offset := 2

	// determine the constant and what to add on both sides
	layerPow2 := int64(1 << uint(len(t.Bags)))

	tare := layerPow2 - ((t.pb.K + 1) % layerPow2)
	tare = tare % layerPow2
	t.Tare = tare
	bTare := Binary(tare)

	// output of sorter in layer $i-1$
	bIn := make([]int, 0)

	finalMapping := make(map[int]int, len(t.Sorter.In))

	for i, layer := range layers {

		offset = layer.Normalize(offset, []int{})
		t.Sorter.Comparators = append(t.Sorter.Comparators, layer.Comparators...)

		t.Sorter.In = append(t.Sorter.In, layer.In...)

		size := len(bIn) + len(layers[i].In)

		mergeIn := make([]int, 0, size)
		mergeIn = append(mergeIn, bIn...)
		mergeIn = append(mergeIn, layer.Out...)

		merger := sorters.CreateSortingNetwork(size, len(bIn), t.typ)
		offset = merger.Normalize(offset, mergeIn)

		// halving circuit:

		odd := 1

		if i < len(bTare) && bTare[i] == 1 {
			odd = 0
			bIn = make([]int, (len(merger.Out)+1)/2)
		} else {
			bIn = make([]int, len(merger.Out)/2)
		}

		// Alternate depending on bTare
		for j, x := range merger.Out {
			if j%2 == odd {
				bIn[j/2] = x
			} else if i < len(layers)-1 { // not in last layer, but else
				finalMapping[x] = -1
			}
		}

		t.Sorter.Comparators = append(t.Sorter.Comparators, merger.Comparators...)

	}

	// outLastLayer identifies the nth output in the last layer
	outLastLayer := ((t.pb.K + 1 + tare) / int64(layerPow2)) - 1

	// debug stuff:
	//glob.D("len last layer:", len(bIn), "kth output in last layer: ", outLastLayer)
	//glob.D("K+1+tar", t.pb.K+1+tare, "n layers", layerPow2)

	idSetToZero := bIn[outLastLayer]

	// and propagate the rest backwards
	setTo := -1 // dont care
	for _, id := range t.Sorter.ComputeOut() {
		if id == idSetToZero {
			setTo = 0
		}
		if _, ok := finalMapping[id]; !ok {
			finalMapping[id] = setTo
		}
	}

	t.Sorter.PropagateBackwards(finalMapping)
	t.Sorter.Normalize(2, []int{})

	//fmt.Println("LitIn", t.LitIn)
	//fmt.Println("final debug: tSorter", t.Sorter)

}

// assumes LE, positive weights
func (t *SortingNetwork) CreateBags() {

	nBags := len(Binary(t.pb.K))
	bins := make([][]int, len(t.pb.Entries))
	bagPos := make([]int, nBags)
	bagSize := make([]int, nBags)

	maxWeight := int64(0)

	for i, e := range t.pb.Entries {
		bins[i] = Binary(e.Weight)

		for j, x := range bins[i] {
			bagSize[j] += x
		}

		if maxWeight < e.Weight {
			maxWeight = e.Weight
		}

	}

	t.Bags = make([][]sat.Literal, len(Binary(maxWeight)))

	for i := range t.Bags {
		t.Bags[i] = make([]sat.Literal, bagSize[i])
	}

	for i, e := range t.pb.Entries {
		for j, x := range bins[i] {
			if x == 1 {
				t.Bags[j][bagPos[j]] = e.Literal
				bagPos[j]++
			}
		}
	}

}
