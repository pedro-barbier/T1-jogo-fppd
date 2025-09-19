// personagem.go - Funções para movimentação e ações do personagem
package main

import (
	"fmt"
	"time"
)

// Atualiza a posição do personagem com base na tecla pressionada (WASD)
func personagemMover(tecla rune, direcao chan (string), jogo *Jogo, lock chan struct{}) {
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
}

// Define o que ocorre quando o jogador pressiona a tecla de interação
// Neste exemplo, apenas exibe uma mensagem de status
// Você pode expandir essa função para incluir lógica de interação com objetos
func personagemAtirar(direcao chan (string), jogo *Jogo, lock chan struct{}) {
	// Atualmente apenas exibe uma mensagem de status
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
			jogo.Mapa[y][x] = Vazio
			if jogo.Mapa[y_ant][x_ant] == Tiro {
				jogo.Mapa[y_ant][x_ant] = Vazio
			}
			break
		} else if jogo.Mapa[y][x] == Vazio {

			jogo.Mapa[y][x] = Tiro
			if jogo.Mapa[y_ant][x_ant] == Tiro {
				jogo.Mapa[y_ant][x_ant] = Vazio
			}
			<-lock
			interfaceDesenharJogo(jogo)
			lock <- struct{}{}

			time.Sleep(100 * time.Millisecond)
		} else {
			jogo.Mapa[y_ant][x_ant] = Vazio
			break
		}
		i++
	}
	jogo.StatusMsg = fmt.Sprintf("Interagindo em (%d, %d)", jogo.PosX, jogo.PosY)
}

// Processa o evento do teclado e executa a ação correspondente
func personagemExecutarAcao(ev EventoTeclado, direcao chan string, jogo *Jogo, lock chan struct{}) bool {
	switch ev.Tipo {
	case "sair":
		// Retorna false para indicar que o jogo deve terminar
		return false
	case "interagir":
		// Executa a ação de interação
		go personagemAtirar(direcao, jogo, lock)

	case "mover":
		// Move o personagem com base na tecla
		personagemMover(ev.Tecla, direcao, jogo, lock)
	}
	return true // Continua o jogo
}
