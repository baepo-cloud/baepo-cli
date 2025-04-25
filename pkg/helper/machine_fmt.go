package helper

import (
	"fmt"

	"github.com/baepo-cloud/baepo-cli/pkg/iostream"
	apiv1pb "github.com/baepo-cloud/baepo-proto/go/baepo/api/v1"
	"github.com/dustin/go-humanize"
)

func MachineArrayFmt() iostream.ArrayOptions {
	return iostream.ArrayOptions{
		Fields: []iostream.FieldDefinition{
			{
				DisplayName: "ID",
				FormatFunc: func(obj any) string {
					if m, ok := obj.(*apiv1pb.Machine); ok {
						return m.Id
					}
					return blank
				},
			},
			{
				DisplayName: "Node ID",
				FormatFunc: func(obj any) string {
					if m, ok := obj.(*apiv1pb.Machine); ok {
						return m.GetNodeId()
					}
					return blank
				},
			},
			{
				DisplayName: "Workspace ID",
				FormatFunc: func(obj any) string {
					if m, ok := obj.(*apiv1pb.Machine); ok {
						return m.GetWorkspaceId()
					}
					return blank
				},
			},
			{
				DisplayName: "Name",
				FormatFunc: func(obj any) string {
					if m, ok := obj.(*apiv1pb.Machine); ok {
						if m.GetName() == "" {
							return blank
						}
						return m.GetName()
					}
					return blank
				},
			},
			{
				DisplayName: "State",
				FormatFunc: func(obj any) string {
					if m, ok := obj.(*apiv1pb.Machine); ok {
						return MachineStateToHumanString(m.GetState())
					}
					return blank
				},
			},
			{
				DisplayName: "Desired State",
				FormatFunc: func(obj any) string {
					if m, ok := obj.(*apiv1pb.Machine); ok {
						return MachineDesiredStateToHumanString(m.GetDesiredState())
					}
					return blank
				},
			},
			{
				DisplayName: "Started At",
				FormatFunc: func(obj any) string {
					if m, ok := obj.(*apiv1pb.Machine); ok {
						return TimestampToHumanString(m.GetStartedAt())
					}
					return blank
				},
			},
			{
				DisplayName: "Expires At",
				FormatFunc: func(obj any) string {
					if m, ok := obj.(*apiv1pb.Machine); ok {
						return TimestampToHumanString(m.GetExpiresAt())
					}
					return blank
				},
			},
			{
				DisplayName: "Terminated At",
				FormatFunc: func(obj any) string {
					if m, ok := obj.(*apiv1pb.Machine); ok {
						return TimestampToHumanString(m.GetTerminatedAt())
					}
					return blank
				},
			},
			{
				DisplayName: "Termination Cause",
				FormatFunc: func(obj any) string {
					if m, ok := obj.(*apiv1pb.Machine); ok {
						if m.TerminatedAt == nil {
							return blank
						}
						return MachineTerminationCauseToHumanString(m.GetTerminationCause())
					}
					return blank
				},
			},
			{
				DisplayName: "Termination Details",
				FormatFunc: func(obj any) string {
					if m, ok := obj.(*apiv1pb.Machine); ok {
						if m.TerminatedAt == nil {
							return blank
						}
						return m.GetTerminationDetails()
					}
					return blank
				},
			},
		},
	}
}

func MachineObjectFmt() iostream.ObjectOptions {
	return iostream.ObjectOptions{
		Fields: []iostream.FieldDefinition{
			{
				DisplayName: "ID",
				FormatFunc: func(obj any) string {
					if m, ok := obj.(*apiv1pb.Machine); ok {
						return m.Id
					}
					return blank
				},
			},
			{
				DisplayName: "Node ID",
				FormatFunc: func(obj any) string {
					if m, ok := obj.(*apiv1pb.Machine); ok {
						return m.GetNodeId()
					}
					return blank
				},
			},
			{
				DisplayName: "Workspace ID",
				FormatFunc: func(obj any) string {
					if m, ok := obj.(*apiv1pb.Machine); ok {
						return m.GetWorkspaceId()
					}
					return blank
				},
			},
			{
				DisplayName: "Name",
				FormatFunc: func(obj any) string {
					if m, ok := obj.(*apiv1pb.Machine); ok {
						if m.GetName() == "" {
							return blank
						}
						return m.GetName()
					}
					return blank
				},
			},
			{
				DisplayName: "Spec / Cores",
				FormatFunc: func(obj any) string {
					if m, ok := obj.(*apiv1pb.Machine); ok {
						return fmt.Sprintf("%d", m.GetSpec().GetCpus())
					}
					return blank
				},
			},
			{
				DisplayName: "Spec / RAM",
				FormatFunc: func(obj any) string {
					if m, ok := obj.(*apiv1pb.Machine); ok {
						return humanize.Bytes(m.GetSpec().GetMemoryMb() * 1024 * 1024)
					}
					return blank
				},
			},
			{
				DisplayName: "State",
				FormatFunc: func(obj any) string {
					if m, ok := obj.(*apiv1pb.Machine); ok {
						return MachineStateToHumanString(m.GetState())
					}
					return blank
				},
			},
			{
				DisplayName: "Desired State",
				FormatFunc: func(obj any) string {
					if m, ok := obj.(*apiv1pb.Machine); ok {
						return MachineDesiredStateToHumanString(m.GetDesiredState())
					}
					return blank
				},
			},
			{
				DisplayName: "Created At",
				FormatFunc: func(obj any) string {
					if m, ok := obj.(*apiv1pb.Machine); ok {
						return TimestampToHumanString(m.GetCreatedAt())
					}
					return blank
				},
			},
			{
				DisplayName: "Started At",
				FormatFunc: func(obj any) string {
					if m, ok := obj.(*apiv1pb.Machine); ok {
						return TimestampToHumanString(m.GetStartedAt())
					}
					return blank
				},
			},
			{
				DisplayName: "Expires At",
				FormatFunc: func(obj any) string {
					if m, ok := obj.(*apiv1pb.Machine); ok {
						return TimestampToHumanString(m.GetExpiresAt())
					}
					return blank
				},
			},
			{
				DisplayName: "Terminated At",
				FormatFunc: func(obj any) string {
					if m, ok := obj.(*apiv1pb.Machine); ok {
						return TimestampToHumanString(m.GetTerminatedAt())
					}
					return blank
				},
			},
			{
				DisplayName: "Termination Cause",
				FormatFunc: func(obj any) string {
					if m, ok := obj.(*apiv1pb.Machine); ok {
						if m.TerminatedAt == nil {
							return blank
						}
						return MachineTerminationCauseToHumanString(m.GetTerminationCause())
					}
					return blank
				},
			},
			{
				DisplayName: "Termination Details",
				FormatFunc: func(obj any) string {
					if m, ok := obj.(*apiv1pb.Machine); ok {
						if m.TerminatedAt == nil {
							return blank
						}
						return m.GetTerminationDetails()
					}
					return blank
				},
			},
		},
	}
}
