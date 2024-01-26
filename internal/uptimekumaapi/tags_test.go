package uptimekumaapi

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTags(t *testing.T) {
	api, err := NewUptimeKumaAPI("http://uptime-api.immaleix.casa", "admin", "admin")
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	// CREATE TAG
	tagName := "test"
	tag, err := api.CreateTag(tagName, "red")
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	assert.Equal(t, "test", tag.Name)
	assert.Equal(t, "red", tag.Color)

	// LIST TAGS
	tags, err := api.GetTags()
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	fmt.Println("List of tags:")
	fmt.Println(*tags)
	if len(*tags) == 0 {
		t.Errorf("Expected at least one tag, got %d", len(*tags))
	}

	// DELETE TAG
	err = api.DeleteTag(tagName)
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}
}
