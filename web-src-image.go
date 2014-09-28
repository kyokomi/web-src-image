package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "web-src-image"
	app.Version = Version
	app.Usage = ""
	app.Author = "kyokomi"
	app.Email = "kyoko1220adword@gmail.com"
	app.Action = doMain
	app.Flags = []cli.Flag {
		cli.StringFlag{
			Name:   "name,n",
		},
		cli.StringFlag{
			Name:   "output-dir,o",
			Value:  "images",
		},
		cli.StringFlag{
			Name:   "url",
		},
		cli.IntFlag{
			Name:   "page-count,p",
			Value: 1,
		},
		cli.BoolFlag{
			Name:   "debug",
		},
	}
	app.Run(os.Args)
}

const searchGalleryDoc = "#gdt div"
const searchImageDoc = "div a"
const readImageDoc = "#i3 a img"

var debugFlag = false
func debugLog(a ...interface{}) {
	if debugFlag {
		fmt.Println(a...)
	}
}

func doMain(c *cli.Context) {

	name := c.String("name")
	path := c.String("output-dir")
	url := c.String("url")
	pageCnt := c.Int("page-count")
	debugFlag = c.Bool("debug")

	if name == "" {
		log.Fatal("writeDirPath found not")
	}

	if path == "" {
		log.Fatal("output-dir-path found not")
	}

	if url == "" {
		log.Fatal("target url found not")
	}

	var wg sync.WaitGroup

	// execute
	for i := 0; i < pageCnt; i++ {

		wg.Add(1)

		go func(pageIdx int) {
			pageUrl := fmt.Sprintf("%s?p=%d", url, pageIdx)
			pageName := fmt.Sprintf("%s-%d", name, pageIdx)
			if err := readImages(pageName, path, pageUrl); err != nil {
				fmt.Println(err)
			}
			wg.Done()
		}(i)
	}

	wg.Wait()
}

func readImages(baseName, writeDirPath, url string) error {
	var doc *goquery.Document
	var err error
	if doc, err = goquery.NewDocument(url); err != nil {
		return err
	}

	if err := createDir(writeDirPath); err != nil {
		return err
	}

	var count int
	var wg sync.WaitGroup
	doc.Find(searchGalleryDoc).Each(func(_ int, s *goquery.Selection) {
		imageURL, hit := s.Find(searchImageDoc).Attr("href")
		if !hit {
			return
		}

		wg.Add(1)
		count++
		go func(cnt int) {
			image, err := readImagePath(readImageDoc, imageURL)
			if err != nil || image == "" {
				wg.Done()
				return
			}
			filePath := fmt.Sprintf("%s/%s-%d", writeDirPath, baseName, cnt)
			writeImage(filePath, image)

			wg.Done()
		}(count)
	})

	wg.Wait()

	return nil
}

func readImagePath(searchQuery, imageURL string) (string, error) {
	var doc *goquery.Document
	var err error
	if doc, err = goquery.NewDocument(imageURL); err != nil {
		return "", err
	}

	srcURL, hit := doc.Find(searchQuery).Attr("src")
	if !hit {
		return "", nil
	}

	debugLog(srcURL)

	return srcURL, nil
}

// すでに存在する場合スルー
func createDir(dirPath string) error {
	// check
	if _, err := ioutil.ReadDir(dirPath); err == nil {
		return nil
	}

	// create dir
	return os.MkdirAll(dirPath, 0755)
}

func writeImage(filePath, url string) {

	// 拡張子
	idx := strings.LastIndex(url, ".")
	ex := url[idx:]
	// ファイル名
	writePath := filePath + ex

	// すでに存在する場合は何もしない
	if isFile(writePath) {
		return
	}

	fmt.Println(writePath)

	res, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	file, err := os.Create(writePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	io.Copy(file, res.Body)

	return
}

func isFile(filePath string) bool {
	f, err := os.Open(filePath)
	if err != nil {
		return false
	}
	defer f.Close()

	return true
}
