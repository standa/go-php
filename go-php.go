/**
 * Run TARGET_FILE phar with parameters after downloading it from an HTTP source.
 *
 * Compile by `go build go-php.php` and run go-php.exe.
 *
 * @author standa
 * @version 0.1
 */
package main

import (
	"archive/zip"
	"fmt"
	"os"
	"net/http"
	"io"
	"os/exec"
	"path/filepath"
	"strings"
	"log"
)

// @todo Parse manifest.json and use it instead of FILE_URL.
// @todo Read ./config.json to override parameters below.

const (
	FILE_URL     = "http://example.com/exchangeDirectory/client.phar"
	MANIFEST_URL = "http://example.com/exchangeDirectory/manifest.json"
	TARGET_DIR   = "runtime"
	TARGET_FILE  = "client.phar"
	PHP_ZIP_URL  = "http://windows.php.net/downloads/releases/archives/php-5.6.5-nts-Win32-VC11-x86.zip"
)


func main() {

	// download client.phar
	//if _, err := exec.LookPath(TARGET_DIR + "/" + TARGET_FILE); err != nil {
	// @todo do not download the phar archive always
		os.MkdirAll(TARGET_DIR, 0755)
		fmt.Print("Downloading ", FILE_URL, " into ", TARGET_DIR + "/" +TARGET_FILE, "... ")
		n, err := download(FILE_URL, TARGET_DIR + "/" + TARGET_FILE)
		if err != nil {
			fmt.Println("Download error: ", err)
		} else {
			fmt.Println(n, "bytes")
		}
	//}


	// download PHP runtime if necessary
//	binary, lookErr := exec.LookPath("php")
//	if lookErr != nil {
		binary, lookErr := exec.LookPath(TARGET_DIR + "/php/php.exe")
		if (lookErr != nil) {
			binary, lookErr = downloadPhp()
			if lookErr != nil {
				fmt.Println("Cannot find PHP executable on PATH and in " + TARGET_DIR + "/php/php.exe")
				return
			}
		}
//	}

	// output PHP version
	phpVersion, _ := exec.Command(binary, "-v").Output()
	fmt.Println("PHP Version:", string(phpVersion))

	// run the command
	fmt.Println("Running", binary, TARGET_DIR + "/" + TARGET_FILE, "start")
	cmd := exec.Command(binary, TARGET_DIR + "/" + TARGET_FILE, "start")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()

}

func downloadPhp() (string, error) {

	fmt.Print("Downloading PHP from ", PHP_ZIP_URL, "... ")
	n, err := download(PHP_ZIP_URL, TARGET_DIR + "/php-dist.zip")
	if err != nil {
		return "", err
	}
	fmt.Println("downloaded ", n, "bytes")

	unzip(TARGET_DIR + "/php-dist.zip", TARGET_DIR + "/php")
	os.Remove(TARGET_DIR + "/php-dist.zip")

	return exec.LookPath(TARGET_DIR + "/php/php.exe")
}

func download(url, dest string) (int64, error) {
	out, err := os.Create(dest)
	if err != nil {
//		log.Fatal("Error creating output file ", url, ": ", err)
		return 0, err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
//		log.Fatal("Error downloading ", url, ": ", err)
		return 0, err
	}
	defer resp.Body.Close()

	n, err := io.Copy(out, resp.Body)
	if err != nil {
//		fmt.Println("Error copying the file from the download directory: ", err)
		return 0, err
	}

//	fmt.Println("Downloaded ", n, " bytes")

	return n, nil
}

func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		fpath := filepath.Join(dest, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, f.Mode())
		} else {
			var fdir string
			if lastIndex := strings.LastIndex(fpath,string(os.PathSeparator)); lastIndex > -1 {
				fdir = fpath[:lastIndex]
			}

			err = os.MkdirAll(fdir, f.Mode())
			if err != nil {
				log.Fatal(err)
				return err
			}
			f, err := os.OpenFile(
				fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer f.Close()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
