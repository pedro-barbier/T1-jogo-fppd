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

	lock := make(chan struct{}, 1)
	lock <- struct{}{}

	direcao := make(chan string, 1)
	direcao <- "Default"

	lim := make(chan struct{}, 1)
	timeout := make(chan bool)

	go func() {
		for {
			time.Sleep(15 * time.Second)
			lim <- struct{}{}
		}
	}()
	// Loop principal de entrada
	for {
		evento := interfaceLerEventoTeclado()
		continuar := personagemExecutarAcao(evento, direcao, &jogo, lock)
		if !continuar {
			break
		}

		select {
		case <-lim:
			go jogoSpawnPowerUp(&jogo, timeout, lock)
		default:
		}

		<-lock
		interfaceDesenharJogo(&jogo)
		lock <- struct{}{}
	}
}
