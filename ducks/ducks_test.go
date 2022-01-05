package ducks

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type duck struct {
	name    string
	flyer   Flyer
	speaker Speaker
}

var ducks = []struct {
	name         string
	duck         duck
	flightOutput string
	quackOutput  string
}{
	{
		name: "mallard",
		duck: duck{
			name:    "mallard",
			flyer:   WingedFlyer(2),
			speaker: Quack(),
		},
		flightOutput: "I am flying with 2 wing(s)\n",
		quackOutput:  "Quack\n",
	},
	{
		name: "donald",
		duck: duck{
			name:    "donald",
			flyer:   GroundedFlyer(""),
			speaker: SpeakerFunc(func() { fmt.Println("Aw, phooey!") }),
		},
		flightOutput: "I am grounded as I cannot fly\n",
		quackOutput:  "Aw, phooey!\n",
	},
}

func captureOutput(f func()) string {
	orig := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	_ = w.Close()
	output, _ := ioutil.ReadAll(r)
	os.Stdout = orig

	return string(output)
}

func TestDuck_Fly(t *testing.T) {
	for _, tt := range ducks {
		t.Run(tt.name, func(t *testing.T) {
			d := New(tt.duck.name).WithFlyer(tt.duck.flyer)
			output := captureOutput(d.Fly)
			_ = assert.Equal(t, tt.flightOutput, output)
		})
	}
}

func TestDuck_Quack(t *testing.T) {
	for _, tt := range ducks {
		t.Run(tt.name, func(t *testing.T) {
			d := New(tt.duck.name).WithSpeaker(tt.duck.speaker)
			output := captureOutput(d.Quack)
			_ = assert.Equal(t, tt.quackOutput, output)
		})
	}
}
