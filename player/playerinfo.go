package player

import (
	"fmt"

	"github.com/cbrookscode/pokedexcli2/internal"
)

type Player struct {
	Party []internal.Pokemon
	Level int
}

func (p *Player) AddPokemonToPlayerParty(pokemon internal.Pokemon) {
	if len(p.Party) < 6 {
		p.Party = append(p.Party, pokemon)
		fmt.Printf("%s has been added to your party!\n", pokemon.Name)
	} else {
		fmt.Println("Your party is at max capacity (6).")
	}
}
