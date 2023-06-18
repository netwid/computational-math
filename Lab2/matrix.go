package main

import (
	"errors"
	"math"
)

const limit = 1000

func Compute(matrix [][]float64, accuracy float64) ([]float64, error) {
	n := len(matrix)

	err := diagonalDominance(&matrix)
	if err != nil {
		return nil, err
	}

	C := make([][]float64, n)
	for i := 0; i < n; i++ {
		C[i] = make([]float64, n)
		for j := 0; j < n; j++ {
			C[i][j] = -matrix[i][j] / matrix[i][i]
		}
		C[i][i] = 0
	}
	d := make([]float64, n)
	for i := 0; i < n; i++ {
		d[i] = matrix[i][n] / matrix[i][i]
	}

	prevX := make([]float64, n)
	copy(prevX, d)
	X := make([]float64, n)

	iterationCount := 0
	for {
		iterate(C, d, prevX, X)
		iterationCount++
		if iterationCount > limit || isExact(prevX, X, accuracy) {
			break
		}
		for i := 0; i < len(C); i++ {
			prevX[i] = X[i]
		}
	}

	if iterationCount > limit {
		return nil, errors.New("limit exceeded")
	}

	return X, nil
}

func diagonalDominance(matrix *[][]float64) error {
	n := len(*matrix)

	maxIndexes := map[int]int{} // column -> row
	cntStrict := 0

	for i := 0; i < n; i++ {
		max := math.Abs((*matrix)[i][0])
		maxIndex := 0
		sum := 0.0
		for j := 0; j < n; j++ {
			sum += math.Abs((*matrix)[i][j])
			if math.Abs((*matrix)[i][j]) > max {
				max = math.Abs((*matrix)[i][j])
				maxIndex = j
			}
		}

		if max < sum-max {
			return errors.New("can't make the diagonal dominance")
		}

		if _, ok := maxIndexes[maxIndex]; ok {
			return errors.New("can't make the diagonal dominance")
		}

		if max > sum-max {
			cntStrict++
		}
		maxIndexes[maxIndex] = i
	}

	if cntStrict == 0 {
		return errors.New("can't make the diagonal dominance")
	}

	newMatrix := make([][]float64, n)
	for i := 0; i < n; i++ {
		newMatrix[i] = make([]float64, n+1)
		newMatrix[i][n] = (*matrix)[i][n]
	}

	for column, row := range maxIndexes {
		for i := 0; i < n; i++ {
			newMatrix[column] = (*matrix)[row]
		}
	}

	*matrix = newMatrix

	return nil
}

func iterate(C [][]float64, d []float64, prevX []float64, X []float64) {
	for i := 0; i < len(C); i++ {
		X[i] = d[i]
		for j := 0; j < len(C); j++ {
			if j < i {
				X[i] += C[i][j] * X[j]
			} else {
				X[i] += C[i][j] * prevX[j]
			}
		}
	}
}

func isExact(x1, x2 []float64, eps float64) bool {
	var max float64 = 0
	for index, item := range x1 {
		if math.Abs(item-x2[index]) > max {
			max = math.Abs(item - x2[index])
		}
	}

	if max < eps {
		return true
	}
	return false
}
