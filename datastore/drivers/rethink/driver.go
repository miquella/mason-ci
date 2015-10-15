package rethink

import (
	"encoding/json"
	r "github.com/dancannon/gorethink"
	"github.com/miquella/mason-ci/datastore"
	"log"
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

func (rd *rethinkDriver) GetJobs() ([]*datastore.Job, error) {
	cursor, err := r.Table("jobs").Run(rd.session)
	if err != nil {
		return nil, err
	}

	jobs := []*datastore.Job{}
	err = cursor.All(&jobs)
	if err != nil {
		return nil, err
	}

	return jobs, nil
}

func (rd *rethinkDriver) NewBuild(jobKey string) (*datastore.Build, error) {
	// atomically update the build number of the job
	result, err := r.Table("jobs").Filter(
		r.Row.Field("key").Eq(jobKey),
	).Update(
		map[string]interface{}{
			"last_build_number": r.Row.Field("last_build_number").Add(1),
		},
		r.UpdateOpts{ReturnChanges: true},
	).RunWrite(rd.session)
	if err != nil {
		return nil, err
	}

	if len(result.Changes) < 1 {
		return nil, datastore.ErrJobNotFound
	} else if len(result.Changes) > 1 {
		log.Print("rethink: somehow updated more than one job getting a build number??")
	}

	// convert the udpated value into a job object
	d, err := json.Marshal(result.Changes[0].NewValue)
	if err != nil {
		return nil, err
	}

	var job datastore.Job
	err = json.Unmarshal(d, &job)
	if err != nil {
		return nil, err
	}

	// save the new build object
	build := &datastore.Build{
		JobId:  job.Id,
		Number: job.LastBuildNumber,
		Config: job.Config,
	}
	id, err := rd.PutBuild(build)
	if err != nil {
		return nil, err
	}

	build.Id = id
	return build, nil
}

func (rd *rethinkDriver) PutBuild(build *datastore.Build) (string, error) {
	result, err := r.Table("builds").Insert(
		r.Table("builds").Filter(
			r.And(
				r.Row.Field("job_id").Eq(build.JobId),
				r.Row.Field("number").Eq(build.Number),
			),
		).CoerceTo("array").Do(func(docs r.Term) interface{} {
			return r.Branch(
				r.Or(docs.IsEmpty(), docs.Field("id").Contains(build.Id)),
				build,
				r.Error("Build with job_id and number exists"),
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
		return build.Id, nil
	}
}

func (rd *rethinkDriver) GetBuild(jobKey string, buildNumber int64) (*datastore.Build, error) {
	cursor, err := r.Table("builds").Filter(
		r.And(
			r.Row.Field("job_id").Eq(jobKey),
			r.Row.Field("number").Eq(buildNumber),
		),
	).Run(rd.session)
	if err != nil {
		return nil, err
	}

	var build datastore.Build
	err = cursor.One(&build)
	if err != nil {
		return nil, err
	}

	return &build, nil
}
