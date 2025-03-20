package mmdbgeo

import (
	"context"
	"testing"
	"time"
)

func TestIsWorking(t *testing.T) {
	GlobalInit(context.Background(), 10*time.Second)
	russianIP := "109.108.32.1"
	kazakhIP := "103.106.3.1"

	dataRu, err := GetDataByIP(russianIP)
	if err != nil {
		t.Error(err)
		return
	}
	dataKz, err := GetDataByIP(kazakhIP)
	if err != nil {
		t.Error(err)
		return
	}

	if dataRu.CountryISO != "RU" {
		t.Error("Expected Russian IP to be in Russia, but got", dataKz.CountryISO)
	}

	if dataKz.CountryISO != "KZ" {
		t.Error("Expected Kazakh IP to be in Kazakhstan, but got", dataRu.CountryISO)
	}
}
