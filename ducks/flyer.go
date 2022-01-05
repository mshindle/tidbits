package ducks

import "fmt"

type Flyer interface {
	Fly()
}

type WingedFlyer int

func (w WingedFlyer) Fly() {
	fmt.Printf("I am flying with %d wing(s)\n", w)
}

type GroundedFlyer string

func (g GroundedFlyer) Fly() {
	fmt.Println("I am grounded as I cannot fly")
}
