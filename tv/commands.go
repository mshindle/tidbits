package tv

// Command represents the instruction being sent to the device
type Command interface {
	Press()
}

type onCommand struct {
	device Device
}

func (c *onCommand) Press() {
	c.device.On()
}

type offCommand struct {
	device Device
}

func (c *offCommand) Press() {
	c.device.Off()
}

type increaseVolumeCommand struct {
	device Device
}

func (c *increaseVolumeCommand) Press() {
	c.device.IncreaseVolume()
}

type decreaseVolumeCommand struct {
	device Device
}

func (c *decreaseVolumeCommand) Press() {
	c.device.DecreaseVolume()
}