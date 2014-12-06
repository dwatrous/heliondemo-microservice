package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
	"fmt"

	"upper.io/db"
	"upper.io/db/mongo"
)

var (
	mongosettings mongo.ConnectionURL
	mongodb       db.Database

	listen       = flag.String("l", os.Getenv("PORT"), "listener address")
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

	// TODO(mischief): insert results into a mongo collection

	// the only possible error here is that the collection doesn't exist.
	// in that case, inserting will create it (for mongo driver)
	//col, _ := mongodb.Collection(surveyname)
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

func main() {
	flag.Parse()

	mongosettings.Database = *mongodbname
	mongosettings.Address = db.Host(*mongoaddress)

	// XXX disabled for now since we don't insert yet.

	//var err error
	//mongodb, err = db.Open(mongo.Adapter, mongosettings)
	//if err != nil {
	//	log.Fatal(err)
	//}

	fs := new(FormServer)

	mux := http.NewServeMux()
	mux.HandleFunc("/survey/", fs.Survey)

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
