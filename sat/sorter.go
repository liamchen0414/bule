import (
	"fmt"
	"github.com/vale1410/bule/sorters"
	"strconv"
)

var sorterClauses int
var sorterType sorters.SortingNetworkType

// which sets the type of clauses translated from sorting networks
// see below CreateEncoding for ids wrt clauses
// 0: false, false, false, true, true, true, false, false
// 1: false, false, false, true, true, true, false, true
// 2: false, true, true, true, true, true, true, false
// 3: false, true, true, true, true, true, true, true
func SetUp(which int, typ sorters.SortingNetworkType) {
	sorterClauses = which
	sorterType = typ
}


// CreateCardinality takes set of literals and creates a sorting network
// encoding.
func CreateCardinality(tag string, input []Literal, k int, cType sorters.CardinalityType) ClauseSet {

	var which [8]bool

	switch sorterClauses {
	case 1:
		which = [8]bool{false, false, false, true, true, true, false, false}
	case 2:
		which = [8]bool{false, false, false, true, true, true, false, true}
	case 3:
		which = [8]bool{false, true, true, true, true, true, true, false}
	case 4:
		which = [8]bool{false, true, true, true, true, true, true, true}
	default:
		panic("sorterClauses in sat module not set")
	}

	sorter := sorters.CreateCardinalityNetwork(len(input), k, cType, sorterType)
	sorter.RemoveOutput()

	uniqueId++
	pred := Pred("cnt" + strconv.Itoa(uniqueId))

	output := make([]Literal, 0)

	return CreateEncoding(input, which, output, tag, pred, sorter)
}

// Create Encoding for Sorting Network
// 0)  Omitted for clarity (ids as in paper)
// 1)  A or -D
// 2)  B or -D
// 3) -A or -B or D
// 4) -A or  C
// 5) -B or  C
// 6)  A or  B or -C
// 7)  C or -D
// -1,0,1 mean dontCare, false, true
func CreateEncoding(input []Literal, which [8]bool, output []Literal, tag string, pred Pred, sorter sorters.Sorter) (cs ClauseSet) {

	cs = make([]Clause, 0, 7*len(sorter.Comparators))

	backup := make(map[int]Literal, len(sorter.Out)+len(sorter.In))

	for i, x := range sorter.In {
		backup[x] = input[i]
	}

	for i, x := range sorter.Out {
		backup[x] = output[i]
	}

	for _, comp := range sorter.Comparators {

		if comp.D == 1 || comp.C == 0 {
			fmt.Println("something is wrong with the comparator", comp)
			panic("something is wrong with the comparator")
		}

		getLit := func(x int) Literal {
			if lit, ok := backup[x]; ok {
				return lit
			} else {
				return Literal{true, Atom{pred, x, 0}}
			}
		}

		a := getLit(comp.A)
		b := getLit(comp.B)
		c := getLit(comp.C)
		d := getLit(comp.D)

		if comp.C == 1 { // 6) A or B
			//if which[6] {
			cs.AddClause(tag, a, b)
			//}
		} else if comp.C > 0 { // 4) 5) 6)
			//4)
			if which[4] {
				cs.AddClause(tag, Neg(a), c)
			}
			//5)
			if which[5] {
				cs.AddClause(tag, Neg(b), c)
			}
			//6)
			if which[6] {
				cs.AddClause(tag, a, b, Neg(c))
			}
		}
		if comp.D == 0 { //3)
			//if which[3] {
			cs.AddClause(tag, Neg(a), Neg(b))
			//}
		} else if comp.D > 0 { // 1) 2) 3)
			//1)
			if which[1] {
				cs.AddClause(tag, a, Neg(d))
			}
			//2)
			if which[2] {
				cs.AddClause(tag, b, Neg(d))
			}
			//3)
			if which[3] {
				cs.AddClause(tag, Neg(a), Neg(b), d)
			}
		}

		if which[7] && (comp.D > 1 || comp.D > 1) { // 7)
			cs.AddClause(tag, c, Neg(d))
		}
	}
	return
}
