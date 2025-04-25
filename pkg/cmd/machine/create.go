package machine

import (
	"encoding/json"
	"fmt"
	"strings"

	"connectrpc.com/connect"
	"github.com/baepo-cloud/baepo-cli/pkg/app"
	"github.com/baepo-cloud/baepo-cli/pkg/baepoerrors"
	"github.com/baepo-cloud/baepo-cli/pkg/helper"
	"github.com/baepo-cloud/baepo-cli/pkg/iostream"
	apiv1pb "github.com/baepo-cloud/baepo-proto/go/baepo/api/v1"
	corev1pb "github.com/baepo-cloud/baepo-proto/go/baepo/core/v1"
	"github.com/spf13/cobra"
)

func newCreateCmd() *cobra.Command {
	var name string
	var cpus uint32
	var memoryMB uint64
	var image string
	var env []string
	var healthPort int32
	var healthPath string
	var healthMethod string
	var containersJSON string
	var start bool

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create machine",
		Example: `
# Create a machine with a single container
baepo machine create --name myapp --cpus 2 --memory 2048 --image nginx:latest --env KEY1=value1 --env KEY2=value2 --start

# Create a machine with a container that has a healthcheck
baepo machine create --name myapp --cpus 2 --memory 2048 --image myapp:latest --health-port 8080 --health-path /health

# Create a machine with multiple containers using JSON
baepo machine create --name mydb --cpus 4 --memory 8192 --containers '[{"image":"postgres:14","env":{"POSTGRES_PASSWORD":"secret"}},{"image":"redis:alpine"}]' --start
		`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			a := app.FromContext(ctx)

			spec := &corev1pb.MachineSpec{
				Cpus:     cpus,
				MemoryMb: memoryMB,
			}

			if spec.Cpus == 0 {
				a.IOStream.Error("CPUs must be greater than 0")
				return baepoerrors.InvalidArgsError
			}

			if spec.MemoryMb == 0 {
				a.IOStream.Error("Memory must be greater than 0")
				return baepoerrors.InvalidArgsError
			}

			// Parse containers
			if containersJSON != "" {
				// Parse the JSON array of containers
				var containersData []map[string]interface{}
				if err := json.Unmarshal([]byte(containersJSON), &containersData); err != nil {
					a.IOStream.Error(fmt.Sprintf("Parsing containers JSON: %v", err))
					return baepoerrors.MachineError
				}

				// Convert to MachineContainerSpec
				for _, containerData := range containersData {
					container := &corev1pb.MachineContainerSpec{}

					// Image
					if img, ok := containerData["image"].(string); ok {
						container.Image = img
					} else {
						a.IOStream.Error("Container is missing 'image' field or it's not a string")
						return baepoerrors.MachineError
					}

					// Environment variables
					if envMap, ok := containerData["env"].(map[string]interface{}); ok {
						container.Env = make(map[string]string)
						for k, v := range envMap {
							if strVal, ok := v.(string); ok {
								container.Env[k] = strVal
							} else {
								container.Env[k] = fmt.Sprintf("%v", v)
							}
						}
					}

					// Command
					if cmd, ok := containerData["command"].([]interface{}); ok {
						container.Command = make([]string, 0, len(cmd))
						for _, c := range cmd {
							if strVal, ok := c.(string); ok {
								container.Command = append(container.Command, strVal)
							}
						}
					}

					// Healthcheck
					if healthData, ok := containerData["healthcheck"].(map[string]interface{}); ok {
						healthcheck := &corev1pb.MachineContainerHealthcheckSpec{}

						if initialDelay, ok := healthData["initial_delay_seconds"].(float64); ok {
							healthcheck.InitialDelaySeconds = int32(initialDelay)
						}

						if periodSeconds, ok := healthData["period_seconds"].(float64); ok {
							healthcheck.PeriodSeconds = int32(periodSeconds)
						}

						if httpData, ok := healthData["http"].(map[string]interface{}); ok {
							httpHealthcheck := &corev1pb.MachineContainerHealthcheckSpec_HttpHealthcheckSpec{}

							if method, ok := httpData["method"].(string); ok {
								httpHealthcheck.Method = method
							} else {
								httpHealthcheck.Method = "GET"
							}

							if path, ok := httpData["path"].(string); ok {
								httpHealthcheck.Path = path
							}

							if port, ok := httpData["port"].(float64); ok {
								httpHealthcheck.Port = int32(port)
							}

							if headers, ok := httpData["headers"].(map[string]interface{}); ok {
								httpHealthcheck.Headers = make(map[string]string)
								for k, v := range headers {
									if strVal, ok := v.(string); ok {
										httpHealthcheck.Headers[k] = strVal
									} else {
										httpHealthcheck.Headers[k] = fmt.Sprintf("%v", v)
									}
								}
							}

							healthcheck.Type = &corev1pb.MachineContainerHealthcheckSpec_Http{
								Http: httpHealthcheck,
							}
						}

						container.Healthcheck = healthcheck
					}

					spec.Containers = append(spec.Containers, container)
				}
			} else if image != "" {
				// Single container from command-line args
				container := &corev1pb.MachineContainerSpec{
					Image: image,
					Env:   make(map[string]string),
				}

				// Parse env variables
				for _, e := range env {
					parts := strings.SplitN(e, "=", 2)
					if len(parts) == 2 {
						container.Env[parts[0]] = parts[1]
					}
				}

				// Add healthcheck if specified
				if healthPort > 0 || healthPath != "" {
					httpHealthcheck := &corev1pb.MachineContainerHealthcheckSpec_HttpHealthcheckSpec{
						Method: healthMethod,
						Path:   healthPath,
						Port:   healthPort,
					}

					container.Healthcheck = &corev1pb.MachineContainerHealthcheckSpec{
						InitialDelaySeconds: 5,
						PeriodSeconds:       10,
						Type: &corev1pb.MachineContainerHealthcheckSpec_Http{
							Http: httpHealthcheck,
						},
					}
				}

				spec.Containers = append(spec.Containers, container)
			} else {
				a.IOStream.Error("Either --image or --containers must be specified")
				return baepoerrors.MachineError
			}

			// Create the machine
			req := connect.NewRequest(&apiv1pb.MachineCreateRequest{
				WorkspaceId: *a.Config.CurrentContext.WorkspaceID,
				Spec:        spec,
				Start:       start,
			})

			if name != "" {
				nameStr := name
				req.Msg.Name = &nameStr
			}

			res, err := a.MachineClient.Create(ctx, req)
			if err != nil {
				a.IOStream.Error(fmt.Sprintf("Creating machine: %v", err))
				return baepoerrors.MachineError
			}

			a.IOStream.Object(res.Msg.Machine, helper.MachineMapping(), iostream.ObjectOptions{Full: true})

			return nil
		},
	}

	// Add flags
	cmd.Flags().StringVar(&name, "name", "", "Name of the machine")
	cmd.Flags().Uint32Var(&cpus, "cpus", 1, "Number of CPUs")
	cmd.Flags().Uint64Var(&memoryMB, "memory", 1024, "Memory in MB")
	cmd.Flags().StringVar(&image, "image", "", "Container image")
	cmd.Flags().StringSliceVar(&env, "env", []string{}, "Environment variables in KEY=VALUE format")
	cmd.Flags().Int32Var(&healthPort, "health-port", 0, "Healthcheck port")
	cmd.Flags().StringVar(&healthPath, "health-path", "/", "Healthcheck path")
	cmd.Flags().StringVar(&healthMethod, "health-method", "GET", "Healthcheck method")
	cmd.Flags().StringVar(&containersJSON, "containers", "", "Container definitions in JSON format")
	cmd.Flags().BoolVar(&start, "start", false, "Start the machine after creation")

	return cmd
}
