package ducks

import "fmt"

type Speaker interface {
	Speak()
}

type SpeakerFunc func()

func (s SpeakerFunc) Speak() {
	s()
}

func Quack() Speaker {
	return SpeakerFunc(func() {
		fmt.Println("Quack")
	})
}

func Mute() Speaker {
	return SpeakerFunc(func() {
		fmt.Println("<< silence >>")
	})
}
