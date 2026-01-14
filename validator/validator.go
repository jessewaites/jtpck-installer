package validator

import (
	"os/exec"
)

// ToolStatus represents the installation status of a tool
type ToolStatus struct {
	Name      string
	Installed bool
	Path      string
}

// CheckTools validates which tools are installed
func CheckTools(tools []string) []ToolStatus {
	var statuses []ToolStatus

	for _, tool := range tools {
		path, err := exec.LookPath(tool)
		status := ToolStatus{
			Name:      tool,
			Installed: err == nil,
			Path:      path,
		}
		statuses = append(statuses, status)
	}

	return statuses
}

// GetInstalledTools returns only the installed tools
func GetInstalledTools(tools []string) []string {
	var installed []string
	for _, status := range CheckTools(tools) {
		if status.Installed {
			installed = append(installed, status.Name)
		}
	}
	return installed
}

// GetMissingTools returns only the missing tools
func GetMissingTools(tools []string) []string {
	var missing []string
	for _, status := range CheckTools(tools) {
		if !status.Installed {
			missing = append(missing, status.Name)
		}
	}
	return missing
}
