package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"image/color"
	"log"
	"math"
	"os"
)

type MethodType int64

const (
	oneStep MethodType = iota
	multiStep
)

type Method struct {
	s             string
	f             func(x []float64, y0, h float64, f func(x, y float64) float64) []float64
	t             MethodType
	accuracyOrder float64
}

type Function struct {
	s     string
	f     func(x, y float64) float64
	exact func(x, x0, y0 float64) float64
}

func main() {
	methods := []Method{
		{
			"Euler",
			Euler,
			oneStep,
			1,
		},
		{
			"Runge-Kutta",
			RungeKutta,
			oneStep,
			4,
		},
		{
			"Adams",
			Adams,
			multiStep,
			0,
		},
	}

	functions := []Function{
		{
			"y' = y + (1 + x) * y^2",
			func(x, y float64) float64 { return y + (1+x)*y*y },
			func(x, x0, y0 float64) float64 {
				C := -math.Pow(math.E, x0)/y0 - x0*math.Pow(math.E, x0)
				return -math.Pow(math.E, x) / (x*math.Pow(math.E, x) + C)
			},
		},
		{
			"y' = y / x",
			func(x, y float64) float64 { return y / x },
			func(x, x0, y0 float64) float64 {
				C := y0 / x0
				return C * x
			},
		},
		{
			"y' = y + x",
			func(x, y float64) float64 { return y + x },
			func(x, x0, y0 float64) float64 {
				C := (y0 + x0 + 1) / math.Pow(math.E, x0)
				return C*math.Pow(math.E, x) - x - 1
			},
		},
	}

	app := &cli.App{
		Name:  "Interpolation",
		Usage: "Interpolation",
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
			var m int
			fmt.Println("Choose method:")
			for i := 0; i < len(methods); i++ {
				fmt.Printf("%d. %s\n", i+1, methods[i].s)
			}
			fmt.Scan(&m)

			if m < 1 || m > len(methods) {
				fmt.Errorf("incorrect choice")
			}

			var f int
			fmt.Println("Choose function:")
			for i := 0; i < len(functions); i++ {
				fmt.Printf("%d. %s\n", i+1, functions[i].s)
			}
			fmt.Scan(&f)

			if f < 1 || f > len(functions) {
				fmt.Errorf("incorrect choice")
			}

			var x0, xn, y0, h, ε float64
			fmt.Print("Enter x0: ")
			fmt.Scan(&x0)
			fmt.Print("Enter xn: ")
			fmt.Scan(&xn)
			fmt.Print("Enter y0: ")
			fmt.Scan(&y0)
			fmt.Print("Enter h: ")
			fmt.Scan(&h)
			fmt.Print("Enter ε: ")
			fmt.Scan(&ε)

			if x0 > xn {
				fmt.Errorf("x0 must be less than xn")
			}

			x := xRange(x0, xn, h)
			yExact := yExactRange(x, x0, y0, functions[f-1].exact)

			ans := methods[m-1].f(x, y0, h, functions[f-1].f)
			if ans == nil {
				fmt.Errorf("error while processing ans")
			}

			iterations := 0
			if methods[m-1].t == oneStep {
				x = xRange(x0, xn, h/2)
				next := methods[m-1].f(x, y0, h/2, functions[f-1].f)
				acc := RungeRule(ans[len(ans)-1], next[len(next)-1], methods[m-1].accuracyOrder)
				fmt.Printf("%f -> %f\n", h/2, next[len(next)-1])
				for acc > ε {
					h /= 2
					ans = next
					x = xRange(x0, xn, h/2)
					next = methods[m-1].f(x, y0, h/2, functions[f-1].f)
					fmt.Printf("%f -> %f\n", h/2, next[len(next)-1])
					acc = RungeRule(ans[len(ans)-1], next[len(next)-1], methods[m-1].accuracyOrder)
					iterations++
				}
				fmt.Printf("Iterations: %d\n", iterations)
				fmt.Printf("Accuracy: %f", acc)
			} else {
				acc := MultiStepAcc(ans, yExact)

				for acc > ε {
					h /= 2
					x = xRange(x0, xn, h)
					yExact = yExactRange(x, x0, y0, functions[f-1].exact)
					ans = methods[m-1].f(x, y0, h, functions[f-1].f)
					fmt.Printf("%f -> %f\n", h, ans[len(ans)-1])
					acc = MultiStepAcc(ans, yExact)
					iterations++
				}

				fmt.Printf("Iterations: %d\n", iterations)
				fmt.Printf("Accuracy: %f", acc)
			}

			x = xRange(x0, xn, h)
			yExact = yExactRange(x, x0, y0, functions[f-1].exact)
			draw(x, ans, yExact)

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func xRange(x0, xn, h float64) []float64 {
	var x []float64
	for i := x0; i <= xn; i += h {
		x = append(x, i)
	}

	return x
}

func yExactRange(x []float64, x0, y0 float64, f func(x, x0, y0 float64) float64) []float64 {
	var yExact []float64
	for _, xi := range x {
		yExact = append(yExact, f(xi, x0, y0))
	}

	return yExact
}

func draw(x, y, yExact []float64) {
	p := plot.New()
	p.Title.Text = "Graphics"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	pts := make(plotter.XYs, len(x))
	for i := range x {
		pts[i].X = x[i]
		pts[i].Y = y[i]
	}
	line, _ := plotter.NewLine(pts)
	line.Color = color.RGBA{R: 255, A: 255}

	p.Add(line)

	ptsExact := make(plotter.XYs, len(x))
	for i := range x {
		ptsExact[i].X = x[i]
		ptsExact[i].Y = yExact[i]
	}
	line, _ = plotter.NewLine(ptsExact)
	line.Color = color.RGBA{G: 255, A: 255}

	p.Add(line)

	xAxis := plotter.NewFunction(func(x float64) float64 { return 0 })
	xAxis.Dashes = []vg.Length{vg.Points(2), vg.Points(2)}
	xAxis.Width = vg.Points(1.5)
	xAxis.Color = color.RGBA{A: 255}
	p.Add(xAxis)

	xMin, xMax, yMin, yMax := x[0], x[len(x)-1], yExact[0], yExact[len(y)-1]
	p.X.Min = xMin - 0.5
	p.X.Max = xMax + 0.5
	p.Y.Min = yMin - 0.5
	p.Y.Max = yMax + 0.5

	if err := p.Save(7*vg.Inch, 7*vg.Inch, "function.png"); err != nil {
		log.Fatal(err)
	}
}
