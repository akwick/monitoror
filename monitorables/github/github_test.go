package github

import (
	"os"
	"testing"

	"github.com/monitoror/monitoror/internal/pkg/monitorable/test"
	"github.com/stretchr/testify/assert"
)

func TestNewMonitorable(t *testing.T) {
	// init Store
	mockRouter, mockRouterGroup, mockConfigManager, s := test.InitMockAndStore()

	// init Env
	// OK
	_ = os.Setenv("MO_MONITORABLE_GITHUB_VARIANT0_TOKEN", "xxx")
	// Missing Token
	_ = os.Setenv("MO_MONITORABLE_GITHUB_VARIANT1_URL", "https://github.example.com/")
	// Url broken
	_ = os.Setenv("MO_MONITORABLE_GITHUB_VARIANT2_URL", "url%sgithub.example.com/")

	// NewMonitorable
	monitorable := NewMonitorable(s)
	assert.NotNil(t, monitorable)

	// GetDisplayName
	assert.NotNil(t, monitorable.GetDisplayName())

	// GetVariants and check
	if assert.Len(t, monitorable.GetVariants(), 4) {
		_, err := monitorable.Validate("variant1")
		assert.Error(t, err)
		_, err = monitorable.Validate("variant2")
		assert.Error(t, err)
	}

	// Enable
	for _, variant := range monitorable.GetVariants() {
		if valid, _ := monitorable.Validate(variant); valid {
			monitorable.Enable(variant)
		}
	}

	// Test calls
	mockRouter.AssertNumberOfCalls(t, "Group", 1)
	mockRouterGroup.AssertNumberOfCalls(t, "GET", 2)
	mockConfigManager.AssertNumberOfCalls(t, "RegisterTile", 3)
	mockConfigManager.AssertNumberOfCalls(t, "EnableTile", 2)
	mockConfigManager.AssertNumberOfCalls(t, "EnableDynamicTile", 1)
}
