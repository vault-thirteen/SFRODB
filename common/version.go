package common

import (
	"fmt"
	"runtime"
)

const (
	LibName    = "SFRODB"
	LibVersion = "0.13.0"

	ProductServer     = "Server"
	ProductTestClient = "Test Client"

	ComponentCache = "Cache"
)

const (
	StartupText       = "%s %s, ver. %s, %s."
	ComponentInfoText = "%s Version: %s."
)

func ShowIntroText(product string) {
	fmt.Println(
		fmt.Sprintf(StartupText, LibName, product, LibVersion, runtime.Version()),
	)
}

func ShowComponentInfoText(componentName string, componentVersion string) {
	fmt.Println(fmt.Sprintf(ComponentInfoText, componentName, componentVersion))
}
