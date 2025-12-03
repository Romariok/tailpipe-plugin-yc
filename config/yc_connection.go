package config

import (
	"fmt"
	"os"
	"strings"

	"log/slog"

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
			slog.Warn("Using key_file from YC_KEY_FILE environment variable")
		} else {
			return fmt.Errorf("key_file is required (path to authorized_key.json)")
		}
	}
	if strings.TrimSpace(c.FolderID) == "" {
		if v := strings.TrimSpace(os.Getenv("YC_FOLDER_ID")); v != "" {
			c.FolderID = v
			slog.Warn("Using folder_id from YC_FOLDER_ID environment variable")
		} else {
			return fmt.Errorf("folder_id is required")
		}
	}
	// Optional endpoint fallback from env
	if c.EndpointUrl == nil || strings.TrimSpace(*c.EndpointUrl) == "" {
		if v := strings.TrimSpace(os.Getenv("YC_ENDPOINT_URL")); v != "" {
			c.EndpointUrl = &v
			slog.Info("Using endpoint_url from YC_ENDPOINT_URL environment variable", "endpoint_url", v)
		}
	}
	return nil
}

func (c *YandexCloudConnection) Identifier() string {
	return PluginName
}
