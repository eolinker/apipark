package organization

import (
	"time"
)

// 组织

type Organization struct {
	Id          int64     `gorm:"column:id;type:BIGINT(20);AUTO_INCREMENT;NOT NULL;comment:id;primary_key;comment:主键ID;"`
	UUID        string    `gorm:"type:varchar(36);not null;column:uuid;uniqueIndex:uuid;comment:UUID;"`
	Name        string    `gorm:"type:varchar(100);not null;column:name;comment:name"`
	CreateAt    time.Time `gorm:"type:timestamp;NOT NULL;DEFAULT:CURRENT_TIMESTAMP;column:create_at;comment:创建时间"`
	UpdateAt    time.Time `gorm:"type:timestamp;NOT NULL;DEFAULT:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;column:update_at;comment:修改时间" json:"update_at"`
	Creator     string    `gorm:"type:varchar(36);not null;column:creator;comment:创建人id" aovalue:"creator"`
	Updater     string    `gorm:"type:varchar(36);not null;column:updater;comment:修改人id" aovalue:"updater"`
	Description string    `gorm:"size:255;not null;column:description;comment:description"`
	Master      string    `gorm:"size:36;not null;column:master;comment:负责人id"` // 负责人id
	Prefix      string    `gorm:"size:255;not null;column:prefix;comment:分区前缀"`
}

func (o *Organization) TableName() string {
	return "organization"
}

func (o *Organization) IdValue() int64 {
	return o.Id
}

type Partition struct {
	Id         int64     `gorm:"type:BIGINT(20);size:20;not null;auto_increment;primary_key;column:id;comment:主键ID;"`
	Oid        string    `gorm:"size:36;not null;column:oid;comment:组织id;uniqueIndex:oid_pid;"`
	Pid        string    `gorm:"size:36;not null;column:pid;comment:分区id;uniqueIndex:oid_pid;"`
	CreateTime time.Time `gorm:"type:timestamp;NOT NULL;DEFAULT:CURRENT_TIMESTAMP;column:create_at;comment:创建时间"`
}

func (p *Partition) TableName() string {
	return "organization_partition"
}

func (p *Partition) IdValue() int64 {
	return p.Id
}
