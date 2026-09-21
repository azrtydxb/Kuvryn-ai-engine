package manifestcheck

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestEngineDockerfilesInstallRuntimeCommands(t *testing.T) {
	root := repoRoot(t)
	matches, err := filepath.Glob(filepath.Join(root, "engines", "*", "*", "Dockerfile"))
	if err != nil {
		t.Fatal(err)
	}
	if len(matches) == 0 {
		t.Fatal("no engine Dockerfiles found")
	}
	for _, path := range matches {
		path := path
		t.Run(filepath.ToSlash(strings.TrimPrefix(path, root+string(os.PathSeparator))), func(t *testing.T) {
			b, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			s := string(b)
			flavor := flavorFromDockerfilePath(path)
			switch {
			case strings.HasPrefix(flavor, "llama-cpp-"):
				requireContains(t, s, "install-llama-cpp.sh")
				requireContains(t, s, "kuvryn-llama-cpp-serve")
			case strings.HasPrefix(flavor, "vllm-"):
				requireContains(t, s, "install-vllm.sh")
				requireContains(t, s, "kuvryn-vllm-serve")
			case strings.HasPrefix(flavor, "sglang-"):
				requireContains(t, s, "install-sglang.sh")
				requireContains(t, s, "kuvryn-sglang-serve")
			case strings.HasPrefix(flavor, "tensorrt-llm-"):
				requireContains(t, s, "verify-tensorrt-llm.sh")
				requireContains(t, s, "kuvryn-trtllm-serve")
			default:
				t.Fatalf("unhandled flavor %q", flavor)
			}
			requireContains(t, s, "kuvryn-runtime-wrapper")
		})
	}
}

func repoRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
}

func flavorFromDockerfilePath(path string) string {
	vendor := filepath.Base(filepath.Dir(path))
	engine := filepath.Base(filepath.Dir(filepath.Dir(path)))
	return engine + "-" + vendor
}

func requireContains(t *testing.T, s, want string) {
	t.Helper()
	if !strings.Contains(s, want) {
		t.Fatalf("Dockerfile does not contain %q", want)
	}
}
