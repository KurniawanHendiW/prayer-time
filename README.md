# prayer-time-calendar

iCal URL supports in google calendar, outlook calendar, and calendar app on macOS.

## Rest API

### Get key calendar
```
[POST] https://prayer-time-calendar.herokuapp.com/prayer-time/get-key
```
Example request Body
```
{
    "city": "jakarta",
    "start_date": "2021-11-10",
    "end_date": "2021-11-10",
    "day": ["Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"],
    "sholat": ["Imsak", "Sunrise", "Dhuhr", "Asr", "Sunset", "Maghrib", "Isha", "Midnight"]
}
```
Example response
```
{
    "url": "https://prayer-time-calendar.herokuapp.com/prayer-time/get?key=1a62ce13-862a-426e-b1bb-a749483d6152",
    "key": "1a62ce13-862a-426e-b1bb-a749483d6152",
    "message": "Url expired in 10 minutes"
}
```

### Get iCal with key
```
[GET] https://prayer-time-calendar.herokuapp.com/prayer-time/get?key=:key
```

### Get available city
```
[GET] https://prayer-time-calendar.herokuapp.com/prayer-time/get-city?name=jakarta
```
Example response
```
[
  {
    "cityCode": "jakarta",
    "cityName": "Jakarta",
    "countryCode": "ID",
    "countryName": "Republic of Indonesia"
  }
]
```

# Credits
1. Waktu shalat https://waktusholat.org/api/docs/today
2. GoLang iCal https://github.com/arran4/golang-ical
