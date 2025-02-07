package io

import (
	"cmp"
	"compress/gzip"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/aileron-gateway/aileron-gateway/kernel/er"
)

var (
	// OpenFile opens a file.
	// This can be replaced when testing.
	OpenFile = func(name string, flag int, perm fs.FileMode) (*os.File, error) {
		f, err := os.OpenFile(name, flag, perm)
		if err != nil {
			return os.OpenFile(name, flag, perm) // Retry once.
		}
		return f, nil
	}
	// Remove removes file.
	// This can be replaced when testing.
	Remove = func(name string) error {
		if err := os.Remove(name); err != nil {
			return os.Remove(name) // Retry once.
		}
		return nil
	}
	// Rename renames file.
	// This can be replaced when testing.
	Rename = func(oldpath string, newpath string) error {
		if err := os.Rename(oldpath, newpath); err != nil {
			return os.Rename(oldpath, newpath) // Retry once.
		}
		return nil
	}

	// idOnlyPattern is the pattern of IDs.
	// This pattern matches to, for example, ".0", ".1", ...
	idOnlyPattern = regexp.MustCompile(`^\.[0-9]+$`)
	// timeIDPattern is the pattern of timestamp and ID.
	// This pattern matched to, for example, ".2024-12-31_23-59-59.1", ...
	timeIDPattern = regexp.MustCompile(`^\..+\.[0-9]+$`)
)

func logErr(err error, args ...any) {
	if err == nil {
		return
	}
	e := (&er.Error{
		Package:     ErrPkg,
		Type:        ErrTypeLF,
		Description: "logging only error",
	}).Wrap(err)
	elems := []any{}
	elems = append(elems, time.Now().Local().Format(time.DateTime))
	elems = append(elems, "ERROR ["+e.Error()+"]")
	elems = append(elems, args...)
	fmt.Fprintln(os.Stderr, elems...)
}

func fatalErr(err error, args ...any) {
	if err == nil {
		return
	}
	stack := make([]byte, 1<<12) // Read max 4kiB stack traces.
	n := runtime.Stack(stack, false)
	args = append(args, "\n\nStackTrace:\n", string(stack[:n]))
	logErr(err, args...)
	fmt.Fprintln(os.Stderr, "Try graceful shutdown.")
	p, err := os.FindProcess(os.Getpid())
	if err != nil {
		return
	}
	_ = p.Signal(os.Interrupt)
}

// newMatchFunc return a new file name match function.
// The function returns <pattern> part if the given file name follows
// the pattern of `<base><pattern><ext>`, otherwise returns an empty string.
// ".*" pattern will be used if the given pattern is nil.
// When the "<pattern>" is an empty string, the match function returns a space " " instead.
//
// example:
//
//	matchFunc := newMatchFunc("test", ".log", regexp.MustCompile(`\.[0-9]+`))
//	println(matchFunc("test.log"))     // ""
//	println(matchFunc("test.0.log"))   // ".0"
//	println(matchFunc("test.foo.log")) // ""
func newMatchFunc(base, ext string, pattern *regexp.Regexp) func(string) string {
	if pattern == nil {
		pattern = regexp.MustCompile(`.*`)
	}
	return func(s string) string {
		if !strings.HasPrefix(s, base) || !strings.HasSuffix(s, ext) {
			return "" // Not matches to `<base><pattern><ext>`
		}
		s = strings.TrimPrefix(s, base)
		s = strings.TrimSuffix(s, ext)
		if !pattern.MatchString(s) {
			return "" // Not matches to `<base><pattern><ext>`
		}
		return cmp.Or(s, " ") // Return `<pattern>` part of `<base><pattern><ext>`
	}
}

