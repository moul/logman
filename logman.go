package logman

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"time"

	"go.uber.org/multierr"
	"moul.io/u"
)

// Manager is a configuration object used to create log files with automatic GC rules.
type Manager struct {
	// Path is the target directory containing the log files.
	// Default is '.'.
	Path string

	// MaxFiles is the maximum number of log files in the directory.
	// If 0, won't automatically GC based on this criteria.
	MaxFiles int

	// MaxFilesWithName is the maximum number of log files with the same app name.
	// If 0, won't automatically GC based on this criteria.
	// TODO: MaxFilesWithName int
}

// File defines a log file with metadata.
type File struct {
	// Full path.
	Path string

	// Size in bytes.
	Size int64

	// Provided name when creating the file.
	Name string

	// Creation date of the file.
	Time time.Time

	// Whether it is the most recent log file for the provided app name or not.
	Latest bool

	// If there were errors when trying to get info about this file.
	Errs error `json:"Errs,omitempty"`
}

// Create a new log file and perform automatic GC of the old log files if needed.
//
// The created log file will looks like:
//    <path/to/log/dir>/<name>-<time>.log
//
// Depending on the provided configuration of Manager, an automatic GC will be run automatically.
func (m Manager) New(name string) (io.WriteCloser, error) {
	// FIXME: validate m.Path
	// FIXME: validate name

	startTime := time.Now().Format(filePatternDateLayout)
	filename := filepath.Join(
		m.Path,
		fmt.Sprintf("%s-%s.log", name, startTime),
	)

	// run gc
	if err := m.gc(); err != nil {
		return nil, fmt.Errorf("auto GC: %w", err)
	}

	// create log dif if missing
	if dir := filepath.Dir(filename); !u.DirExists(dir) {
		err := os.MkdirAll(dir, 0o711)
		if err != nil {
			return nil, fmt.Errorf("create log dir: %w", err)
		}
	}

	// create/open the file and create a WriteCloser
	var writer io.WriteCloser
	if u.FileExists(filename) {
		var err error
		writer, err = os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		if err != nil {
			return nil, fmt.Errorf("open log file: %q: %w", filename, err)
		}
	} else {
		var err error
		writer, err = os.Create(filename)
		if err != nil {
			return nil, fmt.Errorf("create log file: %q: %w", filename, err)
		}
	}

	return writer, nil
}

const filePatternDateLayout = "2006-01-02T15-04-05.000"

var filePatternRegex = regexp.MustCompile(`(?m)^(.*)-(\d{4}-\d{2}-\d{2}T\d{2}-\d{2}-\d{2}.\d{3}).log$`)

// Files returns a list of existing log files.
func (m Manager) Files() ([]File, error) {
	dir := m.dir()
	osFiles, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	files := []File{}
	for _, file := range osFiles {
		sub := filePatternRegex.FindStringSubmatch(file.Name())
		if sub == nil {
			continue
		}
		t, err := time.Parse(filePatternDateLayout, sub[2])
		var errs error
		if err != nil {
			errs = multierr.Append(errs, err)
		}

		files = append(files, File{
			Path: filepath.Join(dir, file.Name()),
			Size: file.Size(),
			Name: sub[1],
			Time: t,
			Errs: errs,
		})
	}

	// compute latest
	if len(files) > 0 {
		var maxTime time.Time
		for _, file := range files {
			if file.Time.After(maxTime) {
				maxTime = file.Time
			}
		}
		for idx, file := range files {
			if file.Time == maxTime {
				files[idx].Latest = true
			}
		}
	}

	return files, nil
}

func (m Manager) dir() string {
	// FIXME: expand homedir, etc
	if m.Path == "" {
		return "."
	}
	return m.Path
}

func (m Manager) gc() error {
	dir := m.dir()

	if !u.DirExists(dir) {
		return nil
	}

	files, err := m.Files()
	if err != nil {
		return err
	}

	if len(files) < m.MaxFiles-1 {
		return nil
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].Time.Before(files[j].Time)
	})

	var errs error
	for i := 0; i <= len(files)-m.MaxFiles; i++ {
		err := os.Remove(files[i].Path)
		if err != nil {
			errs = multierr.Append(errs, err)
		}
	}
	return errs
}

/*
   // GCWithName cleans up old log files matching a specific name.
   func (m Manager) GCWithName(name string) error {
	return nil
	}

   // GC cleans up old log files.
   func (m Manager) GC() error {
	return nil
	}
*/

func (f File) String() string {
	return f.Path
}
