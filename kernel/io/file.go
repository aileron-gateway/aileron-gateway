// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package io

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/aileron-gateway/aileron-gateway/kernel/er"
)

// ReadFiles reads files in the given paths.
// Paths can be absolute or relative, and file or directory.
// ReadFiles returns an empty map and nil error when the given paths has no element.
func ReadFiles(recursive bool, paths ...string) (map[string][]byte, error) {
	results := map[string][]byte{}

	if len(paths) == 0 {
		return results, nil
	}

	files, err := ListFiles(recursive, paths...)
	if err != nil {
		return nil, err // Return err as-is.
	}

	for _, file := range files {
		b, err := os.ReadFile(file)
		if err != nil {
			return nil, (&er.Error{
				Package:     ErrPkg,
				Type:        ErrTypeFile,
				Description: ErrDscReadFile,
			}).Wrap(err)
		}
		results[file] = b
	}

	return results, nil
}

// ListFiles get file paths in the given paths.
// Paths can be directories or file paths.
// If there are no paths in given paths, it returns nil slice and nil error.
// The first argument of
func ListFiles(recursive bool, paths ...string) ([]string, error) {
	var files []string

	for _, path := range paths {
		if path == "" {
			continue // Ignore empty path.
		}

		path = filepath.Clean(path)
		err := filepath.Walk(path, func(pt string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			pt = filepath.Clean(pt)

			if info.IsDir() {
				if !recursive && path != pt {
					return fs.SkipDir // Skip internal dirs when not recursive.
				}
				return nil
			}

			files = append(files, pt)
			return nil
		})

		if err != nil && err != fs.SkipDir {
			return nil, (&er.Error{
				Package:     ErrPkg,
				Type:        ErrTypeFile,
				Description: ErrDscListFile,
			}).Wrap(err)
		}
	}

	return files, nil
}

// SplitMultiDoc splits a documents in []byte format with a given separator.
// If an empty separator is given, the default separator "---\n" is used.
// Empty contents are ignored.
func SplitMultiDoc(in []byte, sep string) [][]byte {
	if sep == "" {
		sep = "---\n"
	}

	var outArr [][]byte
	inArr := bytes.Split(in, []byte(sep))
	for _, b := range inArr {
		// exclude empty documents.
		b = bytes.Trim(b, " \n")
		if len(b) != 0 {
			outArr = append(outArr, b)
		}
	}

	return outArr
}
