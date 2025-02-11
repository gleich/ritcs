package conf

import (
	"bufio"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strings"

	"pkg.mattglei.ch/timber"
)

const IGNORE_FILENAME = ".ritcsignore"

func ReadIgnore() ([]string, error) {
	file, err := os.Open(IGNORE_FILENAME)
	if errors.Is(err, fs.ErrNotExist) {
		return []string{}, nil
	}
	if err != nil {
		return []string{}, fmt.Errorf("%v failed to open %s", err, IGNORE_FILENAME)
	}
	defer file.Close()

	lines := []string{}
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	err = scanner.Err()
	if err != nil {
		return []string{}, fmt.Errorf("%v failed to scan lines from file", err)
	}

	statements := []string{}
	for _, line := range lines {
		if !strings.HasPrefix(strings.TrimLeft(line, " "), "#") {
			statements = append(statements, line)
		}
	}

	if !Config.Silent {
		timber.Done("loaded", IGNORE_FILENAME)
	}

	return statements, nil
}
