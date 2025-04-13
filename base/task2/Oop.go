package main

import (
	"fmt"
	"math"
)

// task1

type Shape interface {
	Area() float64
	Perimeter() float64
}

type Rectangle struct {
	Width  float64
	Height float64
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

type Circle struct {
	Radius float64
}

func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}

func printShapeInfo(s Shape) {
	fmt.Printf("Area: %.2f\n", s.Area())
	fmt.Printf("Perimeter: %.2f\n\n", s.Perimeter())
}

// task2

type Person struct {
	Name string
	Age  int
}

type Employee struct {
	Person     // 组合 Person 结构体
	EmployeeID int
}

func (e *Employee) PrintInfo() {
	fmt.Printf("Employee Details:\n")
	fmt.Printf("• Name: %s\n• Age: %d\n• ID: %d\n",
		e.Name, e.Age, e.EmployeeID)
}
