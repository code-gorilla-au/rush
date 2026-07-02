package teams

import (
	"fmt"
	"math/rand/v2"
	"strings"
)

type TeamNameGenerator struct {
	rng *rand.Rand
}

func NewTeamNameGenerator() *TeamNameGenerator {
	// Using a fixed seed for determinism across runs if desired,
	// or we could pass one in. For now, let's just use a new RNG.
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
)

func (g *TeamNameGenerator) Generate() string {
	templates := []string{
		"The %s %s", // The {Adjective} {Noun}s
		"%s %s",     // {Place} {Noun}s
		"%s %s",     // {Adjective} {Creature}s (reusing Nouns)
		"%s of %s",  // {Noun} of {DoomThing}
		"%s-%s",     // {DockerAdjective}-{DockerNoun}
		"%s %s",     // {Profession} {PluralNoun} (reusing Nouns)
	}

	templateIdx := g.rng.IntN(len(templates))
	template := templates[templateIdx]

	switch templateIdx {
	case 0: // The {Adjective} {Noun}s
		return fmt.Sprintf(template, g.randomString(adjectives), g.randomString(teamNouns))
	case 1: // {Place} {Noun}s
		return fmt.Sprintf(template, g.randomString(places), g.randomString(teamNouns))
	case 2: // {Adjective} {Creature}s
		return fmt.Sprintf(template, g.randomString(adjectives), g.randomString(teamNouns))
	case 3: // {Noun} of {DoomThing}
		// Singularize noun for "X of Y"? Or keep plural? Example "Snacklords of Doom".
		// Let's singularize some nouns if they end in 's'.
		noun := g.randomString(teamNouns)
		if strings.HasSuffix(noun, "s") {
			noun = strings.TrimSuffix(noun, "s")
		}
		return fmt.Sprintf(template, noun, g.randomString(doomThings))
	case 4: // {DockerAdjective}-{DockerNoun}
		return fmt.Sprintf(template, g.randomString(dockerAdjectives), g.randomString(dockerNouns))
	case 5: // {Profession} {PluralNoun}
		return fmt.Sprintf(template, g.randomString(professions), g.randomString(teamNouns))
	default:
		return "Unknown Team"
	}
}

func (g *TeamNameGenerator) randomString(list []string) string {
	return list[g.rng.IntN(len(list))]
}

func (g *TeamNameGenerator) GenerateUnique(count int) []string {
	seen := make(map[string]bool)
	names := make([]string, 0, count)

	for len(names) < count {
		name := g.Generate()
		if !seen[name] {
			seen[name] = true
			names = append(names, name)
		}
	}
	return names
}
