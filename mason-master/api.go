package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/miquella/mason-ci/datastore"
	"log"
	"net/http"
)

func init() {
	api := Router.PathPrefix("/api").Subrouter()

	api.Path("/jobs").Methods("GET").Handler(http.RedirectHandler("jobs/", http.StatusMovedPermanently))
	api.Path("/jobs/").Methods("GET").HandlerFunc(jobsIndexHandler)
	api.Path("/jobs/").Methods("POST").HandlerFunc(jobCreateHandler)
	api.Path("/jobs/{job}").Methods("GET").HandlerFunc(jobGetHandler)
}

func apiJob(job *datastore.Job) map[string]interface{} {
	return map[string]interface{}{
		"key":               job.Key,
		"name":              job.Name,
		"description":       job.Description,
		"last_build_number": job.LastBuildNumber,
		"config": map[string]interface{}{
			"environment": job.Config.Environment,
			"commands":    job.Config.Commands,
		},
	}
}

func jobsIndexHandler(w http.ResponseWriter, r *http.Request) {
	jobs, err := Store.LoadJobs()
	if err != nil {
		log.Printf("Failed to query jobs: %s", err)
		http.Error(w, "failed while querying jobs", http.StatusInternalServerError)
		return
	}

	apiJobs := make([]map[string]interface{}, len(jobs))
	for i, job := range jobs {
		apiJobs[i] = apiJob(job)
	}

	encoder := json.NewEncoder(w)
	err = encoder.Encode(&apiJobs)
	if err != nil {
		log.Printf("Failed to encode jobs data: %s", err)
	}
}

func jobCreateHandler(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header["Content-Type"]
	if len(contentType) != 1 || contentType[0] != "application/json" {
		http.Error(w, "expected application/json entity", http.StatusBadRequest)
		return
	}

	var postedJob *datastore.Job
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&postedJob)
	if err != nil {
		log.Printf("Failed to decode job object: %s", err)
		http.Error(w, "invalid job object provided", http.StatusBadRequest)
		return
	}

	job := datastore.Job{
		Name:        postedJob.Name,
		Description: postedJob.Description,
		Config: datastore.JobConfig{
			Environment: postedJob.Config.Environment,
			Commands:    postedJob.Config.Commands,
		},
	}
	err = Store.SaveJob(&job)
	if err != nil {
		log.Printf("Failed to save job object: %s", err)
		http.Error(w, "failed to save job object", http.StatusInternalServerError)
		return
	}

	encoder := json.NewEncoder(w)
	err = encoder.Encode(apiJob(&job))
	if err != nil {
		log.Printf("Failed to encode new job data: %s", err)
	}
}

func jobGetHandler(w http.ResponseWriter, r *http.Request) {
	jobKey := mux.Vars(r)["job"]
	job, err := Store.LoadJob(jobKey)
	if err == datastore.ErrJobNotFound {
		http.NotFound(w, r)
		return
	} else if err != nil {
		log.Printf("Failed to query job (key: %s): %s", jobKey, err)
		http.Error(w, "error querying job", http.StatusInternalServerError)
		return
	}

	encoder := json.NewEncoder(w)
	err = encoder.Encode(apiJob(job))
	if err != nil {
		log.Printf("Failed to encode job data: %s", err)
	}
}
