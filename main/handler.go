package main

import (
    "encoding/json"
    "fmt"
    "github.com/marni/goigc"
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
       fmt.Fprint(w, "It was I, Judas!")
   }
}

//Redirects requests from root to API
func HandlerAPIRedirect(w http.ResponseWriter, r *http.Request){

    if r.URL.Path == "/paragliding" || r.URL.Path == "/paragliding/"{
        //if there is nothing after paragliding in the URL, redirect to API
        http.Redirect(w, r, "/paragliding/api", http.StatusSeeOther)
    } else {
        //else return 404
        http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
        fmt.Fprint(w, "You're bad and should feel bad (failed to resolve URL)")
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
                //bson.NewObjectId(),
                track.Date,
                track.Pilot,
                track.GliderType,
                track.GliderID,
                TotalDistance(track),
                urlRequest.URL,
                Millisec(),
            }
            fmt.Fprintf(w, "track id:%v",encode.Timestamp) //return the id
            db.Add(encode) //add to the database
        } else {
            http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
        }

    } else if r.Method=="GET" {
        parts := strings.Split(r.URL.Path, "/")
        requestString:=""

        if len(parts)>4 { //this if block prevents  accessing space outside the array
            requestString = parts[4]
        }

        if !isNumeric(requestString) && requestString!="" {
            //check if the ID is numeric (and that the request was not for all tracks
            http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)

        }else{
        requestedID, _ := strconv.ParseInt(requestString, 10, 64)

            //the OR clause in the following if statements are to deal with trailing slashes.
            //We check whether an URL is of the correct length, or if the last segment of the URL is blank.
            if len(parts)==4 || (len(parts)==5 && parts[4]==""){
                //when entire array is requested
                tracks := db.GetAll() //add all tracks to array
                var ids []int64

                for i:= 0; i < len(tracks);i++{
                    ids = append(ids, tracks[i].Timestamp) //append all ids to new array
                }
                http.Header.Add(w.Header(), "content-type", "application/json")
                json.NewEncoder(w).Encode(ids)
                //fmt.Fprint(w, "This space for rent\n")

            } else if len(parts)==5 ||(len(parts)==6 &&parts[5]==""){
                //when single ID is requested
                track, err := db.Get(requestedID) //try to fetch track by ID

                if err==nil { //if that works, return it
                    http.Header.Add(w.Header(), "content-type", "application/json")
                    json.NewEncoder(w).Encode(track)
                }else{ //if this track could not be fetched, throw 404
                    http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
                    fmt.Fprint(w, "\nWhy are we still here, just to suffer?\n", requestedID)
                }

            }else if len(parts)==6|| (len(parts)==7 &&parts[6]==""){
                //when a single field is requested
                //This is an inefficient way of retrieving a field, but the GetField function in database.go
                //was not willing to cooperate, so this is the hack I ended up with.
                track, err :=db.Get(requestedID) //copy the track
                if err != nil{ //if that doesn't work throw 404
                    http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
                    fmt.Fprint(w, "\nEvery night, I feel my leg")
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
                    default: //Throw Bad Request
                        http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
                    }
                }

            }else{
                http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
                fmt.Fprint(w, "\nMy arm, even my fingers")
            }
        }
    }
}



