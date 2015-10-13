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

type JobConfig struct {
	Key         string   `gorethink:"key"`
	Name        string   `gorethink:"name"`
	Description string   `gorethink:"description"`
	Environment []string `gorethink:"environment"`
	Commands    []string `gorethink:"command"`
}

type Datastore interface {
	PutJobConfig(*JobConfig) (string, error)
	GetJobConfig(key string) (*JobConfig, error)
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

func (rd RethinkDatastore) PutJobConfig(config *JobConfig) (string, error) {
	if config.Key == "" {
		config.Key = calculateRethinkKey(config.Name)
	} else if config.Key != calculateRethinkKey(config.Name) {
		return "", MismatchedKeyErr
	}
	r.Table("jobs").Insert(config).RunWrite(rd.rethinkSession)

	return config.Key, nil
}

func (rd RethinkDatastore) GetJobConfig(key string) (*JobConfig, error) {
	cursor, err := r.Table("jobs").Filter(map[string]interface{}{"key": key}).Run(rd.rethinkSession)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()

	var config JobConfig
	cursor.One(&config)

	return &config, nil
}
