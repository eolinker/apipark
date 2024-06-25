package organization_dto

type CreateOrganization struct {
	Id          string   `json:"id,omitempty"`          // Id is the UUID of the organization
	Name        string   `json:"name,omitempty"`        // Name is the name of the organization
	Description string   `json:"description,omitempty"` // Description is the description of the organization
	Master      string   `json:"master,omitempty"`      // Master is the UUID of the organization's master partition
	Partitions  []string `json:"partitions,omitempty"`  // Partition is the list of the organization's partition UUID
	Prefix      string   `json:"prefix,omitempty"`      // Prefix is the prefix of the organization's UUID
}

type EditOrganization struct {
	Name        *string   `json:"name,omitempty"`        // Name is the name of the organization
	Description *string   `json:"description,omitempty"` // Description is the description of the organization
	Master      *string   `json:"master,omitempty"`      // Master is the UUID of the organization's master partition
	Partitions  *[]string `json:"partitions,omitempty"`  // Partition is the
	Prefix      *string   `json:"prefix,omitempty"`      // Prefix is the prefix of the organization's UUID
}
