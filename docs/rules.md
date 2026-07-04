# Rules

Game rules.

## Background

Rush is a five-player lane battle game. Rush takes some elements from risk and battle arena games.

## Match setup

- A match is ten rounds.
- At the battle selection screen, a human coach picks an opponent.
- Human coaches select a coach persona (Vanguard, Bastion, Trickster, or Wildcard).
- Human coaches must select one of their created playbooks before the match starts.
- A playbook defines how the coach assigns players to lanes and how those assignments are prioritised each round.
- Coaches equip exactly three single-use tokens from their persona's allowed list.

## Round flow

Each round is resolved in this order:

1. **Reset players**
   - All players return for the new round, even if they were eliminated in the previous round.

2. **Augmentation and Lane assignment**
   - Before each round, coaches have an option of providing an augmentation from their equipped list.
   - A coach can use a maximum of 3 augmentations per game.
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

### 3) Tactical Tokens

Coaches use tokens to influence the outcome of the match.
- Each token is **single-use** per match.
- Max **3 augmentations per game**.
- An augmentation is selected **before each round**.

#### Token Library

| Token | Timing | Effect |
| :--- | :--- | :--- |
| **Twist of Fate** | Pre-roll | Roll 2d6 and keep the highest. |
| **Overpower** | Pre-roll | Gain +1 to your roll total this duel. |
| **Hamstring** | Pre-roll | Opponent gets -1 to their roll total this duel. |
| **Jamming Signal** | Pre-roll | Cancel the opponent's declared Pre-roll token. |
| **Momentum Surge** | Pre-roll | If you won your previous duel, gain +2 this duel. |
| **Second Chance** | Reaction | Re-roll your own die once; second result replaces the first. |
| **Precision Strike** | Reaction | Add +1 to your revealed total (usually if losing by 1). |
| **Last Stand** | Resolution | Prevent your elimination; lane remains unresolved for this duel. |
| **Ice in Veins** | Resolution | Convert a tie into a win for your side. |

### 4) Coach Personas

At match start, the player selects one coach persona. Each persona restricts which tokens can be equipped for that match.

- **Vanguard Coach (Aggressive)**: Focuses on tempo and lane conversion.
  - *Allowed*: Twist of Fate, Overpower, Momentum Surge, Precision Strike, Ice in Veins.
- **Bastion Coach (Defensive)**: Focuses on attrition, denial, and mistake recovery.
  - *Allowed*: Second Chance, Hamstring, Last Stand, Jamming Signal, Ice in Veins.
- **Trickster Coach (Control)**: Focuses on information and timing traps.
  - *Allowed*: Jamming Signal, Precision Strike, Second Chance, Overpower.
- **Wildcard Coach (Flexible)**: Broad adaptability with reduced extreme effects.
  - *Allowed*: Twist of Fate, Second Chance, Overpower, Hamstring, Precision Strike, Jamming Signal, Ice in Veins.

### 5) Lane outcome value beyond score

- Keep round scoring at `1` point, but track lane-level performance for tie context and future tournament systems:
  - lanes won,
  - survivor differential,
  - clutch re-roll usage.
- These metrics improve strategy feedback without changing the core win condition.