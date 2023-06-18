package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gopkg.in/yaml.v2"
	"image/color"
	"io/ioutil"
	"log"
	"math"
	"os"
)

type data struct {
	EquationOrSystem int     `yaml:"equationOrSystem"`
	Method           int     `yaml:"method"`
	A                float64 `yaml:"a"`
	B                float64 `yaml:"b"`
	Eps              float64 `yaml:"eps"`
	X0               float64 `yaml:"x0"`
	Y0               float64 `yaml:"y0"`
}

type Equation struct {
	s           string
	f           func(x float64) float64
	derivative  func(x float64) float64
	derivative2 func(x float64) float64
}

type Method struct {
	name string
	f    func(e Equation, a float64, b float64, eps float64)
}

type System struct {
	s     []string
	f     []func(x, y float64) float64
	jacob func(x, y float64) [][]float64
}

func main() {
	methods := []Method{
		{"Choord method", ChordMethod},
		{"Newton's method", NewtonMethod},
		{"Simple iteration method", SimpleIterationMethod},
	}

	equations := []Equation{
		{
			"sin(x)",
			func(x float64) float64 { return math.Sin(x) },
			func(x float64) float64 { return math.Cos(x) },
			func(x float64) float64 { return -math.Sin(x) },
		},
		{
			"x^3 - x + 4",
			func(x float64) float64 { return x*x*x - x + 4 },
			func(x float64) float64 { return 3*x*x - 1 },
			func(x float64) float64 { return 6 * x },
		},
		{
			"x^3 - 2x^2 + 4x - 8",
			func(x float64) float64 { return x*x*x - 2*x*x + 4*x - 8 },
			func(x float64) float64 { return 3*x*x - 4*x + 4 },
			func(x float64) float64 { return 6*x - 4 },
		},
	}

	systems := []System{
		{[]string{"x^2 + y^2 = 4", "y = 3x^2"}, []func(x, y float64) float64{
			func(x, y float64) float64 { return x*x + y*y - 4 },
			func(x, y float64) float64 { return y - 3*x*x },
		}, func(x, y float64) [][]float64 {
			return [][]float64{
				{2 * x, 2 * y, 4 - x*x - y*y},
				{-6 * x, 1, 3*x*x - y},
			}
		}},
		{[]string{"y = x^2 - 1", "y = 1"}, []func(x, y float64) float64{
			func(x, y float64) float64 { return x*x - y*y - 1 },
			func(x, y float64) float64 { return y - 1 },
		}, func(x, y float64) [][]float64 {
			return [][]float64{
				{2 * x, 1, -x*x + 1 + y},
				{0, 1, 1 - y},
			}
		}},
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
			var d data

			if cCtx.Bool("console-input") {
				isSystem := false

				for i, method := range methods {
					fmt.Printf("%d. %s\n", i+1, method.name)
				}
				fmt.Println("4. Newton's method (system)")
				fmt.Print("Choose method: ")
				fmt.Scan(&d.Method)
				if d.Method == len(methods)+1 {
					isSystem = true
				}
				if !isSystem && (d.Method < 1 || d.Method > len(methods)) {
					return fmt.Errorf("invalid method")
				}

				if isSystem {
					for i, system := range systems {
						fmt.Printf("%d. %s\n", i+1, system.s)
					}
					fmt.Print("Choose system: ")
					fmt.Scan(&d.EquationOrSystem)
					if d.EquationOrSystem < 1 || d.EquationOrSystem > len(systems) {
						return fmt.Errorf("invalid system")
					}

					fmt.Print("Enter x0: ")
					fmt.Scan(&d.X0)
					fmt.Print("Enter y0: ")
					fmt.Scan(&d.Y0)
					fmt.Print("Enter eps: ")
					fmt.Scan(&d.Eps)
				} else {
					for i, equation := range equations {
						fmt.Printf("%d. %s\n", i+1, equation.s)
					}
					fmt.Print("Choose equation: ")
					fmt.Scan(&d.EquationOrSystem)
					if d.EquationOrSystem < 1 || d.EquationOrSystem > len(equations) {
						return fmt.Errorf("invalid equation")
					}

					fmt.Print("Enter a: ")
					fmt.Scan(&d.A)
					fmt.Print("Enter b: ")
					fmt.Scan(&d.B)
					fmt.Print("Enter eps: ")
					fmt.Scan(&d.Eps)
				}
			} else {
				filename := "data.yml"

				if len(cCtx.String("filename")) > 0 {
					filename = cCtx.String("filename")
				}

				yamlFile, err := ioutil.ReadFile(filename)
				if err != nil {
					return err
				}

				err = yaml.Unmarshal(yamlFile, &d)
				if err != nil {
					return err
				}
			}

			if d.Method >= 1 && d.Method <= len(methods) {
				checkRoots(equations[d.EquationOrSystem-1], d.A, d.B)
				methods[d.Method-1].f(equations[d.EquationOrSystem-1], d.A, d.B, d.Eps)
				drawPlot(equations[d.EquationOrSystem-1], d.A, d.B)
				fmt.Println("Plot saved to function.png")
			} else {
				NewtonMethodSystem(systems[d.EquationOrSystem-1], d.X0, d.Y0, d.Eps)
				drawSystem(systems[d.EquationOrSystem-1], d.X0, d.Y0)
				fmt.Println("Plot saved to system.png")
			}

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func checkRoots(e Equation, a, b float64) {
	if e.f(a)*e.f(b) > 0 {
		log.Fatal("No roots in this interval")
	}

	cnt := 0
	for i := a; i < b; i += 0.01 {
		if e.f(i)*e.f(i+0.01) < 0 {
			if cnt > 0 {
				log.Fatal("More than 1 root in this interval")
			}
			cnt++
		}
	}
}

func drawSystem(s System, x0, y0 float64) {
	p := plot.New()
	p.Title.Text = s.s[0] + " and " + s.s[1]
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	f1 := plotter.NewFunction(func(x float64) float64 {
		if math.Abs(x) < 2 {
			return math.Sqrt(4 - x*x)
		}
		return 0
	})
	f1.Color = color.RGBA{B: 255, A: 255}

	f2 := plotter.NewFunction(func(x float64) float64 { return 3 * x * x })
	f2.Color = color.RGBA{G: 255, A: 255}

	xAxis := plotter.NewFunction(func(x float64) float64 { return 0 })
	xAxis.Dashes = []vg.Length{vg.Points(2), vg.Points(2)}
	xAxis.Width = vg.Points(1.5)
	xAxis.Color = color.RGBA{A: 255}

	p.Add(f1, xAxis, f2)
	p.Legend.Add("f1", f1)
	p.Legend.Add("f2", f2)
	p.Legend.ThumbnailWidth = 1 * vg.Inch

	p.X.Min = x0 - 5
	p.X.Max = x0 + 2
	p.Y.Min = y0 - 5
	p.Y.Max = y0 + 5

	if err := p.Save(7*vg.Inch, 7*vg.Inch, "system.png"); err != nil {
		log.Fatal(err)
	}
}

func drawPlot(e Equation, a, b float64) {
	p := plot.New()
	p.Title.Text = e.s
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	f := plotter.NewFunction(e.f)
	f.Color = color.RGBA{B: 255, A: 255}

	der := plotter.NewFunction(e.derivative)
	der.Color = color.RGBA{G: 255, A: 255}

	xAxis := plotter.NewFunction(func(x float64) float64 { return 0 })
	xAxis.Dashes = []vg.Length{vg.Points(2), vg.Points(2)}
	xAxis.Width = vg.Points(1.5)
	xAxis.Color = color.RGBA{A: 255}

	p.Add(f, xAxis, der)
	p.Legend.Add("function", f)
	p.Legend.Add("derivative", der)
	p.Legend.ThumbnailWidth = 1 * vg.Inch

	p.X.Min = a - 2
	p.X.Max = b + 2
	p.Y.Min = -5
	p.Y.Max = 5

	if err := p.Save(7*vg.Inch, 7*vg.Inch, "function.png"); err != nil {
		log.Fatal(err)
	}
}
