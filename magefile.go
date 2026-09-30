//go:build mage
// +build mage

package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"

	devtools "github.com/elastic/beats/v7/dev-tools/mage"
	"github.com/elastic/beats/v7/dev-tools/mage/target/build"
	"github.com/elastic/beats/v7/dev-tools/mage/target/common"
	"github.com/elastic/beats/v7/dev-tools/mage/target/pkg"
	"github.com/elastic/beats/v7/dev-tools/mage/target/unittest"
)

func init() {
	// dev-tools locates the beats sources with `go list -m`, which returns
	// nothing in vendor mode, and vendor/ lacks the non-Go files it needs
	// (fields.yml, config templates). Point it at the module cache instead.
	out, err := sh.Output("go", "mod", "download", "-json", "github.com/elastic/beats/v7")
	if err != nil {
		panic(err)
	}
	var mod struct{ Dir string }
	if err := json.Unmarshal([]byte(out), &mod); err != nil {
		panic(err)
	}
	devtools.SetElasticBeatsDir(mod.Dir)

	devtools.SetBuildVariableSources(devtools.DefaultBeatBuildVariableSources)

	devtools.BeatDescription = "Countbeat periodically publishes an event carrying an incrementing counter."
	devtools.BeatVendor = "Peter Cesnek"
	devtools.BeatProjectType = devtools.CommunityProject
}

// VendorUpdate updates the vendor dir
func VendorUpdate() error {
	return sh.Run("go", "mod", "vendor")
}

// Package packages the Beat for distribution.
// Use SNAPSHOT=true to build snapshots.
// Use PLATFORMS to control the target platforms.
func Package() {
	start := time.Now()
	defer func() { fmt.Println("package ran for", time.Since(start)) }()

	devtools.MustUsePackaging("community_beat", "dev-tools/packaging/packages.yml")

	mg.Deps(Update)
	mg.Deps(build.CrossBuild)
	mg.SerialDeps(devtools.Package, pkg.PackageTest)
}

// Update updates the generated files (aka make update).
func Update() {
	mg.SerialDeps(Fields, FieldsGo, Config, FieldDocs)
}

// Fields generates a fields.yml for the Beat.
func Fields() error {
	return devtools.GenerateFieldsYAML()
}

// FieldsGo generates include/fields.go, which embeds fields.yml in the binary.
func FieldsGo() error {
	return devtools.GenerateAllInOneFieldsGo()
}

// FieldDocs generates docs/fields.asciidoc from fields.yml.
func FieldDocs() error {
	return devtools.Docs.FieldDocs("fields.yml")
}

// Config generates both the short/reference/docker configs.
func Config() error {
	p := devtools.DefaultConfigFileParams()
	p.Templates = append(p.Templates, devtools.OSSBeatDir("_meta/config/*.tmpl"))
	return devtools.Config(devtools.AllConfigTypes, p, ".")
}

// Clean cleans all generated files and build artifacts.
func Clean() error {
	return devtools.Clean()
}

// Check formats code, updates generated content, check for common errors, and
// checks for any modified files.
func Check() {
	common.Check()
}

// Fmt formats source code (.go and .py) and adds license headers.
func Fmt() {
	common.Fmt()
}

// Test runs all available tests
func Test() {
	mg.Deps(unittest.GoUnitTest)
}

// Build builds the Beat binary.
func Build() error {
	return build.Build()
}

// CrossBuild cross-builds the beat for all target platforms.
func CrossBuild() error {
	return build.CrossBuild()
}

// GolangCrossBuild build the Beat binary inside of the golang-builder.
// Do not use directly, use crossBuild instead.
func GolangCrossBuild() error {
	return build.GolangCrossBuild()
}
