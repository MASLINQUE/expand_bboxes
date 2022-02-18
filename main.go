package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

var scale_w = flag.Float64("sw", 1.0, "Scale of widht")
var scale_h = flag.Float64("sh", 1.0, "Scale of height")
var scale_both = flag.Bool("sb", false, "If true, scale w and h by w value")

func main() {

	// open output file
	flag.Parse()

	// make a write buffer
	logrus.SetLevel(logrus.DebugLevel)
	s := bufio.NewScanner(os.Stdin)
	bufsize := 10 << 20
	buf := make([]byte, bufsize)
	s.Buffer(buf, bufsize)
	for {
		if s.Scan() {
			reqdata := s.Bytes()
			expand(reqdata)

		}
	}

}

func expand(reqdata []byte) {

	items := gjson.ParseBytes(reqdata).Get("items")
	if len(items.Array()) == 0 {
		fmt.Println(string(reqdata))
		return
	}
	responseStringArray := []string{}

	// dets := [][]float64{}
	for _, item := range items.Array() {
		item_str := item.String()
		bbox := []float64{}
		item.Get("bbox").ForEach(func(key, value gjson.Result) bool {
			bbox = append(bbox, value.Num)
			return true
		})

		bbox_res := ResizeFromCenter(bbox, *scale_w, *scale_h, *scale_both)

		item_str, _ = sjson.Set(item_str, "bbox.0", toFixed(bbox_res[0], 3))
		item_str, _ = sjson.Set(item_str, "bbox.1", toFixed(bbox_res[1], 3))
		item_str, _ = sjson.Set(item_str, "bbox.2", toFixed(bbox_res[2], 3))
		item_str, _ = sjson.Set(item_str, "bbox.3", toFixed(bbox_res[3], 3))
		responseStringArray = append(responseStringArray, item_str)
	}

	responseString := fmt.Sprintf("[%s]", strings.Join(responseStringArray, ", "))

	reqdata_st, _ := sjson.SetRaw(string(reqdata), "items", responseString)
	fmt.Println(string(reqdata_st))

}
