package main

import "fmt"

func nod(x, y int) int {
	for y != 0 {
		x, y = y, x%y
	}
	return x
}

func abs(num int) int {
	return (num ^ (num >> 31)) - (num >> 31)
}

type Rational struct {
	numer, denom int
}

func sign(x int) int {
	if x < 0 {
		return -1
	}
	return 1
}

func (frac Rational) reduce() Rational {
	commonDivisor := nod(frac.numer, frac.denom)
	frac.numer /= commonDivisor
	frac.denom /= commonDivisor
	if frac.denom < 0 {
		frac.denom = -frac.denom
		frac.numer = -frac.numer
	}
	return frac
}

func createFraction(numerator, denominator int) Rational {
	reducedFraction := Rational{numerator, denominator}.reduce()
	return reducedFraction
}

func addFractions(f1, f2 Rational) Rational {
	lcmDenom := f1.denom * f2.denom / nod(abs(f1.denom), abs(f2.denom))
	adjustedNumer1 := f1.numer * (lcmDenom / f1.denom)
	adjustedNumer2 := f2.numer * (lcmDenom / f2.denom)
	sumNumer := adjustedNumer1 + adjustedNumer2
	resultFraction := Rational{numer: sumNumer, denom: lcmDenom}
	return resultFraction.reduce()
}

func subFractions(f1, f2 Rational) Rational {
	lcmDenom := f1.denom * f2.denom / nod(abs(f1.denom), abs(f2.denom))
	adjustedNumer1 := f1.numer * (lcmDenom / f1.denom)
	adjustedNumer2 := f2.numer * (lcmDenom / f2.denom)
	sumNumer := adjustedNumer1 - adjustedNumer2
	resultFraction := Rational{numer: sumNumer, denom: lcmDenom}
	return resultFraction.reduce()
}

func multFractions(f1, f2 Rational) Rational {
	prodNumer := f1.numer * f2.numer
	prodDenom := f1.denom * f2.denom
	resultFraction := Rational{numer: prodNumer, denom: prodDenom}
	return resultFraction.reduce()
}

func divideFractions(f1, f2 Rational) Rational {
	inversef2 := Rational{
		numer: f2.denom * sign(f2.numer),
		denom: abs(f2.numer),
	}
	return multFractions(f1, inversef2)
}

func makeAbs(frac Rational) Rational {
	return Rational{
		numer: abs(frac.numer),
		denom: frac.denom,
	}
}

func gauss(f1 [][]Rational, f2 []Rational, size int) []Rational {
	o := make([]Rational, size)
	var ind int
	for n := 0; n < size; n++ {
		helper := makeAbs(f1[n][n])
		ind = n
		for m := n + 1; m < size; m++ {
			if subFractions(makeAbs(f1[m][n]), helper).numer == 0 {
				helper = makeAbs(f1[m][n])
				ind = m
			}
		}
		if helper.numer == 0 {
			return nil
		}
		f1[n], f1[ind] = f1[ind], f1[n]
		f2[n], f2[ind] = f2[ind], f2[n]
		for m := 0; m < size; m++ {
			if m != n {
				factor := divideFractions(f1[m][n], f1[n][n])
				for k := n; k < size; k++ {
					f1[m][k] = subFractions(f1[m][k], multFractions(factor, f1[n][k]))
				}
				f2[m] = subFractions(f2[m], multFractions(factor, f2[n]))
			}
		}
	}
	for n := size - 1; n >= 0; n-- {
		sum := Rational{numer: 0, denom: 1}
		for k := n + 1; k < size; k++ {
			sum = addFractions(sum, multFractions(f1[n][k], o[k]))
		}
		o[n] = divideFractions(subFractions(f2[n], sum), f1[n][n])
	}
	return o
}

func inputFraction() Rational {
	var value int
	fmt.Scan(&value)
	return Rational{numer: value, denom: 1}
}

func main() {
	var size int
	fmt.Scan(&size)
	results := make([]Rational, size)
	matrix := make([][]Rational, size)
	for i := range matrix {
		matrix[i] = make([]Rational, size)
		for j := range matrix[i] {
			matrix[i][j] = inputFraction()
		}
		results[i] = inputFraction()
	}
	solutions := gauss(matrix, results, size)
	if solutions == nil {
		fmt.Println("No solution")
	} else {
		for _, solution := range solutions {
			fmt.Printf("%d/%d\n", solution.numer, solution.denom)
		}
	}
}
