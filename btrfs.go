package main

import (
	"context"
	"fmt"
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
func CreateSnapshot(ctx context.Context, src, dst string, rw bool) error {
	_, err := os.Stat(dst)
	switch {
	case os.IsNotExist(err):

	case err == nil:
		return &SnapshotDestExistsError{dst}
	default:
		return err
	}

	args := make([]string, 0, 5)
	args = append(args, "subvolume", "snapshot")
	if !rw {
		args = append(args, "-r")
	}
	args = append(args, src, dst)

	btrfs := exec.CommandContext(ctx, BtrfsCommand, args...)
	btrfs.Stderr = os.Stderr

	return btrfs.Run()
}

func DeleteSubvol(ctx context.Context, path string) error {
	btrfs := exec.CommandContext(ctx, BtrfsCommand, "subvolume", "delete", "-c", path)
	btrfs.Stderr = os.Stderr

	return btrfs.Run()
}

// A SnapshotDestExistsError is returned by Snapshot if the
// destination already exists.
type SnapshotDestExistsError struct {
	Dest string
}

func (err SnapshotDestExistsError) Error() string {
	return fmt.Sprintf("Snapshot destination %q already exists", err.Dest)
}
