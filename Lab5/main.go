package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gopkg.in/yaml.v2"
	"image/color"
	"io/ioutil"
	"log"
	"math"
	"os"
)

type data struct {
	Polynomial int     `yaml:"polynomial"`
	Points     []Point `yaml:"points"`
}

type Function struct {
	s string
	f func(x float64) float64
}

type Point struct {
	X, Y float64
}

type Polynomial struct {
	name string
	f    func(points []Point, x float64) float64
}

func main() {
	polynomials := []Polynomial{
		{"Lagrange polynomial", Lagrange},
		{"Newton's polynomial", Newton},
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
			var p int
			var points []Point

			if cCtx.Bool("console-input") {
				fmt.Println("Choose polynomial:")
				for i := 0; i < len(polynomials); i++ {
					fmt.Printf("%d. %s\n", i+1, polynomials[i].name)
				}
				fmt.Scan(&p)

				var option int
				fmt.Print("Set or function? [1/2]: ")
				fmt.Scan(&option)

				switch option {
				case 1:
					var num int
					fmt.Println("Enter number of points: ")
					fmt.Scan(&num)

					var x, y float64
					for i := 0; i < num; i++ {
						fmt.Printf("#Point %d\n", i+1)
						fmt.Print("x: ")
						fmt.Scan(&x)
						fmt.Print("y: ")
						fmt.Scan(&y)

						points = append(points, Point{
							X: x,
							Y: y,
						})
					}

				case 2:
					var f int
					fmt.Println("Choose function:")
					for i := 0; i < len(functions); i++ {
						fmt.Printf("%d. %s\n", i+1, functions[i].s)
					}
					fmt.Scan(&f)

					var a, b float64
					var num int
					fmt.Print("Enter a: ")
					fmt.Scan(&a)
					fmt.Print("Enter b: ")
					fmt.Scan(&b)
					fmt.Print("Enter number of intervals: ")
					fmt.Scan(&num)

					step := (b - a) / float64(num)

					for i := a; i <= b; i += step {
						points = append(points, Point{
							X: i,
							Y: functions[f-1].f(i),
						})
					}

				default:
					fmt.Println("Incorrect choice")
					return nil
				}
			} else {
				var d data
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

				p = d.Polynomial
				points = d.Points
			}

			checkPoints(points)

			table := FiniteDifferenceTable(points)
			for i := 0; i < len(table); i++ {
				for j := 0; j < len(table[i]); j++ {
					fmt.Printf("%.2f\t", table[i][j])
				}
				fmt.Println()
			}

			drawPolynomial(polynomials[p-1], points)
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func checkPoints(points []Point) {
	eps := 0.00001
	if len(points) < 2 {
		fmt.Errorf("error: too few points")
	}

	h := points[1].X - points[0].X

	for i := 2; i < len(points); i++ {
		if math.Abs(points[i].X-points[i-1].X-h) > eps {
			fmt.Errorf("error: non-equal nodes")
		}
	}
}

func drawPolynomial(polynomial Polynomial, points []Point) {
	p := plot.New()
	p.Title.Text = polynomial.name
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	f := plotter.NewFunction(func(x float64) float64 {
		return polynomial.f(points, x)
	})
	f.Color = color.RGBA{B: 255, A: 255}
	p.Add(f)

	xAxis := plotter.NewFunction(func(x float64) float64 { return 0 })
	xAxis.Dashes = []vg.Length{vg.Points(2), vg.Points(2)}
	xAxis.Width = vg.Points(1.5)
	xAxis.Color = color.RGBA{A: 255}
	p.Add(xAxis)

	pointsScatter, err := plotter.NewScatter(getPoints(points))
	if err != nil {
		log.Fatal(err)
	}
	pointsScatter.GlyphStyle.Shape = draw.CircleGlyph{}
	pointsScatter.GlyphStyle.Color = color.RGBA{R: 255, A: 255}
	p.Add(pointsScatter)

	xMin, xMax, yMin, yMax := getMinMaxPoints(points)
	p.X.Min = xMin - 0.5
	p.X.Max = xMax + 0.5
	p.Y.Min = yMin - 0.5
	p.Y.Max = yMax + 0.5

	if err := p.Save(7*vg.Inch, 7*vg.Inch, "function.png"); err != nil {
		log.Fatal(err)
	}
}

func getMinMaxPoints(points []Point) (float64, float64, float64, float64) {
	xMin, xMax := points[0].X, points[0].X
	yMin, yMax := points[0].Y, points[0].Y

	for _, pt := range points {
		if pt.X < xMin {
			xMin = pt.X
		}
		if pt.X > xMax {
			xMax = pt.X
		}
		if pt.Y < yMin {
			yMin = pt.Y
		}
		if pt.Y > yMax {
			yMax = pt.Y
		}
	}

	return xMin, xMax, yMin, yMax
}

func getPoints(points []Point) plotter.XYs {
	pts := make(plotter.XYs, len(points))
	for i, pt := range points {
		pts[i].X = pt.X
		pts[i].Y = pt.Y
	}
	return pts
}
