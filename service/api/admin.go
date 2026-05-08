package api

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
	"workhub/config"
	"workhub/db"
	"workhub/middleware"
	"workhub/model"
	"workhub/util"
	"workhub/util/response"

	"github.com/gin-gonic/gin"
)

type AdminApi struct{}

type AdminLoginRequest struct {
	Password string `json:"password"`
}

// Login 管理端登录：
// - 首次请求通常走 /admin/login?token=xxx（你会手动补充 token 后缀）
// - body 只做密码校验
// - 通过后返回 ok（具体会话逻辑后续再加）
func (a *AdminApi) Login(c *gin.Context) {
	var req AdminLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	expectedPwd := strings.TrimSpace(config.AppConfig.Admin.Password)
	if expectedPwd == "" {
		response.FailWithMessage("admin password 未配置", c)
		return
	}
	if strings.TrimSpace(req.Password) != expectedPwd {
		response.NoAuth("密码错误", c)
		return
	}

	// 同时要求携带固定 token（query/header 均可）
	if middleware.GetAdminToken(c) == "" {
		response.NoAuth("缺少 admin token", c)
		return
	}
	// token 正确性由 AdminAuth 统一处理，这里直接复用一次校验，避免漏配
	expectedToken := strings.TrimSpace(config.AppConfig.Admin.Token)
	if middleware.GetAdminToken(c) != expectedToken {
		response.NoAuth("admin token 校验失败", c)
		return
	}

	response.OkWithMessage("登录成功", c)
}

func (a *AdminApi) Ping(c *gin.Context) {
	response.OkWithMessage("admin pong", c)
}

type AdminCreateProjectRequest struct {
	Name        string     `json:"name"`
	FolderName  string     `json:"folderName"`
	IsPublic    *bool      `json:"isPublic"`
	Intro       string     `json:"intro"`
	CodeStartAt *time.Time `json:"codeStartAt"`
	GitRepo     string     `json:"gitRepo"`
	Guide       string     `json:"guide"`
	Tags        string     `json:"tags"`
}

type AdminUpdateProjectRequest struct {
	Name     *string `json:"name"`
	IsPublic *bool   `json:"isPublic"`
	Intro    *string `json:"intro"`
	GitRepo  *string `json:"gitRepo"`
	Guide    *string `json:"guide"`
	Tags     *string `json:"tags"`
}

func (a *AdminApi) CreateProject(c *gin.Context) {
	var req AdminCreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		response.FailWithMessage("name 不能为空", c)
		return
	}

	folder := strings.TrimSpace(req.FolderName)
	if folder == "" {
		folder = slugify(name)
		if folder == "" {
			folder = fmt.Sprintf("project-%d", time.Now().Unix())
		}
	}

	isPublic := true
	if req.IsPublic != nil {
		isPublic = *req.IsPublic
	}

	codeStartAt := time.Now()
	if req.CodeStartAt != nil && !req.CodeStartAt.IsZero() {
		codeStartAt = *req.CodeStartAt
	}

	p := model.Project{
		Name:        name,
		FolderName:  folder,
		IsPublic:    isPublic,
		Intro:       strings.TrimSpace(req.Intro),
		CodeStartAt: codeStartAt,
		GitRepo:     strings.TrimSpace(req.GitRepo),
		Guide:       req.Guide,
		Tags:        strings.TrimSpace(req.Tags),
	}

	if err := db.DB.Create(&p).Error; err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// After DB insert succeeds, ensure filesystem dirs exist:
	// ./resources/<folderName>/{media,zip}
	if err := util.EnsureProjectDirs("./resources", p.FolderName); err != nil {
		// Best-effort rollback to keep DB and filesystem consistent.
		_ = db.DB.Delete(&p).Error
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(p, "创建成功", c)
}

func (a *AdminApi) UpdateProject(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		response.FailWithMessage("invalid project id", c)
		return
	}

	var req AdminUpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	var proj model.Project
	if err := db.DB.First(&proj, id).Error; err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	updates := map[string]interface{}{}
	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if name == "" {
			response.FailWithMessage("name 不能为空", c)
			return
		}
		updates["name"] = name
	}
	if req.IsPublic != nil {
		updates["is_public"] = *req.IsPublic
	}
	if req.Intro != nil {
		updates["intro"] = strings.TrimSpace(*req.Intro)
	}
	if req.GitRepo != nil {
		updates["git_repo"] = strings.TrimSpace(*req.GitRepo)
	}
	if req.Guide != nil {
		updates["guide"] = *req.Guide
	}
	if req.Tags != nil {
		updates["tags"] = strings.TrimSpace(*req.Tags)
	}

	if len(updates) == 0 {
		response.OkWithDetailed(proj, "无更新", c)
		return
	}

	if err := db.DB.Model(&proj).Updates(updates).Error; err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := db.DB.First(&proj, id).Error; err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithDetailed(proj, "更新成功", c)
}

