package helper

import (
	corev1pb "github.com/baepo-cloud/baepo-proto/go/baepo/core/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	blank = "-"
)

func MachineStateToHumanString(s corev1pb.MachineState) string {
	switch s {
	case corev1pb.MachineState_MachineState_Pending:
		return "Pending"
	case corev1pb.MachineState_MachineState_Starting:
		return "Starting"
	case corev1pb.MachineState_MachineState_Running:
		return "Running"
	case corev1pb.MachineState_MachineState_Degraded:
		return "Degraded"
	case corev1pb.MachineState_MachineState_Error:
		return "Error"
	case corev1pb.MachineState_MachineState_Terminating:
		return "Terminating"
	case corev1pb.MachineState_MachineState_Terminated:
		return "Terminated"
	default:
		return blank
	}
}

func MachineDesiredStateToHumanString(s corev1pb.MachineDesiredState) string {
	switch s {
	case corev1pb.MachineDesiredState_MachineDesiredState_Pending:
		return "Pending"
	case corev1pb.MachineDesiredState_MachineDesiredState_Running:
		return "Running"
	case corev1pb.MachineDesiredState_MachineDesiredState_Terminated:
		return "Terminated"
	default:
		return blank
	}
}

func TimestampToHumanString(at *timestamppb.Timestamp) string {
	if at == nil {
		return blank
	}
	return at.AsTime().Format("2006-01-02 15:04:05")
}

func MachineTerminationCauseToHumanString(cause corev1pb.MachineTerminationCause) string {
	switch cause {
	case corev1pb.MachineTerminationCause_MachineTerminationCause_HealthcheckFailed:
		return "Healthcheck Failed"
	case corev1pb.MachineTerminationCause_MachineTerminationCause_ManuallyRequested:
		return "Manually Requested"
	case corev1pb.MachineTerminationCause_MachineTerminationCause_InternalError:
		return "Internal Error"
	case corev1pb.MachineTerminationCause_MachineTerminationCause_NoNodeAvailable:
		return "No Node Available"
	case corev1pb.MachineTerminationCause_MachineTerminationCause_Expired:
		return "Expired"
	default:
		return blank
	}
}
