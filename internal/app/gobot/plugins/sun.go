package gobot

import (
	"fmt"
	"time"

	"github.com/sargon2/gobot/internal/app/gobot"

	"github.com/sixdouglas/suncalc"
)

type Sun struct {
	hub            *gobot.Hub
	locationFinder *gobot.LocationFinder
}

func NewSun(hub *gobot.Hub, locationFinder *gobot.LocationFinder) *Sun {
	ret := &Sun{
		hub:            hub,
		locationFinder: locationFinder,
	}

	hub.RegisterBangHandler("sun", ret.handleMessage)

	return ret
}

func (t *Sun) handleMessage(source *gobot.MessageSource, message string) {
	location, err := t.locationFinder.FindLocation(message)
	if err != nil {
		t.hub.Message(source, err.Error())
		return
	}

	timeResult := suncalc.GetTimes(time.Now(), location.Latitude, location.Longitude)

	t.hub.Message(source, fmt.Sprint(
		location.Description,
		": ",
		timeResult[suncalc.Dawn].Value.In(location.TimeLocation).Round(time.Minute).Format("3:04pm"),
		" - ",
		timeResult[suncalc.Dusk].Value.In(location.TimeLocation).Round(time.Minute).Format("3:04pm MST"),
	))
}
