// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package io

const (
	ErrPkg = "io"

	ErrTypeEnv  = "env"
	ErrTypeFile = "file"
	ErrTypeLF   = "logical file"

	// ErrDscSetEnv is a error description.
	// This description indicates the failure of
	// setting environmental variable.
	ErrDscSetEnv = "setting environmental variable failed."

	// ErrDscLoadEnv is a error description.
	// This description indicates the failure of
	// loading environmental variable.
	ErrDscLoadEnv = "loading environmental variable failed."

	// ErrDscReadFile is a error description.
	// This description indicates the failure of
	// reading file.
	ErrDscReadFile = "reading file failed."

	// ErrDscListFile is a error description.
	// This description indicates the failure of
	// listing directory files.
	ErrDscListFile = "listing files failed."

	// ErrDscLogicalFile is a error description.
	// This description indicates the failure of
	// logical file operation.
	ErrDscLogicalFile = "logical file operation failed."

	// ErrDscFileSys is a error description.
	// This description indicates the failure of
	// some file system operations.
	ErrDscFileSys = "file system operation failed."
)
