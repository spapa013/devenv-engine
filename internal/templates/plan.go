package templates

import (
	"fmt"

	"github.com/nauticalab/devenv-engine/internal/config"
)

var devTemplates = []string{"statefulset", "service", "env-vars", "startup-scripts", "ingress"}

var systemTemplates = []string{"namespace"}

// BuildDevRenderPlan computes the template set from config before rendering.
func BuildDevRenderPlan(cfg *config.DevEnvConfig) ([]string, error) {
	if cfg == nil {
		return nil, fmt.Errorf("BuildDevRenderPlan requires non-nil config")
	}

	templateNames := make([]string, 0, len(devTemplates))
	for _, templateName := range devTemplates {
		if templateName == "ingress" && !cfg.ShouldRenderIngress() {
			continue
		}
		templateNames = append(templateNames, templateName)
	}

	return templateNames, nil
}

// BuildSystemRenderPlan computes the template set for system-level manifests.
func BuildSystemRenderPlan() []string {
	return append([]string{}, systemTemplates...)
}

func DevCleanupScope() []string {
	return append([]string{}, devTemplates...)
}

func SystemCleanupScope() []string {
	return append([]string{}, systemTemplates...)
}