// newParseFunc returns a new timestamp and id parse function.
// The function expects the pattern of `.<id>` or `.<timestamp>.<id>` as an argument.
// The <timestamp> will be parsed with the given layout and location.
// The parse function returns unix seconds of timestamp and id.
// The parse function returns -1,-1 when the timestamp or id was an invalid format.
// time.Local will be used if the given loc is nil.
//
// example:
//
//	parseFunc := newParseFunc("2006-01-02_15-04-05", time.UTC)
//	println(parseFunc(".1"))                     // 0, 1
//	println(parseFunc(".1970-01-01_00-00-00.1")) // 0, 1
//	println(parseFunc(".1970-01-01_00-01-00.1")) // 60, 1
func newParseFunc(layout string, loc *time.Location) func(string) (int64, int) {
	// Expect `.<id>` or `.<timestamp>.<id>` as the argument.
	if loc == nil {
		loc = time.Local
	}
	return func(pattern string) (created int64, id int) {
		switch pos := strings.LastIndex(pattern, "."); pos {
		case -1:
			return -1, -1 // Invalid format.
		case 0:
			n, err := strconv.Atoi(pattern[1:])
			if err != nil {
				return -1, -1 // Invalid format.
			}
			return 0, n
		default:
			t, err := time.ParseInLocation(layout, pattern[1:pos], loc)
			if err != nil {
				return -1, -1 // Invalid format.
			}
			n, err := strconv.Atoi(pattern[pos+1:])
			if err != nil {
				return -1, -1 // Invalid format.
			}
			return t.Unix(), n
		}
	}
}

type LogicalFileConfig struct {
	// Filename is the log filename.
	// For example "application.log"
	// Backup file names will be "application.<timestamp>.<id>.log".
	FileName string

	// TimeLayout is the layout of timestamp.
	// Layout must follow the go time layout.
	// Timestamp will be omitted from the backup file names
	// if this field is set to be empty.
	// See https://pkg.go.dev/time for the layout.
	TimeLayout string
	// TimeZone is the timezone of the timestamp.
	// For example "Asia/Tokyo".
	TimeZone string

	// SrcDir is the source file directory.
	// If empty, log files won't be managed.
	SrcDir string
	// DstDir is the destination file directory
	// where the managed files are located.
	// If empty, log files won't be manages.
	DstDir string

	// RotateSize is the log file size
	// that should be rotated.
	// The unit is byte.
	RotateSize int64

	// // Cron is the cron expression for time based log rotation.
	// // Time based rotation will be disabled is the value is not set.
	// // Format should be "minute hour day month week year".
	// // For example, "0 * * * *" means hourly rotation.
	// Cron string

	// CompressLv is the gzip compression level.
	// This must be valid compression level for gzip.
	// If 0, gzip compression won't be applied.
	CompressLv int

	// MaxBackup is the max number of the backup files.
	// If the number of backup files exceeded this value,
	// backup files will be removed from the older one.
	// If zero or negative, all backups will be removed.
	MaxBackup int

	// MaxAge is the max age of the backup files in second.
	// Files older than this age will be removed.
	// If zero or negative, all backups will be removed.
	MaxAge int64

	// MaxTotalSize is the total log file size.
	// If the total size exceeded this value,
	// backup files are deleted from older one.
	// The unit is byte.
	MaxTotalSize int64
}

func (c *LogicalFileConfig) New() (*LogicalFile, error) {
	c.FileName = cmp.Or(c.FileName, "application.log")
	c.SrcDir = filepath.Clean(cmp.Or(c.SrcDir, "./"))     // Use "./" if empty.
	c.DstDir = filepath.Clean(cmp.Or(c.DstDir, c.SrcDir)) // Use source directory if empty.

	if err := os.MkdirAll(c.SrcDir, os.ModePerm); err != nil {
		return nil, (&er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeFile,
			Description: ErrDscFileSys,
			Detail:      "failed to make directory.",
		}).Wrap(err)
	}
	if err := ReadWriteTest(c.SrcDir); err != nil {
		return nil, err // Return err as-is.
	}
	if err := os.MkdirAll(c.DstDir, os.ModePerm); err != nil {
		return nil, (&er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeFile,
			Description: ErrDscFileSys,
			Detail:      "failed to make directory.",
		}).Wrap(err)
	}
	if err := ReadWriteTest(c.DstDir); err != nil {
		return nil, err // Return err as-is.
	}

	name := filepath.Base(c.FileName)     // Filename of the file. "foo.log" for "aaa/bbb/foo.log".
	ext := filepath.Ext(name)             // Extension of the file. ".log" for "foo.log".
	base := strings.TrimSuffix(name, ext) // Basename of the file. "foo" for "foo.log".

	var pattern *regexp.Regexp
	var timeFunc func() string

	loc := location(c.TimeZone)
	layout := c.TimeLayout
	if c.TimeLayout == "" {
		pattern = idOnlyPattern
	} else {
		pattern = timeIDPattern
		timeFunc = func() string {
			return "." + time.Now().In(loc).Format(layout)
		}
	}

	archiveExt := ext
	if c.CompressLv != gzip.NoCompression {
		archiveExt += ".gz"
	}

	m := &fileManager{
		maxAge:       c.MaxAge,
		maxBackup:    c.MaxBackup,
		maxTotalSize: c.MaxTotalSize,
		targetDir:    c.DstDir,
		matchFunc:    newMatchFunc(base, archiveExt, pattern),
		parseFunc:    newParseFunc(layout, loc),
	}

	f := &LogicalFile{
		curFile:        os.Stderr,
		manageFunc:     m.manage,
		srcMatchFunc:   newMatchFunc(base, ext, pattern),
		parseFunc:      newParseFunc(layout, loc),
		fileBase:       base,
		fileExt:        ext,
		fileArchiveExt: archiveExt,
		timeFunc:       timeFunc,
		srcDir:         c.SrcDir,
		dstDir:         c.DstDir,
		rotateSize:     c.RotateSize,
		compressLv:     c.CompressLv,
	}

	// Create initial file.
	if err := f.swapFile(); err != nil {
		return nil, (&er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeLF,
			Description: ErrDscLogicalFile,
		}).Wrap(err)
	}

	return f, nil
}

