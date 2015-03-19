package config

import (
	"fmt"
)

// global configuration accessible from everywhere

var Complex_flag string
var Timeout_flag int
var MDD_max_flag int
var MDD_redundant_flag bool

func PringConfig() {
	fmt.Println("Configuration")
	fmt.Println("Complex_flag :\t", Complex_flag)
	fmt.Println("Timeout_flag :\t", Timeout_flag)
	fmt.Println("MDD_max_flag :\t", MDD_max_flag)
	fmt.Println("MDD_redundant_flag :\t", MDD_redundant_flag)
}

// These constants are for future implementations

// MDD translation by which clauses? implications over branchens? @ignasis implementation

// Use of what type of sorting networks:

// Use Mergers

var SortersT int
var EquationT int