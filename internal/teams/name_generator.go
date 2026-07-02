package teams

import (
	"math/rand/v2"
	"strings"
)

type TeamNameGenerator struct {
	rng *rand.Rand
}

func NewTeamNameGenerator() *TeamNameGenerator {
	return &TeamNameGenerator{
		rng: rand.New(rand.NewPCG(42, 42)),
	}
}

var (
	adjectives = []string{
		"Rusty", "Brutal", "Grim", "Blessed", "Cursed", "Unstable", "Feral", "Mouldy", "Turbo", "Heretical",
		"Recursive", "Screaming", "Sneaky", "Explosive", "Unwashed", "Violent", "Questionable", "Bloodied", "Rabid", "Malfunctioning",
		"Ancient", "Putrid", "Vengeful", "Drunken", "Radioactive", "Glitchy", "Homicidal", "Desperate", "Forgotten", "Toxic",
	}

	teamNouns = []string{
		"Goblins", "Ogres", "Ratlings", "Bonechewers", "Skullpunters", "Maulers", "Bruisers", "Bashers", "Fumblers", "Blitzers",
		"Misfits", "Snacklords", "Doomrunners", "Mudstompers", "Groinbiters", "Helmetlickers", "Tunnelgrinders", "Chainpullers", "Binfires", "Scraplords",
		"Bonecrushers", "Fleshweavers", "Geargrinders", "Muckrakers", "Voidwalkers", "Sludgehounds", "Ironbellies", "Gutterrunners", "Stormbringers", "Hellraisers",
	}

	places = []string{
		"Krumpgate", "Middenheap", "Ironbog", "Gutterspire", "Blackditch", "Rustvale", "Muckford", "New Grotsburg", "Ashbarrow", "Boltmarsh",
		"Grubchester", "Voidcrate", "Slagwell", "Dockerford", "Meatminster",
	}

	doomThings = []string{
		"The Final Whistle", "The Sacred Boot", "The Broken Helmet", "The Emperor's Lunchbox", "The Rusted Cog",
		"The Spiked Ball", "The Great Bin", "The Unpaid Invoice", "The Ninth Referee", "The Eternal Scrum",
		"The Leaking Barrel", "The Golden Tooth", "The Cursed Whistle", "The Dented Cup", "The Last Rations",
	}

	dockerAdjectives = []string{
		"sleepy", "angry", "clever", "feral", "hungry", "nervous", "chaotic", "brave", "grumpy", "wobbly", "sneaky", "loud",
	}

	dockerNouns = []string{
		"badger", "squid", "goblin", "troll", "weasel", "servo", "bucket", "daemon", "hamster", "squiggle", "penguin", "slug",
	}

	professions = []string{
		"Cook", "Smith", "Baker", "Butcher", "Slayer", "Reaper", "Whacker", "Crusher", "Believer", "Denier",
	}

	namePatterns = []func(*TeamNameGenerator) string{
		func(g *TeamNameGenerator) string {
			return "The " + g.pick(adjectives) + " " + g.pick(teamNouns)
		},
		func(g *TeamNameGenerator) string {
			return g.pick(places) + " " + g.pick(teamNouns)
		},
		func(g *TeamNameGenerator) string {
			return g.pick(adjectives) + " " + g.pick(teamNouns)
		},
		func(g *TeamNameGenerator) string {
			noun := g.pick(teamNouns)
			if singular, ok := strings.CutSuffix(noun, "s"); ok {
				noun = singular
			}

			return noun + " of " + g.pick(doomThings)
		},
		func(g *TeamNameGenerator) string {
			return g.pick(dockerAdjectives) + "-" + g.pick(dockerNouns)
		},
		func(g *TeamNameGenerator) string {
			return g.pick(professions) + " " + g.pick(teamNouns)
		},
	}
)

func (g *TeamNameGenerator) Generate() string {
	pattern := namePatterns[g.rng.IntN(len(namePatterns))]
	return pattern(g)
}

func (g *TeamNameGenerator) pick(list []string) string {
	return list[g.rng.IntN(len(list))]
}

func (g *TeamNameGenerator) GenerateUnique(count int) []string {
	if count <= 0 {
		return []string{}
	}

	seen := make(map[string]struct{}, count)
	names := make([]string, 0, count)

	for len(names) < count {
		name := g.Generate()
		if _, exists := seen[name]; exists {
			continue
		}

		seen[name] = struct{}{}
		names = append(names, name)
	}

	return names
}
