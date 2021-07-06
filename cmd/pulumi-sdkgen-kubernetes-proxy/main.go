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
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pulumi/pulumi/sdk/v3/go/common/tools"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/cmdutil"

	"github.com/pkg/errors"
	pygen "github.com/pulumi/pulumi/pkg/v3/codegen/python"
	"github.com/pulumi/pulumi/pkg/v3/codegen/schema"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: pulumi-sdkgen-kubernetes-proxy <version>\n")
		return
	}

	if err := run(os.Args[1]); err != nil {
		cmdutil.ExitError(err.Error())
	}
}

func run(version string) error {
	spec := schema.PackageSpec{
		Name:              "kubernetes-proxy",
		Version:           version,
		Description:       "A Pulumi package for managing Kubernetes proxies during provisioning.",
		License:           "Apache-2.0",
		Repository:        "https://github.com/benesch/pulumi-kubernetes-proxy",
		PluginDownloadURL: fmt.Sprintf("https://github.com/benesch/pulumi-kubernetes-proxy/releases/download/%s/", version),
		Provider: schema.ResourceSpec{
			InputProperties: map[string]schema.PropertySpec{
				"kubeconfig":  {TypeSpec: schema.TypeSpec{Type: "string"}},
				"namespace":   {TypeSpec: schema.TypeSpec{Type: "string"}},
				"podSelector": {TypeSpec: schema.TypeSpec{Type: "string"}},
				"remotePort":  {TypeSpec: schema.TypeSpec{Type: "integer"}},
				"hostPort":    {TypeSpec: schema.TypeSpec{Type: "integer"}},
			},
			RequiredInputs: []string{
				"kubeconfig",
				"namespace",
				"podSelector",
				"remotePort",
				"hostPort",
			},
		},
		Language: map[string]json.RawMessage{
			"python": json.RawMessage("{}"),
		},
	}
	ppkg, err := schema.ImportSpec(spec, nil)
	if err != nil {
		return errors.Wrap(err, "reading schema")
	}

	toolDescription := "the Pulumi SDK Generator"
	extraFiles := map[string][]byte{}
	files, err := pygen.GeneratePackage(toolDescription, ppkg, extraFiles)
	if err != nil {
		return fmt.Errorf("generating python package: %v", err)
	}

	for path, contents := range files {
		path = filepath.Join("sdk", "python", path)
		if err := tools.EnsureFileDir(path); err != nil {
			return fmt.Errorf("creating directory: %v", err)
		}
		if err := os.WriteFile(path, contents, 0644); err != nil {
			return fmt.Errorf("writing file: %v", err)
		}
	}

	return nil
}
