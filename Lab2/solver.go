package main

import (
	"fmt"
	"log"
	"math"
)

func ChordMethod(e Equation, a float64, b float64, eps float64) {
	fA := e.f(a)
	fB := e.f(b)

	iterations := 0
	for math.Abs(fB-fA) > eps {
		x := (a*fB - b*fA) / (fB - fA)
		fX := e.f(x)

		if fA*fX < 0 {
			b = x
			fB = fX
		} else {
			a = x
			fA = fX
		}
		iterations++
	}
	fmt.Println("X =", (a+b)/2)
	fmt.Println("f(x) =", e.f((a+b)/2))
	fmt.Println("Number of iterations:", iterations)

}

func NewtonMethod(e Equation, a float64, b float64, eps float64) {
	var x0, x float64
	if e.f(a)*e.derivative2(a) > 0 {
		x0 = a
	} else {
		x0 = b
	}

	iterations := 0
	for {
		x = x0 - e.f(x0)/e.derivative(x0)
		if math.Abs(x-x0) < eps {
			break
		}
		x0 = x
		iterations++
	}

	fmt.Println("X =", x)
	fmt.Println("f(x) =", e.f(x))
	fmt.Println("Number of iterations:", iterations)
}

func SimpleIterationMethod(e Equation, a float64, b float64, eps float64) {
	max := e.derivative(a)
	maxAbs := math.Abs(max)
	for i := a; i < b; i += eps {
		if math.Abs(e.derivative(i)) > maxAbs {
			max = e.derivative(i)
			maxAbs = math.Abs(max)
		}
	}

	x0 := a

	var x float64

	iterations := 0
	for {
		lambda := -1 / max
		if e.derivative(x0)*max < 0 {
			lambda = -lambda
		}
		phi := func(x0 float64) float64 {
			return x0 + lambda*e.f(x0)
		}
		x = phi(x0)
		if math.Abs(x-x0) < eps && math.Abs(e.f(x)) < eps {
			break
		}
		x0 = x
		iterations++
	}

	fmt.Println("X =", x)
	fmt.Println("f(x) =", e.f(x))
	fmt.Println("Number of iterations:", iterations)
}

func NewtonMethodSystem(s System, x0 float64, y0 float64, eps float64) {
	var x, y float64

	iterations := 0
	for {
		solve, err := Compute(s.jacob(x0, y0), 0.001)
		if err != nil {
			log.Fatal(err)
		}

		x = x0 + solve[0]
		y = y0 + solve[1]

		if solve[0] < eps && solve[1] < eps && math.Abs(s.f[0](x, y)-s.f[0](x0, y0)) < eps && math.Abs(s.f[1](x, y)-s.f[1](x0, y0)) < eps {
			break
		}

		x0 = x
		y0 = y

		iterations++
	}

	fmt.Println("X =", x)
	fmt.Println("Y =", y)
	fmt.Println("f1(x, y) =", s.f[0](x, y))
	fmt.Println("f2(x, y) =", s.f[1](x, y))
	fmt.Println("Number of iterations:", iterations)
}
