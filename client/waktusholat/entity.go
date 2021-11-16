package waktusholat

type (
	PrayTimeRequest struct {
		City      string `json:"city" binding:"required"`
		StartDate string `json:"start_date" binding:"required"`
		EndDate   string `json:"end_date" binding:"required"`
	}

	PrayTimeResponse struct {
		Code    int                  `json:"code"`
		Status  string               `json:"status"`
		Results DataPrayTimeResponse `json:"results"`
	}

	DataPrayTimeResponse struct {
		Datetime []DateTime `json:"datetime"`
		Location Location   `json:"location"`
	}

	DateTime struct {
		Times map[string]string `json:"times"`
		Date  struct {
			Timestamp int64  `json:"timestamp"`
			Gregorian string `json:"gregorian"`
			Hijri     string `json:"hijri"`
		}
	}

	Location struct {
		City        string `json:"city"`
		Country     string `json:"country"`
		CountryCode string `json:"country_code"`
		Timezone    string `json:"timezone"`
	}

	GetCityByNameResponse struct {
		CityCode    string `json:"cityCode"`
		CityName    string `json:"cityName"`
		CountryCode string `json:"countryCode"`
		CountryName string `json:"countryName"`
	}
)
