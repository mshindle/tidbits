package tv

import "github.com/apex/log"

// Device represents any TV (i.e. Sony, Samsung, LG, etc) which can execute these common TV functions
type Device interface {
	On()
	Off()
	IncreaseVolume()
	DecreaseVolume()
}

type Remote struct {
	On Command
	Off Command
	VolUp Command
	VolDown Command
}

type Television struct {
	isOn bool
	volume int
}

func (t *Television) On() {
	t.isOn = true
	log.WithField("isOn", t.isOn).Info("turning tv on/off")
}

func (t *Television) Off() {
	t.isOn = false
	log.WithField("isOn", t.isOn).Info("turning tv on/off")
}

func (t *Television) IncreaseVolume() {
	if !t.isOn {
		log.Error("Cannot change volume as tv is off")
		return
	}
	if t.volume >= 0 && t.volume < 100 {
		t.volume++
		log.WithField("volume", t.volume).Info("increased volume")
	} else {
		log.WithField("volume", t.volume).Info("at max volume")
	}
}

func (t *Television) DecreaseVolume() {
	if !t.isOn {
		log.Error("Cannot change volume as tv is off")
		return
	}
	if t.volume > 0 && t.volume <= 100 {
		t.volume--
		log.WithField("volume", t.volume).Info("decreased volume")
	} else {
		log.WithField("volume", t.volume).Info("at min volume")
	}
}

func (t *Television) GrabRemote() *Remote {
	return &Remote{
		On: &onCommand{t},
		Off: &offCommand{t},
		VolUp: &increaseVolumeCommand{t},
		VolDown: &decreaseVolumeCommand{t},
	}
}

func New() *Television {
	return &Television{false,10}
}
