package main

import "math"

//ResizeFromCenter resizes a bounding box by a scale factor from its center
func ResizeFromCenter(bbox []float64, scale_w float64, scale_h float64, scale_both bool) []float64 {

	if scale_both {
		scale_w = scale_h
	}

	if scale_h == 0 || scale_w == 0 {
		return bbox
	}
	w := bbox[2]
	h := bbox[3]

	dx := scale_w * w / 2
	dy := scale_h * h / 2
	xc := bbox[0] + w/2
	yc := bbox[1] + h/2

	bbox2 := make([]float64, 4)
	bbox2[0] = math.Max(xc-dx, 0)
	bbox2[1] = math.Max(yc-dy, 0)
	bbox2[2] = math.Min(xc+dx-bbox2[0], 1)
	bbox2[3] = math.Min(yc+dy-bbox2[1], 1)

	return bbox2
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}
