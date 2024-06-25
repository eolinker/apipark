package partition

import "time"

type Partition struct {
	Id       int64     `gorm:"column:id;type:BIGINT(20);AUTO_INCREMENT;NOT NULL;comment:id;primary_key;comment:主键ID;"`
	UUID     string    `gorm:"type:varchar(36);not null;column:uuid;uniqueIndex:uuid;comment:UUID;"`
	Name     string    `gorm:"type:varchar(100);not null;column:name;comment:name"`
	CreateAt time.Time `gorm:"type:timestamp;NOT NULL;DEFAULT:CURRENT_TIMESTAMP;column:create_at;comment:创建时间"`
	UpdateAt time.Time `gorm:"type:timestamp;NOT NULL;DEFAULT:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;column:update_at;comment:修改时间" json:"update_at"`
	Creator  string    `gorm:"type:varchar(36);not null;column:creator;comment:创建人" aovalue:"creator"`
	Updater  string    `gorm:"type:varchar(36);not null;column:updater;comment:修改人" aovalue:"updater"`
	Resume   string    `gorm:"size:255;not null;column:resume;comment:备注"`
	Prefix   string    `gorm:"size:255;not null;column:prefix;comment:分区前缀"`
	Url      string    `gorm:"size:512;not null;column:url;comment:分区访问的地址,格式为 https://{host}/{prefix}"`
	Cluster  string    `gorm:"size:255;not null;column:cluster;comment:所属集群"`
	IsDelete bool      `gorm:"type:tinyint(1);not null;column:is_delete;comment:是否删除"`
}

func (p *Partition) IdValue() int64 {
	return p.Id
}
func (p *Partition) TableName() string {
	return "partition"
}
