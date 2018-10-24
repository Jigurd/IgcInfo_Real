package main

import (
    "encoding/json"
    "fmt"
    "github.com/marni/goigc"
    "log"
    "net/http"
    "os"
    "strconv"
    "strings"
    "time"
)

//JSON Structs

//IDArray ... Array keeping track of all IDs.
type IDArray struct {
    Ids []string `json:"ids"`
}
//Metadata ... stores metadata about app
type Metadata struct {
	Uptime 	string 	`json:"uptime"`
	Info    string 	`json:"info"`
	Version string 	`json:"version"`
}

//Track ... stores metadata about track
type Track struct {
    Hdate time.Time `json:"H_date"`
    Pilot string  `json:"pilot"`
    Glider string `json:"glider"`
    GliderID string `json:"glider_id"`
    TrackLength string `json:"track_length"`
}

//URLRequest ... stores URL request
type URLRequest struct {
    URL string  `json:"url"`
}

//global variables
var apiStruct Metadata //contains meta information
var start = time.Now() //keeps track of uptime
//LastID ... keeps track of used IDs
var LastID int

//arrays
var tracks = []Track{}//make(map[string]Track)
var ids IDArray

// HANDLERS

func handlerAPI(w http.ResponseWriter, r *http.Request) {

    //check that there is no rubbish behind api
    if r.URL.Path == "/paragliding/api" || r.URL.Path == "/paragliding/api/"{

        //finding uptime
        //I only track uptime until the point of days, as I find it unlikely that this service would
        //be running for weeks on end, let alone months or years.
        elapsedTime := time.Since(start)
        apiStruct.Uptime = fmt.Sprintf("P%dD%dH%dM%dS",
            int(elapsedTime.Hours()/24),    //number of days (no Days method available)
            int(elapsedTime.Hours())%24,    //number of hours
            int(elapsedTime.Minutes())%60,  //number of minutes
            int(elapsedTime.Seconds())%60,  //number of seconds
            )
        json.NewEncoder(w).Encode(apiStruct)

    } else {
        // if there is rubbish behind /api/, return 404
        http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
    }
}

func handlerIGC(w http.ResponseWriter, r *http.Request) {
    if (r.Method=="POST"){
      http.Header.Add(w.Header(), "content-type", "application/json")

        var urlRequest URLRequest

        decoder := json.NewDecoder(r.Body)
        decoder.Decode(&urlRequest)



        track, err := igc.ParseLocation(urlRequest.URL)
        if err == nil{
            id := fmt.Sprintf("track id: %d", LastID)
            ids.Ids = append(ids.Ids, id)
            LastID++

            encode := Track{track.Date, track.Pilot, track.GliderType, track.GliderID, "0",}
            encode.TrackLength = totalDistance(track)

            tracks = append(tracks, encode)
            json.NewEncoder(w).Encode(id)

        } else {
            http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
            fmt.Sprintf("%f", track)
        }
    } else if (r.Method=="GET"){
        parts :=strings.Split(r.URL.Path, "/")

        if len(parts)>=4 { //Check whether a specific id is being requested
            requestedID, err := strconv.Atoi(parts[4])

            if requestedID < LastID { //make sure requestedID is not out of bounds
                track := tracks[requestedID]

                if err != nil {
                    //the track does not exist
                    http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
                    json.NewEncoder(w).Encode(requestedID)
                    //fmt.Fprintf(w, "This is the first NotFound block\n")
                }

                if requestedID <= LastID {
                    if len(parts) == 5 {
                        http.Header.Add(w.Header(), "content-type", "application/json")
                        requestedTrack := Track{
                            track.Hdate,
                            track.Pilot,
                            track.Glider,
                            track.GliderID,
                            track.TrackLength,
                        }
                        json.NewEncoder(w).Encode(requestedTrack)
                    } else if len(parts) == 6 {
                        switch strings.ToLower(parts[5]) {
                        case "pilot":
                            fmt.Fprintf(w, track.Pilot)
                        case "glider":
                            fmt.Fprintf(w, track.Glider)
                        case "glider_id":
                            fmt.Fprintf(w, track.GliderID)
                        case "track_length":
                            fmt.Fprintf(w, "%f", track.TrackLength)
                        case "h_date":
                            fmt.Fprintf(w, "%v", track.Hdate)
                        default:
                            http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
                            //fmt.Fprintf(w, "This is the second NotFound block\n")
                        }

                    }
                } else {
                    http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
                    //fmt.Fprintf(w, "This is the third NotFound block\n")
                }
            }else{ //id is out of bounds, does not exist
                http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
                }
        }else{
            //return array of all ids
            http.Header.Add(w.Header(), "content-type", "application/json")
            json.NewEncoder(w).Encode(ids)
        }

    } else { //if
        http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
    }
}

func handlerAPIRedirect(w http.ResponseWriter, r *http.Request){

    if r.URL.Path == "/paragliding" || r.URL.Path == "/paragliding/"{
        //if there is nothing after paragliding in the URL, redirect to API
        http.Redirect(w, r, "/paragliding/api/", http.StatusSeeOther)
    } else {
        //else return 404
        http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
    }

}



//UTILITY FUNCTIONS

func totalDistance(t igc.Track) string {
    track := t
    totalDistance := 0.0
    for i := 0; i < len(track.Points)-1; i++ {
        totalDistance += track.Points[i].Distance(track.Points[i+1])
    }

    return fmt.Sprintf("%f", totalDistance)
}

func main() {
    LastID = 0
    ids = IDArray{make([]string, 0)}

   // set port. if no port, default to 8080 (well not at the moment but you know, in theory)
    port := ":"+os.Getenv("PORT")
    if ( port == ":"){
        port = ":8080";
    }

    apiStruct = Metadata{Uptime: "", Info:"Info for paragliding tracks.", Version: "v1" }
    http.HandleFunc("/paragliding/", handlerAPIRedirect)
	http.HandleFunc("/paragliding/api/", handlerAPI)
    http.HandleFunc("/paragliding/api/igc/", handlerIGC)
	log.Fatal(http.ListenAndServe(port, nil))
}
