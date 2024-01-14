package main

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"hash"
	"io"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/BurntSushi/toml"
)

var modDir = "./mods"

type Packwiz struct {
	Name     string `toml:"name"`
	Filename string `toml:"filename"`
	Side     string `toml:"side"`
	Download struct {
		URL      string `toml:"url"`
		Hash     string `toml:"hash"`
		HashType string `toml:"hash-format"`
	} `toml:"download"`
}

func init() {
	flag.StringVar(&modDir, "dir", "./mods", "directory to scan for mods")
	flag.Parse()
}
func LoadPackwiz(dir string) ([]Packwiz, error) {
	packwizSuffix := ".pw.toml"
	var err error
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var packwizFiles []Packwiz = make([]Packwiz, 0, len(files))
	for _, file := range files {
		if strings.HasSuffix(file.Name(), packwizSuffix) {
			p := Packwiz{}
			f, err := os.ReadFile(path.Join(dir, file.Name()))
			if err != nil {
				return nil, err
			}
			err = toml.Unmarshal(f, &p)
			if err != nil {
				return nil, err
			}
			packwizFiles = append(packwizFiles, p)
		}
	}
	return packwizFiles, nil
}
func HashSum(data []byte, hashType string) (string, error) {
	var hash hash.Hash
	switch hashType {
	case "sha1":
		hash = sha1.New()
	default:
		return "", errors.New("Unsupported hash type: " + hashType)
	}
	hash.Write(data)
	hashInBytes := hash.Sum(nil)
	hashString := hex.EncodeToString(hashInBytes)
	return hashString, nil
}
func Download(packwiz Packwiz, distDir string) error {
	distName := path.Join(distDir, packwiz.Filename)
	// 检查文件存在
	if _, err := os.Stat(distName); err == nil {
		fmt.Printf("%s already exists, skipping download.\n", distName)
		return nil // 文件已存在，无需下载
	}
	req, err := http.NewRequest("GET", packwiz.Download.URL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/octet-stream")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return errors.New("failed to download file")
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	hash, err := HashSum(data, packwiz.Download.HashType)
	if err != nil {
		return err
	}
	if hash != packwiz.Download.Hash {
		return errors.New("hash mismatch")
	}
	err = os.WriteFile(distName, data, 0644)
	if err != nil {
		return err
	}
	return nil
}
func main() {
	ps, err := LoadPackwiz(modDir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, p := range ps {
		err := Download(p, modDir)
		if err != nil {
			fmt.Printf("Failed to download %s: %s\n", p.Name, err)
		} else {
			fmt.Printf("Downloaded %s\n", p.Name)
		}
	}
	fmt.Println("Download complete.")
}
