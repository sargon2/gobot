package gobot

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/jasonwinn/geocoder"
	"github.com/sixdouglas/suncalc"
)

// Intended output: Sunrise is technically __ but there will be light at __.  Sunset is technically __ but it will be dark at __.
// TODO swap output if you query it after sunrise?  And show the following day.

type Sun struct {
	hub Hub
}

func NewSun(hub Hub) (*Sun, error) {
	ret := &Sun{
		hub: hub,
	}
	hub.RegisterBangHandler("sun", ret.handleMessage)

	apiKey := os.Getenv("MAPQUEST_API_KEY")
	if apiKey == "" {
		return nil, errors.New("MAPQUEST_API_KEY must be set")
	}

	geocoder.SetAPIKey(apiKey)

	return ret, nil
}

func (t *Sun) handleMessage(source *MessageSource, message string) {
	if message == "" {
		t.hub.Message(source, "Need a location")
		return
	}

	result, err := geocoder.FullGeocode(message)
	if err != nil {
		t.hub.Message(source, err.Error())
		return
	}
	if result == nil || len(result.Results) == 0 {
		t.hub.Message(source, "No location results")
		return
	}
	locations := result.Results[0].Locations
	// Use the first location, unless any are in the US, in which case use the first of those.
	// This should be in a function, but can't be because of geocoder's strange nested struct declarations.
	locationToUse := locations[0]
	if len(message) == 5 && govalidator.IsInt(message) { // If the user gave a zip code, limit it to the US
		for _, location := range locations {
			if location.AdminArea1 == "US" {
				locationToUse = location
			}
		}
	}
	var parts []string
	for _, s := range []string{locationToUse.AdminArea6, locationToUse.AdminArea5, locationToUse.AdminArea4, locationToUse.AdminArea3, locationToUse.AdminArea1} {
		if s != "" {
			parts = append(parts, s)
		}
	}
	locationStr := strings.Join(parts, ", ")

	timeResult := suncalc.GetTimes(time.Now(), locationToUse.LatLng.Lat, locationToUse.LatLng.Lng)
	t.hub.Message(source, fmt.Sprint(
		locationStr,
		": ",
		timeResult[suncalc.Dawn].Time.Round(time.Minute).Format("3:04pm"),
		" - ",
		timeResult[suncalc.Dusk].Time.Round(time.Minute).Format("3:04pm"),
	))
}
