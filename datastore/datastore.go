package datastore

import (
	"errors"
	"regexp"
	"strings"
)

var (
	ErrInvalidJobKey    = errors.New("invalid job key")
	ErrJobNotFound      = errors.New("job not found")
	ErrMismatchedJobKey = errors.New("mismatched job key")
)

var (
	jobKeyRegex = regexp.MustCompile("[^A-Za-z0-9_]+")
)

type Datastore struct {
	Driver DatastoreDriver
}

type DatastoreDriver interface {
	PutJob(*Job) (id string, err error)
	GetJob(key string) (*Job, error)
	GetJobs() ([]*Job, error)

	NewBuild(jobKey string) (*Build, error)
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

	id, err := d.Driver.PutJob(job)
	if err != nil {
		return err
	}

	job.Id = id
	return nil
}

func (d *Datastore) LoadJob(key string) (*Job, error) {
	return d.Driver.GetJob(key)
}

func (d *Datastore) LoadJobs() ([]*Job, error) {
	return d.Driver.GetJobs()
}

func (d *Datastore) CreateBuild(jobKey string) (*Build, error) {
	return d.Driver.NewBuild(jobKey)
}

func calculateJobKey(name string) string {
	key := strings.ToLower(strings.TrimSpace(name))
	return strings.Trim(jobKeyRegex.ReplaceAllString(key, "-"), "-")
}
