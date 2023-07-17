package main

import (
	"fmt"
	"os"

	"github.com/alexflint/go-arg"
	"github.com/anasrar/bulkdl/downloader"
	"github.com/anasrar/bulkdl/file"
	"github.com/anasrar/bulkdl/utils"
	tea "github.com/charmbracelet/bubbletea"
)

var VERSION = "dev"

type Args struct {
	Yaml                string `arg:"positional, required"`
	MaxParallelDownload int    `arg:"--parallel, -p" help:"Maximum parallel download" default:"3"`
}

func (Args) Version() string {
	return fmt.Sprintf("bulkdl %s", VERSION)
}

var (
	P *tea.Program
	D downloader.Downloader
	E []error
)

func main() {
	var args Args
	arg.MustParse(&args)
	if stat, err := os.Stat(args.Yaml); err != nil || stat.IsDir() {
		e := fmt.Sprintf("%s not found", args.Yaml)
		if err == nil && stat.IsDir() {
			e = fmt.Sprintf("%s is not file", args.Yaml)
		}
		fmt.Println(utils.SyleFgRed.Render(e))
		os.Exit(1)
	}
	y, err := os.ReadFile(args.Yaml)
	if err != nil {
		fmt.Println(utils.SyleFgRed.Render(err.Error()))
		os.Exit(2)
	}
	fileYaml, err := file.Unmarshal(y)
	if err != nil {
		fmt.Println(utils.SyleFgRed.Render(err.Error()))
		os.Exit(3)
	}
	D = downloader.New(fileYaml, args.MaxParallelDownload)
	go func() {
		E = D.Run()
		P.Send(TUIMsgFinish{})
	}()
	P = tea.NewProgram(TUIModel{})
	if _, err := P.Run(); err != nil {
		fmt.Println(TUIStyleFgRed.Render(err.Error()))
		os.Exit(4)
	}
	if *D.ItemsError > 0.0 {
		for _, err := range E {
			fmt.Println(TUIStyleFgRed.Render(err.Error()))
		}
		os.Exit(5)
	}
}
