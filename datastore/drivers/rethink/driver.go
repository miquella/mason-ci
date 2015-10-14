package rethink

import (
	r "github.com/dancannon/gorethink"
	"github.com/miquella/mason-ci/datastore"
)

type rethinkDriver struct {
	session *r.Session
}

func NewRethinkDriver(address string, database string) (*rethinkDriver, error) {
	session, err := r.Connect(r.ConnectOpts{
		Address:  address,
		Database: database,
	})
	if err != nil {
		return nil, err
	}

	r.TableCreate("jobs", r.TableCreateOpts{PrimaryKey: "key"}).RunWrite(session)
	return &rethinkDriver{session: session}, nil
}

func (rd rethinkDriver) PutJob(job *datastore.Job) error {
	r.Table("jobs").Insert(job).RunWrite(rd.session)
	return nil
}

func (rd *rethinkDriver) GetJob(key string) (*datastore.Job, error) {
	cursor, err := r.Table("jobs").Get(key).Run(rd.session)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()

	var job datastore.Job
	err = cursor.One(&job)
	if err != nil {
		return nil, err
	}

	return &job, nil
}
