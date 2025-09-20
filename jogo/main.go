// main.go - Loop principal do jogo
package main

import (
	"os"
	"time"
)

func main() {
	// Inicializa a interface (termbox)
	interfaceIniciar()
	defer interfaceFinalizar()

	// Usa "mapa.txt" como arquivo padrão ou lê o primeiro argumento
	mapaFile := "mapa.txt"
	if len(os.Args) > 1 {
		mapaFile = os.Args[1]
	}

	// Inicializa o jogo
	jogo := jogoNovo()
	if err := jogoCarregarMapa(mapaFile, &jogo); err != nil {
		panic(err)
	}

	// Desenha o estado inicial do jogo
	interfaceDesenharJogo(&jogo)

	// Canais para sincronização e comunicação entre goroutines
	lock := make(chan struct{}, 1)
	lock <- struct{}{}

	direcao := make(chan string, 1)
	direcao <- "Default"

	limPowerup := make(chan struct{}, 1)
	limInimigo := make(chan struct{}, 1)
	dano_confirmation := make(chan bool, 1)
	timeout := make(chan bool)

	// Goroutine para spawnar power-ups periodicamente
	go func() {
		for {
			time.Sleep(15 * time.Second)
			limPowerup <- struct{}{}
		}
	}()

	go func() {
		for {
			time.Sleep(5 * time.Second)
			limInimigo <- struct{}{}
		}
	}()
	// Loop principal de entrada
	for {
		evento := interfaceLerEventoTeclado()
		continuar := personagemExecutarAcao(&jogo, evento, direcao, lock)
		if !continuar {
			break
		}

		select {
		case <-limPowerup:
			go jogoSpawnPowerUp(&jogo, timeout, lock)
		case <-limInimigo:
			go jogoSpawnInimigo(&jogo, dano_confirmation, lock)
		default:
		}

		<-lock
		interfaceDesenharJogo(&jogo)
		lock <- struct{}{}
	}
}
