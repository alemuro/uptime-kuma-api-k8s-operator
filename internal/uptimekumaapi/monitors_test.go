package uptimekumaapi

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMonitors(t *testing.T) {
	api, err := NewUptimeKumaAPI("http://uptime-api.immaleix.casa", "admin", "admin")
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	// CREATE MONITOR
	monitor, err := api.CreateMonitor("test-operator", "http://test.com", 60, []string{"test"})
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	assert.Equal(t, "test-operator", monitor.Name)
	assert.Equal(t, "http://test.com", monitor.URL)
	assert.Equal(t, 60, monitor.Interval)
	assert.Equal(t, "http", monitor.Type)

	// LIST MONITORS
	monitors, err := api.GetMonitors()
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	fmt.Println("List of monitors:")
	fmt.Println(monitors)
	if len(*monitors) == 0 {
		t.Errorf("Expected at least one monitor, got %d", len(*monitors))
	}

	// DELETE MONITOR
	err = api.DeleteMonitor("test-operator")
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}
}
