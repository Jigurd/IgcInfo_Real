package main

import (
    "fmt"
    "github.com/marni/goigc"
    "regexp"
    "time"
)


//Finds the total distance of an IGC track
func TotalDistance(t igc.Track) string {
	track := t
	totalDistance := 0.0
	for i := 0; i < len(track.Points)-1; i++ {
		totalDistance += track.Points[i].Distance(track.Points[i+1])
	}

	return fmt.Sprintf("%f", totalDistance)
}

func isNumeric(s string) bool { //Checks whether given string is numeric
    value, _ := regexp.MatchString("[0-9]+", s)
    return value
}


//returns monotonic time as an int64
func Millisec() int64{
    nowmilli := time.Now().UnixNano()/1000000 //
    return nowmilli
}
