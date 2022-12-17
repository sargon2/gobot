package gobot

import (
	"errors"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/locationservice"
	"github.com/zsefvlol/timezonemapper"
)

type Location struct {
	Latitude     float64
	Longitude    float64
	Description  string
	TimeLocation *time.Location
}

type LocationFinder struct {
	locationService *locationservice.LocationService
}

func NewLocationFinder() *LocationFinder {
	locationService := locationservice.New(session.Must(session.NewSession()), aws.NewConfig().WithRegion("us-east-1"))
	return &LocationFinder{locationService: locationService}
}

func (l *LocationFinder) FindLocation(input string) (*Location, error) {
	if input == "" {
		return nil, errors.New("Need a location")
	}

	maxResults := int64(1) // can't inline these because go sucks
	indexName := "GobotPlaceIndex"
	// Somewhere in the middle of the US
	biasPositionLat := float64(37.0902)
	biasPositionLong := float64(-95.7129)

	biasPosition := []*float64{&biasPositionLong, &biasPositionLat}
	asdf := &locationservice.SearchPlaceIndexForTextInput{
		IndexName:    &indexName,
		Text:         &input,
		BiasPosition: biasPosition,
		MaxResults:   &maxResults,
	}
	output, err := l.locationService.SearchPlaceIndexForText(asdf)
	if err != nil {
		return nil, err
	}

	lat := output.Results[0].Place.Geometry.Point[1]
	lng := output.Results[0].Place.Geometry.Point[0]

	timeLocation, _ := time.LoadLocation(timezonemapper.LatLngToTimezoneString(*lat, *lng))

	return &Location{
		Latitude:     *lat,
		Longitude:    *lng,
		Description:  *output.Results[0].Place.Label,
		TimeLocation: timeLocation,
	}, nil
}
