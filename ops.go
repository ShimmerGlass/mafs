package main

import (
	"fmt"
	"math"
)

func sqrt(in []float64) (float64, error) {
	if len(in) != 1 {
		return 0, fmt.Errorf("this function requires 1 argument")
	}

	return math.Sqrt(in[0]), nil
}

func pow(in []float64) (float64, error) {
	if len(in) != 2 {
		return 0, fmt.Errorf("this function requires 2 arguments")
	}

	return math.Pow(in[0], in[1]), nil
}
