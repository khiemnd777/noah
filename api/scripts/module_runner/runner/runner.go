package runner

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/khiemnd777/noah_api/shared/runtime"
	"github.com/khiemnd777/noah_api/shared/utils"
)

const (
	moduleStartTimeout = 30 * time.Second
	namespaceGuidance  = "Run this command inside the API container: docker compose exec api go run ./scripts/module_runner <cmd> ..."
)

func EnsureRuntimeNamespaceForMutation() error {
	if isInsideContainer() {
		return nil
	}
	running, err := composeAPIContainerRunning()
	if err != nil {
		return err
	}
	if !running {
		return nil
	}

	return fmt.Errorf("runtime namespace mismatch: API container is running. %s", namespaceGuidance)
}

func StartModule(module string) error {
	if err := EnsureRuntimeNamespaceForMutation(); err != nil {
		return err
	}

	entry, err := runtime.EnsureModuleEntry(module)
	if err != nil {
		return err
	}

	if utils.CheckPortOpen(entry.Host, entry.Port) {
		if pid, err := utils.DetectPIDFromPort(entry.Port); err == nil {
			return fmt.Errorf("❌ Port %d on host %s is already in use by PID %d", entry.Port, entry.Host, pid)
		}
		return fmt.Errorf("❌ Port %d on host %s is already in use but listener PID could not be detected", entry.Port, entry.Host)
	}

	fmt.Printf("🚀 Starting module '%s' on %s:%d...\n", module, entry.Host, entry.Port)
	cmd := exec.Command("go", "run", utils.GetFullPath("modules", module, "main.go"))
	cmd.Env = append(os.Environ(), "GATEWAY_MODE=true")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("❌ Failed to start module: %w", err)
	}

	if err := waitForModuleReady(cmd, entry.Host, entry.Port, moduleStartTimeout); err != nil {
		return fmt.Errorf("❌ Module '%s' failed to become ready on %s:%d: %w", module, entry.Host, entry.Port, err)
	}

	status, pid, err := syncModuleState(module, entry)
	if err != nil {
		return err
	}
	if status == "STALE" {
		return fmt.Errorf("❌ Module '%s' opened %s:%d but listener PID could not be detected", module, entry.Host, entry.Port)
	}

	fmt.Printf("✅ Module '%s' is running on %s:%d (PID: %d)\n", module, entry.Host, entry.Port, pid)
	return nil
}

func StartModulesInBatch(modules []string) error {
	var wg sync.WaitGroup
	for _, module := range modules {
		wg.Add(1)
		go func(m string) {
			defer wg.Done()
			if err := StartModule(m); err != nil {
				fmt.Printf("⚠️  Failed to start module [%s]: %v\n", m, err)
			}
		}(module)
	}
	wg.Wait()
	return nil
}

func StopModule(module string) error {
	if err := EnsureRuntimeNamespaceForMutation(); err != nil {
		return err
	}

	registry, err := runtime.LoadRegistry()
	if err != nil {
		return fmt.Errorf("cannot load runtime registry: %w", err)
	}
	entry, ok := registry[module]
	if !ok || entry.Host == "" || entry.Port == 0 {
		return fmt.Errorf("❌ Module '%s' not found in tmp/runtime.json", module)
	}

	if !utils.CheckPortOpen(entry.Host, entry.Port) {
		if err := runtime.RemoveModuleEntry(module); err != nil {
			return fmt.Errorf("failed to clean runtime entry for module '%s': %w", module, err)
		}
		fmt.Printf("🧹 Module '%s' is not listening on %s:%d; runtime entry removed\n", module, entry.Host, entry.Port)
		return nil
	}

	pid, err := utils.DetectPIDFromPort(entry.Port)
	if err != nil || pid <= 0 {
		_ = runtime.UpdateRegistry(func(reg runtime.Registry) {
			current := reg[module]
			current.PID = 0
			reg[module] = current
		})
		return fmt.Errorf("❌ Module '%s' is reachable on %s:%d but listener PID could not be detected; not stopping an unknown process", module, entry.Host, entry.Port)
	}

	proc, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("❌ Could not find process %d: %w", pid, err)
	}

	if err := proc.Kill(); err != nil {
		return fmt.Errorf("❌ Failed to kill PID %d for module '%s': %w", pid, module, err)
	}

	if err := waitForPortClosed(entry.Host, entry.Port, 5*time.Second); err != nil {
		return fmt.Errorf("❌ Module '%s' PID %d was killed but %s:%d is still reachable: %w", module, pid, entry.Host, entry.Port, err)
	}

	if err := runtime.RemoveModuleEntry(module); err != nil {
		return fmt.Errorf("failed to remove runtime entry for module '%s': %w", module, err)
	}

	fmt.Printf("🛑 Stopped module [%s] on %s:%d (PID: %d)\n", module, entry.Host, entry.Port, pid)
	return nil
}

