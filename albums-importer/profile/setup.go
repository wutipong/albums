package profile

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/url"
	"os"
	"path/filepath"

	"github.com/manifoldco/promptui"
	"go.yaml.in/yaml/v4"
)

func setupProfile(profile string) error {
	configMap, err := OpenConfigMap()
	if err != nil {
		slog.Error("unable to open configuration map", slog.String("error", err.Error()))
		return err
	}

	config, ok := configMap[profile]

	if !ok {
		fmt.Printf("Profile does not exists. Creating new profile [%s].\n", profile)
		config = Profile{}
	} else {
		fmt.Printf("Profile [%s] found. Updating existing profile.\n", profile)
	}

	urlPrompt := promptui.Prompt{
		Label: "Server URL",
		Validate: func(s string) error {
			if s == "" {
				return fmt.Errorf("server URL cannot be empty")
			}
			if u, err := url.Parse(s); err == nil {
				if u.Scheme != "http" && u.Scheme != "https" {
					return fmt.Errorf("invalid URL scheme: %s", u.Scheme)
				}
			} else {
				return fmt.Errorf("invalid server URL: %w", err)
			}
			return nil
		},
		Default:   config.URL,
		AllowEdit: true,
	}

	result, err := urlPrompt.Run()
	if err != nil {
		err = fmt.Errorf("prompt failed: %w", err)
		return err
	}
	config.URL = result

	configMap[profile] = config

	err = SaveConfigMap(configMap)
	if err != nil {
		return err
	}

	return nil
}

func CreateConfigPath() (path string, err error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		slog.Error(
			"failed to get user home directory",
			slog.String("error", err.Error()),
		)
		return
	}
	path = filepath.Join(homeDir, ".albums-importer", "config.yaml")

	return
}

func OpenConfigMap() (profileMap map[string]Profile, err error) {
	path, err := CreateConfigPath()
	if err != nil {
		return
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDONLY, 0644)
	if errors.Is(err, os.ErrNotExist) {
		slog.Warn(
			"Configuration file does not exist. Creating empty configuration map.",
			slog.String("path", path),
		)
		profileMap = make(map[string]Profile)
		err = nil
		return
	} else if err != nil {
		err = fmt.Errorf("failed to open configuration file: %w", err)
		return
	}
	defer file.Close()

	profileMap = make(map[string]Profile)
	decoder := yaml.NewDecoder(file)

	err = decoder.Decode(&profileMap)
	if errors.Is(err, io.EOF) {
		slog.Warn(
			"Configuration file is empty. Creating empty configuration map.",
			slog.String("path", path),
		)
		profileMap = make(map[string]Profile)
		return
	} else if err != nil {
		err = fmt.Errorf("failed to parse existing configuration file: %w", err)
		return
	}
	return
}

func SaveConfigMap(configMap map[string]Profile) error {
	path, err := CreateConfigPath()
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Dir(path), 0777)
	if err != nil {
		return fmt.Errorf("failed to create configuration directory: %w", err)
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to save configuration file: %w", err)
	}
	defer file.Close()

	yamlEncoder := yaml.NewEncoder(file)
	defer yamlEncoder.Close()

	err = yamlEncoder.Encode(configMap)
	if err != nil {
		return fmt.Errorf("failed to save configuration file: %w", err)
	}

	return nil
}
