package gobot

import (
	"fmt"
	"time"

	"github.com/sixdouglas/suncalc"
)

// TODO swap output if you query it after sunrise?  And show the following day.

type Sun struct {
	hub            Hub
	locationFinder *LocationFinder
}

func NewSun(hub Hub, locationFinder *LocationFinder) (*Sun, error) {
	ret := &Sun{
		hub:            hub,
		locationFinder: locationFinder,
	}

	hub.RegisterBangHandler("sun", ret.handleMessage)

	return ret, nil
}

func (t *Sun) handleMessage(source *MessageSource, message string) {
	location, err := t.locationFinder.FindLocation(message)
	if err != nil {
		t.hub.Message(source, err.Error())
		return
	}

	timeResult := suncalc.GetTimes(time.Now(), location.Latitude, location.Longitude)

	t.hub.Message(source, fmt.Sprint(
		location.Description,
		": ",
		timeResult[suncalc.Dawn].Time.In(location.TimeLocation).Round(time.Minute).Format("3:04pm"),
		" - ",
		timeResult[suncalc.Dusk].Time.In(location.TimeLocation).Round(time.Minute).Format("3:04pm MST"),
	))
}
