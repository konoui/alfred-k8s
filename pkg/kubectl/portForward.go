package kubectl

import (
	"context"
	"fmt"
	"strconv"
	"strings"
)

const (
	podPortProtoTPL = `{{ range .spec.containers -}}
{{ range .ports -}}
{{.containerPort}}/{{.protocol}}
{{ end -}}
{{ end -}}`
	deploymentPortProtoTPL = `{{ range .spec.template.spec.containers -}}
{{ range .ports -}}
{{.containerPort}}/{{.protocol}}
{{ end -}}
{{ end -}}`
	servicePortProtoTPL = `{{ range .spec.ports -}}
{{.port}}/{{.protocol}}
{{ end -}}`
)

// PortResource represents resource of ports
type PortResource struct {
	Resource  string   `json:"Resource"`
	Namespace string   `json:"Namespace"`
	Name      string   `json:"Name"`
	Ports     []string `json:"Ports"`
	Job
}

// PortForward exec kubectl port-forward command
func (k *Kubectl) PortForward(ctx context.Context, res, name, ns string, ports []string) (resp *PortResource, errCh <-chan error) {
	args := append([]string{
		"port-forward",
		res + "/" + name,
		"--namespace",
		ns,
	}, ports...)
	job, errCh := k.StartJob(ctx, args...)
	resp = &PortResource{
		Resource:  res,
		Namespace: ns,
		Name:      name,
		Ports:     ports,
		Job:       *job,
	}
	return
}

// GetPorts return pod/deployment/service ports
func (k *Kubectl) GetPorts(res, name, ns string) (ports []string) {
	switch res {
	case "pod", "po":
		ports = k.GetPodPorts(name, ns)
	case "deploy", "deployment":
		ports = k.GetDeploymentPorts(name, ns)
	case "svc", "service":
		ports = k.GetServicePorts(name, ns)
	}
	return
}

// GetPodPorts returns TCP ports defined in containers
func (k *Kubectl) GetPodPorts(name, ns string) []string {
	arg := fmt.Sprintf("get pod %s --namespace %s -o go-template --template='%s'", name, ns, podPortProtoTPL)
	resp, err := k.Execute(arg)
	if err != nil {
		return []string{}
	}
	return parsePortsFromResponse(resp)
}

// GetDeploymentPorts returns TCP ports defined in containers
func (k *Kubectl) GetDeploymentPorts(name, ns string) []string {
	arg := fmt.Sprintf("get deployment %s --namespace %s -o go-template --template='%s'", name, ns, deploymentPortProtoTPL)
	resp, err := k.Execute(arg)
	if err != nil {
		return []string{}
	}
	return parsePortsFromResponse(resp)
}

// GetServicePorts returns TCP ports defined in service
func (k *Kubectl) GetServicePorts(name, ns string) []string {
	arg := fmt.Sprintf("get service %s --namespace %s -o go-template --template='%s'", name, ns, servicePortProtoTPL)
	resp, err := k.Execute(arg)
	if err != nil {
		return []string{}
	}
	return parsePortsFromResponse(resp)
}

func parsePortsFromResponse(resp *Response) (ports []string) {
	for portProto := range resp.Readline() {
		port := parsePort(portProto)
		if port == "" {
			continue
		}
		ports = append(ports, port)
	}
	return
}

// parsePort assumes Port/Protocol format
func parsePort(portProto string) string {
	out := strings.Split(portProto, "/")
	if len(out) < 2 {
		return ""
	}
	if strings.EqualFold(out[1], "UDP") {
		return ""
	}

	port := out[0]
	if _, err := strconv.Atoi(out[0]); err != nil {
		return ""
	}
	return port
}
