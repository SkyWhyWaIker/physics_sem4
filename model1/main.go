package main

import (
	"fmt"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"image/color"
	"math"
	_ "os"
)

const (
	mu0  = 4 * math.Pi * 1e-7
	e    = 1.602e-19
	eMas = 9.109e-31
)

func speed(U, e, m float64) float64 {
	return math.Sqrt(2 * U * e / m)
}

func field(m, v, e, r float64) float64 {
	if r == 0 {
		return 0
	}
	return m * v / (e * r)
}

func main() {
	var l int
	var Rk, Ra, U float64

	fmt.Print("Кол-во витков на единицу длины (витков/м): ")
	fmt.Scan(&l)
	if l <= 0 {
		fmt.Println("Ошибка: количество витков на единицу длины должно быть больше 0")
		return
	}

	fmt.Print("Радиус катода (м): ")
	fmt.Scan(&Rk)
	if Rk < 0 {
		fmt.Println("Ошибка: радиус катода не может быть отрицательным")
		return
	}

	fmt.Print("Радиус анода (м): ")
	fmt.Scan(&Ra)
	if Ra < 0 {
		fmt.Println("Ошибка: радиус анода не может быть отрицательным")
		return
	}

	if Ra <= Rk {
		fmt.Println("Ошибка: радиус анода должен быть больше радиуса катода")
		return
	}

	fmt.Print("Напряжение (В): ")
	fmt.Scan(&U)
	if U <= 0 {
		fmt.Println("Ошибка: напряжение должно быть больше 0")
		return
	}

	v := speed(U, e, eMas)
	r := (Ra - Rk) / 2
	B := field(eMas, v, e, r)
	Ic := B / (mu0 * float64(l))

	p1 := plot.New()
	p1.Title.Text = "Ic(U)"
	p1.X.Label.Text = "U, В"
	p1.Y.Label.Text = "Ic, А"
	p1.Add(plotter.NewGrid())

	pts := make(plotter.XYs, 200)
	for i := 0; i < 200; i++ {
		Uval := 1.0 + float64(i)*(100.0-1.0)/199.0
		vVal := speed(Uval, e, eMas)
		BVal := field(eMas, vVal, e, r)
		IcVal := BVal / (mu0 * float64(l))
		pts[i].X = Uval
		pts[i].Y = IcVal
	}

	line, err := plotter.NewLine(pts)
	if err != nil {
		fmt.Println(err)
		return
	}
	p1.Add(line)
	p1.Legend.Add("Ic(U)", line)

	p2 := plot.New()
	p2.X.Label.Text = "x, м"
	p2.Y.Label.Text = "y, м"
	p2.Add(plotter.NewGrid())

	trajectory := make(plotter.XYs, 100)
	for i := 0; i < 100; i++ {
		alpha := float64(i) * 2 * math.Pi / 99
		trajectory[i].X = r * math.Cos(alpha)
		trajectory[i].Y = r * math.Sin(alpha)
	}

	path, err := plotter.NewLine(trajectory)
	if err != nil {
		fmt.Println(err)
		return
	}
	p2.Add(path)
	p2.Legend.Add(fmt.Sprintf("Траектория электрона при U = %.0fВ и Ic = %.2fA", U, Ic), path)

	center := plotter.XYs{{X: 0, Y: 0}}
	scatter, err := plotter.NewScatter(center)
	if err != nil {
		fmt.Println(err)
		return
	}
	scatter.GlyphStyle.Color = color.RGBA{R: 255, G: 0, B: 0, A: 255} // Красный цвет
	p2.Add(scatter)
	p2.Legend.Add("Центр окружности", scatter)

	if err := p1.Save(5*vg.Inch, 6*vg.Inch, "ic_u.png"); err != nil {
		fmt.Println(err)
		return
	}
	if err := p2.Save(5*vg.Inch, 6*vg.Inch, "trajectory.png"); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Графики сохранены как ic_u.png и trajectory.png")
}
