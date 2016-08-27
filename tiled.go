package tiled

import (
	"encoding/hex"
	"encoding/xml"
	"errors"
	"io"
	"strconv"
	"strings"
)

func Read(r io.Reader) (m Map, err error) {
	err = xml.NewDecoder(r).Decode(&m)
	return
}

type Map struct {
	XMLName         xml.Name    `xml:"map"`
	Version         string      `xml:"version,attr"`
	Orientation     Orientation `xml:"orientation,attr"`
	RenderOrder     RenderOrder `xml:"renderorder,attr"`
	Width           int         `xml:"width,attr"`
	Height          int         `xml:"height,attr"`
	TileWidth       int         `xml:"tilewidth,attr"`
	TileHeight      int         `xml:"tileheight,attr"`
	HexSideLength   int         `xml:"hexsidelength,attr"`
	StaggerAxis     int         `xml:"staggeraxis,attr"`
	StaggerIndex    int         `xml:"staggerindex,attr"`
	BackgroundColor Color       `xml:"backgroundcolor,attr"`
	NextObjectID    int         `xml:"nextobjectid,attr"`
	TileSets        []TileSet   `xml:"tileset"`
	Layers          []Layer     `xml:"layer"`
}

type Orientation string

const (
	Orthogonal Orientation = "orthogonal"
	Isometric              = "isometric"
	Staggered              = "staggered"
	Hexagonal              = "hexagonal"
)

func (o *Orientation) UnmarshalXMLAttr(attr xml.Attr) error {
	*o = Orientation(attr.Value)
	if len(*o) == 0 {
		*o = Orthogonal
	}
	return nil
}

type RenderOrder string

const (
	RighDown RenderOrder = "right-down"
	RightUp              = "right-up"
	LeftDown             = "left-down"
	LeftUp               = "left-up"
)

func (o *RenderOrder) UnmarshalXMLAttr(attr xml.Attr) error {
	*o = RenderOrder(attr.Value)
	if len(*o) == 0 {
		*o = RighDown
	}
	return nil
}

type Color struct{ R, G, B, A uint8 }

func (c *Color) UnmarshalXMLAttr(attr xml.Attr) error {
	if len(attr.Value) == 0 {
		return nil
	}
	bytes, err := hex.DecodeString(attr.Value[1:])
	if err != nil {
		return err
	}
	n := len(bytes) - 1
	if n > 0 {
		c.B = bytes[n]
	}
	if n > 1 {
		c.G = bytes[n-1]
	}
	if n > 2 {
		c.R = bytes[n-2]
	}
	if n > 3 {
		c.A = bytes[n-3]
	}
	return nil
}

type TileSet struct {
	XMLName       xml.Name     `xml:"tileset"`
	FirstGlobalID int          `xml:"firstgid,attr"`
	Source        string       `xml:"source,attr"`
	Name          string       `xml:"name,attr"`
	TileWidth     int          `xml:"tilewidth,attr"`
	TileHeight    int          `xml:"tileheight,attr"`
	Spacing       int          `xml:"spacing,attr"`
	Margin        int          `xml:"margin,attr"`
	TileCount     int          `xml:"tilecount,attr"`
	Columns       int          `xml:"columns,attr"`
	Image         Image        `xml:"image"`
	TerrainTypes  TerrainTypes `xml:"terraintypes"`
	Tiles         []Tile       `xml:"tile"`
	// TODO TileOffset
}

// general TODO: all properties

type Image struct {
	XMLName xml.Name `xml:"image"`
	Source  string   `xml:"source,attr"`
	Width   int      `xml:"width,attr"`
	Height  int      `xml:"height,attr"`
}

type TerrainTypes struct {
	XMLName  xml.Name  `xml:"terraintypes"`
	Terrains []Terrain `xml:"terrain"`
}

type Terrain struct {
	XMLName xml.Name `xml:"terrain"`
	Name    string   `xml:"name,attr"`
	Tile    int      `xml:"tile,attr"`
}

type Tile struct {
	XMLName xml.Name    `xml:"tile"`
	ID      int         `xml:"id,attr"`
	Terrain TerrainList `xml:"terrain,attr"`
}

type TerrainList struct {
	Valid   bool
	Corners [4]int
}

const (
	NoTerrain = -1

	CornerTopLeft     = 0
	CornerTopRight    = 1
	CornerBottomLeft  = 2
	CornerBottomRight = 3
)

func (t *TerrainList) UnmarshalXMLAttr(attr xml.Attr) error {
	parts := strings.Split(attr.Value, ",")
	if len(parts) != 4 {
		return errors.New("invlaid terrain list, must contain four comma-separated parts: " + attr.Value)
	}
	for i, part := range parts {
		t.Corners[i] = NoTerrain

		if len(part) > 0 {
			n, err := strconv.Atoi(part)
			if err != nil {
				return err
			}
			t.Corners[i] = n
		}
	}
	t.Valid = true
	return nil
}

type Layer struct {
	XMLName xml.Name  `xml:"layer"`
	Name    string    `xml:"name,attr"`
	Width   int       `xml:"width,attr"`
	Height  int       `xml:"height,attr"`
	Data    LayerData `xml:"data"`
}

// TODO make LayerData a 2D grid and decode it when reading it, this is more
// convenient

type LayerData struct {
	XMLName     xml.Name `xml:"data"`
	Encoding    string   `xml:"encoding,attr"`
	Compression string   `xml:"compression,attr"`
	Text        string   `xml:",chardata"`
}
