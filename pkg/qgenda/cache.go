package qgenda

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type CacheConfig struct {
	Dir           string
	ValidDuration time.Duration
	Enable        bool
}

// NewCacheConfig returns a *CacheConfig that defaults to os.UserCacheDir
// with subDir appended and a 30 day valid duration
func NewCacheConfig(subDir string) (*CacheConfig, error) {
	dir, err := os.UserCacheDir()
	if err != nil {
		return nil, err
	}
	dir = filepath.Join(dir, subDir)
	dur, err := time.ParseDuration("1h")
	if err != nil {
		return nil, err
	}
	cf := &CacheConfig{
		Dir:           dir,
		ValidDuration: (dur * 24 * 30),
		Enable:        true,
	}
	return cf, nil
}

// DefaultCacheConfig is a simple wrapper for NewCacheConfig with a
// default directory
func DefaultCacheConfig() (*CacheConfig, error) {
	return NewCacheConfig("qgenda")
}

// CacheFile helps manage caches
type CacheFile struct {
	Name      string    // name of the file
	Timestamp time.Time // meant to capture cache create/recreate time
	Expires   time.Time
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
func NewCacheFile(file string, subDir string, cfg *CacheConfig) (*CacheFile, error) {
	switch {
	case cfg == nil:
		return nil, fmt.Errorf("CacheConfig missing")
	case file == "":
		return nil, fmt.Errorf("cache file name missing")
	case !cfg.Enable:
		return nil, fmt.Errorf("cache disabled by CacheFile.Enable = false")
	default:
		cfg.Dir = filepath.Join(cfg.Dir, subDir)
		cf := &CacheFile{
			Name:        file,
			Timestamp:   time.Now().UTC(),
			CacheConfig: *cfg,
		}

		if cfg.ValidDuration > 0 {
			cf.Expires = cf.Timestamp.Add(cfg.ValidDuration)
		}
		return cf, nil
	}
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

// Read is a wrapper method for ReadCacheFile
func (cf *CacheFile) Read() ([]byte, error) {
	return ReadCacheFile(cf)
}

func (cf *CacheFile) Write(data []byte) error {
	return WriteCacheFile(cf, data)
}

func (cf *CacheFile) Valid() bool {
	return ValidCacheFile(cf)
}
func CreateCacheDir(cf *CacheFile) error {
	if cf.Dir == "" || !cf.Enable {
		return fmt.Errorf("cache disabled by CacheFile.Enable = false")
	}
	if err := os.MkdirAll(cf.Dir, 0777); err != nil {
		return err
	}
	return nil
}

// CreateCacheFile checks if the file exists and creates it if it doesn't
func CreateCacheFile(cf *CacheFile) error {
	if !cf.Enable {
		return fmt.Errorf("cache disabled by CacheFile.Enable = false")
	}
	if err := CreateCacheDir(cf); err != nil {
		return err
	}
	if cf.FileExists() {
		return nil
	}
	f, err := os.Create(cf.String())
	if err != nil {
		return err
	}
	fi, err := f.Stat()
	if err != nil {
		return err
	}
	cf.Timestamp = fi.ModTime().UTC()
	if cf.ValidDuration > 0 {
		cf.Expires = cf.Timestamp.Add(cf.ValidDuration)
	}
	return nil
}

// ReadCacheFile is a wrapper for os.ReadFile. It doesn't parse anything,
// just consumes the file and returns it as []byte. Note that, in
// keeping with os, it doesn't use a context.Context
func ReadCacheFile(cf *CacheFile) ([]byte, error) {

	data, err := os.ReadFile(cf.String())
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(cf.String())
	if err != nil {
		return nil, err
	}
	cf.Timestamp = fi.ModTime().UTC()
	return data, nil
}

// WriteCacheFile writes to the file and updates the CacheFile.Timestamp
func WriteCacheFile(cf *CacheFile, data []byte) error {
	if err := CreateCacheFile(cf); err != nil {
		return err
	}
	os.WriteFile(cf.String(), data, 0666)
	cf.Timestamp = time.Now().UTC()
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

func CacheFileTimestamp(cf *CacheFile) time.Time {
	fi, err := os.Stat(cf.String())
	if err != nil {
		panic(err)
	}
	return fi.ModTime().UTC()
}
