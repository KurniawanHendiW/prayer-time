package waktusholat

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/prayer-time/util"
)

type Service interface {
	GetPrayTimes(req PrayTimeRequest) (PrayTimeResponse, error)
	GetCityByName(name string) ([]GetCityByNameResponse, error)
}

type service struct {
	waktuSholatHost string
	apiPrayZoneHost string
	debugLog        bool
}

func NewService(waktuSholatHost, apiPrayZoneHost string, debugLog bool) Service {
	return &service{
		waktuSholatHost: waktuSholatHost,
		apiPrayZoneHost: apiPrayZoneHost,
		debugLog:        debugLog,
	}
}

func (s *service) GetCityByName(name string) ([]GetCityByNameResponse, error) {
	uri := url.URL{}
	uri.Path = fmt.Sprintf("/api/docs/ajax/cities/%s", name)

	opts := util.ReqOpts{
		Host:        s.waktuSholatHost,
		Method:      http.MethodGet,
		RelativeURL: uri.String(),
	}

	var resp []GetCityByNameResponse
	if err := util.Call(&resp, opts, s.debugLog); err != nil {
		return resp, err
	}

	return resp, nil
}

func (s *service) GetPrayTimes(req PrayTimeRequest) (PrayTimeResponse, error) {
	timeStart, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return PrayTimeResponse{}, err
	}

	timeEnd, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return PrayTimeResponse{}, err
	}

	if timeStart.After(timeEnd) {
		return PrayTimeResponse{}, fmt.Errorf("end date must be greater than start date")
	}

	uri := url.URL{}
	uri.Path = "/v2/times/dates.json"

	school := SchoolMap[req.CountryCode]

	query := uri.Query()
	query.Set("city", req.City)
	query.Set("start", req.StartDate)
	query.Set("end", req.EndDate)
	query.Set("school", fmt.Sprint(school.ID))

	uri.RawQuery = query.Encode()

	opts := util.ReqOpts{
		Host:        s.apiPrayZoneHost,
		Method:      http.MethodGet,
		RelativeURL: uri.String(),
	}

	var resp PrayTimeResponse
	if err := util.Call(&resp, opts, s.debugLog); err != nil {
		return resp, err
	}

	return resp, nil
}
