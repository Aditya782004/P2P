package main

import "fmt"

type Shape interface {
	CalculateShape()
	CalculatePerimeter()
}

type Rectangle struct {
	Width  int
	Length int
}
type Circle struct {
	Radius float32
	Pie    float32
}

func (r *Rectangle) CalculateShape() {
	result := r.Width * r.Length
	fmt.Println("The Result Is: ", result)
}
func (c *Circle) CalculateShape() {
	result := c.Pie * (c.Radius * c.Radius)
	fmt.Println("The Result Is: ", result)
}
func (r *Rectangle) CalculatePerimeter() {
	result := 2 * (r.Width + r.Length)
	fmt.Println("The Result Is: ", result)
}
func (c *Circle) CalculatePerimeter() {
	result := 2 * (c.Radius + c.Pie)
	fmt.Println("The Result Is: ", result)
}

func calculateCaller(s Shape) {
	s.CalculateShape()

}

func main() {
	circle := new(Circle)
	circle = &Circle{
		Radius: 2,
		Pie:    3.14,
	}
	rectangle := &Rectangle{
		Width:  2,
		Length: 4,
	}
	calculateCaller(circle)
	calculateCaller(rectangle)

}
