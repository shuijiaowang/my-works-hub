package api

import (
	"net/http"
	"path/filepath"
	"strconv"
	"workhub/db"
	"workhub/model"
	"workhub/util"
	"workhub/util/response"

	"github.com/gin-gonic/gin"
)

type ProjectApi struct{}

type ProjectAllItem struct {
	Project  model.Project `json:"project"`
	CoverUrl string        `json:"coverUrl"`
}

func coverUrlFromMediaIndex(proj model.Project) string {
	mediaDir := util.ProjectMediaDir("./resources", proj.FolderName)
	items, err := util.LoadMediaIndex(mediaDir)
	if err != nil || len(items) == 0 {
		return ""
	}

	// "index1" 优先，其次取最小 order；优先图片。
	var best *util.MediaItem
	for i := range items {
		it := &items[i]
		if best == nil {
			best = it
			continue
		}
		if it.Order == 1 && best.Order != 1 {
			best = it
			continue
		}
		if best.Order == 1 && it.Order != 1 {
			continue
		}
		if it.Order < best.Order {
			best = it
			continue
		}
		if it.Order == best.Order && best.Kind != util.MediaKindImage && it.Kind == util.MediaKindImage {
			best = it
			continue
		}
	}
	if best == nil || best.FileName == "" {
		return ""
	}
	return filepath.ToSlash(filepath.Join("/api/resources", proj.FolderName, "media", best.FileName))
}

// All returns ALL projects, plus:
// - coverUrl: project media cover (prefer index1)
//
// Both normal users and admins can call it (no auth guard).
func (p *ProjectApi) All(c *gin.Context) {
	var projects []model.Project
	if err := db.DB.Order("id desc").Find(&projects).Error; err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	items := make([]ProjectAllItem, 0, len(projects))
	for _, proj := range projects {
		items = append(items, ProjectAllItem{
			Project:  proj,
			CoverUrl: coverUrlFromMediaIndex(proj),
		})
	}

	response.OkWithData(gin.H{
		"projects": items,
	}, c)
}

// Detail returns ONE project by id.
// It does not require auth.
func (p *ProjectApi) Detail(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		response.FailWithMessage("invalid project id", c)
		return
	}

	var proj model.Project
	if err := db.DB.First(&proj, id).Error; err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	item, err := buildProjectItem(proj)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Response{Code: response.ERROR, Data: nil, Msg: err.Error()})
		return
	}

	response.OkWithData(gin.H{
		"project": item,
	}, c)
}

type ProjectDetailItem struct {
	Project    model.Project `json:"project"`
	MediaDir   string        `json:"mediaDir"`
	ZipDir     string        `json:"zipDir"`
	MediaFiles []string      `json:"mediaFiles"`
	ZipFiles   []string      `json:"zipFiles"`
}

func buildProjectItem(proj model.Project) (ProjectDetailItem, error) {
	mediaDir := util.ProjectMediaDir("./resources", proj.FolderName)
	zipDir := util.ProjectZipDir("./resources", proj.FolderName)

	mediaFiles, err := util.ListFiles(mediaDir)
	if err != nil {
		return ProjectDetailItem{}, err
	}
	zipFiles, err := util.ListFiles(zipDir)
	if err != nil {
		return ProjectDetailItem{}, err
	}

	return ProjectDetailItem{
		Project:    proj,
		MediaDir:   mediaDir,
		ZipDir:     zipDir,
		MediaFiles: mediaFiles,
		ZipFiles:   zipFiles,
	}, nil
}

type ProjectMediaItem struct {
	ID       string         `json:"id"`
	Kind     util.MediaKind `json:"kind"`
	Order    int            `json:"order"`
	Url      string         `json:"url"`
	FileName string         `json:"fileName"`
}

// MediaList returns ordered media items for preview.
// It does not require auth.
func (p *ProjectApi) MediaList(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		response.FailWithMessage("invalid project id", c)
		return
	}

	var proj model.Project
	if err := db.DB.First(&proj, id).Error; err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	mediaDir := util.ProjectMediaDir("./resources", proj.FolderName)
	items, err := ensureAndLoadMediaIndex(mediaDir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Response{
			Code: response.ERROR,
			Data: nil,
			Msg:  err.Error(),
		})
		return
	}

	out := make([]ProjectMediaItem, 0, len(items))
	for _, it := range items {
		out = append(out, ProjectMediaItem{
			ID:       it.ID,
			Kind:     it.Kind,
			Order:    it.Order,
			Url:      filepath.ToSlash(filepath.Join("/api/resources", proj.FolderName, "media", it.FileName)),
			FileName: it.FileName,
		})
	}

	response.OkWithData(gin.H{
		"items": out,
	}, c)
}
