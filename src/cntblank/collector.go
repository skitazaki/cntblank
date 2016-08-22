package main

import (
	"errors"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
)

// FileCollector object.
type FileCollector struct {
	files      []TargetFile
	recursive  bool
	extentions []string
}

// TargetFile object.
type TargetFile struct {
	path      string
	size      int64
	modTime   time.Time
	filename  string
	extention string
	md5sum    string
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
	t := TargetFile{
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
	} else if strings.HasPrefix(p, ".") {
		// Skip hidden file.
		return false
	}
	ext := strings.ToLower(path.Ext(p))
	for _, fmt := range c.extentions {
		if ext == fmt {
			return true
		}
	}
	return false
}

func newFileCollector() (c *FileCollector, err error) {
	c = new(FileCollector)
	c.recursive = true
	c.extentions = []string{
		".csv",
		".tsv",
		".txt",
		".xlsx",
	}
	return c, nil
}
