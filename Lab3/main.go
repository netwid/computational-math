package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"math"
	"os"
)

type Function struct {
	s string
	f func(x float64) float64
}

type Method struct {
	name string
	f    func(a, b float64, n int, f func(x float64) float64) float64
	num  float64
}

func main() {
	methods := []Method{
		{"Left rectangle method", LeftRectangleMethod, 3},
		{"Middle rectangle method", MiddleRectangleMethod, 3},
		{"Right rectangle method", RightRectangleMethod, 3},
		{"Trapezoidal method", TrapezoidalMethod, 3},
		{"Simpson method", SimpsonMethod, 15},
	}

	functions := []Function{
		{
			"sin(x)",
			func(x float64) float64 { return math.Sin(x) },
		},
		{
			"x^3 - x + 4",
			func(x float64) float64 { return x*x*x - x + 4 },
		},
		{
			"x^3 - 2x^2 + 4x - 8",
			func(x float64) float64 { return x*x*x - 2*x*x + 4*x - 8 },
		},
	}

	app := &cli.App{
		Name:  "Computation",
		Usage: "Solve equations",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "console-input",
				Aliases: []string{"i"},
				Usage:   "Use console input",
			},
			&cli.StringFlag{
				Name:    "filename",
				Aliases: []string{"f"},
				Usage:   "Filename if not console input (default: data.yml)",
			},
		},
		Action: func(cCtx *cli.Context) error {
			var method, function int
			var a, b, eps float64

			for i, method := range methods {
				fmt.Printf("%d. %s\n", i+1, method.name)
			}
			fmt.Scan(&method)
			if method < 1 || method > len(methods) {
				return fmt.Errorf("invalid method")
			}

			for i, function := range functions {
				fmt.Printf("%d. %s\n", i+1, function.s)
			}
			fmt.Print("Choose function: ")
			fmt.Scan(&function)
			if function < 1 || function > len(functions) {
				return fmt.Errorf("invalid function")
			}

			fmt.Print("Enter a: ")
			fmt.Scan(&a)
			fmt.Print("Enter b: ")
			fmt.Scan(&b)
			fmt.Print("Enter eps: ")
			fmt.Scan(&eps)

			Solve(a, b, eps, methods[method-1], functions[function-1].f)

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