var nonWord = regexp.MustCompile(`[^a-zA-Z0-9_-]+`)

func slugify(s string) string {
	s = strings.TrimSpace(strings.ToLower(s))
	s = strings.ReplaceAll(s, " ", "-")
	s = nonWord.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	return s
}

type AdminMediaItem struct {
	ID           string         `json:"id"`
	Kind         util.MediaKind `json:"kind"`
	Order        int            `json:"order"`
	Url          string         `json:"url"`
	FileName     string         `json:"fileName"`
	OriginalName string         `json:"originalName"`
}

func (a *AdminApi) ListProjectMedia(c *gin.Context) {
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

	out := make([]AdminMediaItem, 0, len(items))
	for _, it := range items {
		out = append(out, AdminMediaItem{
			ID:           it.ID,
			Kind:         it.Kind,
			Order:        it.Order,
			Url:          filepath.ToSlash(filepath.Join("/api/resources", proj.FolderName, "media", it.FileName)),
			FileName:     it.FileName,
			OriginalName: it.OriginalName,
		})
	}
	response.OkWithData(gin.H{"items": out}, c)
}

func (a *AdminApi) UploadProjectMedia(c *gin.Context) {
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

	fh, err := c.FormFile("file")
	if err != nil {
		response.FailWithMessage("缺少上传文件字段 file", c)
		return
	}

	mediaDir := util.ProjectMediaDir("./resources", proj.FolderName)
	if err := os.MkdirAll(mediaDir, 0o755); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	origName := util.SanitizeFileBaseName(fh.Filename)
	if origName == "" {
		origName = "upload"
	}

	mediaId, fileName := util.NewMediaFileName(fh.Filename)
	dst := filepath.Join(mediaDir, fileName)

	if err := c.SaveUploadedFile(fh, dst); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	items, err := ensureAndLoadMediaIndex(mediaDir)
	if err != nil {
		_ = os.Remove(dst)
		response.FailWithMessage(err.Error(), c)
		return
	}

	newItem := util.MediaItem{
		ID:           mediaId,
		FileName:     fileName,
		OriginalName: fh.Filename,
		Kind:         util.DetectMediaKindByName(fileName),
		Order:        len(items),
		CreatedAt:    time.Now(),
	}
	items = append(items, newItem)
	items = util.NormalizeAndSortMedia(items)

	if err := util.SaveMediaIndex(mediaDir, items); err != nil {
		_ = os.Remove(dst)
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithData(AdminMediaItem{
		ID:           newItem.ID,
		Kind:         newItem.Kind,
		Order:        newItem.Order,
		Url:          filepath.ToSlash(filepath.Join("/api/resources", proj.FolderName, "media", newItem.FileName)),
		FileName:     newItem.FileName,
		OriginalName: newItem.OriginalName,
	}, c)
}

func (a *AdminApi) DeleteProjectMedia(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		response.FailWithMessage("invalid project id", c)
		return
	}
	mediaId := strings.TrimSpace(c.Param("mediaId"))
	if mediaId == "" {
		response.FailWithMessage("invalid media id", c)
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
		response.FailWithMessage(err.Error(), c)
		return
	}

	var target *util.MediaItem
	next := make([]util.MediaItem, 0, len(items))
	for _, it := range items {
		if it.ID == mediaId {
			tmp := it
			target = &tmp
			continue
		}
		next = append(next, it)
	}
	if target == nil {
		response.FailWithMessage("media not found", c)
		return
	}

	_ = os.Remove(filepath.Join(mediaDir, target.FileName))
	next = util.NormalizeAndSortMedia(next)
	if err := util.SaveMediaIndex(mediaDir, next); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

type AdminMoveMediaRequest struct {
	Direction string `json:"direction"` // left|right
}

func (a *AdminApi) MoveProjectMedia(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		response.FailWithMessage("invalid project id", c)
		return
	}
	mediaId := strings.TrimSpace(c.Param("mediaId"))
	if mediaId == "" {
		response.FailWithMessage("invalid media id", c)
		return
	}

	var req AdminMoveMediaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	dir := strings.ToLower(strings.TrimSpace(req.Direction))
	if dir != "left" && dir != "right" {
		response.FailWithMessage("direction must be left|right", c)
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
		response.FailWithMessage(err.Error(), c)
		return
	}
	items = util.NormalizeAndSortMedia(items)

	pos := -1
	for i := range items {
		if items[i].ID == mediaId {
			pos = i
			break
		}
	}
	if pos < 0 {
		response.FailWithMessage("media not found", c)
		return
	}

	if dir == "left" {
		if pos == 0 {
			// already first, no-op
		} else {
			items[pos-1], items[pos] = items[pos], items[pos-1]
		}
	} else {
		if pos == len(items)-1 {
			// already last, no-op
		} else {
			items[pos+1], items[pos] = items[pos], items[pos+1]
		}
	}

	// IMPORTANT: swap slice positions alone is not enough because NormalizeAndSortMedia
	// sorts by the persisted Order field. Re-assign Order based on the new slice order first.
	for i := range items {
		items[i].Order = i
	}
	items = util.NormalizeAndSortMedia(items)
	if err := util.SaveMediaIndex(mediaDir, items); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	out := make([]AdminMediaItem, 0, len(items))
	for _, it := range items {
		out = append(out, AdminMediaItem{
			ID:           it.ID,
			Kind:         it.Kind,
			Order:        it.Order,
			Url:          filepath.ToSlash(filepath.Join("/api/resources", proj.FolderName, "media", it.FileName)),
			FileName:     it.FileName,
			OriginalName: it.OriginalName,
		})
	}
	response.OkWithData(gin.H{"items": out}, c)
}

