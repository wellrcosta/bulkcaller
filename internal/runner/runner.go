package runner

import (
	"github.com/wellrcosta/bulkcaller/internal/config"
	"github.com/wellrcosta/bulkcaller/internal/reader"
)

func Run(cfg *config.Config) error {
	if cfg.FilePath == "" {
		return nil
	}
	_, err := reader.ReadCSV(cfg.FilePath)
	return err
}
