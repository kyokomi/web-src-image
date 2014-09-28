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
	app.Run(os.Args)
}

const searchGalleryDoc = "#gdt div"
const searchImageDoc = "div a"
const readImageDoc = "#i3 a img"

func doMain(c *cli.Context) {

	writeDirPath := c.Args().Get(0)
	url := c.Args().Get(1)

	if writeDirPath == "" {
		log.Fatal("writeDirPath found not")
	}

	if url == "" {
		log.Fatal("args 2 url found not")
	}

	readImages(writeDirPath, url)
}

func readImages(writeDirPath, url string) error {
	var doc *goquery.Document
	var err error
	if doc, err = goquery.NewDocument(url); err != nil {
		return err
	}

	var wg sync.WaitGroup
	doc.Find(searchGalleryDoc).Each(func(_ int, s *goquery.Selection) {
		imageURL, hit := s.Find(searchImageDoc).Attr("href")
		if !hit {
			return
		}

		wg.Add(1)

		go func() {
			image, err := readImagePath(readImageDoc, imageURL)
			if err != nil || image == "" {
				wg.Done()
				return
			}
			writeImage(writeDirPath, image)
			wg.Done()
		}()
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
	fmt.Println(srcURL)

	return srcURL, nil
}

func writeImage(writeDir, url string) {

	idx := strings.LastIndex(url, "/")
	fileName := strings.Join([]string{writeDir, url[idx+1:]}, "/")

	_, err := ioutil.ReadFile(fileName)
	if err == nil {
		return
	}
	fmt.Println("fileName ", fileName)

	res, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	io.Copy(file, res.Body)

	return
}
