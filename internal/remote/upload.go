package remote

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/pkg/sftp"
	"pkg.mattglei.ch/ritcs/internal/conf"
	"pkg.mattglei.ch/ritcs/internal/util"
	"pkg.mattglei.ch/timber"
)

func UploadCWD(
	sftpClient *sftp.Client,
	ignoreStatements []string,
	remoteTarPath string,
) error {
	start := time.Now()
	if !conf.Config.Silent {
		fmt.Println()
		timber.Info("uploading files")
	}

	remoteFile, err := sftpClient.Create(remoteTarPath)
	if err != nil {
		return fmt.Errorf("%v failed to create remote file at %s", err, remoteTarPath)
	}
	defer remoteFile.Close()

	var (
		gw = gzip.NewWriter(remoteFile)
		tw = tar.NewWriter(gw)
	)
	defer gw.Close()
	defer tw.Close()

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("%v failed to get working directory", err)
	}

	filesUploaded := 0
	err = filepath.Walk(cwd, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(cwd, path)
		if err != nil {
			return fmt.Errorf("%v failed to get relative path for %s", err, path)
		}

		ignore := false
		for _, statement := range ignoreStatements {
			match, err := doublestar.Match(statement, relPath)
			if err != nil {
				return fmt.Errorf(
					"%v failed to check match with %s for %s",
					err,
					statement,
					relPath,
				)
			}
			if match {
				ignore = true
				break
			}
		}
		if ignore {
			return nil
		}

		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return fmt.Errorf("%v failed to create tar header", err)
		}
		header.Name = filepath.ToSlash(relPath)
		err = tw.WriteHeader(header)
		if err != nil {
			return fmt.Errorf("%v failed to write header", err)
		}

		if info.Mode().IsRegular() {
			file, err := os.Open(path)
			if err != nil {
				return fmt.Errorf("%v failed to open %s", err, path)
			}
			defer file.Close()

			_, err = io.Copy(tw, file)
			if err != nil {
				return fmt.Errorf("%v failed to copy %s to tar writer", err, path)
			}

			if !conf.Config.Silent {
				timber.Done("uploaded", relPath)
				filesUploaded++
			}
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("%v failed to walk directory %s", err, cwd)
	}
	if !conf.Config.Silent {
		timber.Done("uploaded", filesUploaded, "files in", util.FormatDuration(time.Since(start)))
	}
	return nil
}