func RestartModule(module string) error {
	if err := StopModule(module); err != nil {
		return err
	}
	return StartModule(module)
}

func StopAllModules() error {
	if err := EnsureRuntimeNamespaceForMutation(); err != nil {
		return err
	}

	registry, err := runtime.LoadRegistry()
	if err != nil {
		return fmt.Errorf("cannot load runtime registry: %w", err)
	}

	names := make([]string, 0, len(registry))
	for name := range registry {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		if err := StopModule(name); err != nil {
			fmt.Printf("⚠️  Failed to stop module [%s]: %v\n", name, err)
		}
	}
	return nil
}

func SyncRunningModules() error {
	if err := EnsureRuntimeNamespaceForMutation(); err != nil {
		return err
	}

	entries, err := os.ReadDir(utils.GetFullPath("modules"))
	if err != nil {
		return fmt.Errorf("failed to read modules directory: %w", err)
	}

	runningCount := 0
	staleCount := 0
	stoppedCount := 0

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		moduleEntry, err := runtime.EnsureModuleEntry(name)
		if err != nil {
			fmt.Printf("⚠️  Skipping module '%s': %v\n", name, err)
			continue
		}

		status, pid, err := syncModuleState(name, moduleEntry)
		if err != nil {
			return err
		}

		switch status {
		case "RUNNING":
			runningCount++
			fmt.Printf("✅ Module '%s' is running on %s:%d (PID: %d)\n", name, moduleEntry.Host, moduleEntry.Port, pid)
		case "STALE":
			staleCount++
			fmt.Printf("⚠️  Module '%s' is STALE on %s:%d (PID: unknown)\n", name, moduleEntry.Host, moduleEntry.Port)
		default:
			stoppedCount++
			fmt.Printf("🛑 Module '%s' is not running on %s:%d\n", name, moduleEntry.Host, moduleEntry.Port)
		}
	}

	fmt.Printf("🔁 Synced runtime registry: %d running, %d stale, %d stopped\n", runningCount, staleCount, stoppedCount)
	return nil
}

func ShowStatus() error {
	entries, err := os.ReadDir(utils.GetFullPath("modules"))
	if err != nil {
		return fmt.Errorf("failed to read modules directory: %w", err)
	}

	type moduleRow struct {
		Name   string
		Host   string
		Port   int
		PID    int
		RunAt  string
		Status string
		Color  string
	}
	var rows []moduleRow

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		moduleEntry, err := runtime.ResolveModuleEntry(name)
		if err != nil {
			rows = append(rows, moduleRow{
				Name: name, Host: "-", Port: 0, PID: -1, RunAt: "-", Status: "CONFIG ERROR", Color: "\033[33m",
			})
			continue
		}

		status := "STOPPED"
		pid := -1
		color := "\033[31m"
		runAt := "-"
		if !moduleEntry.RunAt.IsZero() {
			runAt = moduleEntry.RunAt.Format("2006-01-02 15:04:05")
		}

		if utils.CheckPortOpen(moduleEntry.Host, moduleEntry.Port) {
			detectedPID, err := utils.DetectPIDFromPort(moduleEntry.Port)
			if err != nil || detectedPID <= 0 {
				status = "STALE"
				pid = 0
				color = "\033[33m"
			} else {
				status = "RUNNING"
				pid = detectedPID
				color = "\033[32m"
			}
		}

		rows = append(rows, moduleRow{
			Name: name, Host: moduleEntry.Host, Port: moduleEntry.Port, PID: pid, RunAt: runAt, Status: status, Color: color,
		})
	}

	sort.SliceStable(rows, func(i, j int) bool {
		return strings.ToLower(rows[i].Name) < strings.ToLower(rows[j].Name)
	})

	fmt.Println("\n📦 Module Status:")
	fmt.Println("----------------------------------------------------------------------------------------")
	fmt.Printf("%-20s | %-15s | %-5s | %-6s | %-20s | %-12s\n", "Module", "Host", "Port", "PID", "RunAt", "Status")
	fmt.Println("----------------------------------------------------------------------------------------")
	for _, r := range rows {
		fmt.Printf("%-20s | %-15s | %-5d | %-6d | %-20s | %s%-12s\033[0m\n",
			r.Name, r.Host, r.Port, r.PID, r.RunAt, r.Color, r.Status)
	}
	fmt.Println("----------------------------------------------------------------------------------------")
	return nil
}

