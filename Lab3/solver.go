package main

import (
	"fmt"
	"math"
)

func Solve(a, b, eps float64, m Method, f func(x float64) float64) {
	n := 4
	res0 := m.f(a, b, n, f)
	res1 := m.f(a, b, n*2, f)

	for math.Abs(res1-res0)/m.num > eps {
		n *= 2
		res1, res0 = m.f(a, b, n*2, f), res1
	}

	fmt.Println("I =", res1)
	fmt.Println("Eps =", math.Abs(res1-res0)/m.num)
	fmt.Println("N =", n*2)
}

func LeftRectangleMethod(a, b float64, n int, f func(x float64) float64) float64 {
	h := (b - a) / float64(n)
	ans := float64(0)

	for i := a; i <= b; i += h {
		ans += f(i) * h
	}

	return ans
}

func MiddleRectangleMethod(a, b float64, n int, f func(x float64) float64) float64 {
	h := (b - a) / float64(n)
	ans := float64(0)

	for i := a; i <= b; i += h {
		ans += f(i+h/2) * h
	}

	return ans
}

func RightRectangleMethod(a, b float64, n int, f func(x float64) float64) float64 {
	h := (b - a) / float64(n)
	ans := float64(0)

	for i := a; i <= b; i += h {
		ans += f(i+h) * h
	}

	return ans
}

func TrapezoidalMethod(a, b float64, n int, f func(x float64) float64) float64 {
	h := (b - a) / float64(n)
	ans := float64(0)

	for i := a + h; i < b; i += h {
		ans += f(i)
	}

	return h * ((f(a)+f(b))/2 + ans)
}

func SimpsonMethod(a, b float64, n int, f func(x float64) float64) float64 {
	h := (b - a) / float64(n)
	ans := f(a) + f(b)

	for i := a + h; i < b; i += 2 * h {
		ans += 4 * f(i)
	}
	for i := a + 2*h; i < b; i += 2 * h {
		ans += 2 * f(i)
	}

	return h / 3 * ans
}
