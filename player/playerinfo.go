package player

import (
	"fmt"

	"github.com/cbrookscode/pokedexcli2/internal"
)

type Player struct {
	Party []internal.Pokemon
	Level int
}

func (p *Player) AddPokemonToPlayerParty(pokemon internal.Pokemon, pokedex *internal.Pokedex) error {
	// make sure pokemon to add is in pokedex
	_, exists := pokedex.Entries[pokemon.Name]
	if !exists {
		return fmt.Errorf("pokemon is not in your pokedex yet, and so cannot be added to your party")
	}

	// make sure pokemon is not already in the players party
	for _, party_pokies := range p.Party {
		if pokemon.Name == party_pokies.Name {
			return fmt.Errorf("pokemon already is in your party")
		}
	}
	// make sure party size isn't max
	if len(p.Party) < 6 {
		// adjust current stats based on player level
		p.UpdatePokemonCurrentStatsToFull(&pokemon)
		p.Party = append(p.Party, pokemon)
		return nil
	} else {
		return fmt.Errorf("your party is at max capacity (6)")
	}
}

func (p *Player) UpdatePokemonCurrentStatsToFull(pokemon *internal.Pokemon) {
	for _, statstruct := range pokemon.Stats {
		switch statstruct.Stat.Name {
		case "hp":
			pokemon.Current_stats.Hp = p.Level * statstruct.BaseStat
			pokemon.Current_health = pokemon.Current_stats.Hp
		case "attack":
			pokemon.Current_stats.Attack = p.Level * statstruct.BaseStat
		case "defense":
			pokemon.Current_stats.Defense = p.Level * statstruct.BaseStat
		case "special-attack":
			pokemon.Current_stats.Special_attack = p.Level * statstruct.BaseStat
		case "special-defense":
			pokemon.Current_stats.Special_defense = p.Level * statstruct.BaseStat
		case "speed":
			pokemon.Current_stats.Speed = p.Level * statstruct.BaseStat
		}
	}
}