type LogicalFile struct {
	// mu protects curFile.
	mu sync.RWMutex

	// curFile is the current physical file.
	curFile *os.File
	// curSize is the current file size in Byte.
	curSize atomic.Int64
	// swapping is the flag that the file is swapping.
	swapping atomic.Bool

	// manageFunc manages archived files.
	// This function should be called after physical file
	// was swapped, or file rotated.
	manageFunc func() error

	// fileBase is the base name of log file.
	// e.g. "application" for "application.log".
	fileBase string
	// fileExt is the file extension of log file.
	// e.g. ".log" for "application.log".
	fileExt string
	// fileArchiveExt is the file extension of archived log file.
	// e.g. ".log.gz" for "application.log.gz".
	fileArchiveExt string

	// srcDir is the source directory of log files.
	// curFile will be placed in this directory.
	srcDir string
	// dstDir is the destination directory of log backup files.
	// Archived files will be located in this directory.
	dstDir string

	// timeFunc is the function that returns
	// string expression of a timestamp.
	timeFunc func() string

	srcMatchFunc func(string) string
	parseFunc    func(string) (int64, int)

	// compressLv is the gzip compression level.
	// Backup files won't be compressed if 0.
	compressLv int

	// rotateSize is the max log file size in Byte.
	// curFile will be rotated when the curSize
	// exceeded this size.
	// This cannot be zero or negative.
	rotateSize int64
}

func (f *LogicalFile) Write(b []byte) (int, error) {
	f.mu.RLock() // Protect f.curFile. It should not be swapped now.
	defer f.mu.RUnlock()
	n, err := f.curFile.Write(b)
	if f.curSize.Add(int64(n)) >= f.rotateSize && !f.swapping.Swap(true) {
		go func() { _ = f.swapFile() }()
	}
	return n, err
}

func (f *LogicalFile) Close() error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.curFile == nil {
		return nil
	}
	if err := f.curFile.Close(); err != nil {
		return err
	}
	name := f.issueArchiveFileName()
	err := Rename(f.curFile.Name(), filepath.Join(f.srcDir, name))
	logErr(err)
	if err := compressFiles(f.srcDir, f.dstDir, f.compressLv, f.srcMatchFunc); err != nil {
		return err
	}
	// We don't manage archived files because
	// the current file should better not to be contained
	// in the case it does not contain enough data to count in.
	return nil
}

func (f *LogicalFile) SwapFile() error {
	if f.swapping.Swap(true) {
		return nil // Swapping under progress.
	}
	return f.swapFile()
}

