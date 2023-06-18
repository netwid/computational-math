package main

func Lagrange(points []Point, x float64) float64 {
	ans := float64(0)

	for i := 0; i < len(points); i++ {
		basis := 1.0

		for j := 0; j < len(points); j++ {
			if i != j {
				basis *= (x - points[j].X) / (points[i].X - points[j].X)
			}
		}

		ans += points[i].Y * basis
	}

	return ans
}

func FiniteDifferenceTable(points []Point) [][]float64 {
	table := make([][]float64, len(points))

	table[0] = make([]float64, len(points))
	for i := 0; i < len(points); i++ {
		table[0][i] = points[i].Y
	}

	for i := 1; i < len(table); i++ {
		table[i] = make([]float64, len(points)-i)

		for j := 0; j < len(table[i]); j++ {
			table[i][j] = table[i-1][j+1] - table[i-1][j]
		}
	}

	return table
}

func Newton(points []Point, x float64) float64 {
	finiteDifferenceTable := FiniteDifferenceTable(points)

	middle := len(points) / 2

	h := points[1].X - points[0].X

	ans := float64(0)
	if x < points[middle].X {
		i := 0
		for points[i+1].X < x {
			i++
		}

		t := (x - points[i].X) / h
		y := points[i].Y

		ans += y
		numerator := t
		factorial := 1
		denominator := factorial
		for j := 1; j < len(finiteDifferenceTable)-i; j++ {
			ans += numerator / float64(denominator) * finiteDifferenceTable[j][i]
			t--
			numerator *= t
			factorial++
			denominator *= factorial
		}
	} else {
		ans += points[len(points)-1].Y

		t := (x - points[len(points)-1].X) / h

		numerator := t
		factorial := 1
		denominator := factorial
		for j := len(points) - 2; j >= 0; j-- {
			ans += numerator / float64(denominator) * finiteDifferenceTable[len(finiteDifferenceTable)-j-1][j]
			t++
			numerator *= t
			factorial++
			denominator *= factorial
		}
	}

	return ans
}
