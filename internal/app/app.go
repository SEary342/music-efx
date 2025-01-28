package app

import "music-efx/internal/menu"

type GlobalModel struct {
	CurrentView string
	SharedData  map[string]interface{}
}

type Model struct {
	Global *GlobalModel
	Menu   *menu.MenuModel
}
