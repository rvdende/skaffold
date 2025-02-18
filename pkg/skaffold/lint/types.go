/*
Copyright 2021 The Skaffold Authors

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

package lint

import (
	"fmt"

	"go.lsp.dev/protocol"
	"sigs.k8s.io/kustomize/kyaml/yaml"

	"github.com/GoogleContainerTools/skaffold/pkg/skaffold/parser"
)

// Options holds flag values for the various `skaffold lint` commands
type Options struct {
	// Filename is the `skaffold.yaml` file path
	Filename string
	// RepoCacheDir is the directory for the remote git repository cache
	RepoCacheDir string
	// OutFormat is the output format. One of: json
	OutFormat string
	// Modules is the module filter for specific commands
	Modules []string
	// Profiles is the slice of profile names to activate.
	Profiles []string
}

type Rule struct {
	RuleID               RuleID
	RuleType             RuleType
	ExplanationTemplate  string
	Severity             protocol.DiagnosticSeverity
	Filter               interface{}
	ExplanationPopulator func(InputParams) (explanationInfo, error) `json:"-"`
	LintConditions       []func(InputParams) bool                   `json:"-"`
}

type explanationInfo struct {
	FieldMap map[string]interface{}
	// FieldMap map[string]string
}

type Result struct {
	Rule        *Rule
	AbsFilePath string
	RelFilePath string
	Explanation string
	Line        int
	Column      int
}

type YamlFieldFilter struct {
	Filter      yaml.Filter
	FieldOnly   string
	InvertMatch bool
}

type ConfigFile struct {
	AbsPath string
	RelPath string
	Text    string
}

type RuleType int

const (
	RegExpLintLintRule RuleType = iota
	YamlFieldLintRule
)

func (a RuleType) String() string {
	return [...]string{"RegExpLintLintRule", "YamlFieldLintRule"}[a]
}

type RuleID int

const (
	DummyRuleIDForTesting RuleID = iota

	SkaffoldYamlAPIVersionOutOfDate
	SkaffoldYamlUseStaticPort
	SkaffoldYamlSyncPython
)

func (a RuleID) String() string {
	return fmt.Sprintf("ID%06d", a)
}

type InputParams struct {
	ConfigFile         ConfigFile
	DockerfileToDepMap map[string][]string
	SkaffoldConfig     *parser.SkaffoldConfigEntry
}

type Linter interface {
	Lint(InputParams, *[]Rule) (*[]Result, error)
}
