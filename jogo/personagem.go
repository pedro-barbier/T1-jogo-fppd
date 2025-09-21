// personagem.go - Funções para movimentação e ações do personagem
package main

import (
	"fmt"
	"time"
)

// Atualiza a posição do personagem com base na tecla pressionada (WASD)
func personagemMover(jogo *Jogo, tecla rune, direcao chan (string), timeout chan struct{}, lock chan struct{}) {
	dx, dy := 0, 0
	switch tecla {
	case 'w':
		dy = -1 // Move para cima
		<-direcao
		direcao <- "N"
	case 'a':
		dx = -1 // Move para a esquerda
		<-direcao
		direcao <- "O"
	case 's':
		dy = 1 // Move para baixo
		<-direcao
		direcao <- "S"
	case 'd':
		dx = 1 // Move para a direita
		<-direcao
		direcao <- "L"
	}

	nx, ny := jogo.PosX+dx, jogo.PosY+dy
	// Verifica se o movimento é permitido e realiza a movimentação
	if jogoPodeMoverPara(jogo, nx, ny) {
		<-lock
		jogoMoverElemento(jogo, jogo.PosX, jogo.PosY, dx, dy)
		jogo.PosX, jogo.PosY = nx, ny
		lock <- struct{}{}
	}
	if jogo.UltimoVisitado == Powerup {
		<-lock
		jogo.UltimoVisitado = Vazio
		jogo.StatusMsg = "Você coletou um power-up! Vida restaurada."
		interfaceDesenharJogo(jogo)
		lock <- struct{}{}
		timeout <- struct{}{}
	}
}

// Define o que ocorre quando o jogador pressiona a tecla de interação
// Neste exemplo, apenas exibe uma mensagem de status
// Você pode expandir essa função para incluir lógica de interação com objetos
func personagemAtirar(jogo *Jogo, direcao chan (string), lock chan struct{}) {
	// Atualmente apenas exibe uma mensagem de status
	<-lock
	jogo.StatusMsg = fmt.Sprintf("Atirando em (%d, %d)", jogo.PosX, jogo.PosY)
	lock <- struct{}{}
	x, y := jogo.PosX, jogo.PosY
	x_ant, y_ant := x, y
	dir := <-direcao
	direcao <- dir

	i := 0
	for {
		switch dir {
		case "N":
			y--
			if i > 0 {
				y_ant--
			}
		case "S":
			y++
			if i > 0 {
				y_ant++
			}
		case "L":
			x++
			if i > 0 {
				x_ant++
			}
		case "O":
			x--
			if i > 0 {
				x_ant--
			}
		}

		if jogo.Mapa[y][x] == Inimigo {
			<-lock
			jogo.Mapa[y][x] = Vazio
			interfaceDesenharJogo(jogo)
			lock <- struct{}{}
			if jogo.Mapa[y_ant][x_ant] == Tiro {
				<-lock
				jogo.Mapa[y_ant][x_ant] = Vazio
				interfaceDesenharJogo(jogo)
				lock <- struct{}{}
			}
			break
		} else if jogo.Mapa[y][x] == Vazio {
			<-lock
			jogo.Mapa[y][x] = Tiro
			if jogo.Mapa[y_ant][x_ant] == Tiro {
				jogo.Mapa[y_ant][x_ant] = Vazio
			}
			interfaceDesenharJogo(jogo)
			lock <- struct{}{}

			time.Sleep(100 * time.Millisecond)
		} else {
			<-lock
			jogo.Mapa[y_ant][x_ant] = Vazio
			interfaceDesenharJogo(jogo)
			lock <- struct{}{}
			break
		}
		i++
	}
}

// Processa o evento do teclado e executa a ação correspondente
func personagemExecutarAcao(jogo *Jogo, ev EventoTeclado, direcao chan string, timeout chan struct{}, lock chan struct{}) bool {
	switch ev.Tipo {
	case "sair":
		// Retorna false para indicar que o jogo deve terminar
		return false
	case "interagir":
		// Executa a ação de interação
		go personagemAtirar(jogo, direcao, lock)

	case "mover":
		// Move o personagem com base na tecla
		personagemMover(jogo, ev.Tecla, direcao, timeout, lock)
	}
	return true // Continua o jogo
}
