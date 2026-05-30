package model

import (
	"time"

	"gorm.io/gorm"
)

// ProjectLink 项目外链，如 GitHub / Gitee / 官网 / 插件下载等。
type ProjectLink struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// Project 作品/项目（用于主页展示与下载聚合）。
// 部分字段均可为空
type Project struct {
	gorm.Model
	Name          string        `gorm:"type:varchar(255);not null;comment:'标题/项目名'"`
	IsPublic      bool          `gorm:"default:true;comment:'是否公共（默认公开）'"`
	Intro         string        `gorm:"type:varchar(512);comment:'简介'"`
	CodeStartAt   time.Time     `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP;comment:'项目开始写代码的时间'"`
	Links         []ProjectLink `gorm:"serializer:json;type:json;comment:'外链列表'"`
	Guide         string        `gorm:"type:text;comment:'文字教程'"`
	Tags          string        `gorm:"type:varchar(255);comment:'标签（如：网站,浏览器脚本）'"` //逗号分隔
	FolderName    string        `gorm:"type:varchar(255);not null;comment:'项目文件夹名'"`
	ViewCount     int           `gorm:"default:0;comment:'观看次数'"`
	DownloadCount int           `gorm:"default:0;comment:'下载次数'"`
}
