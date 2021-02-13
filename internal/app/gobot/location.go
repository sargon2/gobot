package gobot

import (
	"errors"
	"os"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/jasonwinn/geocoder"
)

type Location struct {
	Latitude    float64
	Longitude   float64
	Description string
}

type LocationFinder struct{}

func NewLocationFinder() (*LocationFinder, error) {
	apiKey := os.Getenv("MAPQUEST_API_KEY")
	if apiKey == "" {
		return nil, errors.New("MAPQUEST_API_KEY must be set")
	}

	geocoder.SetAPIKey(apiKey)

	return &LocationFinder{}, nil
}

func (*LocationFinder) FindLocation(input string) (*Location, error) {
	if input == "" {
		return nil, errors.New("Need a location")
	}

	result, err := geocoder.FullGeocode(input)
	if err != nil {
		return nil, err
	}
	if result == nil || len(result.Results) == 0 {
		return nil, errors.New("No location results")
	}
	locations := result.Results[0].Locations
	// Choose a location.
	// This should be in a function, but can't be because of geocoder's strange nested struct declarations.
	locationToUse := locations[0]
	if len(input) == 5 && govalidator.IsInt(input) { // If the user gave a zip code, limit it to the US
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

	return &Location{
		Latitude:    locationToUse.LatLng.Lat,
		Longitude:   locationToUse.LatLng.Lng,
		Description: strings.Join(parts, ", "),
	}, nil
}
