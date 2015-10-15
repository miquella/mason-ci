package datastore

type Job struct {
	Id              string `json:"id,omitempty" gorethink:"id,omitempty"`
	Key             string `json:"key" gorethink:"key"`
	Name            string `json:"name" gorethink:"name"`
	Description     string `json:"description" gorethink:"description"`
	LastBuildNumber int64  `json:"last_build_number" gorethink:"last_build_number"`

	Config JobConfig `json:"config" gorethink:"config"`
}

type JobConfig struct {
	Environment []string `json:"environment" gorethink:"environment"`
	Commands    []string `json:"commands" gorethink:"commands"`
}

type Build struct {
	Id     string `json:"id,omitempty" gorethink:"id,omitempty"`
	JobId  string `json:"job_id" gorethink:"job_id"`
	Number int64  `json:"number" gorethink:"number"`

	Config JobConfig `json:"config" gorethink:"config"`
}
