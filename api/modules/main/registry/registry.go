package registry

import (
	"log/slog"
	"sort"
	"sync"

	"github.com/khiemnd777/noah_api/modules/main/config"
	"github.com/khiemnd777/noah_framework/shared/metadata/customfields"
	"github.com/khiemnd777/noah_framework/shared/module"
	frameworkhttp "github.com/khiemnd777/noah_framework/pkg/http"
)

type Registrar interface {
	ID() string
	Priority() int // smaller = earlier
	Register(r frameworkhttp.Router, deps *module.ModuleDeps[config.ModuleConfig], cfMgr *customfields.Manager) error
}

var (
	mu         sync.RWMutex
	registrars = make(map[string]Registrar)
)

func Register(reg Registrar) {
	mu.Lock()
	defer mu.Unlock()
	if _, ok := registrars[reg.ID()]; ok {
		slog.Warn("Duplicate feature registration", "id", reg.ID())
	}
	registrars[reg.ID()] = reg
}

type InitOptions struct {
	EnabledIDs []string // if empty => enable all
}

func Init(r frameworkhttp.Router, deps *module.ModuleDeps[config.ModuleConfig], cfMgr *customfields.Manager, opts InitOptions) {
	mu.RLock()
	defer mu.RUnlock()

	var list []Registrar
	if len(opts.EnabledIDs) == 0 {
		for _, rg := range registrars {
			list = append(list, rg)
		}
	} else {
		enabled := map[string]struct{}{}
		for _, id := range opts.EnabledIDs {
			enabled[id] = struct{}{}
		}
		for id, rg := range registrars {
			if _, ok := enabled[id]; ok {
				list = append(list, rg)
			}
		}
	}

	sort.Slice(list, func(i, j int) bool {
		pi, pj := list[i].Priority(), list[j].Priority()
		if pi == 0 {
			pi = 100
		}
		if pj == 0 {
			pj = 100
		}
		return pi < pj
	})

	for _, rg := range list {
		if err := rg.Register(r, deps, cfMgr); err != nil {
			slog.Error("Feature registration failed", "id", rg.ID(), "err", err)
		} else {
			slog.Info("Feature registered", "id", rg.ID())
		}
	}
}
