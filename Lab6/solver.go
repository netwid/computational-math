package main

import "math"

func RungeRule(yh, yh2, p float64) float64 {
	return math.Abs(yh-yh2) / (math.Pow(2, p) - 1)
}

func MultiStepAcc(ans, yExact []float64) float64 {
	max := math.Abs(yExact[0] - ans[0])
	for idx, y := range ans {
		diff := math.Abs(yExact[idx] - y)
		if max < diff {
			max = diff
		}
	}
	return max
}

func Euler(x []float64, y0, h float64, f func(x, y float64) float64) []float64 {
	y := y0
	var ys []float64
	for _, xi := range x {
		ys = append(ys, y)
		y += h * f(xi, y)
	}

	return ys
}

func RungeKutta(x []float64, y0, h float64, f func(x, y float64) float64) []float64 {
	y := y0
	var ys []float64

	for _, xi := range x {
		ys = append(ys, y)

		k1 := h * f(xi, y)
		k2 := h * f(xi+h/2, y+k1/2)
		k3 := h * f(xi+h/2, y+k2/2)
		k4 := h * f(xi+h, y+k3)

		y += (k1 + 2*k2 + 2*k3 + k4) / 6
	}

	return ys
}

func Adams(x []float64, y0, h float64, f func(x, y float64) float64) []float64 {
	n := len(x)
	y := make([]float64, n)
	y[0] = y0

	if n < 4 {
		return nil
	}

	rk := RungeKutta(x, y0, h, f)
	for i := 0; i < 4; i++ {
		y[i+1] = rk[i]
		x[i+1] = x[i] + h
	}

	for i := 4; i < n; i++ {
		k1 := f(x[i-1], y[i-1]) - f(x[i-2], y[i-2])
		k2 := k1 + f(x[i-3], y[i-3]) - f(x[i-2], y[i-2])
		k3 := k2 + 2*f(x[i-3], y[i-3]) - f(x[i-2], y[i-2]) - f(x[i-4], y[i-4])

		y[i] = y[i-1] + h*f(x[i-1], y[i-1]) + (math.Pow(h, 2)/2)*k1 +
			(5*math.Pow(h, 3)/12)*k2 + (3*math.Pow(h, 4)/8)*k3
	}

	return y
}
