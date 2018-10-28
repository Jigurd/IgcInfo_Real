package main

import (
    "encoding/json"
    "net/http"
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
    requestStartTime := Millisec()

    tracks := db.GetAll() //pull all tracks from db

    var page []int64
    
    for i:=0;i < Min(pagesize, len(tracks));i++{ //put up to 5 elements in the page array
        page = append(page, tracks[i].Timestamp)
    }


    var ticker = Ticker{
        tracks[len(tracks)-1].Timestamp,
        tracks[0].Timestamp,
        tracks[pagesize-1].Timestamp,
        page,
        Millisec()-requestStartTime,
    }

    http.Header.Add(w.Header(), "content-type", "application/json")
    json.NewEncoder(w).Encode(ticker)
}

