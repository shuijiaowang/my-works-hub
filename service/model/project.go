package model

import (
	"gorm.io/gorm"
)

// Project 作品/项目（用于主页展示与下载聚合）。
// 部分字段均可为空
type Project struct {
	gorm.Model
	Name          string `gorm:"type:varchar(255);not null;comment:'标题/项目名'"`
	GitRepo       string `gorm:"type:varchar(255);comment:'Git开源地址'"`
	Guide         string `gorm:"type:text;comment:'文字教程'"`
	Tags          string `gorm:"type:varchar(255);comment:'标签（如：网站,浏览器脚本）'"` //逗号分隔
	FolderName    string `gorm:"type:varchar(255);not null;comment:'项目文件夹名'"`
	ViewCount     int    `gorm:"default:0;comment:'观看次数'"`
	DownloadCount int    `gorm:"default:0;comment:'下载次数'"`
}
