package irate

import "golang.org/x/time/rate"

var gameScenes map[string]*rate.Limiter

func init() {
	gameScenes = make(map[string]*rate.Limiter)
}

func NewRate(key string, speed rate.Limit, capacity int) *rate.Limiter {
	if d, ok := gameScenes[key]; ok {
		return d
	} else {
		gameScene := rate.NewLimiter(speed, capacity)
		gameScenes[key] = gameScene
		return gameScene
	}
}
