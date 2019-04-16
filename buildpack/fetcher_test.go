package buildpack_test

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/fatih/color"
	"github.com/onsi/gomega/ghttp"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"

	"github.com/buildpack/pack/buildpack"
	h "github.com/buildpack/pack/testhelpers"
)

func TestBuildpackFetcher(t *testing.T) {
	h.RequireDocker(t)
	color.NoColor = true
	if runtime.GOOS == "windows" {
		t.Skip("create builder is not implemented on windows")
	}
	spec.Run(t, "BuildpackFetcher", testBuildpackFetcher, spec.Parallel(), spec.Report(report.Terminal{}))
}

type emptyLogger struct {
}

func (e *emptyLogger) Verbose(format string, a ...interface{}) {
}

func testBuildpackFetcher(t *testing.T, when spec.G, it spec.S) {
	when("#FetchBuildpack", func() {
		var (
			err      error
			tmpDir   string
			cacheDir string
			subject  *buildpack.Fetcher
		)

		it.Before(func() {
			tmpDir, err = ioutil.TempDir("", "")
			h.AssertNil(t, err)

			cacheDir, err = ioutil.TempDir("", "")
			h.AssertNil(t, err)

			subject = buildpack.NewFetcher(&emptyLogger{}, cacheDir)
		})

		it.After(func() {
			os.RemoveAll(tmpDir)
			os.RemoveAll(cacheDir)
		})

		it("fetches from a relative directory", func() {
			out, err := subject.FetchBuildpack(filepath.Join("testdata", "buildpack"))
			h.AssertNil(t, err)
			h.AssertEq(t, out.ID, "bp.one")
			h.AssertEq(t, out.Version, "some-buildpack-version")
			h.AssertEq(t, out.Stacks[0].ID, "some.stack.id")
			h.AssertEq(t, out.Stacks[1].ID, "other.stack.id")
			h.AssertNotEq(t, out.Dir, "")
			h.AssertDirContainsFileWithContents(t, out.Dir, "bin/detect", "I come from a directory\n")
			h.AssertDirContainsFileWithContents(t, out.Dir, "bin/build", "I come from a directory\n")
		})

		it("fetches from a relative tgz", func() {
			out, err := subject.FetchBuildpack(filepath.Join("testdata", "buildpack.tgz"))
			h.AssertNil(t, err)
			h.AssertEq(t, out.ID, "bp.one")
			h.AssertEq(t, out.Version, "some-buildpack-version")
			h.AssertEq(t, out.Stacks[0].ID, "some.stack.id")
			h.AssertEq(t, out.Stacks[1].ID, "other.stack.id")
			h.AssertNotEq(t, out.Dir, "")
			h.AssertDirContainsFileWithContents(t, out.Dir, "bin/detect", "I come from an archive\n")
			h.AssertDirContainsFileWithContents(t, out.Dir, "bin/build", "I come from an archive\n")
		})

		it("fetches from an absolute directory", func() {
			absPath, err := filepath.Abs(filepath.Join("testdata", "buildpack"))
			h.AssertNil(t, err)

			out, err := subject.FetchBuildpack(absPath)
			h.AssertNil(t, err)
			h.AssertEq(t, out.ID, "bp.one")
			h.AssertEq(t, out.Version, "some-buildpack-version")
			h.AssertEq(t, out.Stacks[0].ID, "some.stack.id")
			h.AssertEq(t, out.Stacks[1].ID, "other.stack.id")
			h.AssertNotEq(t, out.Dir, "")
			h.AssertDirContainsFileWithContents(t, out.Dir, "bin/detect", "I come from a directory\n")
			h.AssertDirContainsFileWithContents(t, out.Dir, "bin/build", "I come from a directory\n")
		})

		it("fetches from an absolute tgz", func() {
			absPath, err := filepath.Abs(filepath.Join("testdata", "buildpack.tgz"))
			h.AssertNil(t, err)

			out, err := subject.FetchBuildpack(absPath)
			h.AssertNil(t, err)
			h.AssertEq(t, out.ID, "bp.one")
			h.AssertEq(t, out.Version, "some-buildpack-version")
			h.AssertEq(t, out.Stacks[0].ID, "some.stack.id")
			h.AssertEq(t, out.Stacks[1].ID, "other.stack.id")
			h.AssertNotEq(t, out.Dir, "")
			h.AssertDirContainsFileWithContents(t, out.Dir, "bin/detect", "I come from an archive\n")
			h.AssertDirContainsFileWithContents(t, out.Dir, "bin/build", "I come from an archive\n")
		})

		it("fetches from a 'file://' URI directory", func() {
			absPath, err := filepath.Abs(filepath.Join("testdata", "buildpack"))
			h.AssertNil(t, err)

			out, err := subject.FetchBuildpack("file://" + absPath)
			h.AssertNil(t, err)
			h.AssertEq(t, out.ID, "bp.one")
			h.AssertEq(t, out.Version, "some-buildpack-version")
			h.AssertEq(t, out.Stacks[0].ID, "some.stack.id")
			h.AssertEq(t, out.Stacks[1].ID, "other.stack.id")
			h.AssertNotEq(t, out.Dir, "")
			h.AssertDirContainsFileWithContents(t, out.Dir, "bin/detect", "I come from a directory\n")
			h.AssertDirContainsFileWithContents(t, out.Dir, "bin/build", "I come from a directory\n")
		})

		it("fetches from a 'file://' URI tgz", func() {
			absPath, err := filepath.Abs(filepath.Join("testdata", "buildpack.tgz"))
			h.AssertNil(t, err)

			out, err := subject.FetchBuildpack("file://" + absPath)
			h.AssertNil(t, err)
			h.AssertEq(t, out.ID, "bp.one")
			h.AssertEq(t, out.Version, "some-buildpack-version")
			h.AssertEq(t, out.Stacks[0].ID, "some.stack.id")
			h.AssertEq(t, out.Stacks[1].ID, "other.stack.id")
			h.AssertNotEq(t, out.Dir, "")
			h.AssertDirContainsFileWithContents(t, out.Dir, "bin/detect", "I come from an archive\n")
			h.AssertDirContainsFileWithContents(t, out.Dir, "bin/build", "I come from an archive\n")
		})

		it("fetches from a 'http(s)://' URI tgz", func() {
			server := ghttp.NewServer()
			server.AppendHandlers(func(w http.ResponseWriter, r *http.Request) {
				path := filepath.Join("testdata", r.URL.Path)
				http.ServeFile(w, r, path)
			})
			defer server.Close()

			out, err := subject.FetchBuildpack(server.URL() + "/buildpack.tgz")
			h.AssertNil(t, err)
			h.AssertEq(t, out.ID, "bp.one")
			h.AssertEq(t, out.Version, "some-buildpack-version")
			h.AssertEq(t, out.Stacks[0].ID, "some.stack.id")
			h.AssertEq(t, out.Stacks[1].ID, "other.stack.id")
			h.AssertNotEq(t, out.Dir, "")
			h.AssertDirContainsFileWithContents(t, out.Dir, "bin/detect", "I come from an archive\n")
			h.AssertDirContainsFileWithContents(t, out.Dir, "bin/build", "I come from an archive\n")
		})
	})
}
