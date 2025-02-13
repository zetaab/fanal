package json

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"

	"golang.org/x/xerrors"

	"github.com/aquasecurity/fanal/analyzer"
	"github.com/aquasecurity/fanal/types"
)

const version = 1

var (
	requiredExt   = ".json"
	excludedFiles = []string{types.NpmPkgLock, types.NuGetPkgsLock, types.NuGetPkgsConfig}
)

type ConfigAnalyzer struct {
	filePattern *regexp.Regexp
}

func NewConfigAnalyzer(filePattern *regexp.Regexp) ConfigAnalyzer {
	return ConfigAnalyzer{
		filePattern: filePattern,
	}
}

func (a ConfigAnalyzer) Analyze(_ context.Context, input analyzer.AnalysisInput) (*analyzer.AnalysisResult, error) {
	var parsed interface{}
	if err := json.NewDecoder(input.Content).Decode(&parsed); err != nil {
		return nil, xerrors.Errorf("unable to decode JSON (%s): %w", input.FilePath, err)
	}

	return &analyzer.AnalysisResult{
		Configs: []types.Config{
			{
				Type:     types.JSON,
				FilePath: input.FilePath,
				Content:  parsed,
			},
		},
	}, nil
}

func (a ConfigAnalyzer) Required(filePath string, _ os.FileInfo) bool {
	if a.filePattern != nil && a.filePattern.MatchString(filePath) {
		return true
	}

	filename := filepath.Base(filePath)
	for _, excludedFile := range excludedFiles {
		if filename == excludedFile {
			return false
		}
	}

	return filepath.Ext(filePath) == requiredExt
}

func (ConfigAnalyzer) Type() analyzer.Type {
	return analyzer.TypeJSON
}

func (ConfigAnalyzer) Version() int {
	return version
}
