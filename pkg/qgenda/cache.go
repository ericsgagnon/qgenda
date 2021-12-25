package qgenda

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"time"
)

// consolidate these somewhere
func typeName(a any) string {
	return reflect.TypeOf(a).Name()
}

type CacheConfig struct {
	Dir           string
	ValidDuration time.Duration
	Enable        bool
}

func NewCacheConfig() *CacheConfig {
	dir, err := os.UserCacheDir()
	if err != nil {
		panic(err)
	}
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

type CacheFile struct {
	Name string
	CacheConfig
}

func (cf *CacheFile) String() string {
	file := filepath.Join(cf.Dir, cf.Name)
	file = filepath.Clean(file)
	return file
}

func NewCacheFile(file string, cfg *CacheConfig) *CacheFile {
	if file == "" || !cfg.Enable {
		return nil
	}
	if cfg == nil {
		cfg = NewCacheConfig()
	}
	return &CacheFile{
		Name:        file,
		CacheConfig: *cfg,
	}
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

func (cf *CacheFile) FileExists() bool {
	_, err := os.Stat(cf.String())
	return !os.IsNotExist(err)
}

// CheckCacheFile compares the cache files metadata with the given parameters
// to determine if it is 'valid'. It will always return false if it is unable
// to stat the file, if CacheFile.Enable == false, or the file size is 0.
func ValidCacheFile(cf *CacheFile) bool {
	if !cf.Enable {
		return false
	}
	fs, err := os.Stat(cf.String())
	if err != nil {
		return false
	}
	if time.Now().Sub(fs.ModTime()) > cf.ValidDuration {
		return false
	}
	if fs.Size() == 0 {
		return false
	}
	return true
}

// ReadCacheFile is a wrapper for os.ReadFile
func ReadCacheFile(cf *CacheFile, data []byte) error {
	data, err := os.ReadFile(cf.String())
	return err
}

// WriteCacheFile is a wrapper for os.WriteFile with standard defaults
func WriteCacheFile(cf *CacheFile, data []byte) error {
	os.WriteFile(cf.String(), data, 0666)
	return nil
}



// ReadCache reads the AuthToken from a file cache
func (t *AuthToken) xReadCache(ctx context.Context, filename string) error {

	log.Printf("Read cached AuthToken from file")
	// get the file to write cache to
	filename, err := t.CacheFile(ctx, filename)
	if err != nil {
		return err
	}

	f, err := os.Open(filename)
	if err != nil {
		log.Printf("Error opening AuthToken cache file %v: %v\n", filename, err)
		return err
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Printf("Error reading AuthToken cache file %v: %v\n", filename, err)
		return err
	}

	tkn := &AuthToken{
		Token:   &http.Header{},
		Expires: time.Time{},
	}
	if err := json.Unmarshal(b, tkn); err != nil {
		log.Printf("Error unmarshalling cached AuthToken: %v", err)
	}

	if tkn.Expires.UTC().Before(t.Expires.UTC()) || tkn.Expires.UTC().Before(time.Now().UTC()) {
		m := fmt.Sprintf("Cached AuthToken expired %v", tkn.Expires.UTC().String())
		log.Printf(m)
		return errors.New(m)
	}
	t.Token = tkn.Token
	t.Expires = tkn.Expires
	t = tkn

	log.Printf("Cached AuthToken appears valid until %v", t.Expires.UTC().String())
	return nil

}
