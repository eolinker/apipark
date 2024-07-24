package service

import "time"

type Service struct {
	Id          int64     `gorm:"column:id;type:BIGINT(20);AUTO_INCREMENT;NOT NULL;comment:id;primary_key;comment:主键ID;"`
	UUID        string    `gorm:"type:varchar(36);not null;column:uuid;uniqueIndex:uuid;comment:UUID;"`
	Name        string    `gorm:"type:varchar(100);not null;column:name;comment:name"`
	Description string    `gorm:"type:varchar(255);not null;column:description;comment:描述"`
	Logo        string    `gorm:"type:text;column:logo;comment:logo"`
	ServiceType string    `gorm:"type:varchar(36);not null;column:service_type;comment:服务类型"`
	Project     string    `gorm:"type:varchar(36);not null;column:project;comment:项目"`
	Team        string    `gorm:"type:varchar(36);not null;column:team;comment:团队"`
	Catalogue   string    `gorm:"type:varchar(36);not null;column:catalogue;comment:目录"`
	Status      string    `gorm:"type:varchar(36);not null;column:status;comment:状态"`
	Tag         string    `gorm:"type:varchar(255);not null;column:tag;comment:标签"`
	CreateAt    time.Time `gorm:"type:timestamp;NOT NULL;DEFAULT:CURRENT_TIMESTAMP;column:create_at;comment:创建时间"`
	UpdateAt    time.Time `gorm:"type:timestamp;NOT NULL;DEFAULT:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;column:update_at;comment:修改时间" json:"update_at"`
	Creator     string    `gorm:"type:varchar(36);not null;column:creator;comment:创建者" aovalue:"creator"`
	Updater     string    `gorm:"type:varchar(36);not null;column:updater;comment:修改者" aovalue:"updater"`
	IsDelete    int       `gorm:"type:tinyint(1);not null;column:is_delete;comment:是否删除 0:未删除 1:已删除"`
}

func (o *Service) TableName() string {
	return "service"
}

func (o *Service) IdValue() int64 {
	return o.Id
}

type Partition struct {
	Id  int64  `gorm:"type:BIGINT(20);size:20;not null;auto_increment;primary_key;column:id;comment:主键ID;"`
	Sid string `gorm:"size:36;not null;column:sid;comment:服务id;uniqueIndex:sid_pid; index:sid;"`
	Pid string `gorm:"size:36;not null;column:pid;comment:分区id;uniqueIndex:sid_pid;index:pid;"`
}

func (p *Partition) IdValue() int64 {
	return p.Id
}
func (p *Partition) TableName() string {
	return "server_partition"
}

type Tag struct {
	Id  int64  `gorm:"column:id;type:BIGINT(20);AUTO_INCREMENT;NOT NULL;comment:id;primary_key;comment:主键ID;"`
	Tid string `gorm:"size:36;not null;column:tid;comment:标签id;uniqueIndex:sid_tid;index:tid;"`
	Sid string `gorm:"size:36;not null;column:sid;comment:服务id;uniqueIndex:sid_tid;index:sid;"`
}

func (t *Tag) IdValue() int64 {
	return t.Id
}

func (t *Tag) TableName() string {
	return "server_tag"
}

type Doc struct {
	Id       int64     `gorm:"column:id;type:BIGINT(20);AUTO_INCREMENT;NOT NULL;comment:id;primary_key;comment:主键ID;"`
	Sid      string    `gorm:"size:36;not null;column:sid;comment:服务id;uniqueIndex:unique_sid;"`
	Doc      string    `gorm:"type:text;column:content;comment:内容"`
	CreateAt time.Time `gorm:"type:timestamp;NOT NULL;DEFAULT:CURRENT_TIMESTAMP;column:create_at;comment:创建时间"`
	UpdateAt time.Time `gorm:"type:timestamp;NOT NULL;DEFAULT:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;column:update_at;comment:修改时间" json:"update_at"`
	Creator  string    `gorm:"type:varchar(36);not null;column:creator;comment:创建者"`
	Updater  string    `gorm:"type:varchar(36);not null;column:updater;comment:修改者"`
}

func (d *Doc) IdValue() int64 {
	return d.Id
}

func (d *Doc) TableName() string {
	return "server_doc"
}

type Api struct {
	Id   int64  `gorm:"type:BIGINT(20);size:20;not null;auto_increment;primary_key;column:id;comment:主键ID;"`
	Sid  string `gorm:"size:36;not null;column:sid;uniqueIndex:sid_api;index:sid;comment:服务id;"`
	Aid  string `gorm:"size:36;not null;column:aid;uniqueIndex:sid_api;comment:api id;index:api;"`
	Sort int    `gorm:"type:int;not null;column:sort;comment:排序"`
}

func (a *Api) IdValue() int64 {
	return a.Id
}
func (a *Api) TableName() string {
	return "service_api"
}
