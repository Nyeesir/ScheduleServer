package configLoader

import (
	"errors"
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
)

type AppConfig struct {
	ScraperUrl string `yaml:"scraper_url"`
	ServerPort string `yaml:"server_port"`
}

var (
	ErrFileNotFound           = errors.New("plik konfiguracyjny nie został znaleziony")
	ErrReadFailed             = errors.New("nie udało się odczytać pliku konfiguracyjnego")
	ErrUnmarshallFailed       = errors.New("wystąpił błąd podczas deserializacji YAML")
	ErrTemplateCreationFailed = errors.New("wystąpił błąd podczas tworzenia pliku konfiguracji")
	ErrFileSaveFailed         = errors.New("wystąpił błąd podczas zapisu pliku")
	ErrMarshallFailed         = errors.New("wystąpił błąd podczas konwersji do YAML")
)

func LoadYamlConfig[T any](filePath string, config *T) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("%w: %s", ErrFileNotFound, err)
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrReadFailed, err)
	}

	err = yaml.Unmarshal(data, config)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrUnmarshallFailed, err)
	}
	return nil
}

func LoadOrCreateYamlConfig[T any](filePath string, config *T, createIfNotExist bool) error {
	err := LoadYamlConfig[T](filePath, config)
	if errors.Is(err, ErrFileNotFound) {
		if createIfNotExist {
			err := CreateYamlConfigTemplate(filePath, config)
			if err != nil {
				return fmt.Errorf("%w: %s", ErrTemplateCreationFailed, err)
			}
		}
	} else if err != nil {
		return err
	}
	return nil
}

func CreateYamlConfigTemplate[T any](filePath string, config *T) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrMarshallFailed, err)
	}

	err = os.WriteFile(filePath, data, 0666)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrFileSaveFailed, err)
	}

	return nil
}
