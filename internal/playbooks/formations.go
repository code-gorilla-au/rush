package playbooks

import "slices"

var formations = []Formation{
	{
		Name:        "balanced-right",
		Description: "each lane has at least one player, with centre and right having two players in each lane.",
		Lane1:       1,
		Lane2:       2,
		Lane3:       2,
	},
	{
		Name:        "balanced-left",
		Description: "each lane has at least one player, with centre and left having two players in each lane.",
		Lane1:       2,
		Lane2:       2,
		Lane3:       1,
	},
	{
		Name:        "strong-right",
		Description: "A strong formation with more players on the right side.",
		Lane1:       1,
		Lane2:       1,
		Lane3:       3,
	},
	{
		Name:        "strong-left",
		Description: "A strong formation with more players on the left side.",
		Lane1:       3,
		Lane2:       1,
		Lane3:       1,
	},
	{
		Name:        "strong-centre",
		Description: "Concentrated strength in the middle lane.",
		Lane1:       1,
		Lane2:       3,
		Lane3:       1,
	},
	{
		Name:        "split-balanced",
		Description: "Players are split to stack team members on the outside lanes.",
		Lane1:       2,
		Lane2:       1,
		Lane3:       2,
	},
	{
		Name:        "split-right",
		Description: "Players are stacked on the outside lanes, leaving the centre open, with a strong presence on the right flank.",
		Lane1:       2,
		Lane2:       0,
		Lane3:       3,
	},
	{
		Name:        "split-left",
		Description: "Players are stacked on the outside lanes, leaving the centre open, with a strong presence on the left flank.",
		Lane1:       3,
		Lane2:       0,
		Lane3:       2,
	},
	{
		Name:        "overload-right",
		Description: "Maximum pressure on the right flank, leaving the left open.",
		Lane1:       0,
		Lane2:       2,
		Lane3:       3,
	},
	{
		Name:        "overload-left",
		Description: "Maximum pressure on the left flank, leaving the right open.",
		Lane1:       3,
		Lane2:       2,
		Lane3:       0,
	},
	{
		Name:        "overload-centre-left",
		Description: "Heavy presence in the centre and left, leaving the right open.",
		Lane1:       2,
		Lane2:       3,
		Lane3:       0,
	},
	{
		Name:        "overload-centre-right",
		Description: "Heavy presence in the centre and right, leaving the left open.",
		Lane1:       0,
		Lane2:       3,
		Lane3:       2,
	},
	{
		Name:        "single-lane-left",
		Description: "All players are concentrated in the left lane.",
		Lane1:       5,
		Lane2:       0,
		Lane3:       0,
	},
	{
		Name:        "single-lane-centre",
		Description: "All players are concentrated in the center lane.",
		Lane1:       0,
		Lane2:       5,
		Lane3:       0,
	},
	{
		Name:        "single-lane-right",
		Description: "All players are concentrated in the right lane.",
		Lane1:       0,
		Lane2:       0,
		Lane3:       5,
	},
}

func Formations() []Formation {
	cloned := slices.Clone(formations)
	return cloned
}
