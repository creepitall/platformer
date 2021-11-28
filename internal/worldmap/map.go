package worldmap

import (
	"github.com/creepitall/test_pixel/internal/domain"
	tmx "github.com/creepitall/test_pixel/internal/pkg/tmx"
)

var currentScenePhys []frontObjectPhys

type levelSettings struct {
	width      int
	height     int
	tileWidth  int
	tileHeight int
}

type frontObjectPhys struct {
	Min coor
	Max coor
}

type coor struct {
	X, Y float64
}

func CreateNewMap() []frontObjectPhys {
	levelmap, err := tmx.ReadFile(domain.ReturnFilePath("assets/start.tmx"))
	if err != nil {
		panic(err)
	}
	decodeMap(levelmap)

	return currentScenePhys
}

// Подсчет Y идет сверху вниз
// Подсчет X идет слева направо
// Координата нижнего левого угла - 0,0
// верхнего - 0,1344
func decodeMap(levelmap *tmx.Map) {
	ls := levelSettings{
		width:      levelmap.TileWidth * levelmap.Width,
		height:     levelmap.TileHeight * levelmap.Height,
		tileWidth:  levelmap.TileWidth,
		tileHeight: levelmap.TileHeight,
	}

	currentScenePhys = make([]frontObjectPhys, 0)
	for _, mapObjectsGr := range levelmap.ObjectGroups {
		for _, object := range mapObjectsGr.Objects {

			min := coor{
				X: float64(object.X),
				Y: float64(ls.height) - float64(object.Y),
			}
			max := coor{
				X: min.X + float64(object.Width),
				Y: min.Y - float64(object.Height),
			}

			fop := frontObjectPhys{
				Min: min,
				Max: max,
			}

			currentScenePhys = append(currentScenePhys, fop)
		}
	}

}
