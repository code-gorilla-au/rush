# Rules

Game rules.

## Background

Rush is a five-player lane battle game. Rush takes some elements from risk and battle arena games.

## Match setup

- A match is ten rounds.
- At the battle selection screen, a human coach picks an opponent.
- Human coaches must select one of their created playbooks before the match starts.
- A playbook defines how the coach assigns players to lanes and how those assignments are prioritised each round.

## Round flow

Each round is resolved in this order:

1. **Reset players**
   - All players return for the new round, even if they were eliminated in the previous round.

2. **Lane assignment**
   - Coaches determine where their players go across the three lanes.
   - Assignments are made according to the selected playbook.

3. **Lane duels**
   - Players in the same lane duel and roll `1d6`. Highest score wins.
   - If dice rolls are equal, immediately re-roll.
   - The losing player is eliminated for the round.
   - Lane duels continue until no opposing players remain in that lane.

4. **Round completion and scoring**
   - The round only ends when **all lanes** have completed their duels.
   - If the total number of remaining players on each team is equal, the round is a draw.
   - The team with the higher number of remaining players gets `1` point for the round.

## Win condition

- After ten rounds, the coach with the most points wins.
- If both teams have the same number of points, the game is a draw.

## Strategic depth mechanics

To keep the game simple but add meaningful decisions, playbooks and coaches should include the following strategic layers:

### 1) Lane priority planning

- Each playbook should define a lane priority profile (for example: balanced, left-heavy, right-heavy, center-control).
- Coaches can shift commitments between rounds, but every shift creates tradeoffs in the other lanes.

### 2) Tempo vs. control choices

- Coaches should choose between:
  - **Tempo play**: stack one lane to secure fast eliminations.
  - **Control play**: spread players to contest all lanes and reduce variance.
- This creates different risk profiles across rounds instead of repeating one optimal pattern.

### 3) Tactical resources per match

- Give each coach a small tactical resource pool (for example, `2` tokens per match).
- A token can be spent on effects such as:
  - `+1` to a single duel roll (declared before rolling), or
  - one re-roll in a chosen lane.
- Limited resources force long-term planning across ten rounds.

### 4) Lane outcome value beyond score

- Keep round scoring at `1` point, but track lane-level performance for tie context and future tournament systems:
  - lanes won,
  - survivor differential,
  - clutch re-roll usage.
- These metrics improve strategy feedback without changing the core win condition.