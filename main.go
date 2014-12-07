package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"upper.io/db"
	"upper.io/db/mongo"
)

func GetenvDefault(key, dfault string) string {
	r := os.Getenv("PORT")
	if r != "" {
		return r
	}

	return dfault
}

var (
	mongosettings mongo.ConnectionURL
	mongodb       db.Database

	listen       = flag.String("l", ":"+GetenvDefault("PORT", "8080"), "listener address")
	mongoaddress = flag.String("mongo", "127.0.0.1", "mongo address")
	mongodbname  = flag.String("dbname", "surveys", "mongo database name")
)

func SurveyFile(name string) string {
	return "survey/" + name + ".json"
}

type FormServer struct {
}

func (fs *FormServer) SurveyPost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	// strip leading /survey/
	surveyname := r.URL.Path[8:]

	// make sure this survey exists by looking for the source json file
	_, err := os.Stat(SurveyFile(surveyname))
	if err != nil {
		http.Error(w, "non-existent survey", 400)
		return
	}

	// the only possible error here is that the collection doesn't exist.
	// in that case, inserting will create it (for mongo driver)
	col, err := mongodb.Collection("results")
	//if err != db.ErrCollectionDoesNotExist {
	if err != nil {
		http.Error(w, "retrieving mongo collection: "+err.Error(), 400)
		return
	}

	doc := map[string]interface{}{}

	doc["survey"] = surveyname
	doc["results"] = r.Form

	col.Append(doc)
	if err != nil {
		http.Error(w, "insertion error: "+err.Error(), 400)
		return
	}
}

func (fs *FormServer) SurveyRender(w http.ResponseWriter, r *http.Request) {
	var f []FormElement

	// strip leading /survey/
	surveyname := r.URL.Path[8:]
	log.Printf("loading survey %q", surveyname)

	surveypath := "survey/" + surveyname + ".json"

	js, err := ioutil.ReadFile(surveypath)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	err = json.Unmarshal(js, &f)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	err = formTmpl.Execute(w, f)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
}

func (fs *FormServer) Survey(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		fs.SurveyRender(w, r)
	case "POST":
		fs.SurveyPost(w, r)
	default:
		http.Error(w, "bad method", 400)
	}
}

func (fs *FormServer) Result(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
	default:
		http.Error(w, "bad method", 400)
		return
	}

	// strip leading /result/
	surveyname := r.URL.Path[8:]

	col, err := mongodb.Collection("results")
	//if err != db.ErrCollectionDoesNotExist {
	if err != nil {
		http.Error(w, "retrieving mongo collection: "+err.Error(), 400)
		return
	}

	var res db.Result

	res = col.Find(db.Cond{"survey": surveyname})

	var results []map[string]interface{}

	// Query all results and fill the birthdays variable with them.
	err = res.All(&results)
	if err != nil {
		http.Error(w, "retrieving survey results: "+err.Error(), 400)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func main() {
	flag.Parse()

	mongosettings.Database = *mongodbname
	mongosettings.Address = db.Host(*mongoaddress)

	// XXX disabled for now since we don't insert yet.
	var err error
	mongodb, err = db.Open(mongo.Adapter, mongosettings)
	if err != nil {
		log.Fatal(err)
	}

	fs := new(FormServer)

	mux := http.NewServeMux()
	mux.HandleFunc("/survey/", fs.Survey)
	mux.HandleFunc("/result/", fs.Result)

	s := &http.Server{
		Addr:           *listen,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	fmt.Printf("Listening on %v\n", *listen)
	log.Fatal(s.ListenAndServe())
}
