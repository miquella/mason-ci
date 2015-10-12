package main

import (
	r "github.com/dancannon/gorethink"
)

type JobConfig struct {
}

type Datastore interface {
	PutJobConfig(*JobConfig)
	GetJobConfig(name string) *JobConfig
}

type RethinkDatastore struct {
	rethinkSession *r.Session
}

func (rd RethinkDatastore) PutJobConfig(config *JobConfig) {

}

func (rd RethinkDatastore) GetJobConfig(name string) *JobConfig {
	return nil
}
