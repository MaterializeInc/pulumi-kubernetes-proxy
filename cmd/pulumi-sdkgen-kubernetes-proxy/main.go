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
	"io/ioutil"
	"os"
	"path"

	"github.com/pulumi/pulumi/sdk/v3/go/common/tools"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/contract"

	"github.com/pkg/errors"
	pygen "github.com/pulumi/pulumi/pkg/v3/codegen/python"
	pschema "github.com/pulumi/pulumi/pkg/v3/codegen/schema"
)

// emitPackage emits an entire package pack into the configured output directory
// with the configured settings.
func emitPackage(schemaPath, targetSdkFolder string) error {
	schemaBytes, err := ioutil.ReadFile(schemaPath)
	if err != nil {
		return errors.Wrap(err, "reading schema")
	}

	var spec pschema.PackageSpec
	err = json.Unmarshal(schemaBytes, &spec)
	if err != nil {
		return errors.Wrap(err, "reading schema")
	}

	ppkg, err := pschema.ImportSpec(spec, nil)
	if err != nil {
		return errors.Wrap(err, "reading schema")
	}

	toolDescription := "the Pulumi SDK Generator"
	extraFiles := map[string][]byte{}

	files, err := pygen.GeneratePackage(toolDescription, ppkg, extraFiles)
	if err != nil {
		return errors.Wrapf(err, "generating python package")
	}

	for f, contents := range files {
		if err := emitFile(path.Join(targetSdkFolder, "python"), f, contents); err != nil {
			return errors.Wrapf(err, "emitting file %v", f)
		}
	}

	return nil
}

func emitFile(outDir, relPath string, contents []byte) error {
	p := path.Join(outDir, relPath)
	if err := tools.EnsureDir(path.Dir(p)); err != nil {
		return errors.Wrap(err, "creating directory")
	}

	f, err := os.Create(p)
	if err != nil {
		return errors.Wrap(err, "creating file")
	}
	defer contract.IgnoreClose(f)

	_, err = f.Write(contents)
	return err
}

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("Usage: pulumi-sdkgen-kubernetes-proxy <schema-file> <target-sdk-folder>\n")
		return
	}

	schemaPath := os.Args[1]
	targetSdkFolder := os.Args[2]

	err := emitPackage(schemaPath, targetSdkFolder)
	if err != nil {
		fmt.Printf("Failed: %s", err.Error())
	}
}
