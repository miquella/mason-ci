package main

import (
	"encoding/json"
	"github.com/miquella/mason-ci/datastore"
	"log"
	"net/http"
)

func init() {
	api := Router.PathPrefix("/api").Subrouter()

	jobsApi := api.Path("/jobs").Subrouter()
	jobsApi.Methods("GET").HandlerFunc(jobsIndexHandler)
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
