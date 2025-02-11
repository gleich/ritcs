package host

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"io/fs"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"

	"github.com/bmatcuk/doublestar/v4"
	"pkg.mattglei.ch/ritcs/internal/conf"
	"pkg.mattglei.ch/timber"
)

func CreateTarball(ignoreStatements []string) (string, error) {
	outPath := filepath.Join(
		os.TempDir(),
		"ritcs",
		fmt.Sprintf("%s.tar.gz", strconv.Itoa(rand.Int())),
	)
	err := os.MkdirAll(filepath.Dir(outPath), os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("%v failed to create parent directory for %s", err, outPath)
	}
	outFile, err := os.Create(outPath)
	if err != nil {
		return "", fmt.Errorf("%v failed to create tar file at %s", err, outPath)
	}
	defer outFile.Close()

	var (
		gw = gzip.NewWriter(outFile)
		tw = tar.NewWriter(gw)
	)
	defer gw.Close()
	defer tw.Close()

	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("%v failed to get working directory", err)
	}

	outputtedNewline := false

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
				return fmt.Errorf("%v failed to check match with %s for %s", err, statement, path)
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
			return fmt.Errorf("%v could not obtain tar header for %s", err, path)
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
				return fmt.Errorf("%v failed to copy %s into tar writer", err, path)
			}

			if !conf.Config.Silent {
				if !outputtedNewline {
					fmt.Println()
					outputtedNewline = true
				}
				timber.Done("uploaded", relPath)
			}
		}

		return nil
	})
	if err != nil {
		return "", fmt.Errorf("%v failed to walk directory", err)
	}

	return outPath, nil
}
