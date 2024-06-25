package partition

import "time"

type Partition struct {
	UUID       string
	Name       string
	Resume     string
	Prefix     string
	Url        string
	Updater    string
	UpdateTime time.Time
	Creator    string
	CreateTime time.Time
	Cluster    string
}

type CreatePartition struct {
	Uuid    string
	Name    string
	Resume  string
	Prefix  string
	Url     string
	Cluster string
}

type EditPartition struct {
	Uuid   string
	Name   *string
	Resume *string
	Prefix *string
	Url    *string
}
