package main

import (
	"os"
	"os/exec"
)

const (
	// BtrfsCommand is the name of the btrfs command-line tool binary.
	BtrfsCommand = "btrfs"
)

// Snapshot creates a Btrfs snapshot of src at dst, optionally making
// it writable. If dst already exists, the snapshot is not created and
// an error is returned.
func Snapshot(src, dst string, rw bool) error {
	args := make([]string, 0, 6)
	args = append(args, "btrfs", "subvolume", "snapshot")
	if !rw {
		args = append(args, "-r")
	}
	args = append(args, src, dst)

	btrfs := exec.Command(BtrfsCommand, args...)
	btrfs.Stdin = os.Stdin
	btrfs.Stdout = os.Stdout
	btrfs.Stderr = os.Stderr

	return btrfs.Run()
}
