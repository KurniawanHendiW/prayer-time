package prayerTime

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/prayer-time/client/redis"
	"github.com/prayer-time/client/waktusholat"

	ics "github.com/arran4/golang-ical"
	uuid "github.com/satori/go.uuid"
)

type Service interface {
	GetDataPrayerTime(req DataPrayerTimeRequest) (DataPrayerTimeResponse, error)
	GetKeyPrayerTime(req KeyPrayerTimeRequest) (KeyPrayerTimeResponse, error)
	GetCityByName(name string) ([]waktusholat.GetCityByNameResponse, error)
}

type service struct {
	waktuSholatSvc waktusholat.Service
	redisSvc       redis.Service
	serviceHost    string
	expiredKey     int64
}

func NewService(waktuSholatSvc waktusholat.Service, redisSvc redis.Service, serviceHost string, expiredKey int64) Service {
	return &service{
		waktuSholatSvc: waktuSholatSvc,
		redisSvc:       redisSvc,
		serviceHost:    serviceHost,
		expiredKey:     expiredKey,
	}
}

func (s *service) GetKeyPrayerTime(req KeyPrayerTimeRequest) (KeyPrayerTimeResponse, error) {
	byteData, err := json.Marshal(req)
	if err != nil {
		return KeyPrayerTimeResponse{}, err
	}

	serial := uuid.NewV4()

	redisKey := fmt.Sprintf("prayer-time-%s", serial.String())
	if err = s.redisSvc.Set(redisKey, string(byteData)).Error; err != nil {
		return KeyPrayerTimeResponse{}, err
	}

	if err = s.redisSvc.Expire(redisKey, s.expiredKey).Error; err != nil {
		return KeyPrayerTimeResponse{}, err
	}

	return KeyPrayerTimeResponse{
		Key:     serial.String(),
		Url:     fmt.Sprintf("%s/prayer-time/get?key=%s", s.serviceHost, serial.String()),
		Message: fmt.Sprintf("Url expired in %d minutes", s.expiredKey/60),
	}, nil
}

func (s *service) GetDataPrayerTime(req DataPrayerTimeRequest) (DataPrayerTimeResponse, error) {
	redisData, err := s.redisSvc.Get(fmt.Sprintf("prayer-time-%s", req.Key)).String()
	if err != nil {
		if redis.IsErrorNil(err) {
			return DataPrayerTimeResponse{}, fmt.Errorf("URL is expired")
		}

		return DataPrayerTimeResponse{}, err
	}

	keyPrayerTimeRequest := KeyPrayerTimeRequest{}
	if err = json.Unmarshal([]byte(redisData), &keyPrayerTimeRequest); err != nil {
		return DataPrayerTimeResponse{}, err
	}

	requestSholat := map[string]bool{}
	for _, sholat := range keyPrayerTimeRequest.Sholat {
		requestSholat[sholat] = true
	}
	if len(requestSholat) == 0 {
		requestSholat = MapSholat
	}

	requestDay := map[string]bool{}
	for _, day := range keyPrayerTimeRequest.Day {
		requestDay[day] = true
	}
	if len(requestDay) == 0 {
		requestDay = MapDay
	}

	resp, err := s.waktuSholatSvc.GetPrayTimes(waktusholat.PrayTimeRequest{
		City:      keyPrayerTimeRequest.City,
		StartDate: keyPrayerTimeRequest.StartDate,
		EndDate:   keyPrayerTimeRequest.EndDate,
	})
	if err != nil {
		return DataPrayerTimeResponse{}, err
	}

	timeNow := time.Now()

	cal := ics.NewCalendar()
	cal.SetColor("#009688")
	cal.SetMethod(ics.MethodPublish)
	cal.SetXWRCalName(fmt.Sprintf("Prayer time for %s - %s", resp.Results.Location.City, resp.Results.Location.Country))
	cal.SetXWRCalDesc(fmt.Sprintf("Prayer time for %s - %s", resp.Results.Location.City, resp.Results.Location.Country))
	for _, datetime := range resp.Results.Datetime {
		for day, prayTime := range datetime.Times {
			if _, ok := requestSholat[day]; !ok {
				continue
			}

			timeStart, err := time.Parse("2006-01-02 15:04", fmt.Sprintf("%s %s", datetime.Date.Gregorian, prayTime))
			if err != nil {
				return DataPrayerTimeResponse{}, err
			}

			if _, ok := requestDay[timeStart.Weekday().String()]; !ok {
				continue
			}

			event, err := s.addEventCalendar(cal, datetime, resp.Results.Location, day, timeNow)
			if err != nil {
				return DataPrayerTimeResponse{}, err
			}

			s.addAlarm(event, day)
		}
	}

	return DataPrayerTimeResponse{Data: cal.Serialize()}, nil
}

func (s *service) addEventCalendar(cal *ics.Calendar, datetime waktusholat.DateTime, location waktusholat.Location, day string, timeNow time.Time) (*ics.VEvent, error) {
	prayTime := datetime.Times[day]

	timeStart, err := time.Parse("2006-01-02 15:04", fmt.Sprintf("%s %s", datetime.Date.Gregorian, prayTime))
	if err != nil {
		return nil, err
	}

	event := cal.AddEvent(fmt.Sprintf("%d-%s", datetime.Date.Timestamp, day))
	event.SetDtStampTime(timeNow)
	event.SetProperty("X-GOOGLE-CALENDAR-CONTENT-TITLE", fmt.Sprintf("Time for %s", day))
	event.SetProperty("X-MICROSOFT-CDO-BUSYSTATUS", "TRUE")
	event.SetSummary(fmt.Sprintf("Time for %s (Pray)", day))
	event.SetDescription(fmt.Sprintf("Time for %s", day))
	event.SetLocation(fmt.Sprintf("%s - %s", location.City, location.Country))
	event.SetProperty(ics.ComponentPropertyCategories, "Prayer")
	//event.SetClass(ics.ClassificationPublic)
	event.SetTimeTransparency(ics.TransparencyTransparent)

	timeStart = timeStart.Add(time.Hour * -7)
	event.SetStartAt(timeStart)

	timeEnd := timeStart.Add(time.Minute * 30)
	event.SetEndAt(timeEnd)

	return event, nil
}

func (s *service) addAlarm(event *ics.VEvent, day string) {
	alarm := event.AddAlarm()
	alarm.SetTrigger("PT0S")
	alarm.SetAction(ics.ActionDisplay)
	alarm.SetProperty("DESCRIPTION", fmt.Sprintf("Time for %s", day))
}

func (s *service) GetCityByName(name string) ([]waktusholat.GetCityByNameResponse, error) {
	respGetCity, err := s.waktuSholatSvc.GetCityByName(name)
	if err != nil {
		return respGetCity, err
	}

	resp := []waktusholat.GetCityByNameResponse{}
	for _, city := range respGetCity {
		if strings.Contains(city.CityCode, "-") {
			continue
		}

		resp = append(resp, city)
	}

	return resp, nil
}
