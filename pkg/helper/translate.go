package helper

import (
	"fmt"
	"strings"

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

func EnvToHumanString(env map[string]string) string {
	if len(env) == 0 {
		return blank
	}
	envStr := ""
	for k, v := range env {
		envStr += k + "=" + v + " "
	}
	return envStr
}

func MachineContainerHealthcheckSpecToHumanString(hc *corev1pb.MachineContainerHealthcheckSpec) string {
	if hc == nil {
		return blank
	}

	parts := []string{}

	if hc.InitialDelaySeconds != 0 {
		parts = append(parts, fmt.Sprintf("Initial Delay: %ds", hc.InitialDelaySeconds))
	}

	if hc.PeriodSeconds != 0 {
		parts = append(parts, fmt.Sprintf("Period: %ds", hc.PeriodSeconds))
	}

	if http := hc.GetHttp(); http != nil {
		if http.Method != "" {
			parts = append(parts, fmt.Sprintf("Method: %s", http.Method))
		}

		if http.Port != 0 {
			parts = append(parts, fmt.Sprintf("Port: %d", http.Port))
		}

		if http.Path != "" {
			parts = append(parts, fmt.Sprintf("Path: %s", http.Path))
		}

		if len(http.Headers) > 0 {
			headerParts := []string{}
			for k, v := range http.Headers {
				headerParts = append(headerParts, fmt.Sprintf("%s=%s", k, v))
			}
			parts = append(parts, fmt.Sprintf("Headers: %s", strings.Join(headerParts, ", ")))
		}
	}

	if len(parts) == 0 {
		return blank
	}

	return strings.Join(parts, ", ")
}
