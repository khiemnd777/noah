package runtime

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"

	frameworkmodule "github.com/khiemnd777/noah_framework/pkg/module"
)

type ModuleLoader struct {
	roots []frameworkmodule.DiscoveryRoot
}

func NewModuleLoader(roots []frameworkmodule.DiscoveryRoot) *ModuleLoader {
	copied := make([]frameworkmodule.DiscoveryRoot, len(roots))
	copy(copied, roots)
	return &ModuleLoader{roots: copied}
}

func DiscoverModules(roots []frameworkmodule.DiscoveryRoot) ([]frameworkmodule.Descriptor, error) {
	return NewModuleLoader(roots).Discover()
}

func (l *ModuleLoader) Discover() ([]frameworkmodule.Descriptor, error) {
	descriptors := make([]frameworkmodule.Descriptor, 0)
	seenPaths := make(map[string]struct{})
	seenIDs := make(map[string]string)

	for _, root := range l.roots {
		rootPath, err := filepath.Abs(root.Path)
		if err != nil {
			return nil, fmt.Errorf("resolve root %q: %w", root.Path, err)
		}

		info, err := os.Stat(rootPath)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, fmt.Errorf("stat root %q: %w", rootPath, err)
		}
		if !info.IsDir() {
			return nil, fmt.Errorf("module root %q is not a directory", rootPath)
		}

		rootDescriptors, err := discoverRoot(root, rootPath)
		if err != nil {
			return nil, err
		}

		for _, descriptor := range rootDescriptors {
			if _, ok := seenPaths[descriptor.Path]; ok {
				continue
			}
			if existingPath, ok := seenIDs[descriptor.ID]; ok && existingPath != descriptor.Path {
				return nil, fmt.Errorf("duplicate module id %q found in %q and %q", descriptor.ID, existingPath, descriptor.Path)
			}
			seenPaths[descriptor.Path] = struct{}{}
			seenIDs[descriptor.ID] = descriptor.Path
			descriptors = append(descriptors, descriptor)
		}
	}

	return descriptors, nil
}

func discoverRoot(root frameworkmodule.DiscoveryRoot, rootPath string) ([]frameworkmodule.Descriptor, error) {
	descriptors := make([]frameworkmodule.Descriptor, 0)

	err := filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			return nil
		}

		descriptor, ok, err := describeModuleDir(root, rootPath, path)
		if err != nil {
			return err
		}
		if ok {
			descriptors = append(descriptors, descriptor)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walk root %q: %w", rootPath, err)
	}

	sort.Slice(descriptors, func(i, j int) bool {
		return descriptors[i].Path < descriptors[j].Path
	})

	return descriptors, nil
}

func describeModuleDir(root frameworkmodule.DiscoveryRoot, rootPath, dir string) (frameworkmodule.Descriptor, bool, error) {
	configPath := filepath.Join(dir, "config.yaml")
	if _, err := os.Stat(configPath); err != nil {
		if os.IsNotExist(err) {
			return frameworkmodule.Descriptor{}, false, nil
		}
		return frameworkmodule.Descriptor{}, false, err
	}

	entryPath := filepath.Join(dir, "main.go")
	if _, err := os.Stat(entryPath); err != nil {
		if os.IsNotExist(err) {
			return frameworkmodule.Descriptor{}, false, nil
		}
		return frameworkmodule.Descriptor{}, false, err
	}

	rel, err := filepath.Rel(rootPath, dir)
	if err != nil {
		return frameworkmodule.Descriptor{}, false, err
	}

	moduleID := filepath.ToSlash(rel)
	if moduleID == "." {
		moduleID = filepath.Base(dir)
	}

	return frameworkmodule.Descriptor{
		ID:         moduleID,
		Name:       filepath.Base(dir),
		RootName:   root.Name,
		RootPath:   rootPath,
		Path:       dir,
		ConfigPath: configPath,
		EntryPath:  entryPath,
	}, true, nil
}
