package waktusholat

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/prayer-time/util"
)

type Service interface {
	GetPrayTimes(req PrayTimeRequest) (PrayTimeResponse, error)
	GetCityByName(name string) ([]GetCityByNameResponse, error)
}

type service struct {
	waktuSholatHost string
	apiPrayZoneHost string
}

func NewService(waktuSholatHost, apiPrayZoneHost string) Service {
	return &service{
		waktuSholatHost: waktuSholatHost,
		apiPrayZoneHost: apiPrayZoneHost,
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
	if err := util.Call(&resp, opts); err != nil {
		return resp, err
	}

	return resp, nil
}

func (s *service) GetPrayTimes(req PrayTimeRequest) (PrayTimeResponse, error) {
	uri := url.URL{}
	uri.Path = "/v2/times/dates.json"

	query := uri.Query()
	query.Set("city", req.City)
	query.Set("start", req.StartDate)
	query.Set("end", req.EndDate)
	query.Set("school", "3")

	uri.RawQuery = query.Encode()

	opts := util.ReqOpts{
		Host:        s.apiPrayZoneHost,
		Method:      http.MethodGet,
		RelativeURL: uri.String(),
	}

	var resp PrayTimeResponse
	if err := util.Call(&resp, opts); err != nil {
		return resp, err
	}

	return resp, nil
}
