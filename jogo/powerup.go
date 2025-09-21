package main

import (
	"math/rand"
	"time"
)

func powerUpSpawnar(jogo *Jogo, timeout chan struct{}, heal_confirmation chan bool, lock chan struct{}) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	x, y := 0, 0

	for {
		y = r.Intn(30)
		x = r.Intn(83)
		if jogo.Mapa[y][x] == Vazio {
			<-lock
			jogo.Mapa[y][x] = Powerup
			interfaceDesenharJogo(jogo)
			lock <- struct{}{}
			break
		}
	}
	select {
	case <-timeout:
		heal_confirmation <- true
	case <-time.After(10 * time.Second):
		<-lock
		jogo.Mapa[y][x] = Vazio
		interfaceDesenharJogo(jogo)
		lock <- struct{}{}
	}
}
