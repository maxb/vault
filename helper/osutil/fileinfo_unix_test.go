// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:build !windows

package osutil

import (
	"os"
	"syscall"
	"testing"
)

func TestFileUIDEqual(t *testing.T) {
	uid := syscall.Geteuid()

	testCases := []struct {
		uid      int
		expected bool
	}{
		{
			uid:      uid,
			expected: true,
		},
		{
			uid:      uid + 1,
			expected: false,
		},
	}

	for _, tc := range testCases {
		err := os.Mkdir("testFile", 0o777)
		if err != nil {
			t.Fatal(err)
		}
		info, err := os.Stat("testFile")
		if err != nil {
			t.Errorf("error stating %q: %v", "testFile", err)
		}

		result := FileUIDEqual(info, tc.uid)
		if result != tc.expected {
			t.Errorf("invalid result. expected %t for uid %v", tc.expected, tc.uid)
		}
		err = os.RemoveAll("testFile")
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestFileGIDEqual(t *testing.T) {
	gid := syscall.Getegid()

	testCases := []struct {
		gid      int
		expected bool
	}{
		{
			gid:      gid,
			expected: true,
		},
		{
			gid:      gid + 1,
			expected: false,
		},
	}

	for _, tc := range testCases {
		err := os.Mkdir("testFile", 0o777)
		if err != nil {
			t.Fatal(err)
		}
		info, err := os.Stat("testFile")
		if err != nil {
			t.Errorf("error stating %q: %v", "testFile", err)
		}

		result := FileGIDEqual(info, tc.gid)
		if result != tc.expected {
			t.Errorf("invalid result. expected %t for gid %v", tc.expected, tc.gid)
		}
		err = os.RemoveAll("testFile")
		if err != nil {
			t.Fatal(err)
		}
	}
}