func (f *LogicalFile) swapFile() error {
	f.swapping.Store(true)
	defer f.swapping.Store(false)
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.curFile != nil && (uintptr(unsafe.Pointer(f.curFile)) != uintptr(unsafe.Pointer(os.Stderr))) {
		// It is rare that Close returns error.
		// But os.ErrClosed can happen if the file was remove by others.
		// When the err is non-nil, we can't do anything except for ignoring the file.
		err := f.curFile.Close()
		if err == nil {
			name := f.issueArchiveFileName()
			err := Rename(f.curFile.Name(), filepath.Join(f.srcDir, name))
			logErr(err) // Log output only.
		}
		go func() {
			err := compressFiles(f.srcDir, f.dstDir, f.compressLv, f.srcMatchFunc)
			logErr(err) // Log output only.
			err = f.manageFunc()
			logErr(err) // Log output only.
		}()
	}

	f.curSize.Store(0)
	activeLog := filepath.Join(f.srcDir, f.fileBase+f.fileExt)
	ff, err := OpenFile(activeLog, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		f.curFile = os.Stderr
		fatalErr(err) // Log error and graceful shutdown.
		return err
	}
	info, err := ff.Stat()
	if err != nil {
		f.curFile = os.Stderr
		fatalErr(err) // Log error and graceful shutdown.
		return err
	}
	f.curSize.Store(info.Size())
	f.curFile = ff
	return nil
}

func (f *LogicalFile) issueArchiveFileName() string {
	filename := f.fileBase
	if f.timeFunc != nil {
		filename += f.timeFunc()
	}

	var values []int
	srcMatchFunc := newMatchFunc(filename, f.fileExt, idOnlyPattern)
	srcFilePaths, _ := ListFiles(false, f.srcDir)
	for _, path := range srcFilePaths {
		base := filepath.Base(path)
		if s := srcMatchFunc(base); s != "" {
			_, n := f.parseFunc(s)
			values = append(values, n)
		}
	}
	dstMatchFunc := newMatchFunc(filename, f.fileArchiveExt, idOnlyPattern)
	dstFilePaths, _ := ListFiles(false, f.dstDir)
	for _, path := range dstFilePaths {
		base := filepath.Base(path)
		if s := dstMatchFunc(base); s != "" {
			_, n := f.parseFunc(s)
			values = append(values, n)
		}
	}

	id := 0
	if len(values) > 0 {
		sort.Sort(sort.Reverse(sort.IntSlice(values)))
		id = values[0] + 1
	}

	return filename + "." + strconv.Itoa(id) + f.fileExt
}

// fileInfo is the information of archived files.
type fileInfo struct {
	// path is the file path.
	path string
	// size is the file size in bytes.
	size int64
	// age is the file age in seconds.
	age int64
	// id is the incremented number.
	id int
}

type fileManager struct {
	// MaxAge is the max age of the backup files in second.
	// Files older than this age will be removed.
	// If zero or negative, all backups will be removed.
	maxAge int64

	// MaxBackup is the max number of the backup files.
	// If the number of backup files exceeded this value,
	// backup files will be removed from the older one.
	// If zero or negative, all backups will be removed.
	maxBackup int

	// maxTotalSize is the maximum file size in byte
	// of sum of the all backup files.
	// There will be no limit when set to 0.
	maxTotalSize int64

	targetDir string
	matchFunc func(string) string
	parseFunc func(string) (int64, int)
}

// Manage manages backup files.
// Calling this method may take time because
// of the file operation such as compression.
func (m *fileManager) manage() error {
	filePaths, _ := ListFiles(false, m.targetDir)
	fs := []*fileInfo{}
	now := time.Now().Unix()
	for _, path := range filePaths {
		info, err := os.Stat(path)
		if err != nil {
			continue // Ignore file that cannot get file info.
		}

		name := info.Name()
		if s := m.matchFunc(name); s != "" {
			// s is `.<id>` or `.<timestamp>.<id>`
			created, id := m.parseFunc(s)
			if created < 0 || id < 0 {
				continue // timestamp or id is invalid.
			}
			fs = append(fs, &fileInfo{
				path: path,
				size: info.Size(),
				age:  now - created,
				id:   id,
			})
		}
	}

	sort.SliceStable(fs, func(i, j int) bool {
		if fs[i].age != fs[j].age {
			return fs[i].age < fs[j].age
		}
		return fs[i].id > fs[j].id
	})

	totalSize := int64(0)
	for i, f := range fs {
		if m.maxBackup > 0 && i >= m.maxBackup {
			logErr(Remove(f.path))
		}
		if m.maxAge > 0 && f.age > m.maxAge {
			logErr(Remove(f.path))
		}
		totalSize += f.size
		if m.maxTotalSize > 0 && totalSize > m.maxTotalSize {
			logErr(Remove(f.path))
		}
	}

	return nil
}

