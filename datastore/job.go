package datastore

type Job struct {
	Key         string `gorethink:"key"`
	Name        string `gorethink:"name"`
	Description string `gorethink:"description"`

	JobConfig JobConfig `gorethink:"config"`
}

type JobConfig struct {
	Environment []string `gorethink:"environment"`
	Commands    []string `gorethink:"commands"`
}
