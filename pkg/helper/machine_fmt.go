package helper

import (
	"fmt"
	"strings"

	"github.com/baepo-cloud/baepo-cli/pkg/iostream"
	apiv1pb "github.com/baepo-cloud/baepo-proto/go/baepo/api/v1"
	corev1pb "github.com/baepo-cloud/baepo-proto/go/baepo/core/v1"
)

// MachineArrayConfig returns the declarative configuration for mapping Machine arrays.
func MachineMapping() []any {
	return []any{
		iostream.FieldConfig{
			DisplayName: "ID",
			FormatFunc: func(obj *apiv1pb.Machine) string {
				return fmt.Sprintf("%s", obj.GetId())
			},
		},
		iostream.FieldConfig{
			DisplayName: "Node ID",
			FormatFunc: func(obj *apiv1pb.Machine) string {
				return obj.GetNodeId()
			},
		},
		iostream.FieldConfig{
			DisplayName: "Workspace ID",
			FormatFunc: func(obj *apiv1pb.Machine) string {
				return obj.GetWorkspaceId()
			},
		},
		iostream.FieldConfig{
			DisplayName: "Name",
			FormatFunc: func(obj *apiv1pb.Machine) string {
				if obj.GetName() == "" {
					return ""
				}
				return obj.GetName()
			},
		},
		iostream.ObjectConfig{
			Path:        "Spec",
			DisplayName: "Spec",
			Full:        true,
			Fields: []any{
				iostream.FieldConfig{
					DisplayName: "CPUs",
					FormatFunc: func(obj *corev1pb.MachineSpec) string {
						return fmt.Sprint(obj.Cpus)
					},
				},
				iostream.ArrayConfig{
					Path:        "Containers",
					DisplayName: "Containers",
					ObjectConfig: &iostream.ObjectConfig{
						Fields: []any{
							iostream.FieldConfig{
								DisplayName: "Image",
								FormatFunc: func(obj *corev1pb.MachineContainerSpec) string {
									return fmt.Sprint(obj.Image)
								},
							},
							iostream.FieldConfig{
								DisplayName: "Command",
								FormatFunc: func(obj *corev1pb.MachineContainerSpec) string {
									if len(obj.Command) == 0 {
										return "-"
									}
									return strings.Join(obj.Command, " ")
								},
							},
						},
					},
				},
			},
		},
		iostream.FieldConfig{
			DisplayName: "State",
			FormatFunc: func(obj *apiv1pb.Machine) string {
				return MachineStateToHumanString(obj.GetState())
			},
		},
		iostream.FieldConfig{
			DisplayName: "Desired State",
			FormatFunc: func(obj *apiv1pb.Machine) string {
				return MachineDesiredStateToHumanString(obj.GetDesiredState())
			},
		},
		iostream.FieldConfig{
			DisplayName: "Started At",
			FormatFunc: func(obj *apiv1pb.Machine) string {
				return TimestampToHumanString(obj.GetStartedAt())
			},
		},
		iostream.FieldConfig{
			DisplayName: "Expires At",
			FormatFunc: func(obj *apiv1pb.Machine) string {
				return TimestampToHumanString(obj.GetExpiresAt())
			},
		},
		iostream.FieldConfig{
			DisplayName: "Terminated At",
			FormatFunc: func(obj *apiv1pb.Machine) string {
				return TimestampToHumanString(obj.GetTerminatedAt())
			},
		},
		iostream.FieldConfig{
			DisplayName: "Termination Cause",
			FormatFunc: func(obj *apiv1pb.Machine) string {
				if obj.TerminatedAt == nil {
					return ""
				}
				return MachineTerminationCauseToHumanString(obj.GetTerminationCause())
			},
		},
		iostream.FieldConfig{
			DisplayName: "Termination Details",
			FormatFunc: func(obj *apiv1pb.Machine) string {
				if obj.TerminatedAt == nil {
					return ""
				}
				return obj.GetTerminationDetails()
			},
		},
	}
}
