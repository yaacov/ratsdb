package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/yaacov/ratsdb/models"
)

// Helper function to send an error message to client
func SendErr(w http.ResponseWriter, message models.ErrorResponse) {
	output, _ := json.Marshal(message)
	http.Error(w, string(output), http.StatusInternalServerError)
}

// Handler function for index page request
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	// index page is not implemented
	message := models.ErrorResponse{Status: "running", Message: "unknown request"}
	SendErr(w, message)
}

// Handler function for samples list request
func GetSamplesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	key := r.FormValue("key")
	start := r.FormValue("start")
	end := r.FormValue("end")
	labels := r.FormValue("labels")
	bucket := r.FormValue("bucket")

	if bucket == `` {
		// get samples
		samples := DataList(key, start, end, labels)
		if samples == nil {
			w.Write([]byte("[]"))
			return
		}

		// output json data to user
		if err := json.NewEncoder(w).Encode(samples); err != nil {
			message := models.ErrorResponse{Status: "running", Message: err.Error()}
			SendErr(w, message)
			return
		}
	} else {
		// get buckets
		buckets := DataBuckets(key, start, end, labels, bucket)
		if buckets == nil {
			w.Write([]byte("[]"))
			return
		}

		// output json data to user
		if err := json.NewEncoder(w).Encode(buckets); err != nil {
			message := models.ErrorResponse{Status: "running", Message: err.Error()}
			SendErr(w, message)
			return
		}
	}
}

// Handler function for one smaple request
func GetOneSampleHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["sampleId"], 0, 64)

	// get one sample
	sample := DataFind(int(id))
	if sample.Id == 0 {
		w.Write([]byte("{}"))
		return
	}

	// output json data to user
	if err := json.NewEncoder(w).Encode(sample); err != nil {
		message := models.ErrorResponse{Status: "running", Message: err.Error()}
		SendErr(w, message)
		return
	}
}

// Handler function for one smaple request
func PostSampleHandler(w http.ResponseWriter, r *http.Request) {
	var sample models.Sample
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)

	// get request body
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		message := models.ErrorResponse{Status: "running", Message: err.Error()}
		SendErr(w, message)
		return
	}
	if err := r.Body.Close(); err != nil {
		message := models.ErrorResponse{Status: "running", Message: err.Error()}
		SendErr(w, message)
		return
	}

	// parse request body
	if err := json.Unmarshal(body, &sample); err != nil {
		message := models.ErrorResponse{Status: "running", Message: err.Error()}
		SendErr(w, message)
		return
	}

	// check key and label
	keyOk, _ := regexp.MatchString("^[a-zA-Z0-9-]+$", sample.Key)
	labelsOk, _ := regexp.MatchString("^[a-zA-Z0-9-,]*$", sample.Labels)
	if !keyOk || !labelsOk {
		message := models.ErrorResponse{Status: "running", Message: "bad key or labels"}
		SendErr(w, message)
		return
	}

	// create a new sample
	t := DataCreate(sample)
	if err := json.NewEncoder(w).Encode(t); err != nil {
		message := models.ErrorResponse{Status: "running", Message: err.Error()}
		SendErr(w, message)
		return
	}
}

// Handler function for deletion of one smaple request
func DeleteOneSampleHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["sampleId"], 0, 64)

	// delete one sample
	DataDelete(int(id))
	w.Write([]byte("{}"))
}
