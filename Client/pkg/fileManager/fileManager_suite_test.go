package fileManager

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestFileManager(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "FileManager Suite")
}
