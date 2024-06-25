package apinto

import "time"

type PartitionConfig struct {
	Id         int64     `gorm:"column:id;type:BIGINT(20);AUTO_INCREMENT;NOT NULL;comment:id;primary_key;comment:主键ID;"`
	Partition  string    `gorm:"type:varchar(36);not null;column:partition;comment:partition id;uniqueIndex:partition_key"`
	Key        string    `gorm:"type:varchar(36);not null;column:key;comment:key; uniqueIndex:partition_key"`
	Data       string    `gorm:"type:text;not null;column:data;comment:data;"`
	UpdateTime time.Time `gorm:"type:timestamp;NOT NULL;DEFAULT:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;column:update_at;comment:修改时间"`
}
