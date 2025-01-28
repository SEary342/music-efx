package app

import (
	"music-efx/internal/menu"
	"music-efx/internal/player"
)

type GlobalModel struct {
	CurrentView string
	SharedData  map[string]interface{}
}

type Model struct {
	Global *GlobalModel
	Menu   *menu.MenuModel
	Player *player.PlayerModel
}
