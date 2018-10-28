package main

import (
    "log"
    "net/http"
    "os"
    "time"
)
//JSON STRUCTS
//Metadata ... stores metadata about app
type Metadata struct {
	Uptime 	string 	`json:"uptime"`
	Info    string 	`json:"info"`
	Version string 	`json:"version"`
}

//Track ... stores metadata about track
type Track struct {
    //Id        bson.ObjectId `bson:"_id,omitempty"`
    Hdate time.Time `json:"H_date"`
    Pilot string  `json:"pilot"`
    Glider string `json:"glider"`
    GliderID string `json:"glider_id"`
    TrackLength string `json:"track_length"`
    TrackURL string `json:"track_src_url"`
    Timestamp int64 `json:"timestamp"` //also used as a database ID

}

//Database ... Used to keep track of location of database and database accesories
type DBInfo struct {
    DBurl string
    DBname string
    TrackCollection string

}

//URLRequest ... stores URL request
type URLRequest struct {
    URL string  `json:"url"`
}

//global variables
var apiStruct Metadata //contains meta information
var start = time.Now() //keeps track of uptime


//global structs
var db DBInfo //help struct that contains info about the Database




//MAIN
func main() {
    db.DBname = "trackdb"
    db.TrackCollection = "tracks"
    db.DBurl = "mongodb://admin1:admin1@ds141813.mlab.com:41813/trackdb"


   // set port. if no port, default to 8080 (well not at the moment but you know, in theory)
    port := ":"+os.Getenv("PORT")
    if ( port == ":"){
        port = ":8080";
    }
    db.Init()

    apiStruct = Metadata{Uptime: "", Info:"Info for paragliding tracks.", Version: "v1" }
    http.HandleFunc("/paragliding/api/track/", HandlerTrack)
    http.HandleFunc("/paragliding/api/ticker", HandlerTicker)
	http.HandleFunc("/paragliding/api/", HandlerAPI)
    http.HandleFunc("/paragliding/", HandlerAPIRedirect)
	log.Fatal(http.ListenAndServe(port, nil))
}
