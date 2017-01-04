package main

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
)

// FileCollector object.
type FileCollector struct {
	files      []File
	recursive  bool
	extentions []string
}

// File object.
type File struct {
	path      string
	size      int64
	modTime   time.Time
	filename  string
	extention string
}

// CollectAll collects all files in list of paths.
func (c *FileCollector) CollectAll(paths []string) error {
	if len(paths) == 0 {
		return nil
	}
	log.Debugf("start to collect %d paths", len(paths))
	for i, path := range paths {
		err := c.Collect(path)
		if err != nil {
			log.Errorf("path #%d: %v", i+1, err)
		}
	}
	log.Debugf("finish to collect %d files", len(c.files))
	return nil
}

// Collect collects all files in given one path.
func (c *FileCollector) Collect(path string) error {
	if len(path) == 0 {
		return errors.New("path is empty to collect files")
	}
	fileInfo, err := os.Stat(path)
	if err != nil {
		return err
	}
	if fileInfo.IsDir() {
		log.Infof("walk directory: %s", path)
		err := filepath.Walk(path, func(p string, info os.FileInfo, e error) error {
			if e != nil {
				return e
			}
			if info.IsDir() {
				if !c.recursive && path != p {
					return filepath.SkipDir
				}
			} else if c.isTarget(p) {
				return c.dispatch(p, info)
			}
			return nil
		})
		if err != nil {
			return err
		}
	} else if c.isTarget(path) {
		return c.dispatch(path, fileInfo)
	}
	return nil
}

func (c *FileCollector) dispatch(p string, fileInfo os.FileInfo) error {
	filename := fileInfo.Name()
	t := File{
		path:      p,
		size:      fileInfo.Size(),
		modTime:   fileInfo.ModTime(),
		filename:  filename,
		extention: strings.ToLower(path.Ext(filename)),
	}
	c.files = append(c.files, t)
	return nil
}

func (c *FileCollector) isTarget(p string) bool {
	if len(p) == 0 {
		return false
	} else if strings.HasPrefix(path.Base(p), ".") {
		// Skip hidden file.
		return false
	}
	if len(c.extentions) == 0 {
		return true
	}
	// If file extentions are set, filter them.
	ext := strings.ToLower(path.Ext(p))
	for _, fmt := range c.extentions {
		if ext == fmt {
			return true
		}
	}
	return false
}

// Checksum returns MD5 checksum as hex string.
func (f *File) Checksum() (string, error) {
	if f.path == "" {
		return "", fmt.Errorf("path is empty")
	}
	fp, err := os.Open(f.path)
	if err != nil {
		return "", err
	}
	defer fp.Close()
	hasher := md5.New()
	if _, err := io.Copy(hasher, fp); err != nil {
		return "", err
	}
	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func newFileCollector(recursive bool, extentions []string) *FileCollector {
	c := new(FileCollector)
	c.recursive = recursive
	c.extentions = extentions
	return c
}
