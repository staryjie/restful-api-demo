package version_test

import (
	"fmt"
	"testing"

	"github.com/staryjie/restful-api-demo/version"
)

func TestVersion(t *testing.T) {
	fmt.Println(version.FullVersion())
}
