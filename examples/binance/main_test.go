package main

import (
	"fmt"
	"testing"
)

func BenchmarkRsi(b *testing.B) {
	closes := make([]float64, 15)
	xx := []float64{0.0224, 0.0239, 0.0226, 0.0236, 0.0219, 0.0221, 0.0224, 0.0241, 0.0238, 0.0243, 0.0262, 0.0261, 0.0280, 0.0269, 0.0334}
	for i := 0; i < len(closes); i++ {
		closes[i] = xx[i]
	}

	for n := 0; n < b.N; n++ {
		Rsi(closes, 14)
	}
}

// Benchmark_ReplacingFirstItemUsingAppend-16    	30355455	        34.60 ns/op
func Benchmark_ReplacingFirstItemUsingAppend(b *testing.B) {
	for n := 0; n < b.N; n++ {
		closes := []float64{0.0224, 0.0239, 0.0226, 0.0236, 0.0219, 0.0221, 0.0224, 0.0241, 0.0238, 0.0243, 0.0262, 0.0261, 0.0280, 0.0269, 0.0334}
		closes = append(closes[1:], 0.42)
	}
}

// Benchmark_ReplacingFirstItemByShifting-16    	268561569	         4.451 ns/op
func Benchmark_ReplacingFirstItemByShifting(b *testing.B) {
	for n := 0; n < b.N; n++ {
		closes := [15]float64{0.0224, 0.0239, 0.0226, 0.0236, 0.0219, 0.0221, 0.0224, 0.0241, 0.0238, 0.0243, 0.0262, 0.0261, 0.0280, 0.0269, 0.0334}
		for i := 0; i < len(closes)-1; i++ {
			closes[i] = closes[i+1]
		}
		closes[len(closes)-1] = 0.420
	}
}
func TestShit(t *testing.T) {
	closes := [15]float64{0.0224, 0.0239, 0.0226, 0.0236, 0.0219, 0.0221, 0.0224, 0.0241, 0.0238, 0.0243, 0.0262, 0.0261, 0.0280, 0.0269, 0.0334}
	for i := 0; i < len(closes)-1; i++ {
		closes[i] = closes[i+1]
	}
	closes[len(closes)-1] = 0.420
	fmt.Println(closes)
}

// Rsi - Relative strength index
func Rsi(inReal []float64, inTimePeriod int) []float64 {

	outReal := make([]float64, len(inReal))

	if inTimePeriod < 2 {
		return outReal
	}

	// variable declarations
	tempValue1 := 0.0
	tempValue2 := 0.0
	outIdx := inTimePeriod
	today := 0
	prevValue := inReal[today]
	prevGain := 0.0
	prevLoss := 0.0
	today++

	for i := inTimePeriod; i > 0; i-- {
		tempValue1 = inReal[today]
		today++
		tempValue2 = tempValue1 - prevValue
		prevValue = tempValue1
		if tempValue2 < 0 {
			prevLoss -= tempValue2
		} else {
			prevGain += tempValue2
		}
	}

	prevLoss /= float64(inTimePeriod)
	prevGain /= float64(inTimePeriod)

	if today > 0 {

		tempValue1 = prevGain + prevLoss
		if !((-0.00000000000001 < tempValue1) && (tempValue1 < 0.00000000000001)) {
			outReal[outIdx] = 100.0 * (prevGain / tempValue1)
		} else {
			outReal[outIdx] = 0.0
		}
		outIdx++

	} else {

		for today < 0 {
			tempValue1 = inReal[today]
			tempValue2 = tempValue1 - prevValue
			prevValue = tempValue1
			prevLoss *= float64(inTimePeriod - 1)
			prevGain *= float64(inTimePeriod - 1)
			if tempValue2 < 0 {
				prevLoss -= tempValue2
			} else {
				prevGain += tempValue2
			}
			prevLoss /= float64(inTimePeriod)
			prevGain /= float64(inTimePeriod)
			today++
		}
	}

	for today < len(inReal) {

		tempValue1 = inReal[today]
		today++
		tempValue2 = tempValue1 - prevValue
		prevValue = tempValue1
		prevLoss *= float64(inTimePeriod - 1)
		prevGain *= float64(inTimePeriod - 1)
		if tempValue2 < 0 {
			prevLoss -= tempValue2
		} else {
			prevGain += tempValue2
		}
		prevLoss /= float64(inTimePeriod)
		prevGain /= float64(inTimePeriod)
		tempValue1 = prevGain + prevLoss
		if !((-0.00000000000001 < tempValue1) && (tempValue1 < 0.00000000000001)) {
			outReal[outIdx] = 100.0 * (prevGain / tempValue1)
		} else {
			outReal[outIdx] = 0.0
		}
		outIdx++
	}

	return outReal
}
