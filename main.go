package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/packwiz/packwiz/cmdshared"
	"github.com/packwiz/packwiz/core"
	"github.com/spf13/viper"

	_ "github.com/packwiz/packwiz/curseforge"
	_ "github.com/packwiz/packwiz/github"
	_ "github.com/packwiz/packwiz/migrate"
	_ "github.com/packwiz/packwiz/modrinth"
)

var modPack = ""
var modDir = ""

func init() {
	flag.StringVar(&modPack, "p", "pack.toml", "Path to the Packwiz modpack file")
	flag.StringVar(&modDir, "o", "downloads", "Directory for saving mods")
	flag.Parse()
	viper.Set("pack-file", modPack)
}

func TryFatal(err error) {
	if err == nil {
		return
	}
	fmt.Println(err)
	os.Exit(1)
}
func main() {
	if modDir == "" {
		TryFatal(errors.New("need a -o option"))
	}
	pack, err := core.LoadPack()
	TryFatal(err)
	index, err := pack.LoadIndex()
	TryFatal(err)
	mods, err := index.LoadAllMods()
	TryFatal(err)
	session, err := core.CreateDownloadSession(mods, []string{index.HashFormat})
	TryFatal(err)
	session.GetManualDownloads()
	cmdshared.ListManualDownloads(session)
	TryFatal(CreateDir(modDir))
	for dl := range session.StartDownloads() {
		err := WriteMod(dl, modDir)
		TryFatal(err)
		TryFatal(os.Remove(dl.File.Name()))
	}
}
func WriteMod(dl core.CompletedDownload, output string) error {
	f, err := os.Create(filepath.Join(output, dl.Mod.FileName))
	if err != nil {
		return err
	}
	_, err = io.Copy(f, dl.File)
	return err
}

func CreateDir(dir string) error {
	dir_, err := os.Stat(dir)
	if errors.Is(err, fs.ErrNotExist) {
		return os.MkdirAll(dir, 0766)
	} else if err != nil {
		return err
	}
	if !dir_.IsDir() {
		return errors.New("output must be a directory")
	}
	return nil
}
