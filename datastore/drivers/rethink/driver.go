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

	err = createSchema(session)
	if err != nil {
		return nil, err
	}

	return &rethinkDriver{session: session}, nil
}

func getTableList(s *r.Session) (map[string]bool, error) {
	cursor, err := r.TableList().Run(s)
	if err != nil {
		return nil, err
	}

	tableList := []string{}
	err = cursor.All(&tableList)
	if err != nil {
		return nil, err
	}

	tables := map[string]bool{}
	for _, table := range tableList {
		tables[table] = true
	}

	return tables, nil
}

func createSchema(s *r.Session) error {
	tables, err := getTableList(s)

	if err == nil && !tables["jobs"] {
		err = r.TableCreate("jobs").Exec(s)
	}

	if err == nil && !tables["builds"] {
		err = r.TableCreate("builds").Exec(s)
	}

	return err
}

func (rd rethinkDriver) PutJob(job *datastore.Job) (string, error) {
	result, err := r.Table("jobs").Insert(
		r.Table("jobs").Filter(r.Row.Field("key").Eq(job.Key)).CoerceTo("array").Do(func(docs r.Term) interface{} {
			return r.Branch(
				r.Or(docs.IsEmpty(), docs.Field("id").Contains(job.Id)),
				job,
				r.Error("Job with key exists"),
			)
		}),
		r.InsertOpts{Conflict: "replace"},
	).RunWrite(rd.session)
	if err != nil {
		return "", err
	}

	if len(result.GeneratedKeys) > 0 {
		return result.GeneratedKeys[0], nil
	} else {
		return job.Id, nil
	}
}

func (rd *rethinkDriver) GetJob(key string) (*datastore.Job, error) {
	cursor, err := r.Table("jobs").Filter(r.Row.Field("key").Eq(key)).Run(rd.session)
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
