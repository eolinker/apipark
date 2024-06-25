package external_app

import "time"

type ExternalApp struct {
	Id       int64     `gorm:"column:id;type:BIGINT(20);AUTO_INCREMENT;NOT NULL;comment:id;primary_key;comment:主键ID;"`
	UUID     string    `gorm:"type:varchar(36);not null;column:uuid;uniqueIndex:uuid;comment:UUID;"`
	Name     string    `gorm:"type:varchar(100);not null;column:name;comment:name"`
	Token    string    `gorm:"column:token;type:VARCHAR(36);NOT NULL;comment:鉴权token"`
	Desc     string    `gorm:"column:desc;type:TEXT;comment:应用描述;"`
	Tags     []string  `gorm:"column:tags;type:TEXT;comment:应用标签;serializer:json"`
	Enable   bool      `gorm:"column:enable;type:TINYINT(1); comment:启用状态"`
	IsDelete bool      `gorm:"column:is_delete;type:TINYINT(1); comment:是否删除"`
	Creator  string    `gorm:"type:varchar(36);not null;column:creator;comment:creator" aovalue:"creator"`
	Updater  string    `gorm:"type:varchar(36);not null;column:updater;comment:updater" aovalue:"updater"`
	CreateAt time.Time `gorm:"type:timestamp;NOT NULL;DEFAULT:CURRENT_TIMESTAMP;column:create_at;comment:创建时间"`
	UpdateAt time.Time `gorm:"type:timestamp;NOT NULL;DEFAULT:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;column:update_at;comment:修改时间" json:"update_at"`
}

func (c *ExternalApp) IdValue() int64 {
	return c.Id
}
func (c *ExternalApp) TableName() string {
	return "external_app"
}
