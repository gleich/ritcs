package remote

import (
	"crypto/md5"
	"math/big"
	"path/filepath"

	"go.mattglei.ch/ritcs/internal/conf"
)

func ProjectPath(cwd string) string {
	hash := md5.Sum([]byte(cwd))
	n := new(big.Int)
	n.SetBytes(hash[:])
	hashStr := n.Text(36)
	return filepath.Join(conf.Config.Home, ".ritcs", hashStr)
}
