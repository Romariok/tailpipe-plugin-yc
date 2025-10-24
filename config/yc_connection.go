package config

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcl/v2"
)

const PluginName = "yc"

type YandexCloudConnection struct {
	ConnectionName string `hcl:"connection_name"`
	KeyFile string `hcl:"key_file"`
	FolderID string `hcl:"folder_id"`
	Remain hcl.Body `hcl:",remain"`
}

func (c *YandexCloudConnection) Validate() error {
	if strings.TrimSpace(c.ConnectionName) == "" {
		return fmt.Errorf("connection_name is required")
	}
	if strings.TrimSpace(c.KeyFile) == "" {
		return fmt.Errorf("key_file is required (path to authorized_key.json)")
	}
	if strings.TrimSpace(c.FolderID) == "" {
		return fmt.Errorf("folder_id is required")
	}
	return nil
}

func (c *YandexCloudConnection) Identifier() string {
	return PluginName
}
