package remote

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/sftp"
	"pkg.mattglei.ch/ritcs/internal/conf"
	"pkg.mattglei.ch/ritcs/internal/util"
	"pkg.mattglei.ch/timber"
)

func DownloadFromTarball(sftpClient *sftp.Client, remoteTarPath string) error {
	start := time.Now()
	if !conf.Config.Silent {
		timber.Info("downloading files")
	}

	remoteFile, err := sftpClient.Open(remoteTarPath)
	if err != nil {
		return fmt.Errorf("%v failed to open remote file: %s", err, remoteTarPath)
	}
	defer remoteFile.Close()

	gzReader, err := gzip.NewReader(remoteFile)
	if err != nil {
		return fmt.Errorf("%v failed to create new gzip reader", err)
	}
	defer gzReader.Close()

	tarReader := tar.NewReader(gzReader)

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("%v failed to get working directory", err)
	}

	filesDownloaded := 0
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("%v failed to iterate through tar reader", err)
		}

		targetPath := filepath.Join(cwd, header.Name)
		if !strings.HasPrefix(filepath.Clean(targetPath), cwd) {
			return fmt.Errorf("%v illegal file path detected", err)
		}

		switch header.Typeflag {
		case tar.TypeDir:
			err := os.MkdirAll(targetPath, os.FileMode(header.Mode))
			if err != nil {
				return fmt.Errorf("%v failed to create directory %s", err, targetPath)
			}

		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
				return fmt.Errorf("failed to create parent directory for %q: %w", targetPath, err)
			}
			outFile, err := os.OpenFile(
				targetPath,
				os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
				os.FileMode(header.Mode),
			)
			if err != nil {
				return fmt.Errorf("%v failed to create file %s", err, targetPath)
			}
			_, err = io.Copy(outFile, tarReader)
			if err != nil {
				outFile.Close()
				return fmt.Errorf("%v failed to write file %s", err, targetPath)
			}
			if !conf.Config.Silent {
				relPath, err := filepath.Rel(cwd, targetPath)
				if err != nil {
					return fmt.Errorf("%v failed to get relative path for %s", err, targetPath)
				}
				timber.Done("downloaded", relPath)
				filesDownloaded++
			}
			outFile.Close()

		case tar.TypeSymlink:
			if err := os.Symlink(header.Linkname, targetPath); err != nil {
				return fmt.Errorf(
					"failed to create symlink %q -> %q: %w",
					targetPath,
					header.Linkname,
					err,
				)
			}

		default:
			timber.Warning("skipping unsupported file type in %q\n", header.Name)
		}
	}

	if !conf.Config.Silent {
		timber.Done(
			"downloaded",
			filesDownloaded,
			"files in",
			util.FormatDuration(time.Since(start)),
		)
	}

	return nil
}
