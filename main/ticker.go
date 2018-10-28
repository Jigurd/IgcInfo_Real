package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "strconv"
    "strings"
)

const pagesize = 5 //sets number of ids returned per page

type Ticker struct{
    Tlatest int64 `json:"t_latest"` //last timestamp in db
    Tstart int64 `json:"t_start"`   //first timestamp on this page
    Tstop int64 `json:"t_stop"`     //last timestamp on this page
    Tracks []int64 `json:"tracks"`  //all timestamps on page
    Processing int64 `json:"processing"` //processing time in ms
}

func HandlerTicker(w http.ResponseWriter, r *http.Request){
    //This is in an incredibly inefficient way to do this kind of operation, but I couldn't for the life of me get
    //the GetField function in database.go to cooperate. So, we're stuck with this.

    requestStartTime := Millisec()  //mark start of function
    tracks := db.GetAll()           //pull all tracks from db
    var page []int64                //define page array (which will be returned)
    var requestedID int64           //define variable as it is assigned in another scope
    var reqIndex int
    var AllOK = true

    parts := strings.Split(r.URL.Path, "/")
    requestString:="0"

    if len(parts)>4 { //this if block prevents accessing space outside the array
        requestString = parts[4]
    }

    if !isNumeric(requestString) && requestString!="0" {
        //check if the ID is numeric (and that the request was not for all tracks
        http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
        fmt.Fprint(w, "The body I've lost")

    }else{
        var err error
        requestedID, err = strconv.ParseInt(requestString, 10, 64)

        if requestedID>0{   //if there is a valid ID, then find the index for that. Otherwise use index 0.
            reqIndex, AllOK = FindIndex(tracks, requestedID)
        } else{
            reqIndex=0
        }

        if err!=nil{ //if there is an error
            http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
            fmt.Fprint(w, "The comrades I have lost (error !=nil)")
        }else if !IsInSlice(tracks, requestedID)&&requestedID!=0{ //If the ID is not in slice and not default
            http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
            fmt.Fprint(w,"It won't stop hurting. (Not in slice)")
        }else if !AllOK{ //if there was some error in finding the index
            http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
            fmt.Fprint(w, "It's like they're all still here (AllOK)")
        }else{ //if nothing has gone wrong, commence function

            for i:=0;i < pagesize;i++ { //put up to 5 elements in the page array
                if (i + reqIndex) < len(tracks) { //make sure we're not exceeding the length of the slice
                    page = append(page, tracks[reqIndex+i].Timestamp)
                }
            }

            var ticker = Ticker{
                tracks[len(tracks)-1].Timestamp,
                tracks[reqIndex].Timestamp,
                tracks[Min(reqIndex+pagesize-1, pagesize-1)].Timestamp,//set Tstop to be Tstart+5 or the end
                page,                                                    //of the slice, whichever is smaller
                Millisec()-requestStartTime,
            }

            http.Header.Add(w.Header(), "content-type", "application/json")
            fmt.Fprint(w, "reqIndex: ", reqIndex, "\n")
            json.NewEncoder(w).Encode(ticker)
        }
    }
}
