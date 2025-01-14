/*
Copyright 2015 The Kubernetes Authors.

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
/*
Copyright 2019, 2021 The Multi-Cluster App Dispatcher Authors.

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
// deepcopy-gen is a tool for auto-generating DeepCopy functions.
//
// Given a list of input directories, it will generate functions that
// efficiently perform a full deep-copy of each type.  For any type that
// offers a `.DeepCopy()` method, it will simply call that.  Otherwise it will
// use standard value assignment whenever possible.  If that is not possible it
// will try to call its own generated copy function for the type, if the type is
// within the allowed root packages.  Failing that, it will fall back on
// `conversion.Cloner.DeepCopy(val)` to make the copy.  The resulting file will
// be stored in the same directory as the processed source package.
//
// Generation is governed by comment tags in the source.  Any package may
// request DeepCopy generation by including a comment in the file-comments of
// one file, of the form:
//   // +k8s:deepcopy-gen=package
//
// Packages can request that the generated DeepCopy functions be registered
// with an `init()` function call to `Scheme.AddGeneratedDeepCopyFuncs()` by
// changing the tag to:
//   // +k8s:deepcopy-gen=package,register
//
// DeepCopy functions can be generated for individual types, rather than the
// entire package by specifying a comment on the type definition of the form:
//   // +k8s:deepcopy-gen=true
//
// When generating for a whole package, individual types may opt out of
// DeepCopy generation by specifying a comment on the of the form:
//   // +k8s:deepcopy-gen=false
//
// Note that registration is a whole-package option, and is not available for
// individual types.
package main

import (
	"path/filepath"

	"k8s.io/gengo/args"
	"k8s.io/gengo/examples/deepcopy-gen/generators"

	"github.com/spf13/pflag"
	"k8s.io/klog/v2"
)

func main() {
	klog.InitFlags(nil)
	arguments := args.Default()

	// Override defaults.
	arguments.OutputFileBaseName = "deepcopy_generated"
	current_loc, _ := filepath.Abs(".")
	arguments.GoHeaderFilePath = filepath.Join(current_loc, "hack/boilerplate/boilerplate.go.txt")

	// Custom args.
	customArgs := &generators.CustomArgs{}
	pflag.CommandLine.StringSliceVar(&customArgs.BoundingDirs, "bounding-dirs", customArgs.BoundingDirs,
		"Comma-separated list of import paths which bound the types for which deep-copies will be generated.")
	arguments.CustomArgs = customArgs

	// Run it.
	if err := arguments.Execute(
		generators.NameSystems(),
		generators.DefaultNameSystem(),
		generators.Packages,
	); err != nil {
		klog.Fatalf("Error: %v", err)
	}
	klog.V(2).Info("Completed successfully.")
}
