package main

import (
	"encoding/json"
	"flag"
	"strconv"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"upper.io/db"
	"upper.io/db/mongo"

	"github.com/cloudfoundry-community/go-cfenv"
)

// string to port number
func atop(s string) uint {
	p, _ := strconv.Atoi(s)
	return uint(p)
}

func GetenvDefault(key, dfault string) string {
	r := os.Getenv("PORT")
	if r != "" {
		return r
	}

	return dfault
}

var (
	appenv        *cfenv.App
	mongosettings mongo.ConnectionURL
	mongodb       db.Database
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
	if err != nil {
		if err != db.ErrCollectionDoesNotExists {
			http.Error(w, "error retrieving collection: "+err.Error(), 400)
			return
		}
	} else {
		err = col.Truncate()

		if err != nil {
			http.Error(w, "error truncating collection: "+err.Error(), 400)
			return
		}
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
	if err != nil {
		if err != db.ErrCollectionDoesNotExists {
			http.Error(w, "error retrieving collection: "+err.Error(), 400)
			return
		}
	} else {
		err = col.Truncate()

		if err != nil {
			http.Error(w, "error truncating collection: "+err.Error(), 400)
			return
		}
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
	var err error

	flag.Parse()

	appenv, err = cfenv.Current()
	if err != nil {
		log.Fatalf("can't initialize app environment: %s", err)
	}

	srvname := "mongo-heliondemo"
	mongosrv, err := appenv.Services.WithName(srvname)
	if err != nil {
		log.Fatalf("missing service %q: %s", srvname, err)
	}

	mongosettings.Database = mongosrv.Credentials["database"]
	mongosettings.Address = db.HostPort(mongosrv.Credentials["host"], atop(mongosrv.Credentials["port"]))
	mongosettings.User = mongosrv.Credentials["user"]
	mongosettings.Password = mongosrv.Credentials["pass"]

	mongodb, err = db.Open(mongo.Adapter, mongosettings)
	if err != nil {
		log.Fatalf("failed to connect to mongo at %s: %s", mongosettings.Address, err)
	}

	fs := new(FormServer)

	mux := http.NewServeMux()
	mux.HandleFunc("/survey/", fs.Survey)
	mux.HandleFunc("/result/", fs.Result)

	listen := ":" + GetenvDefault("PORT", "8080")
	s := &http.Server{
		Addr:           listen,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	fmt.Printf("Listening on %v\n", listen)
	log.Fatal(s.ListenAndServe())
}
