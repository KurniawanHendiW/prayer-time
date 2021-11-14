package prayerTime

var (
	MapDay = map[string]bool{
		"Sunday":    true,
		"Monday":    true,
		"Tuesday":   true,
		"Wednesday": true,
		"Thursday":  true,
		"Friday":    true,
		"Saturday":  true,
	}

	MapSholat = map[string]bool{
		"Imsak":    true,
		"Sunrise":  true,
		"Dhuhr":    true,
		"Asr":      true,
		"Sunset":   true,
		"Maghrib":  true,
		"Isha":     true,
		"Midnight": true,
	}
)

type (
	KeyPrayerTimeRequest struct {
		City      string   `json:"city" binding:"required"`
		StartDate string   `json:"start_date" binding:"required"`
		EndDate   string   `json:"end_date" binding:"required"`
		Day       []string `json:"day" binding:"required"`
		Sholat    []string `json:"sholat" binding:"required"`
	}

	KeyPrayerTimeResponse struct {
		Url     string `json:"url"`
		Key     string `json:"key"`
		Message string `json:"message"`
	}

	DataPrayerTimeRequest struct {
		Key string `json:"key"`
	}

	DataPrayerTimeResponse struct {
		Data string `json:"data"`
	}
)
