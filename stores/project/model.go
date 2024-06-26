package project

import "time"

type Project struct {
	Id       int64     `gorm:"column:id;type:BIGINT(20);AUTO_INCREMENT;NOT NULL;comment:id;primary_key;comment:主键ID;"`
	UUID     string    `gorm:"type:varchar(36);not null;column:uuid;uniqueIndex:uuid;comment:UUID;"`
	Name     string    `gorm:"type:varchar(100);not null;column:name;comment:name"`
	CreateAt time.Time `gorm:"type:timestamp;NOT NULL;DEFAULT:CURRENT_TIMESTAMP;column:create_at;comment:创建时间"`
	UpdateAt time.Time `gorm:"type:timestamp;NOT NULL;DEFAULT:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;column:update_at;comment:修改时间"`

	Description  string `gorm:"size:255;not null;column:description;comment:description"`
	Prefix       string `gorm:"size:255;not null;column:prefix;comment:前缀"`
	Team         string `gorm:"size:36;not null;column:team;comment:团队id;index:team"`                 // 团队id
	Organization string `gorm:"size:36;not null;column:organization;comment:组织id;index:organization"` // 组织
	Master       string `gorm:"size:36;not null;column:master;comment:负责人id"`                         // 负责人id
	IsDelete     int    `gorm:"type:tinyint(1);not null;column:is_delete;comment:是否删除"`
	AsServer     bool   `gorm:"type:tinyint(1);not null;column:as_server;comment:是否为服务端项目"`
	AsApp        bool   `gorm:"type:tinyint(1);not null;column:as_app;comment:是否为应用项目"`
}

func (p *Project) IdValue() int64 {
	return p.Id
}
func (p *Project) TableName() string {
	return "project"
}

type MemberRole struct {
	Id         int64     `gorm:"type:BIGINT(20);size:20;not null;auto_increment;primary_key;column:id;comment:主键ID;"`
	Pid        string    `gorm:"size:36;not null;column:pid;index:pid;comment:项目id;uniqueIndex:pid_uid_rid;"`
	Uid        string    `gorm:"size:36;not null;column:uid;index:uid;comment:用户id;uniqueIndex:pid_uid_rid;"`
	Rid        string    `gorm:"size:36;not null;column:rid;index:rid;comment:角色id;uniqueIndex:pid_uid_rid;"`
	CreateTime time.Time `gorm:"type:timestamp;NOT NULL;DEFAULT:CURRENT_TIMESTAMP;column:create_at;comment:创建时间"`
}

func (m *MemberRole) IdValue() int64 {
	return m.Id
}
func (m *MemberRole) TableName() string {
	return "project_role"
}

type Authorization struct {
	Id             int64     `gorm:"type:BIGINT(20);size:20;not null;auto_increment;primary_key;column:id;comment:主键ID;"`
	UUID           string    `gorm:"size:36;not null;column:uuid;uniqueIndex:uuid;comment:UUID;"`
	Name           string    `gorm:"size:100;not null;column:name;comment:名称"`
	Project        string    `gorm:"size:36;not null;column:project;comment:项目id;index:project"`
	Type           string    `gorm:"size:100;not null;column:type;comment:类型"`
	Position       string    `gorm:"size:100;not null;column:position;comment:位置"`
	TokenName      string    `gorm:"size:100;not null;column:token_name;comment:token名称"`
	Config         string    `gorm:"type:text;not null;column:config;comment:配置"`
	Creator        string    `gorm:"size:36;not null;column:creator;comment:创建者" aovalue:"creator"`
	Updater        string    `gorm:"size:36;not null;column:updater;comment:修改者" aovalue:"updater"`
	ExpireTime     int64     `gorm:"type:BIGINT(20);not null;column:expire_time;comment:过期时间"`
	CreateAt       time.Time `gorm:"type:timestamp;NOT NULL;DEFAULT:CURRENT_TIMESTAMP;column:create_at;comment:创建时间"`
	UpdateAt       time.Time `gorm:"type:timestamp;NOT NULL;DEFAULT:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;column:update_at;comment:修改时间"`
	HideCredential bool      `gorm:"type:tinyint(1);not null;column:hide_credential;comment:隐藏凭证"`
}

func (a *Authorization) IdValue() int64 {
	return a.Id
}

func (a *Authorization) TableName() string {
	return "project_authorization"
}

type Partition struct {
	Id         int64     `gorm:"type:BIGINT(20);size:20;not null;auto_increment;primary_key;column:id;comment:主键ID;"`
	Project    string    `gorm:"size:36;not null;column:project;comment:系统id;uniqueIndex:oid_pid;"`
	Partition  string    `gorm:"size:36;not null;column:partition;comment:分区id;uniqueIndex:oid_pid;"`
	CreateTime time.Time `gorm:"type:timestamp;NOT NULL;DEFAULT:CURRENT_TIMESTAMP;column:create_at;comment:创建时间"`
}

func (p *Partition) TableName() string {
	return "project_partition"
}

func (p *Partition) IdValue() int64 {
	return p.Id
}
