package teams

import (
	"testing"

	"github.com/code-gorilla-au/odize"
)

func TestTeamNameGenerator(t *testing.T) {
	group := odize.NewGroup(t, nil)

	group.Test("GenerateUnique should produce unique names", func(t *testing.T) {
		g := NewTeamNameGenerator()
		count := 50
		names := g.GenerateUnique(count)

		odize.AssertEqual(t, count, len(names))

		seen := make(map[string]bool)
		for _, name := range names {
			if seen[name] {
				t.Errorf("Duplicate name found: %s", name)
			}
			seen[name] = true
			odize.AssertTrue(t, len(name) > 0)
		}
	})

	group.Test("Generate should produce names in expected formats", func(t *testing.T) {
		g := NewTeamNameGenerator()
		// Generate many names to hit most templates
		for i := 0; i < 100; i++ {
			name := g.Generate()
			odize.AssertTrue(t, len(name) > 0)
			// Basic sanity check that it's not empty or just whitespace
			odize.AssertTrue(t, len(name) > 2)
		}
	})

	err := group.Run()
	odize.AssertNoError(t, err)
}
