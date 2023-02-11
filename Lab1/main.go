package main

import (
	"errors"
	"fmt"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

type data struct {
	Accuracy float64     `yaml:"accuracy"`
	Matrix   [][]float64 `yaml:"matrix"`
}

func main() {
	app := &cli.App{
		Name:  "Computation",
		Usage: "Gauss-Seidel method",
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
			if cCtx.Bool("console-input") {
				var acc float64
				var n int
				var matrix [][]float64

				fmt.Print("Enter accuracy: ")
				fmt.Scan(&acc)

				fmt.Print("Enter matrix size (n): ")
				fmt.Scan(&n)

				fmt.Println("Enter matrix with D as last column (n x n+1): ")
				matrix = make([][]float64, n)
				for i := 0; i < n; i++ {
					matrix[i] = make([]float64, n+1)
					for j := 0; j < n+1; j++ {
						fmt.Scan(&matrix[i][j])
					}
				}

				err := Compute(matrix, acc)
				if err != nil {
					return err
				}
				return nil
			}

			filename := "data.yml"

			if len(cCtx.String("filename")) > 0 {
				filename = cCtx.String("filename")
			}

			yamlFile, err := ioutil.ReadFile(filename)
			if err != nil {
				return err
			}

			var d data
			err = yaml.Unmarshal(yamlFile, &d)
			if err != nil {
				return err
			}

			if len(d.Matrix) != len(d.Matrix[0])-1 {
				return errors.New("invalid matrix size")
			}

			err = Compute(d.Matrix, d.Accuracy)
			if err != nil {
				return err
			}

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
