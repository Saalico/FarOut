package camera

type coordinates struct {
	x, y float64
}

type vector struct {
	dx, dy float64
}

type Camera struct {
	offset vector
	zoom   float64
	center coordinates
}

type cameraOptions struct {
	tetherRange float64
	zoom        float64
	locked      bool
}

func (c *Camera) getOffset() (float64, float64) {
	return c.offset.dx, c.offset.dy
}

func (c *Camera) setOffset(v vector) {
	c.offset = v
}

func (c *Camera) snapToCoordinates(xy coordinates) {
	c.center = xy
}
