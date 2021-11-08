package waktusholat

import (
	"net/http"
	"net/url"

	"github.com/prayer-time/util"
)

type Service interface {
	GetPrayTimes(req PrayTimeRequest) (PrayTimeResponse, error)
}

type service struct {
	Host string
}

func NewService(host string) Service {
	return &service{Host: host}
}

func (s *service) GetPrayTimes(req PrayTimeRequest) (PrayTimeResponse, error) {
	uri := url.URL{}
	uri.Path = "/v2/times/dates.json"

	query := uri.Query()
	query.Set("city", req.City)
	query.Set("start", req.StartDate)
	query.Set("end", req.EndDate)

	uri.RawQuery = query.Encode()

	opts := util.ReqOpts{
		Host:        s.Host,
		Method:      http.MethodGet,
		RelativeURL: uri.String(),
	}

	var resp PrayTimeResponse
	if err := util.Call(&resp, opts); err != nil {
		return resp, err
	}

	return resp, nil
}
