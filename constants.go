package main

// Wait represents a unit of time to wait
type Wait struct {
	units    string
	quantity int
	seconds  int
}

var timeUnits = []string{"second", "minute", "hour", "day", "week", "month"}
var unitMap = map[string]int{
	"second": 1,
	"minute": 60,
	"hour":   3600,
	"day":    86400,
	"week":   604800,
	"month":  2419200000,
}