func sanitizeZipUploadName(original string) (string, error) {
	name := strings.TrimSpace(original)
	if name == "" {
		return "", fmt.Errorf("zip file name is empty")
	}
	name = filepath.Base(name)
	// extra safety: disallow any separators after Base()
	if strings.ContainsAny(name, `/\`) || strings.Contains(name, "..") {
		return "", fmt.Errorf("invalid zip file name")
	}
	if strings.ToLower(filepath.Ext(name)) != ".zip" {
		return "", fmt.Errorf("only .zip is allowed")
	}
	// keep original base name but sanitize strange chars (avoid creating weird paths)
	base := strings.TrimSuffix(name, filepath.Ext(name))
	base = util.SanitizeFileBaseName(base)
	if base == "" {
		base = "package"
	}
	return base + ".zip", nil
}

// UploadProjectZip uploads (or overwrites) a zip package under:
// ./resources/<folderName>/zip/<originalName>.zip
func (a *AdminApi) UploadProjectZip(c *gin.Context) {
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

	fh, err := c.FormFile("file")
	if err != nil {
		response.FailWithMessage("缺少上传文件字段 file", c)
		return
	}

	saveName, err := sanitizeZipUploadName(fh.Filename)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	zipDir := util.ProjectZipDir("./resources", proj.FolderName)
	if err := os.MkdirAll(zipDir, 0o755); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	dst := filepath.Join(zipDir, saveName)
	// SaveUploadedFile will overwrite if exists (same name => update)
	if err := c.SaveUploadedFile(fh, dst); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithData(gin.H{
		"fileName":     saveName,
		"originalName": fh.Filename,
		"url":          filepath.ToSlash(filepath.Join("/api/admin/projects", strconv.Itoa(id), "zip", saveName)),
	}, c)
}

func (a *AdminApi) DownloadProjectZip(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		response.FailWithMessage("invalid project id", c)
		return
	}

	fileName, err := sanitizeZipUploadName(c.Param("fileName"))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	var proj model.Project
	if err := db.DB.First(&proj, id).Error; err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	zipDir := util.ProjectZipDir("./resources", proj.FolderName)
	p := filepath.Join(zipDir, fileName)
	if _, err := os.Stat(p); err != nil {
		if os.IsNotExist(err) {
			response.FailWithMessage("zip not found", c)
			return
		}
		response.FailWithMessage(err.Error(), c)
		return
	}

	c.FileAttachment(p, fileName)
}

func (a *AdminApi) DeleteProjectZip(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		response.FailWithMessage("invalid project id", c)
		return
	}

	fileName, err := sanitizeZipUploadName(c.Param("fileName"))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	var proj model.Project
	if err := db.DB.First(&proj, id).Error; err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	zipDir := util.ProjectZipDir("./resources", proj.FolderName)
	p := filepath.Join(zipDir, fileName)
	if err := os.Remove(p); err != nil {
		if os.IsNotExist(err) {
			response.FailWithMessage("zip not found", c)
			return
		}
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithMessage("删除成功", c)
}
