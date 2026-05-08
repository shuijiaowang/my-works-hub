package util

import (
	"reflect"
	"strings"
	"workhub/model"
)

type TableColumnSchema struct {
	Field   string `json:"field"`
	GoType  string `json:"goType"`
	GormTag string `json:"gormTag"`
	Comment string `json:"comment"`
}

// ProjectTableSchema returns a model-based "table schema" view for Project.
// It is derived from Go struct field types and gorm tag comments.
func ProjectTableSchema() []TableColumnSchema {
	t := reflect.TypeOf(model.Project{})
	out := make([]TableColumnSchema, 0, t.NumField())

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)

		// Skip embedded gorm.Model fields but keep Project's own fields.
		if f.Anonymous {
			continue
		}

		gormTag := strings.TrimSpace(f.Tag.Get("gorm"))
		out = append(out, TableColumnSchema{
			Field:   f.Name,
			GoType:  f.Type.String(),
			GormTag: gormTag,
			Comment: extractGormComment(gormTag),
		})
	}
	return out
}

func extractGormComment(tag string) string {
	// Example: "type:varchar(255);not null;comment:'标题/项目名'"
	parts := strings.Split(tag, ";")
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if strings.HasPrefix(p, "comment:") {
			v := strings.TrimPrefix(p, "comment:")
			v = strings.TrimSpace(v)
			v = strings.Trim(v, `"'`)
			return v
		}
	}
	return ""
}
