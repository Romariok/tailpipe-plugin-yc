package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/hcl/v2"
)

const PluginName = "yc"

type YandexCloudConnection struct {
	ConnectionName        string   `hcl:"connection_name"`
	KeyFile               string   `hcl:"key_file"`
	FolderID              string   `hcl:"folder_id"`
	MaxErrorRetryAttempts *int     `hcl:"max_error_retry_attempts"`
	MinErrorRetryDelay    *int     `hcl:"min_error_retry_delay"`
	EndpointUrl           *string  `hcl:"endpoint_url"`
	Remain                hcl.Body `hcl:",remain"`
}

func (c *YandexCloudConnection) Validate() error {
	if strings.TrimSpace(c.ConnectionName) == "" {
		return fmt.Errorf("connection_name is required")
	}
	if strings.TrimSpace(c.KeyFile) == "" {
		if v := strings.TrimSpace(os.Getenv("YC_KEY_FILE")); v != "" {
			c.KeyFile = v
		} else {
			return fmt.Errorf("key_file is required (path to authorized_key.json)")
		}
	}
	if strings.TrimSpace(c.FolderID) == "" {
		if v := strings.TrimSpace(os.Getenv("YC_FOLDER_ID")); v != "" {
			c.FolderID = v
		} else {
			return fmt.Errorf("folder_id is required")
		}
	}
	// Optional endpoint fallback from env
	if c.EndpointUrl == nil || strings.TrimSpace(*c.EndpointUrl) == "" {
		if v := strings.TrimSpace(os.Getenv("YC_ENDPOINT_URL")); v != "" {
			c.EndpointUrl = &v
		}
	}
	return nil
}

func (c *YandexCloudConnection) Identifier() string {
	return PluginName
}
