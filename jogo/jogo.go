// jogo.go - Funções para manipular os elementos do jogo, como carregar o mapa e mover o personagem
package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"
)

// Elemento representa qualquer objeto do mapa (parede, personagem, vegetação, etc)
type Elemento struct {
	simbolo  rune
	cor      Cor
	corFundo Cor
	tangivel bool // Indica se o elemento bloqueia passagem
}

// Jogo contém o estado atual do jogo
type Jogo struct {
	Mapa           [][]Elemento // grade 2D representando o mapa
	StatusMsg      string       // mensagem para a barra de status
	PosX, PosY     int          // posição atual do personagem
	UltimoVisitado Elemento     // elemento que estava na posição do personagem antes de mover
}

// Elementos visuais do jogo
var (
	Personagem = Elemento{'☺', CorCinzaEscuro, CorPadrao, true}
	Inimigo    = Elemento{'☠', CorVermelho, CorPadrao, true}
	Parede     = Elemento{'▤', CorParede, CorFundoParede, true}
	Parede2    = Elemento{'░', CorParede, CorPadrao, true}
	Vegetacao  = Elemento{'♣', CorVerde, CorPadrao, false}
	Powerup    = Elemento{'★', CorAmarela, CorPadrao, false}
	Coracao    = Elemento{'♥', CorVermelho, CorPadrao, true}
	Tiro       = Elemento{'✳', CorRoxa, CorPadrao, true}
	Zero       = Elemento{'0', CorTexto, CorPadrao, true}
	Vazio      = Elemento{' ', CorPadrao, CorPadrao, false}
)

// Cria e retorna uma nova instância do jogo
func jogoNovo() Jogo {
	// O ultimo elemento visitado é inicializado como vazio
	// pois o jogo começa com o personagem em uma posição vazia
	return Jogo{UltimoVisitado: Vazio}
}

// Lê um arquivo texto linha por linha e constrói o mapa do jogo
func jogoCarregarMapa(nome string, jogo *Jogo) error {
	arq, err := os.Open(nome)
	if err != nil {
		return err
	}
	defer arq.Close()

	scanner := bufio.NewScanner(arq)
	y := 0
	for scanner.Scan() {
		linha := scanner.Text()
		var linhaElems []Elemento
		for x, ch := range linha {
			e := Vazio
			switch ch {
			case Parede.simbolo:
				e = Parede
			case Parede2.simbolo:
				e = Parede2
			case Inimigo.simbolo:
				e = Inimigo
			case Vegetacao.simbolo:
				e = Vegetacao
			case Powerup.simbolo:
				e = Powerup
			case Coracao.simbolo:
				e = Coracao
			case Zero.simbolo:
				e = Zero
			case Personagem.simbolo:
				jogo.PosX, jogo.PosY = x, y // registra a posição inicial do personagem
			}
			linhaElems = append(linhaElems, e)
		}
		jogo.Mapa = append(jogo.Mapa, linhaElems)
		y++
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func jogoSpawnPowerUp(jogo *Jogo, timeout chan bool, lock chan struct{}) {
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
		fmt.Print("Pegou PowerUp")
	case <-time.After(10 * time.Second):
		<-lock
		jogo.Mapa[y][x] = Vazio
		interfaceDesenharJogo(jogo)
		lock <- struct{}{}
	}
}

func jogoSpawnInimigo(jogo *Jogo, danoConfirmation chan bool, lock chan struct{}) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	side := r.Intn(4)

	pos_spawns_y, pos_spawns_x := 0, 0
	x, y := 0, 0

	// Define a posição inicial do inimigo com base no lado sorteado
spawnLoop:
	for {
		pos_spawns_y = r.Intn(4) + 13
		pos_spawns_x = r.Intn(16) + 35
		switch side {
		case 0: // Top
			if jogo.Mapa[1][pos_spawns_x] == Vazio {
				<-lock
				jogo.Mapa[1][pos_spawns_x] = Inimigo
				interfaceDesenharJogo(jogo)
				lock <- struct{}{}
				y = 1
				x = pos_spawns_x
				break spawnLoop
			}
		case 1: // Right
			if jogo.Mapa[pos_spawns_y][81] == Vazio {
				<-lock
				jogo.Mapa[pos_spawns_y][81] = Inimigo
				interfaceDesenharJogo(jogo)
				lock <- struct{}{}
				y = pos_spawns_y
				x = 81
				break spawnLoop
			}
		case 2: // Bottom
			if jogo.Mapa[29][pos_spawns_x] == Vazio {
				<-lock
				jogo.Mapa[29][pos_spawns_x] = Inimigo
				interfaceDesenharJogo(jogo)
				lock <- struct{}{}
				y = 29
				x = pos_spawns_x
				break spawnLoop
			}
		case 3: // Left
			if jogo.Mapa[pos_spawns_y][1] == Vazio {
				<-lock
				jogo.Mapa[pos_spawns_y][1] = Inimigo
				interfaceDesenharJogo(jogo)
				lock <- struct{}{}
				y = pos_spawns_y
				x = 1
				break spawnLoop
			}
		}
	}

	// Move o inimigo em direção ao personagem
	for {
		time.Sleep(500 * time.Millisecond)
		dx, dy := 0, 0
		if jogo.PosX > x {
			dx = 1
		} else if jogo.PosX < x {
			dx = -1
		} else if jogo.PosY > y {
			dy = 1
		} else if jogo.PosY < y {
			dy = -1
		}
		if jogo.Mapa[y][x] != Inimigo {
			break
		}
		if jogo.Mapa[y+dy][x+dx] == Personagem {
			<-lock
			jogo.Mapa[y][x] = Vazio
			interfaceDesenharJogo(jogo)
			lock <- struct{}{}
			danoConfirmation <- true
			break
		}
		<-lock
		jogo.Mapa[y+dy][x+dx] = Inimigo
		jogo.Mapa[y][x] = Vazio
		interfaceDesenharJogo(jogo)
		lock <- struct{}{}
		y += dy
		x += dx
	}
}

// Verifica se o personagem pode se mover para a posição (x, y)
func jogoPodeMoverPara(jogo *Jogo, x, y int) bool {
	// Verifica se a coordenada Y está dentro dos limites verticais do mapa
	if y < 0 || y >= len(jogo.Mapa) {
		return false
	}

	// Verifica se a coordenada X está dentro dos limites horizontais do mapa
	if x < 0 || x >= len(jogo.Mapa[y]) {
		return false
	}

	// Verifica se o elemento de destino é tangível (bloqueia passagem)
	if jogo.Mapa[y][x].tangivel {
		return false
	}

	// Pode mover para a posição
	return true
}

// Move um elemento para a nova posição
func jogoMoverElemento(jogo *Jogo, x, y, dx, dy int) {
	nx, ny := x+dx, y+dy

	// Obtem elemento atual na posição
	elemento := jogo.Mapa[y][x] // guarda o conteúdo atual da posição

	jogo.Mapa[y][x] = jogo.UltimoVisitado   // restaura o conteúdo anterior
	jogo.UltimoVisitado = jogo.Mapa[ny][nx] // guarda o conteúdo atual da nova posição
	jogo.Mapa[ny][nx] = elemento            // move o elemento
}
