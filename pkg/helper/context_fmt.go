package helper

import (
	"github.com/baepo-cloud/baepo-cli/pkg/config"
	"github.com/baepo-cloud/baepo-cli/pkg/iostream"
)

type ContextFmt struct {
	Name  string
	Value config.Context
}

func ContextFmtMapping() []any {
	return []any{
		iostream.FieldConfig{
			DisplayName: "Name",
			FormatFunc: func(obj *ContextFmt) string {
				return obj.Name
			},
		},
		iostream.FieldConfig{
			DisplayName: "Secret Key",
			FormatFunc: func(obj *ContextFmt) string {
				if obj.Value.SecretKey == nil {
					return blank
				}
				return *obj.Value.SecretKey
			},
		},
		iostream.FieldConfig{
			DisplayName: "Workspace ID",
			FormatFunc: func(obj *ContextFmt) string {
				if obj.Value.WorkspaceID == nil {
					return blank
				}
				return *obj.Value.WorkspaceID
			},
		},
		iostream.FieldConfig{
			DisplayName: "User ID",
			FormatFunc: func(obj *ContextFmt) string {
				if obj.Value.UserID == nil {
					return blank
				}
				return *obj.Value.UserID
			},
		},
		iostream.FieldConfig{
			DisplayName: "URL",
			FormatFunc: func(obj *ContextFmt) string {
				if obj.Value.URL == "" {
					return blank
				}
				return obj.Value.URL
			},
		},
	}
}
