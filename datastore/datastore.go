package datastore

import (
	"errors"
	"regexp"
	"strings"
)

var (
	ErrInvalidJobKey    = errors.New("invalid job key")
	ErrMismatchedJobKey = errors.New("mismatched job key")
)

var (
	jobKeyRegex = regexp.MustCompile("[^A-Za-z0-9_]+")
)

type Datastore struct {
	Driver DatastoreDriver
}

type DatastoreDriver interface {
	PutJob(*Job) error
	GetJob(key string) (*Job, error)
}

func NewDatastore(driver DatastoreDriver) *Datastore {
	return &Datastore{
		Driver: driver,
	}
}

func (d *Datastore) SaveJob(job *Job) error {
	if job.Key == "" {
		job.Key = calculateJobKey(job.Name)
	} else if job.Key != calculateJobKey(job.Name) {
		return ErrMismatchedJobKey
	}

	if job.Key == "" {
		return ErrInvalidJobKey
	}

	return d.Driver.PutJob(job)
}

func (d *Datastore) LoadJob(key string) (*Job, error) {
	return d.Driver.GetJob(key)
}

func calculateJobKey(name string) string {
	key := strings.ToLower(strings.TrimSpace(name))
	return strings.Trim(jobKeyRegex.ReplaceAllString(key, "-"), "-")
}
