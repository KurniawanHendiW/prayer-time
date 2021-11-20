# prayer-time-calendar

## Rest API

Get key calendar
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

Get iCal with key
```
[GET] https://prayer-time-calendar.herokuapp.com/prayer-time/get?key=:key
```

Get available city
```
[GET] GET https://prayer-time-calendar.herokuapp.com/prayer-time/get-city/:city-name
```