func syncModuleState(module string, entry runtime.RunningModule) (string, int, error) {
	if !utils.CheckPortOpen(entry.Host, entry.Port) {
		err := runtime.UpdateRegistry(func(reg runtime.Registry) {
			current := reg[module]
			current.PID = 0
			current.Host = entry.Host
			current.Port = entry.Port
			current.Route = entry.Route
			current.External = entry.External
			reg[module] = current
		})
		return "STOPPED", 0, err
	}

	pid, err := utils.DetectPIDFromPort(entry.Port)
	if err != nil || pid <= 0 {
		updateErr := runtime.UpdateRegistry(func(reg runtime.Registry) {
			current := reg[module]
			current.PID = 0
			current.Host = entry.Host
			current.Port = entry.Port
			current.Route = entry.Route
			current.External = entry.External
			reg[module] = current
		})
		return "STALE", 0, updateErr
	}

	updateErr := runtime.UpdateRegistry(func(reg runtime.Registry) {
		current := reg[module]
		current.PID = pid
		current.Host = entry.Host
		current.Port = entry.Port
		current.Route = entry.Route
		current.External = entry.External
		current.RunAt = time.Now()
		reg[module] = current
	})
	return "RUNNING", pid, updateErr
}

func waitForModuleReady(cmd *exec.Cmd, host string, port int, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		if utils.CheckPortOpen(host, port) {
			return nil
		}

		if cmd.ProcessState != nil && cmd.ProcessState.Exited() {
			return fmt.Errorf("process exited before opening port")
		}

		time.Sleep(250 * time.Millisecond)

		if err := cmd.Process.Signal(syscall.Signal(0)); err != nil {
			if utils.CheckPortOpen(host, port) {
				return nil
			}
			return fmt.Errorf("process exited before opening port: %w", err)
		}
	}

	return fmt.Errorf("timeout after %s", timeout)
}

func waitForPortClosed(host string, port int, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if !utils.CheckPortOpen(host, port) {
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
	return fmt.Errorf("timeout after %s", timeout)
}

func isInsideContainer() bool {
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true
	}

	data, err := os.ReadFile("/proc/1/cgroup")
	if err != nil {
		return false
	}
	content := string(data)
	return strings.Contains(content, "docker") ||
		strings.Contains(content, "containerd") ||
		strings.Contains(content, "kubepods")
}

func composeAPIContainerRunning() (bool, error) {
	if _, err := exec.LookPath("docker"); err != nil {
		return false, nil
	}

	output, err := exec.Command("docker", "compose", "ps", "-q", "api").CombinedOutput()
	if err != nil {
		message := strings.TrimSpace(string(output))
		if message == "" {
			message = err.Error()
		}
		return false, fmt.Errorf("cannot verify Docker Compose API namespace: %s. %s", message, namespaceGuidance)
	}

	for _, id := range strings.Fields(string(output)) {
		inspectOutput, err := exec.Command("docker", "inspect", "-f", "{{.State.Running}}", id).CombinedOutput()
		if err != nil {
			message := strings.TrimSpace(string(inspectOutput))
			if message == "" {
				message = err.Error()
			}
			return false, fmt.Errorf("cannot verify Docker Compose API namespace: %s. %s", message, namespaceGuidance)
		}
		if strings.TrimSpace(string(inspectOutput)) == "true" {
			return true, nil
		}
	}
	return false, nil
}
