/*
Copyright 2021 TriggerMesh Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
)

const (
	cmdPath       = "cmd"
	configPath    = "config"
	templatesPath = "scaffolding"
	adapterPath   = "pkg/adapter"
)

type component struct {
	Kind          string
	LowercaseKind string
	UppercaseKind string
}

func main() {
	kind := flag.String("kind", "TestAdapter", "Specify the Adapter kind")
	cfgDir := flag.String("config", "../../",
		"Path of the directory containing the TriggerMesh deployment manifests")
	flag.Parse()
	*cfgDir = path.Clean(*cfgDir)
	temp := &component{
		Kind:          *kind,
		LowercaseKind: strings.ToLower(*kind),
		UppercaseKind: strings.ToUpper(*kind),
	}

	// make cmd directory
	cmddir := filepath.Join(*cfgDir, temp.LowercaseKind, cmdPath)
	mustMkdirAll(cmddir)

	// make adapter directory
	adpaterdir := filepath.Join(*cfgDir, temp.LowercaseKind, "pkg", "adapter")
	mustMkdirAll(adpaterdir)

	// populate cmd directory
	// read main.go and replace the template variables
	if err := temp.replaceTemplates(
		filepath.Join(templatesPath, cmdPath, "newtarget-adapter", "main.go"),
		filepath.Join(*cfgDir, temp.LowercaseKind, cmdPath+"/main.go"),
	); err != nil {
		log.Fatalf("failed creating the cmd templates: %v", err)
	}

	// populate adapter directory
	// read adapter.go
	if err := temp.replaceTemplates(
		filepath.Join(templatesPath, adapterPath, "adapter.go"),
		filepath.Join(*cfgDir, temp.LowercaseKind, adapterPath, "/adapter.go"),
	); err != nil {
		log.Fatalf("failed creating the adapter templates: %v", err)
	}

	fmt.Println("done")
}

func (a *component) replaceTemplates(filename, outputname string) error {
	tmp, err := template.ParseFiles(filename)
	if err != nil {
		return err
	}

	file, err := os.Create(outputname)
	if err != nil {
		return err
	}
	defer file.Close()

	return tmp.Execute(file, a)
}

func mustMkdirAll(path string) {
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		log.Fatalf("failed creating directory: %v", err)
	}
}
