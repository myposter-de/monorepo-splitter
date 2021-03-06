package action

import (
	"fmt"
	"splitter/pkg"
	"strings"
)

type SetPackagesDependencies struct{}

func (s SetPackagesDependencies) Act(collection *pkg.PackageCollection) {
	versionString := collection.Conf.Semver.CaretedMinorVersion()
	for _, singlePkg := range collection.Packages {
		for name := range singlePkg.Composer.Items.Require {
			if _, ok := collection.Packages[name]; ok {
				singlePkg.Composer.Items.Require[name] = versionString
			} else if currentVersion, ok := collection.RootPackage.Composer.Items.Require[name]; ok {
				singlePkg.Composer.Items.Require[name] = currentVersion
			} else if strings.HasPrefix(name, "ext-") || strings.HasPrefix(name, "php") {
				continue
			} else {
				panic(fmt.Sprintf("package %s not found locally or in root", name))
			}
		}
	}
}

func (s SetPackagesDependencies) Description() string {
	return "set versions of mutual dependencies to current version"
}

func (s SetPackagesDependencies) String() string {
	return "set-packages-dependencies"
}
