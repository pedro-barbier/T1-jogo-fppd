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

	damage_confirmation := make(chan bool, 5)
	heal_confirmation := make(chan bool, 5)
	gameOver := make(chan bool, 1)
	timeout := make(chan struct{}, 1)

	// Goroutine para administrar vida do personagem
	go vidaAdm(&jogo, heal_confirmation, damage_confirmation, gameOver, lock)
	// Goroutine para spawnar power-ups periodicamente
	go func() {
		for {
			time.Sleep(15 * time.Second)
			go powerUpSpawnar(&jogo, timeout, heal_confirmation, lock)
		}
	}()

	// Goroutine para spawnar inimigos periodicamente
	go func() {
		for {
			time.Sleep(2 * time.Second)
			go inimigoSpawnar(&jogo, damage_confirmation, lock)
		}
	}()

	// Loop principal de entrada
	for {
		evento := interfaceLerEventoTeclado()
		continuar := personagemExecutarAcao(&jogo, evento, direcao, timeout, lock)
		if !continuar || len(gameOver) > 0 {
			break
		}

		<-lock
		interfaceDesenharJogo(&jogo)
		lock <- struct{}{}
	}
}
