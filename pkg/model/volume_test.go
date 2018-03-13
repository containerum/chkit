package model

import "testing"

func TestVolumeTableRendering(test *testing.T) {
	volume := new(Volume)
	headers := volume.TableHeaders()
	row := volume.TableRows()[0]
	if len(headers) != len(row) {
		test.Logf("\nHeaders: %v\nRow: %v", headers, row)
		test.Fatalf("num of headers and len of row are not equal!")
	}
}
