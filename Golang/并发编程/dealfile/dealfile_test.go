package dealfile_test

import (
	"happyladysauce/dealfile"
	"testing"
)

func TestDealFile(t *testing.T) {
	dealfile.DealMassFile("./data")
}