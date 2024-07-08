package open_meteo_parser

import "fmt"

func CalculateAQI(pm25, pm10, o3, no2, so2, co float64) int {
	pm25AQI := calculatePM25AQI(pm25)
	pm10AQI := calculatePM10AQI(pm10)
	o3AQI := calculateO3AQI(o3)
	no2AQI := calculateNO2AQI(no2)
	so2AQI := calculateSO2AQI(so2)
	coAQI := calculateCOAQI(co)

	fmt.Println("PM25 AQI: ", pm25AQI)
	fmt.Println("PM10 AQI: ", pm10AQI)
	fmt.Println("O3 AQI: ", o3AQI)
	fmt.Println("NO2 AQI: ", no2AQI)
	fmt.Println("SO2 AQI: ", so2AQI)
	fmt.Println("CO AQI: ", coAQI)

	return int(max(pm25AQI, pm10AQI, o3AQI, no2AQI, so2AQI, coAQI))
}

func calculatePM25AQI(pm25 float64) float64 {
	var iLow, iHigh, bLow, bHigh, cLow, cHigh float64
	if pm25 >= 0 && pm25 <= 12 {
		iHigh = 50
		iLow = 0
		bHigh = 12
		bLow = 0
	} else if pm25 > 12 && pm25 <= 35.4 {
		iHigh = 100
		iLow = 51
		bHigh = 35.4
		bLow = 12
	} else if pm25 > 35.4 && pm25 <= 55.4 {
		iHigh = 150
		iLow = 101
		bHigh = 55.4
		bLow = 35.4
	} else if pm25 > 55.4 && pm25 <= 150.4 {
		iHigh = 200
		iLow = 151
		bHigh = 150.4
		bLow = 55.4
	} else if pm25 > 150.4 && pm25 <= 250.4 {
		iHigh = 300
		iLow = 201
		bHigh = 250.4
		bLow = 150.4
	} else if pm25 > 250.4 && pm25 <= 350.4 {
		iHigh = 400
		iLow = 301
		bHigh = 350.4
		bLow = 250.4
	} else if pm25 > 350.4 && pm25 <= 500.4 {
		iHigh = 500
		iLow = 401
		bHigh = 500.4
		bLow = 350.4
	} else {
		return 500
	}

	return calculateAQI(pm25, iLow, iHigh, bLow, bHigh, cLow, cHigh, 500)
}

func calculatePM10AQI(pm10 float64) float64 {
	var iLow, iHigh, bLow, bHigh, cLow, cHigh float64

	if pm10 >= 0 && pm10 <= 54 {
		iHigh = 50
		iLow = 0
		bHigh = 54
		bLow = 0
	} else if pm10 > 54 && pm10 <= 154 {
		iHigh = 100
		iLow = 51
		bHigh = 154
		bLow = 54
	} else if pm10 > 154 && pm10 <= 254 {
		iHigh = 150
		iLow = 101
		bHigh = 254
		bLow = 154
	} else if pm10 > 254 && pm10 <= 354 {
		iHigh = 200
		iLow = 151
		bHigh = 354
		bLow = 254
	} else if pm10 > 354 && pm10 <= 424 {
		iHigh = 300
		iLow = 201
		bHigh = 424
		bLow = 354
	} else if pm10 > 424 && pm10 <= 504 {
		iHigh = 400
		iLow = 301
		bHigh = 504
		bLow = 424
	} else if pm10 > 504 && pm10 <= 604 {
		iHigh = 500
		iLow = 401
		bHigh = 604
		bLow = 504
	} else {
		return 500
	}
	return calculateAQI(pm10, iLow, iHigh, bLow, bHigh, cLow, cHigh, 500)
}

func calculateO3AQI(o3 float64) float64 {
	var iLow, iHigh, bLow, bHigh, cLow, cHigh float64

	// convert o3 from μg/m3 to ppb
	o3 = convertμgToPpb(o3, 48)

	if o3 >= 0 && o3 <= 54 {
		iHigh = 50
		iLow = 0
		bHigh = 54
		bLow = 0
	} else if o3 > 54 && o3 <= 70 {
		iHigh = 100
		iLow = 51
		bHigh = 70
		bLow = 54
	} else if o3 > 70 && o3 <= 85 {
		iHigh = 150
		iLow = 101
		bHigh = 85
		bLow = 70
	} else if o3 > 85 && o3 <= 105 {
		iHigh = 200
		iLow = 151
		bHigh = 105
		bLow = 85
	} else if o3 > 105 && o3 <= 200 {
		iHigh = 300
		iLow = 201
		bHigh = 200
		bLow = 105
	} else if o3 > 200 && o3 <= 504 {
		iHigh = 400
		iLow = 301
		bHigh = 504
		bLow = 200
	} else if o3 > 504 && o3 <= 604 {
		iHigh = 500
		iLow = 401
		bHigh = 604
		bLow = 504
	} else {
		return 500
	}
	return calculateAQI(o3, iLow, iHigh, bLow, bHigh, cLow, cHigh, 500)
}

