package types

type PixelFormat string

var PixelMorse PixelFormat = "morse"
var PixelSolid PixelFormat = "solid"

type Pixel struct {
	ID     int         `json:"id"`
	Red    int         `json:"red"`
	Green  int         `json:"green"`
	Blue   int         `json:"blue"`
	Format PixelFormat `json:"format"`
	Value  int64       `json:"value"`
}
