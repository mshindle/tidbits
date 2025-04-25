package ducks

import "fmt"

type Duck struct {
	Name    string
	Flyer   Flyer
	Speaker Speaker
}

func New(name string) *Duck {
	return &Duck{Name: name}
}

func (d *Duck) WithFlyer(f Flyer) *Duck {
	d.Flyer = f
	return d
}

func (d *Duck) WithSpeaker(s Speaker) *Duck {
	d.Speaker = s
	return d
}

func (d *Duck) Display() {
	fmt.Printf("I am a duck - a %s to be precise.\n", d.Name)
}

func (d *Duck) Fly() {
	d.Flyer.Fly()
}

func (d *Duck) Quack() {
	d.Speaker.Speak()
}

func (d *Duck) Swim() {
	fmt.Printf("All ducks can swim - including me a %s\n", d.Name)
}
