package infinitegen

import (
	"azul3d.org/engine/tmx"
	"go-game/mapgen"
	"go-game/mapgen/infinitegen/forestmapper"
	"go-game/mapgen/infinitegen/types"
	"go-game/mapgen/infinitegen/noisegen"
	"os"
	"io/ioutil"
	"image"
	// "path/filepath"
	"image/draw"
	// if this is NOT imported, the image format won't be recognized
	_ "image/png"
	"fmt"
	"path/filepath"
	"azul3d.org/engine/gfx"
	"go-game/util"
)

const emptyMapPath = "data/empty.tmx"
const tilesDir = "data/tiles"

type objCoord struct {
	x int64
	z int64
}
type objMapMap = map[objCoord]*mapgen.ObjectMap

type MapLoader struct {
	noiseGen   noisegen.NoiseGen
	resolution int
	chunks     objMapMap
	mapper     types.FieldMapper
	initialMap *tmx.Map
	tsImages   map[string]*image.RGBA
}

func New(seed int64, resolution int) (*MapLoader, error) {
	chunks := make(objMapMap)
	noiseGen := noisegen.New(seed)
	var fieldMapper types.FieldMapper = forestmapper.New()
	// 1) we create an empty *Map (=: m)
	initialMap, err := getInitialMap()
	initialMap.Height = resolution
	initialMap.Width = resolution
	if err != nil {
		return nil, fmt.Errorf("Infinite generator initial map problem: %s", err)
	}
	// 4) we load the tsImages (as done in LoadFile)
	tsImages, err := loadImages(initialMap.Tilesets)
	if err != nil {
		return nil, fmt.Errorf("Infinite generator image loading problem: %s", err)
	}
	return &MapLoader{chunks:chunks, noiseGen:noiseGen, resolution:resolution, mapper: fieldMapper, initialMap: initialMap, tsImages: tsImages}, nil
}

func (m *MapLoader) GenerateMap(x int64, z int64) (mapgen.ObjectMap, error) {
	coord := objCoord{x,z}
	objMap, ok := m.chunks[coord]
	if !ok {
		chunk := newChunk(x, z, m.noiseGen, m.resolution)
		// 2) we generate the layers using a mapper
		// 3) we add the layers to m
		// Layers are always changed
		m.initialMap.Layers = m.mapper.MapFields(m.noiseGen, chunk.terrain, chunk.resolution)
		// 5) we use Load with m (*Map) and the tsImages
		t := tmx.Load(m.initialMap, nil, m.tsImages)
		fmt.Printf("Objects: %v", t)
		// reverseMeshes(t)
		objMap = &t
		m.chunks[coord] = objMap
	}

	return *objMap, nil
}

func reverseMeshes(objects map[string]map[string]*gfx.Object) {
	for _, layer := range objects {
		for _, obj := range layer {
			util.ReverseSlice(obj.Meshes[0].Vertices)
			util.ReverseSlice(obj.Meshes[0].TexCoords)
		}
	}
}

func getInitialMap() (*tmx.Map, error) {
	f, err := os.Open(emptyMapPath)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return tmx.Parse(data)
}

func loadImages(tilesets []*tmx.Tileset) (map[string]*image.RGBA, error) {
	tsImages := make(map[string]*image.RGBA)
	for _, ts := range tilesets {
		// Name of the tileset image file
		tsImage := filepath.Base(ts.Image.Source)

		// Open tileset image
		f, err := os.Open(filepath.Join(tilesDir, tsImage))
		if err != nil {
			return nil, fmt.Errorf("Infinite generator tileset opening error: %s (%s, %s)", err, ts.Image.Source, ts.Source)
		}

		// Decode the image
		src, _, err := image.Decode(f)
		if err != nil {
			return nil, fmt.Errorf("Infinite generator image decoding error: %s, %s, %s, %s", err, ts.Name, f.Name())
		}

		// If need be, convert to RGBA
		rgba, ok := src.(*image.RGBA)
		if !ok {
			b := src.Bounds()
			rgba = image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
			draw.Draw(rgba, rgba.Bounds(), src, b.Min, draw.Src)
		}

		// Put into the tileset images map
		tsImages[tsImage] = rgba
	}
	return tsImages, nil
}