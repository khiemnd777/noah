package runtime

import (
	"os"
	"path/filepath"
	"testing"

	frameworkmodule "github.com/khiemnd777/noah_framework/pkg/module"
)

func TestDiscoverModulesSupportsMultipleRoots(t *testing.T) {
	tmpDir := t.TempDir()

	mustWriteModule(t, filepath.Join(tmpDir, "framework", "modules", "observability"))
	mustWriteModule(t, filepath.Join(tmpDir, "api", "modules", "auth"))
	mustWriteModule(t, filepath.Join(tmpDir, "api", "modules", "main"))
	mustWriteModule(t, filepath.Join(tmpDir, "api", "modules", "main", "department"))

	descriptors, err := DiscoverModules([]frameworkmodule.DiscoveryRoot{
		{Name: "framework", Path: filepath.Join(tmpDir, "framework", "modules")},
		{Name: "api-main", Path: filepath.Join(tmpDir, "api", "modules", "main")},
		{Name: "api", Path: filepath.Join(tmpDir, "api", "modules")},
	})
	if err != nil {
		t.Fatalf("discover modules: %v", err)
	}

	got := make(map[string]string, len(descriptors))
	for _, descriptor := range descriptors {
		got[descriptor.ID] = descriptor.RootName
	}

	want := map[string]string{
		"observability": "framework",
		"auth":          "api",
		"main":          "api-main",
		"department":    "api-main",
	}

	if len(got) != len(want) {
		t.Fatalf("unexpected module count: got=%d want=%d descriptors=%v", len(got), len(want), descriptors)
	}
	for id, rootName := range want {
		if got[id] != rootName {
			t.Fatalf("module %q root mismatch: got=%q want=%q", id, got[id], rootName)
		}
	}
}

func TestDiscoverModulesRejectsDuplicateIDsAcrossDifferentPaths(t *testing.T) {
	tmpDir := t.TempDir()

	mustWriteModule(t, filepath.Join(tmpDir, "framework", "modules", "auth"))
	mustWriteModule(t, filepath.Join(tmpDir, "api", "modules", "auth"))

	_, err := DiscoverModules([]frameworkmodule.DiscoveryRoot{
		{Name: "framework", Path: filepath.Join(tmpDir, "framework", "modules")},
		{Name: "api", Path: filepath.Join(tmpDir, "api", "modules")},
	})
	if err == nil {
		t.Fatal("expected duplicate module id error")
	}
}

func mustWriteModule(t *testing.T, dir string) {
	t.Helper()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", dir, err)
	}
	if err := os.WriteFile(filepath.Join(dir, "config.yaml"), []byte("server:\n  host: 127.0.0.1\n  port: 1\n  route: /x\nexternal: true\n"), 0o644); err != nil {
		t.Fatalf("write config %s: %v", dir, err)
	}
	if err := os.WriteFile(filepath.Join(dir, "main.go"), []byte("package main\nfunc main(){}\n"), 0o644); err != nil {
		t.Fatalf("write main %s: %v", dir, err)
	}
}
