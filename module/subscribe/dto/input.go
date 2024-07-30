package subscribe_dto

type AddSubscriber struct {
	Uuid    string `json:"uuid"`
	Service string `json:"service" aocheck:"service"`
	Project string `json:"subscriber" aocheck:"project"`
	Applier string `json:"applier" aocheck:"user"`
	//Cluster []string `json:"partition" aocheck:"partition"`
}

type Approve struct {
	//Cluster []string `json:"partition" aocheck:"partition"`
	Opinion string `json:"opinion"`
	Operate string `json:"operate"`
}
