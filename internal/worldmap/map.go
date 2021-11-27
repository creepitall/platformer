package worldmap

import (
	"fmt"
	"os"

	"github.com/creepitall/test_pixel/internal/domain"
	tmx "github.com/creepitall/test_pixel/internal/pkg/tmx"
)


type levelSettings struct {
	width int
	height int
	tileWidth int
	tileHeight int
	layers layerSettings
}

type layerSettings {
	idTles int	
	isEmpty true
	vec coor
}

type coor struct {
	x, y float64
}

func CreateNewMap() {
	levelmap, err := tmx.ReadFile(domain.ReturnFilePath("assets/start.tmx"))
	if err != nil {
		panic(err)	
	}
}

// Подсчет Y идет сверху вниз 
// Подсчет X идет слева направо
func decodeMap(levelmap *tmx.Map) {
	ls := levelSettings {
		width: levelmap.Width,
		height: levelmap.height,
		tileWidth: levelmap.TileWidth,
		tileHeight: levelmap.TileHeight,
	}

	var w, h int = 1, 1
	for id, layer := range levelmap.Layers {
		с := coor {
			x: w * float64(ls.tileWidth),
			y: y * float64(ls.tileHeight),
		}

		layers := layerSettings {

		}

		ls.layers = append(ls.layers, layers)
	}

}
