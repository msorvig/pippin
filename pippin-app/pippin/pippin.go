package pippin

import (
	"appengine"
	"appengine/datastore"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

// Struct Pi represents a registered raspberry pi
type Pi struct {
	Name      string `datastore:"name" json:"name"`
	Ip        string `datastore:"ip" json:"ip"`
	LastSeen  string `datastore:"lastSeen" json:"lastSeen"`
	PingCount int    `datastore:"pingCount" json:"pingCount"`
}

var piListKind = "pi-list"

var router = mux.NewRouter()

func init() {
	// forward all requests to the mux router
	http.Handle("/", router)

	// split into different handler funcs
	router.HandleFunc("/", index).Methods("GET")
	router.HandleFunc("/pi", postPi).Methods("POST")
	router.HandleFunc("/pi/{piname}", putPi).Methods("PUT")
	router.HandleFunc("/pi/{piname}", patchPi).Methods("PATCH")
	router.HandleFunc("/pi/{piname}", getPi).Methods("GET")
	router.HandleFunc("/pi/{piname}/{piattribute}", getPiAttribute).Methods("GET")
	router.HandleFunc("/test", getTest).Methods("GET")
}

func index(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	// Get all registered pis from the data store
	q := datastore.NewQuery(piListKind)
	var pis []Pi
	_, err := q.GetAll(c, &pis)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Print list of all pis
	fmt.Fprint(w, "Registered Raspberry Pis:\n")
	for _, pi := range pis {
		fmt.Fprint(w, pi.Name, " -> ", pi.Ip, pi.LastSeen, " ", pi.PingCount, "\n")
	}
}

func populatePiStruct(r *http.Request) Pi {
	r.ParseForm()
	pingCount, _ := strconv.Atoi(r.Form.Get("pingCount"))

	return Pi{
		Name:      r.Form.Get("name"),
		Ip:        r.Form.Get("ip"),
		LastSeen:  r.Form.Get("lastSeen"),
		PingCount: pingCount,
	}
}

/*
// Queries the datastore for a Pi entry by name
func queryPiStruct(name string, c appengine.Context) (Pi, error) {
	q := datastore.NewQuery(piListKind).Filter("name =", name)
	t := q.Run(c)

	var pi Pi
	_, err := t.Next(&pi)

	if err == datastore.Done {
		return pi, err
	}
	if err != nil {
		return pi, http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return pi, nil
}
*/

// Register a pi by POSTing the pi properties
// curl --data "name=KjokkenPi&ip=192.168.0.1" http://localhost:8080/pi
func postPi(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	r.ParseForm()

	pi := populatePiStruct(r)

	// Store pi object in data store
	_, err := datastore.Put(c, datastore.NewKey(c, piListKind, pi.Name, 0, nil), &pi)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, "{ pi: \"", pi.Name, "\" } ")
}

// Register a pi by PUTing the pi properties
func putPi(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	vars := mux.Vars(r)
	name := vars["piname"]

	// Populate the a Pi struct with form data.
	// Get the name from the PUT url.
	pi := populatePiStruct(r)
	pi.Name = name

	// Store pi object in data store
	c := appengine.NewContext(r)
	_, err := datastore.Put(c, datastore.NewKey(c, piListKind, pi.Name, 0, nil), &pi)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// PATCH one or more pi attributes
func patchPi(w http.ResponseWriter, r *http.Request) {
	// Get pi name from request
	vars := mux.Vars(r)
	name := vars["piname"]

	// Retrieve pi object from data store
	c := appengine.NewContext(r)
	q := datastore.NewQuery(piListKind).Filter("name =", name)
	t := q.Run(c)
	var pi Pi
	_, err := t.Next(&pi)
	if err == datastore.Done {
		http.Error(w, "404 Not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set Pi object property
	r.ParseForm()

	// Updating the name is not allowed
	formName := r.Form.Get("name")
	if len(formName) != 0 {
		http.Error(w, "404 Not found", http.StatusNotFound)
		return
	}
	ip := r.Form.Get("ip")
	if len(ip) != 0 {
		pi.Ip = ip
	}
	lastSeen := r.Form.Get("lastSeen")
	if len(lastSeen) != 0 {
		pi.LastSeen = lastSeen
	}
	pingCount := r.Form.Get("pingCount")
	if len(pingCount) != 0 {
		pi.PingCount, _ = strconv.Atoi(r.Form.Get("pingCount"))
	}

	//	fmt.Fprint(w, "name ", , "\n")
	fmt.Fprint(w, "pingCount ", r.Form.Get("pingCount"), " ", pi.PingCount, "\n")

	// Store pi object in data store
	_, err = datastore.Put(c, datastore.NewKey(c, piListKind, name, 0, nil), &pi)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Error(w, "200 OK", http.StatusOK)
	return
}

// GET a complete Pi structire
func getPi(w http.ResponseWriter, r *http.Request) {
	// Get pi name from request
	vars := mux.Vars(r)
	piname := vars["piname"]

	// Retrieve pi object from data store
	c := appengine.NewContext(r)
	q := datastore.NewQuery(piListKind).Filter("name =", piname)
	t := q.Run(c)
	var pi Pi
	_, err := t.Next(&pi)
	if err == datastore.Done {
		http.Error(w, "404 Not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return JSON representatikon of Pi object
	buffer, _ := json.MarshalIndent(pi, "", " ")
	fmt.Fprint(w, string(buffer))
}

// GET a Pi attribute
func getPiAttribute(w http.ResponseWriter, r *http.Request) {
	// Get pi name and property/attribute from request
	vars := mux.Vars(r)
	piname := vars["piname"]
	piproperty := vars["piattribute"]

	// Get pi entry from data store
	c := appengine.NewContext(r)
	q := datastore.NewQuery(piListKind).Filter("name =", piname)
	t := q.Run(c)
	var pi Pi
	_, err := t.Next(&pi)
	if err == datastore.Done {
		http.Error(w, "404 Not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Print attribute value in plain text
	w.Header().Set("Content-Type", "text/plain")
	switch piproperty {
	case "name":
		fmt.Fprint(w, pi.Name)
	case "ip":
		fmt.Fprint(w, pi.Ip)
	case "lastSeen":
		fmt.Fprint(w, pi.LastSeen)
	case "pingCount":
		fmt.Fprint(w, pi.PingCount)
	default:
		http.Error(w, "404 Not found", http.StatusNotFound)
	}
}

// GET a test string (minmum latency testing etc)
func getTest(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "test")
}
