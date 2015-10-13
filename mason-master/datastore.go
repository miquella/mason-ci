package main

import (
	"errors"
	"regexp"
	"strings"

	r "github.com/dancannon/gorethink"
)

var (
	MismatchedKeyErr = errors.New("mismatched job key")
)

type Job struct {
	Key         string   `gorethink:"key"`
	Name        string   `gorethink:"name"`
	Description string   `gorethink:"description"`
	Environment []string `gorethink:"environment"`
	Commands    []string `gorethink:"command"`
}

type Datastore interface {
	PutJob(*Job) (string, error)
	GetJob(key string) (*Job, error)
}

type RethinkDatastore struct {
	rethinkSession *r.Session
}

func NewRethinkDatastore(rethinkAddress string, rethinkDatabase string) (*RethinkDatastore, error) {
	rethinkSession, err := r.Connect(r.ConnectOpts{
		Address:  rethinkAddress,
		Database: rethinkDatabase,
	})
	if err != nil {
		return nil, err
	}
	r.TableCreate("jobs", r.TableCreateOpts{PrimaryKey: "key"}).RunWrite(rethinkSession)
	return &RethinkDatastore{
		rethinkSession: rethinkSession,
	}, nil
}

func calculateRethinkKey(name string) string {
	re := regexp.MustCompile("[^a-z0-9]+")
	return strings.Trim(re.ReplaceAllString(strings.ToLower(name), "-"), "-")
}

func (rd RethinkDatastore) PutJob(job *Job) (string, error) {
	if job.Key == "" {
		job.Key = calculateRethinkKey(job.Name)
	} else if job.Key != calculateRethinkKey(job.Name) {
		return "", MismatchedKeyErr
	}
	r.Table("jobs").Insert(job).RunWrite(rd.rethinkSession)

	return job.Key, nil
}

func (rd RethinkDatastore) GetJob(key string) (*Job, error) {
	cursor, err := r.Table("jobs").Filter(map[string]interface{}{"key": key}).Run(rd.rethinkSession)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()

	var job Job
	cursor.One(&job)

	return &job, nil
}
