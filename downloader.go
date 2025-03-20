package mmdbgeo

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
)

const (
	// he deletes releasesm replace todo
	MMDBASNDownloadLink     = "https://github.com/P3TERX/GeoLite.mmdb/releases/download/2025.02.25/GeoLite2-ASN.mmdb"
	MMDBCityDownloadLink    = "https://github.com/P3TERX/GeoLite.mmdb/releases/download/2025.02.25/GeoLite2-City.mmdb"
	MMDBCountryDownloadLink = "https://github.com/P3TERX/GeoLite.mmdb/releases/download/2025.02.25/GeoLite2-Country.mmdb"

	MMDBASNName     = "GeoLite2-ASN.mmdb"
	MMDBCityName    = "GeoLite2-City.mmdb"
	MMDBCountryName = "GeoLite2-Country.mmdb"
	DBsPath         = "assets/GeoLite2"

	DownloadEnabled = false
)

func checkDBs(senderr chan error) {

	if err := checkAssetsDir(); err != nil {
		senderr <- err
		return
	}

	if DownloadEnabled {
		var wg sync.WaitGroup
		wg.Add(3)

		// removed root dir
		go checkAndDownload(fmt.Sprintf("%s/%s/%s", rootDir, DBsPath, MMDBASNName), MMDBASNDownloadLink, &wg, senderr)
		go checkAndDownload(fmt.Sprintf("%s/%s/%s", rootDir, DBsPath, MMDBCityName), MMDBCityDownloadLink, &wg, senderr)
		go checkAndDownload(fmt.Sprintf("%s/%s/%s", rootDir, DBsPath, MMDBCountryName), MMDBCountryDownloadLink, &wg, senderr)
		wg.Wait()
	} else {
		log.Printf("[INFO] Download geolite dbs has been disabled.")
	}
}

func checkAndDownload(path string, link string, wg *sync.WaitGroup, senderr chan<- error) {
	defer wg.Done()

	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Printf("[INFO] Downloading %s...\n", path)
		response, err := http.Get(link)
		if err != nil {
			senderr <- err
			return
		}
		defer response.Body.Close()
		file, err := os.Create(path)
		if err != nil {
			senderr <- err
			return
		}
		defer file.Close()
		_, err = io.Copy(file, response.Body)
		if err != nil {
			senderr <- err
			return
		}
	} else {
		log.Printf("[INFO] %s already exists\n", path)
	}
}

func checkAssetsDir() error {
	if _, err := os.Stat(DBsPath); os.IsNotExist(err) {
		log.Printf("[INFO] Creating assets directory: %s\n", DBsPath)
		err := os.MkdirAll(DBsPath, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}
