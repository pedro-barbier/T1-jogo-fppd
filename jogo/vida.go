package main

import (
	"time"
)

func vidaAdm(jogo *Jogo, heal_confirmation chan bool, damage_confirmation chan bool, gameOver chan bool, lock chan struct{}) {
	vida := 3
	for {
		select {
		case cura := <-heal_confirmation:
			if cura && vida < 3 {
				vida++
				switch vida {
				case 3:
					<-lock
					jogo.Mapa[2][4] = Coracao
					jogo.StatusMsg = "Vida cheia!"
					interfaceDesenharJogo(jogo)
					lock <- struct{}{}
				case 2:
					<-lock
					jogo.Mapa[2][6] = Coracao
					jogo.StatusMsg = "Você recuperou vida!"
					interfaceDesenharJogo(jogo)
					lock <- struct{}{}
				}
			}

		case dano := <-damage_confirmation:
			if dano {
				vida--
				switch vida {
				case 2:
					<-lock
					jogo.Mapa[2][4] = CoracaoFerido
					jogo.StatusMsg = "Você tomou dano!"
					interfaceDesenharJogo(jogo)
					lock <- struct{}{}
				case 1:
					<-lock
					jogo.Mapa[2][6] = CoracaoFerido
					jogo.StatusMsg = "Você tomou dano!"
					interfaceDesenharJogo(jogo)
					lock <- struct{}{}
				case 0:
					<-lock
					jogo.Mapa[2][8] = CoracaoFerido
					jogo.StatusMsg = "Você morreu! Fim de jogo."
					interfaceDesenharJogo(jogo)
					lock <- struct{}{}
					time.Sleep(2 * time.Second)
					gameOver <- true
					return
				}
			}
		}
	}
}