// location returns time location.
// If an invalid timezone string was given, local timezone will be returned.
func location(tz string) *time.Location {
	if tz, err := time.LoadLocation(tz); err == nil {
		return tz
	}
	return time.Local
}

// compressFiles compresses files which matched to the given matchFunc
// in the srcDir and save it to dstDir.
// This function will panic if the matchFunc is nil.
func compressFiles(srcDir, dstDir string, level int, matchFunc func(string) string) error {
	if srcDir == dstDir && level == gzip.NoCompression {
		return nil
	}
	filePaths, err := ListFiles(false, srcDir)
	if err != nil {
		return err
	}
	var errs []error
	for _, path := range filePaths {
		base := filepath.Base(path)
		if s := matchFunc(base); s != "" {
			err := gzipFile(srcDir, dstDir, base, level)
			errs = append(errs, err)
		}
	}
	// errors.Join ignores nil error.
	// If no errors found, nil wil be returned.
	return errors.Join(errs...)
}

// gzipFile compresses the file.
// This function compress the file srcDir/filename and
// place it to dstDir/filename+".gz".
// This function does not compress the file when the level is 0
// and just rename the srcDir/filename to dstDir/filename.
// This function does not compress the file if the given filename
// has the suffix ".gz" and moves srcDir/filename to dstDir/filename.
func gzipFile(srcDir, dstDir, filename string, level int) error {
	srcFile := filepath.Join(srcDir, filename) // source file. For example, "srcDir/foo.log"
	dstFile := filepath.Join(dstDir, filename) // destination file. For example, "dstDir/foo.log"
	if level != gzip.NoCompression {
		dstFile += ".gz" // Add extension if compression enabled. For example, "dstDir/foo.log.gz"
	}

	if srcFile == dstFile {
		return nil
	}

	if level == gzip.NoCompression {
		if err := Rename(srcFile, dstFile); err != nil {
			return (&er.Error{
				Package:     ErrPkg,
				Type:        ErrTypeFile,
				Description: ErrDscFileSys,
				Detail:      "rename file " + srcFile + " to " + dstFile,
			}).Wrap(err)
		}
		return nil
	}

	src, err := OpenFile(srcFile, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return (&er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeFile,
			Description: ErrDscFileSys,
			Detail:      "open file " + srcFile,
		}).Wrap(err)
	}
	defer src.Close()

	// Overwrite the destination file if exists by passing the O_TRUNC flag.
	dst, err := OpenFile(dstFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return (&er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeFile,
			Description: ErrDscFileSys,
			Detail:      "open file " + dstFile,
		}).Wrap(err)
	}
	defer dst.Close()

	if level < gzip.HuffmanOnly {
		level = gzip.HuffmanOnly
	}
	if level > gzip.BestCompression {
		level = gzip.BestCompression
	}

	// No error returned here because we have
	// restricted the compression level to correct range.
	gw, _ := gzip.NewWriterLevel(dst, level)
	defer gw.Close() // gw must be closed before dst file closed.

	if _, err := CopyBuffer(gw, src); err != nil {
		gw.Close() // gw must be closed before dst file closed.
		dst.Close()
		logErr(Remove(dstFile)) // Remove unused destination file.
		return (&er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeFile,
			Description: ErrDscFileSys,
		}).Wrap(err)
	}

	// Remove the source file.
	src.Close()
	if err := Remove(srcFile); err != nil {
		return (&er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeFile,
			Description: ErrDscFileSys,
			Detail:      "removing " + srcDir,
		}).Wrap(err)
	}

	return nil
}

// ReadWriteTest tests directory permission.
// This functions checks if the application has read and write permission.
// Calling this function creates temporary file named "permission-check.txt"
// in the target folder.
func ReadWriteTest(dir string) error {
	name := filepath.Clean(filepath.Join(dir, "permission-check.txt"))
	f, err := OpenFile(name, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		return (&er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeFile,
			Description: ErrDscFileSys,
			Detail:      "directory permission check failed.",
		}).Wrap(err)
	}
	f.Close()
	if err := Remove(name); err != nil {
		return (&er.Error{
			Package:     ErrPkg,
			Type:        ErrTypeFile,
			Description: ErrDscFileSys,
			Detail:      "directory permission check failed.",
		}).Wrap(err)
	}
	return nil
}
