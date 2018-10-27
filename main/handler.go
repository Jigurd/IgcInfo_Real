package main

import (
    "encoding/json"
    "fmt"
    "github.com/marni/goigc"
    "gopkg.in/mgo.v2/bson"
    "net/http"
    "strconv"
    "strings"
    "time"
)


//returns metadata about the API
func HandlerAPI(w http.ResponseWriter, r *http.Request) {

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

//Redirects requests from root to API
func HandlerAPIRedirect(w http.ResponseWriter, r *http.Request){

    if r.URL.Path == "/paragliding" || r.URL.Path == "/paragliding/"{
        //if there is nothing after paragliding in the URL, redirect to API
        http.Redirect(w, r, "/paragliding/api/", http.StatusSeeOther)
    } else {
        //else return 404
        http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
    }

}

////test function pls ignore
//func HandlerTest(w http.ResponseWriter, r *http.Request){
//    test:=Track {"Test","Case"}
//    fmt.Fprintf(w, "%f", db.Add(test))
//}

//Handles POST and GET requests of tracks
func HandlerTrack(w http.ResponseWriter, r *http.Request) {
    if r.Method=="POST"{
        http.Header.Add(w.Header(), "content-type", "application/json")
        var urlRequest URLRequest
//
        decoder := json.NewDecoder(r.Body)
        decoder.Decode(&urlRequest)
        track, err := igc.ParseLocation(urlRequest.URL)

        if err == nil{ //if there is no problem with the parse

            encode := Track{ //encode track JSON
                bson.NewObjectId(),
                track.Date,
                track.Pilot,
                track.GliderType,
                track.GliderID,
                TotalDistance(track),
                urlRequest.URL,
                Millisec(),
            }
            fmt.Sprintf("track id:%v",encode.Timestamp) //return the id
            db.Add(encode) //add to the database
        } else {
            http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
        }

    } else if r.Method=="GET" {

        parts := strings.Split(r.URL.Path, "/")

        requestString := parts[4]
        //fmt.Fprint(w, requestedID)

        if !isNumeric(requestString) && requestString!="" {
            //check if the ID is numeric (and that the request was not for all tracks
            http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
            fmt.Fprint(w, "\nThis is the first bad request block!")

        }else{
        requestedID, _ := strconv.ParseInt(requestString, 10, 64)

            switch len(parts) {
            case 4: //if the entire array is requested
                fmt.Fprint(w, "This space for rent\n")

            case 5: //if a single track is requested
                track, err := db.Get(requestedID) //try to fetch track by ID
                track = Track{
                    track.Id,
                    track.Hdate,
                    track.Pilot,
                    track.Glider,
                    track.GliderID,
                    track.TrackLength,
                    track.TrackURL,
                    track.Timestamp,
                 }

                if err==nil { //if that works, return it
                    http.Header.Add(w.Header(), "content-type", "application/json")
                    json.NewEncoder(w).Encode(requestedID)
                }else{ //if this track could not be fetched, throw 404
                    http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
                }

            case 6: //if a single field is requested
                track, err :=db.Get(requestedID) //copy the track
                track = Track{
                    track.Id,
                    track.Hdate,
                    track.Pilot,
                    track.Glider,
                    track.GliderID,
                    track.TrackLength,
                    track.TrackURL,
                    track.Timestamp,
                }
                if err != nil{ //if that doesn't work throw 404
                    http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
                }else{  //if it does, return selected field
                    switch strings.ToLower(parts[5]) {
                        case "glider":
                            fmt.Fprintf(w, track.Glider)
                        case "glider_id":
                            fmt.Fprintf(w, track.GliderID)
                        case "h_date":
                            fmt.Fprintf(w, "%v", track.Hdate)
                        case "pilot":
                            fmt.Fprintf(w, track.Pilot)
                        case "timestamp":
                            fmt.Fprintf(w, "%v", track.Timestamp)
                        case "track_length":
                            fmt.Fprintf(w, "%v", track.TrackLength)
                        case "track_src_url":
                            fmt.Fprintf(w, "%v", track.TrackURL)
                        default: //
                            http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
                            fmt.Fprint(w, "\nSomething messed up here too tbh")
                        }
                }
            default:
                http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
            }
            //sjekk om vi har å gjøre med get array, get track, eller get field
            //hvis get array
            //hent timestamp for alle objekter fra db
            //returner array med alle IDer
            //hvis get track
            //request track fra db
            //returner den som json
            //hvis get $FIELD
            //request $FIELD fra db
            //returner $FIELD
        }
    }
//===================================================================
//        parts :=strings.Split(r.URL.Path, "/")
//
//        if len(parts)>=4 { //Check whether a specific id is being requested
//            requestedID, err := strconv.Atoi(parts[4])
//
//            if requestedID < LastID { //make sure requestedID is not out of bounds
//                track := tracks[requestedID]
//
//                if err != nil {
//                    //the track does not exist
//                    http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
//                    json.NewEncoder(w).Encode(requestedID)
//                    //fmt.Fprintf(w, "This is the first NotFound block\n")
//                }
//
//                if requestedID <= LastID {
//                    if len(parts) == 5 {
//                        http.Header.Add(w.Header(), "content-type", "application/json")
//                        requestedTrack := Track{
//                            track.Id,
//                            track.Hdate,
//                            track.Pilot,
//                            track.Glider,
//                            track.GliderID,
//                            track.TrackLength,
//                            track.TrackURL,
//                            track.Timestamp,
//                        }
//                        json.NewEncoder(w).Encode(requestedTrack)
//                    } else if len(parts) == 6 {
//                        switch strings.ToLower(parts[5]) {
//                        case "pilot":
//                            fmt.Fprintf(w, track.Pilot)
//                        case "glider":
//                            fmt.Fprintf(w, track.Glider)
//                        case "glider_id":
//                            fmt.Fprintf(w, track.GliderID)
//                        case "track_length":
//                            fmt.Fprintf(w, "%f", track.TrackLength)
//                        case "h_date":
//                            fmt.Fprintf(w, "%v", track.Hdate)
//                        case "track_src_url":
//                            fmt.Fprintf(w, "%f", track.TrackURL)
//                        case "timestamp":
//                            fmt.Fprintf(w, "%f", track.Timestamp)
//                        default:
//                            http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
//                            //fmt.Fprintf(w, "This is the second NotFound block\n")
//                        }
//                    }
//                } else {
//                    http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
//                    //fmt.Fprintf(w, "This is the third NotFound block\n")
//                }
//            }else{ //id is out of bounds, does not exist
//                http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
//            }
//        }else{
//            //return array of all ids
//            http.Header.Add(w.Header(), "content-type", "application/json")
//            json.NewEncoder(w).Encode(ids)
//        }
//
//    } else { //if
//        http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
}



