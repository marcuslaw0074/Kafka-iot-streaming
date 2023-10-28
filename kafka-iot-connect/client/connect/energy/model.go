package energy

type EnergyᚖMills struct {
	D []TagValuePair `json:"d"`
}

type TagValuePair struct {
	Tag   string  `json:"tag"`
	Value float64 `json:"value"`
}

type TimeSeriesPoint struct {
	Id    string  `json:"id"`
	Value float64 `json:"value"`
	Time  string  `json:"time"`
}