func calculateNO2AQI(no2 float64) float64 {
	var iLow, iHigh, bLow, bHigh, cLow, cHigh float64

	if no2 >= 0 && no2 <= 53 {
		iHigh = 50
		iLow = 0
		bHigh = 53
		bLow = 0
	} else if no2 > 53 && no2 <= 100 {
		iHigh = 100
		iLow = 51
		bHigh = 100
		bLow = 53
	} else if no2 > 100 && no2 <= 360 {
		iHigh = 150
		iLow = 101
		bHigh = 360
		bLow = 100
	} else if no2 > 360 && no2 <= 649 {
		iHigh = 200
		iLow = 151
		bHigh = 649
		bLow = 360
	} else if no2 > 649 && no2 <= 1249 {
		iHigh = 300
		iLow = 201
		bHigh = 1249
		bLow = 649
	} else if no2 > 1249 && no2 <= 1649 {
		iHigh = 400
		iLow = 301
		bHigh = 1649
		bLow = 1249
	} else if no2 > 1649 && no2 <= 2049 {
		iHigh = 500
		iLow = 401
		bHigh = 2049
		bLow = 1649
	} else {
		return 500
	}
	return calculateAQI(no2, iLow, iHigh, bLow, bHigh, cLow, cHigh, 500)
}

func calculateSO2AQI(so2 float64) float64 {
	var iLow, iHigh, bLow, bHigh, cLow, cHigh float64

	// convert so2 from μg/m3 to ppb
	so2 = convertμgToPpb(so2, 64.066)

	if so2 >= 0 && so2 <= 35 {
		iHigh = 50
		iLow = 0
		bHigh = 35
		bLow = 0
	} else if so2 > 35 && so2 <= 75 {
		iHigh = 100
		iLow = 51
		bHigh = 75
		bLow = 35
	} else if so2 > 75 && so2 <= 185 {
		iHigh = 150
		iLow = 101
		bHigh = 185
		bLow = 75
	} else if so2 > 185 && so2 <= 304 {
		iHigh = 200
		iLow = 151
		bHigh = 304
		bLow = 185
	} else if so2 > 304 && so2 <= 604 {
		iHigh = 300
		iLow = 201
		bHigh = 604
		bLow = 304
	} else if so2 > 604 && so2 <= 804 {
		iHigh = 400
		iLow = 301
		bHigh = 804
		bLow = 604
	} else if so2 > 804 && so2 <= 1004 {
		iHigh = 500
		iLow = 401
		bHigh = 1004
		bLow = 804
	} else {
		return 500
	}
	return calculateAQI(so2, iLow, iHigh, bLow, bHigh, cLow, cHigh, 500)
}

func convertμgToPpb(value float64, molecularWeight float64) float64 {
	return value * (24.45 / molecularWeight)
}

func convertμgToPpm(value float64) float64 {
	return value / 1000
}

func calculateCOAQI(co float64) float64 {
	var iLow, iHigh, bLow, bHigh, cLow, cHigh float64

	// convert co from μg/m3 to ppm
	co = co / 1000
	co = 24.45 * (co / 28.01)

	if co >= 0 && co <= 4.4 {
		iHigh = 50
		iLow = 0
		bHigh = 4.4
		bLow = 0
	} else if co > 4.4 && co <= 9.4 {
		iHigh = 100
		iLow = 51
		bHigh = 9.4
		bLow = 4.4
	} else if co > 9.4 && co <= 12.4 {
		iHigh = 150
		iLow = 101
		bHigh = 12.4
		bLow = 9.4
	} else if co > 12.4 && co <= 15.4 {
		iHigh = 200
		iLow = 151
		bHigh = 15.4
		bLow = 12.4
	} else if co > 15.4 && co <= 30.4 {
		iHigh = 300
		iLow = 201
		bHigh = 30.4
		bLow = 15.4
	} else if co > 30.4 && co <= 40.4 {
		iHigh = 400
		iLow = 301
		bHigh = 40.4
		bLow = 30.4
	} else if co > 40.4 && co <= 50.4 {
		iHigh = 500
		iLow = 401
		bHigh = 50.4
		bLow = 40.4
	} else {
		return 500
	}
	return calculateAQI(co, iLow, iHigh, bLow, bHigh, cLow, cHigh, 500)
}

func calculateAQI(c, iLow, iHigh, bLow, bHigh, cLow, cHigh, iHighAQI float64) float64 {
	return ((iHigh-iLow)/(bHigh-bLow))*(c-bLow) + iLow
}

func max(a, b, c, d, e, f float64) float64 {
	max := a
	if b > max {
		max = b
	}
	if c > max {
		max = c
	}
	if d > max {
		max = d
	}
	if e > max {
		max = e
	}
	if f > max {
		max = f
	}
	return max
}
