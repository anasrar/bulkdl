package downloader

import (
	"fmt"
	"time"

	"github.com/anasrar/bulkdl/file"
	"github.com/go-zoox/fetch"
	"github.com/zenthangplus/goccm"
)

type DownloaderItemStatus int

const (
	DownloaderItemStatusWait DownloaderItemStatus = iota
	DownloaderItemStatusRunning
	DownloaderItemStatusFinish
)

type DownloaderItem struct {
	file.File
	Status      DownloaderItemStatus
	Progress    float64
	FetchConfig fetch.Config
}

type Downloader struct {
	Max                 int
	Items               []DownloaderItem
	ItemsError          *float64
	ItemsFinish         *float64
	TOTAL_ITEMS_INT     int
	TOTAL_ITEMS_FLOAT   float64
	FINISH_STRING_WIDTH int
}

func (d Downloader) Run() (errs []error) {
	c := goccm.New(d.Max)
	for index := range d.Items {
		c.Wait()
		go func(index int) {
			d.Items[index].Status = DownloaderItemStatusRunning
			f, err := fetch.New().
				Download(d.Items[index].Url, d.Items[index].Filename, &d.Items[index].FetchConfig).
				SetProgressCallback(func(percent int64, current int64, total int64) {
					d.Items[index].Progress = (float64(current) / float64(total)) * 100
				}).
				Execute()
			if err != nil {
				*d.ItemsError += 1
				errs = append(errs, fmt.Errorf("%s: %s", d.Items[index].Filename, err.Error()))
			}
			if err == nil {
				sC := f.StatusCode()
				if sC != 200 {
					*d.ItemsError += 1
					errs = append(errs, fmt.Errorf("%s: Error Status Code, Expect 200 Got %d", d.Items[index].Filename, sC))
				}
			}
			*d.ItemsFinish += 1
			d.Items[index].Status = DownloaderItemStatusFinish
			c.Done()
		}(index)
	}
	c.WaitAllDone()
	return
}

func New(fileYaml file.Yaml, max int) Downloader {
	itemError, itemFinish := 0.0, 0.0
	d := Downloader{
		Max:                 max,
		Items:               []DownloaderItem{},
		ItemsError:          &itemError,
		ItemsFinish:         &itemFinish,
		TOTAL_ITEMS_INT:     len(fileYaml.Files),
		TOTAL_ITEMS_FLOAT:   float64(len(fileYaml.Files)),
		FINISH_STRING_WIDTH: len(fmt.Sprintf("%d", len(fileYaml.Files))),
	}
	for index, file := range fileYaml.Files {
		s := DownloaderItemStatusWait
		if index < max {
			s = DownloaderItemStatusRunning
		}
		config := fetch.Config{
			Method:  fileYaml.Config.Method,
			Headers: fileYaml.Config.Headers,
			Proxy:   fileYaml.Config.Proxy,
			Timeout: time.Second * time.Duration(fileYaml.Config.Timeout),
		}
		config.Merge(&fetch.Config{
			Method:  file.Config.Method,
			Headers: file.Config.Headers,
			Proxy:   file.Config.Proxy,
			Timeout: time.Second * time.Duration(file.Config.Timeout),
		})
		d.Items = append(d.Items, DownloaderItem{
			File:        file,
			Status:      s,
			Progress:    0,
			FetchConfig: config,
		})
	}
	return d
}
