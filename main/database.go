package main

import (
    "fmt"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
)

func (db *DBInfo) Init() {
    session, err := mgo.Dial(db.DBurl)
    if err != nil {
        panic(err)
    }
    defer session.Close()

    index := mgo.Index{
        Key:        []string{"timestamp"},
        Unique:     true,
        DropDups:   true,
        Background: true,
        Sparse:     true,
    }

    err = session.DB(db.DBname).C(db.TrackCollection).EnsureIndex(index)
    if err != nil {
        panic(err)
    }
}


//adds track to storage
func (db *DBInfo) Add(s Track) error {
    session, err := mgo.Dial(db.DBurl)
    if err != nil {
        panic(err)
    }
    defer session.Close()

    err = session.DB(db.DBname).C(db.TrackCollection).Insert(s)

    if err != nil {
        fmt.Printf("error in Insert(): %v", err.Error())
        return err
    }

    return nil
}

//counts number of tracks in storage
func (db *DBInfo) Count() int {
    session, err := mgo.Dial(db.DBurl)
    if err != nil {
        panic(err)
    }
    defer session.Close()

    // handle to "db"
    count, err := session.DB(db.DBname).C(db.TrackCollection).Count()
    if err != nil {
        fmt.Printf("error in Count(): %v", err.Error())
        return -1
    }

    return count
}

//returns one track
func (db *DBInfo) Get(keyID int64) (Track, error) {
    session, err := mgo.Dial(db.DBurl)
    if err != nil {
        panic(err)
    }
    defer session.Close()

    track := Track{}


    err = session.DB(db.DBname).C(db.TrackCollection).Find(bson.M{"timestamp": keyID}).One(&track)

    return track, err
}


//returns one field from a track as a string
func (db *DBInfo) GetField(keyID int64, field string) (string, bool) {
    session, err := mgo.Dial(db.DBurl)
    if err != nil {
        panic(err)
    }
    defer session.Close()

    var returnfield string

    allWasGood := true

    err = session.DB(db.DBname).C(db.TrackCollection).Find(bson.M{"timestamp": keyID}).Select(bson.M{field: 1}).One(&returnfield)
    if err != nil {
        allWasGood = false
    }

    return returnfield, allWasGood
}


//returns slice with all Tracks
func (db *DBInfo) GetAll() []Track {
    session, err := mgo.Dial(db.DBurl)
    if err != nil {
        panic(err)
    }
    defer session.Close()

    var all []Track

    err = session.DB(db.DBname).C(db.TrackCollection).Find(bson.M{}).All(&all)
    if err != nil {
        return []Track{}
    }

    return all
}