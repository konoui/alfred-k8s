package rootcmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Collector interface {
	Collect() error
}

type User interface {
	Use(s string) error
}

type Deleter interface {
	Delete(s string) error
}

func (cfg *Config) CollectOutput(c Collector, query, cacheKey string) error {
	if err := cfg.Awf().Cache(cacheKey).MaxAge(cfg.cacheTTL).LoadItems().Err(); err == nil {
		cfg.Awf().Filter(query).Output()
		return nil
	}

	if err := c.Collect(); err != nil {
		return err
	}

	cfg.Awf().Cache(cacheKey).StoreItems().Workflow().Filter(query).Output()
	return nil
}

func (cfg *Config) UseOutput(u User, query string) (err error) {
	err = cfg.deleteAllCaches()
	if err != nil {
		return
	}

	err = u.Use(query)
	if err != nil {
		fmt.Fprintf(cfg.Stdout(), "Failed due to %s\n", err)
		return
	}

	fmt.Fprintf(cfg.Stdout(), "Success!!\n")
	return
}

func (cfg *Config) DeleteOutput(d Deleter, query string) (err error) {
	err = cfg.deleteAllCaches()
	if err != nil {
		return
	}

	err = d.Delete(query)
	if err != nil {
		fmt.Fprintf(cfg.Stdout(), "Failed due to %s\n", err)
		return
	}

	fmt.Fprintf(cfg.Stdout(), "Success!!\n")
	return
}

// deleteCache delete all resources for current namespace/context resources.
// 1. list pods in current ns
// 2. switch ns
// 3. list pods in switched current ns
func (cfg *Config) deleteAllCaches() error {
	cacheDir := cfg.cacheDir
	files, err := ioutil.ReadDir(cacheDir)
	if err != nil {
		return fmt.Errorf("invalid cache directory %s", cacheDir)
	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}
		if !strings.HasSuffix(f.Name(), cacheSuffix) {
			continue
		}

		path := filepath.Join(cacheDir, f.Name())
		if err := os.Remove(path); err != nil {
			return fmt.Errorf("failed to delete %s", path)
		}
	}
	return nil
}
