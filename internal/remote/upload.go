package remote

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/pkg/sftp"
	"go.mattglei.ch/ritcs/internal/conf"
	"go.mattglei.ch/timber"
)

func UploadCWD(
	sftpClient *sftp.Client,
	ignoreStatements []string,
	remoteTarPath string,
) (int, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return 0, fmt.Errorf("%v failed to get working directory", err)
	}

	dirContent, err := os.ReadDir(cwd)
	if err != nil {
		return 0, fmt.Errorf("%v failed to read %s", err, cwd)
	}
	if len(dirContent) == 0 {
		return 0, nil
	}

	if !conf.Config.Silent {
		fmt.Println()
		timber.Info("uploading files")
	}

	remoteFile, err := sftpClient.Create(remoteTarPath)
	if err != nil {
		return 0, fmt.Errorf("%v failed to create remote file at %s", err, remoteTarPath)
	}
	defer remoteFile.Close()

	var (
		gw = gzip.NewWriter(remoteFile)
		tw = tar.NewWriter(gw)
	)
	defer gw.Close()
	defer tw.Close()

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
			}
			filesUploaded++
		}
		return nil
	})
	if err != nil {
		return 0, fmt.Errorf("%v failed to walk directory %s", err, cwd)
	}
	return filesUploaded, nil
}
