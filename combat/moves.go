package combat

import "github.com/cbrookscode/pokedexcli2/internal"

// physical, special or status are possible damage classes. Special uses special attack and defense values, physical uses standard ones.
func SkullBash(pokemon *internal.Pokemon, enemy *internal.Pokemon) {
	power := 100
	Damage := ((2*pokemon.Level/5+2)*power*(pokemon.Current_stats.Attack/enemy.Current_stats.Defense))/50 + 2
	enemy.Current_health -= Damage
}

func FireBlast(pokemon *internal.Pokemon, enemy *internal.Pokemon) {
	power := 110
	Damage := ((2*pokemon.Level/5+2)*power*(pokemon.Current_stats.Special_attack/enemy.Current_stats.Special_defense))/50 + 2
	enemy.Current_health -= Damage
}

func MetalClaw(pokemon *internal.Pokemon, enemy *internal.Pokemon) {
	power := 50
	Damage := ((2*pokemon.Level/5+2)*power*(pokemon.Current_stats.Attack/enemy.Current_stats.Defense))/50 + 2
	enemy.Current_health -= Damage
}

func Cut(pokemon *internal.Pokemon, enemy *internal.Pokemon) {
	power := 50
	Damage := ((2*pokemon.Level/5+2)*power*(pokemon.Current_stats.Attack/enemy.Current_stats.Defense))/50 + 2
	enemy.Current_health -= Damage
}
