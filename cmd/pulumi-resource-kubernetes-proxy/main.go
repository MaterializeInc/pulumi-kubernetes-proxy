// Copyright 2016-2020, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	pbempty "github.com/golang/protobuf/ptypes/empty"
	"github.com/pulumi/pulumi/pkg/v3/resource/provider"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/cmdutil"
	rpc "github.com/pulumi/pulumi/sdk/v3/proto/go"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"
)

const version = "0.1.1"

func main() {
	err := provider.Main("kubernetes-proxy", func(host *provider.HostClient) (rpc.ResourceProviderServer, error) {
		return &kubernetesProxyProvider{
			host: host,
		}, nil
	})
	if err != nil {
		cmdutil.ExitError(err.Error())
	}
}

type kubernetesProxyProvider struct {
	host *provider.HostClient
}

func (k *kubernetesProxyProvider) CheckConfig(ctx context.Context, req *rpc.CheckRequest) (*rpc.CheckResponse, error) {
	return &rpc.CheckResponse{Inputs: req.GetNews()}, nil
}

func (k *kubernetesProxyProvider) DiffConfig(ctx context.Context, req *rpc.DiffRequest) (*rpc.DiffResponse, error) {
	return &rpc.DiffResponse{}, nil
}

func (k *kubernetesProxyProvider) Configure(ctx context.Context, req *rpc.ConfigureRequest) (*rpc.ConfigureResponse, error) {
	vars := req.GetVariables()
	kubeconfig := vars["kubernetes-proxy:config:kubeconfig"]
	namespace := vars["kubernetes-proxy:config:namespace"]
	podSelector := vars["kubernetes-proxy:config:podSelector"]
	hostPort := vars["kubernetes-proxy:config:hostPort"]
	remotePort := vars["kubernetes-proxy:config:remotePort"]

	config, err := clientcmd.RESTConfigFromKubeConfig([]byte(kubeconfig))
	if err != nil {
		return nil, fmt.Errorf("unable to load kubeconfig: %v", err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("unable to construct kubernetes client: %v", err)
	}
	pods, err := clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{
		LabelSelector: podSelector,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to list pods: %v", err)
	}
	if len(pods.Items) == 0 {
		return nil, fmt.Errorf("no pods matching selector %s", podSelector)
	}
	pod := pods.Items[0]
	transport, upgrader, err := spdy.RoundTripperFor(config)
	if err != nil {
		return nil, fmt.Errorf("unable to construct transport: %v", err)
	}
	dialer := spdy.NewDialer(
		upgrader,
		&http.Client{Transport: transport},
		"POST",
		clientset.CoreV1().RESTClient().Post().Resource("pods").Namespace(namespace).Name(pod.Name).SubResource("portforward").URL(),
	)
	ports := []string{fmt.Sprintf("%s:%s", hostPort, remotePort)}
	stopCh := make(chan struct{}, 1)
	readyCh := make(chan struct{})
	fw, err := portforward.New(dialer, ports, stopCh, readyCh, os.Stdout, os.Stderr)
	if err != nil {
		return nil, err
	}
	go func() {
		if err := fw.ForwardPorts(); err != nil {
			fmt.Fprintf(os.Stderr, "failed forwarding ports: %v", err)
			os.Exit(1)
		}
	}()
	<-readyCh
	return &rpc.ConfigureResponse{}, nil
}

func (k *kubernetesProxyProvider) Invoke(_ context.Context, req *rpc.InvokeRequest) (*rpc.InvokeResponse, error) {
	tok := req.GetTok()
	return nil, fmt.Errorf("Unknown Invoke token '%s'", tok)
}

func (k *kubernetesProxyProvider) StreamInvoke(req *rpc.InvokeRequest, server rpc.ResourceProvider_StreamInvokeServer) error {
	tok := req.GetTok()
	return fmt.Errorf("Unknown StreamInvoke token '%s'", tok)
}

func (k *kubernetesProxyProvider) Check(ctx context.Context, req *rpc.CheckRequest) (*rpc.CheckResponse, error) {
	urn := resource.URN(req.GetUrn())
	return nil, fmt.Errorf("Unknown resource type '%s'", urn.Type())
}

func (k *kubernetesProxyProvider) Diff(ctx context.Context, req *rpc.DiffRequest) (*rpc.DiffResponse, error) {
	urn := resource.URN(req.GetUrn())
	return nil, fmt.Errorf("Unknown resource type '%s'", urn.Type())
}

func (k *kubernetesProxyProvider) Create(ctx context.Context, req *rpc.CreateRequest) (*rpc.CreateResponse, error) {
	urn := resource.URN(req.GetUrn())
	return nil, fmt.Errorf("Unknown resource type '%s'", urn.Type())
}

func (k *kubernetesProxyProvider) Read(ctx context.Context, req *rpc.ReadRequest) (*rpc.ReadResponse, error) {
	urn := resource.URN(req.GetUrn())
	return nil, fmt.Errorf("Unknown resource type '%s'", urn.Type())
}

func (k *kubernetesProxyProvider) Update(ctx context.Context, req *rpc.UpdateRequest) (*rpc.UpdateResponse, error) {
	urn := resource.URN(req.GetUrn())
	return nil, fmt.Errorf("Unknown resource type '%s'", urn.Type())
}

func (k *kubernetesProxyProvider) Delete(ctx context.Context, req *rpc.DeleteRequest) (*pbempty.Empty, error) {
	urn := resource.URN(req.GetUrn())
	return nil, fmt.Errorf("Unknown resource type '%s'", urn.Type())
}

func (k *kubernetesProxyProvider) Construct(_ context.Context, _ *rpc.ConstructRequest) (*rpc.ConstructResponse, error) {
	panic("Construct not implemented")
}

func (k *kubernetesProxyProvider) GetPluginInfo(context.Context, *pbempty.Empty) (*rpc.PluginInfo, error) {
	return &rpc.PluginInfo{
		Version: version,
	}, nil
}

func (k *kubernetesProxyProvider) GetSchema(ctx context.Context, req *rpc.GetSchemaRequest) (*rpc.GetSchemaResponse, error) {
	return &rpc.GetSchemaResponse{}, nil
}

func (k *kubernetesProxyProvider) Cancel(context.Context, *pbempty.Empty) (*pbempty.Empty, error) {
	return &pbempty.Empty{}, nil
}
