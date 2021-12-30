package qgenda

import (
	"errors"
	"os"
	"path/filepath"
	"time"
)

type CacheConfig struct {
	Dir           string
	ValidDuration time.Duration
	Enable        bool
}

// NewCacheConfig returns a *CacheConfig that defaults
// to os.UserCacheDir with subDir appended
// and a 30 day valid duration
func NewCacheConfig(subDir string) *CacheConfig {
	dir, err := os.UserCacheDir()
	if err != nil {
		panic(err)
	}
	dir = filepath.Join(dir, subDir)
	dur, err := time.ParseDuration("1h")
	if err != nil {
		panic(err)
	}
	return &CacheConfig{
		Dir:           dir,
		ValidDuration: (dur * 24 * 30),
		Enable:        true,
	}
}

// CacheFile helps manage caches
type CacheFile struct {
	Name    string    // name of the file
	Created time.Time // meant to capture cache create/recreate time
	Expires time.Time
	CacheConfig
}

// String returns the path/to/cache/file
func (cf *CacheFile) String() string {
	file := filepath.Join(cf.Dir, cf.Name)
	file = filepath.Clean(file)
	return file
}

// NewCacheFile returns a CacheFile based on *CacheConfig or
// the default CacheConfig if none is provided
func NewCacheFile(file string, subDir string, cfg *CacheConfig) *CacheFile {
	if file == "" || !cfg.Enable {
		return nil
	}
	if cfg == nil {
		cfg = NewCacheConfig(subDir)
	}

	cf := &CacheFile{
		Name:        file,
		Created:     time.Now().UTC(),
		CacheConfig: *cfg,
	}

	if cfg.ValidDuration > 0 {
		cf.Expires = cf.Created.Add(cfg.ValidDuration)
	}
	return cf
}

func (cf *CacheFile) CreateDir() error {
	return CreateCacheDir(cf)
}

// Create handles creation of both the directory (if needed)
// and the file
func (cf *CacheFile) Create() error {
	if err := cf.CreateDir(); err != nil {
		return err
	}
	return CreateCacheFile(cf)
}

func (cf *CacheFile) Valid() bool {
	return ValidCacheFile(cf)
}

// Read is a wrapper method for ReadCacheFile
func (cf *CacheFile) Read(data []byte) error {
	return ReadCacheFile(cf, data)
}

func (cf *CacheFile) Write(data []byte) error {
	return WriteCacheFile(cf, data)
}

func CreateCacheDir(cf *CacheFile) error {
	if cf.Dir == "" || !cf.Enable {
		return errors.New("Cache is disabled by CacheFile.Enable = false")
	}

	if err := os.MkdirAll(cf.Dir, 0777); err != nil {
		// log.Printf("Error making directory %v: %#v", cf.Dir, err)
		return err
	}
	return nil
}

// CreateCacheFile checks if the file exists and creates it if it doesn't
func CreateCacheFile(cf *CacheFile) error {
	if !cf.Enable {
		return errors.New("Cache is disabled by CacheFile.Enable = false")
	}

	if err := CreateCacheDir(cf); err != nil {
		return err
	}

	if !cf.FileExists() {
		if _, err := os.Create(cf.String()); err != nil {
			return err
		}
	}

	return nil
}

// FileExists only checks if the cache file exists
func (cf *CacheFile) FileExists() bool {
	_, err := os.Stat(cf.String())
	return !os.IsNotExist(err)
}

// CheckCacheFile compares the cache files metadata with the given parameters
// to determine if it is 'valid'. It will always return false if it is unable
// to stat the file, if CacheFile.Enable == false, or the file size is 0.
func ValidCacheFile(cf *CacheFile) bool {
	now := time.Now().UTC()
	fs, err := os.Stat(cf.String())
	if err != nil {
		return false
	}

	switch {
	case !cf.Enable:
		return false
	case fs.Size() == 0:
		return false
	case now.Sub(fs.ModTime()) > cf.ValidDuration:
		return false
	case now.After(cf.Expires):
		return false
	}
	return true
}

// ReadCacheFile is a wrapper for os.ReadFile. It doesn't parse anything, just consumes
// the file and writes it to data []byte. Note that, in keeping with os, it doesn't use
// a context.Context
func ReadCacheFile(cf *CacheFile, data []byte) error {
	data, err := os.ReadFile(cf.String())
	return err
}

// WriteCacheFile is a wrapper for os.WriteFile with standard defaults
func WriteCacheFile(cf *CacheFile, data []byte) error {
	os.WriteFile(cf.String(), data, 0666)
	return nil
}
