### Proposal: Advanced Tactical Mechanics (Rules-Ready Draft)

This proposal introduces two strategic layers to **Rush**: tactical tokens and Attack/Defense stats.
This revision defines exact timing, constraints, and probability framing for simultaneous resolution.

---

### 1. Tactical Resources (Single-Use Tokens)

Each coach receives exactly two tokens per game:
*   `Second Chance` ×1
*   `Twist of Fate` ×1

Token constraints:
*   Max **1 token per coach per duel**.
*   Tokens do **not** stack on the same check.
*   Both coaches may use a token in the same duel.
*   If both coaches use tokens, both effects apply using the duel sequence below.

#### Second Chance (Reactive Re-roll)
*   **Timing**: Declare in the Reaction Window, after the first reveal, only if your check missed.
*   **Effect**: Re-roll your own attack check once. The second result replaces the first.
*   **Strategy**: A "Safety Net" for high-value lanes when the first roll fails.
*   **Probability Note**: In a baseline 50% hit check, one re-roll on miss raises hit chance to **75%**.

#### Twist of Fate (Advantage)
*   **Timing**: Declare in the Pre-roll Window.
*   **Effect**: Roll `2d6` for your attack check and keep the highest die.
*   **Strategy**: A "Power Surge." Best for aggressive tempo play to secure an early elimination and gain momentum.
*   **Probability Note**: In a baseline 50% hit check, advantage raises hit chance to **75%**.

#### Additional Token Mechanics (Expansion Candidates)
These are optional candidates for future playtests and are **not** part of the current baseline token set.
To preserve pacing, test with **2 tokens per coach total** by replacing an existing token, not adding extra total uses.

##### Power Play (Flat Attack Boost)
*   **Timing**: Declare in the Pre-roll Window.
*   **Effect**: Gain `+1 ATK` for this cycle only.
*   **Strategy**: Reliable pressure tool; lower variance than advantage.
*   **Probability Note**: In a baseline 50% hit check, `+1 ATK` raises hit chance to **66%**.

##### Brace (Temporary Defense Boost)
*   **Timing**: Declare in the Pre-roll Window.
*   **Effect**: Gain `+1 DEF` for this cycle only.
*   **Strategy**: Defensive stabilization in high-value lanes; strongest as a denial token.
*   **Probability Note**: Against a baseline 50% enemy hit check, `+1 DEF` lowers enemy hit chance to **33%**.

##### Precision Strike (Near-Miss Conversion)
*   **Timing**: Declare in the Reaction Window, after reveal, only if your total missed by exactly `1`.
*   **Effect**: Add `+1` to your revealed total (turning that near miss into a hit).
*   **Strategy**: High-information token that rewards patient use and matchup awareness.

##### Jamming Signal (Token Counter)
*   **Timing**: Declare in the Pre-roll Window.
*   **Effect**: If opponent declared a Pre-roll token this cycle, cancel that token's effect.
*   **Strategy**: Counter-tempo tool that punishes predictable token usage.
*   **Constraint Note**: Cannot cancel Reaction Window tokens.

##### Last Stand (Elimination Prevention)
*   **Timing**: Declare in Resolution, after outcomes are known, only if your player would be eliminated and the opponent would survive.
*   **Effect**: Your player is not eliminated this cycle; lane remains unresolved and proceeds under normal stalemate rules.
*   **Strategy**: Comeback tool to preserve board presence in a critical lane.
*   **Constraint Note**: Cannot prevent elimination in `both eliminated` outcomes.

#### Duel Resolution Sequence (Authoritative)
1. **Pre-roll Window**: Each coach may declare `Twist of Fate`.
2. **Initial Check**: Both players roll attack checks simultaneously (`1d6 + ATK`, or advantage if declared).
3. **Reveal**: Both totals are revealed.
4. **Reaction Window**: Each coach that missed may declare `Second Chance` (if available and no token already used by that coach in this duel).
5. **Re-roll Step**: Any declared `Second Chance` re-rolls are made and replace prior totals.
6. **Resolution**: Apply simultaneous hit outcomes (both hit / one hits / neither hits).
7. **Lane Update**: Apply eliminations and lane state.

---

### 2. Attack / Defense (AC Style) Stats

To add depth to roster building, players are assigned **Attack (ATK)** and **Defense (DEF)** ratings.

#### The Mechanic
*   **Attack Check**: Attacker Rolls `1d6 + ATK`.
*   **Hit Condition**: If the total is **equal to or greater than** the opponent's **DEF**, the opponent is eliminated.
*   **Simultaneous Resolution**: Both players roll at the same time.
    *   **Both Hit**: Both are eliminated.
    *   **One Hits**: The winner survives, the loser is eliminated.
    *   **Neither Hits**: Stalemate.

#### Stalemate Rule (Global)
To keep pacing deterministic and avoid endless loops:
*   A duel can run for up to **3 cycles total** (initial cycle + up to 2 stalemate re-cycles).
*   If still unresolved after cycle 3, both players remain in lane and the lane state is `Contested` for the round.

#### Player Archetypes

| Archetype | ATK | DEF | Profile |
| :--- | :---: | :---: | :--- |
| **Tank** | -1 | 5 | Hard to kill (33% hit chance), but low threat (33% hit chance). |
| **Standard** | +0 | 4 | Balanced (50% hit chance). The baseline for all lanes. |
| **Striker** | +1 | 3 | Glass Cannon (83% hit chance against strikers), high threat (66%+ hit chance). |

#### Single-Check Hit Probabilities (`1d6 + ATK >= DEF`)

| Attacker ↓ vs Defender → | Tank (DEF 5) | Standard (DEF 4) | Striker (DEF 3) |
| :--- | :---: | :---: | :---: |
| **Tank** (ATK -1) | 16% | 33% | 50% |
| **Standard** (ATK 0) | 33% | 50% | 66% |
| **Striker** (ATK +1) | 50% | 66% | 83% |

These are **hit chances per check**, not duel win rates.

#### Simultaneous Outcome Probabilities (Single Cycle)
Given `pA` = attacker hit chance and `pD` = defender hit chance:
*   Attacker survives: `pA * (1 - pD)`
*   Defender survives: `pD * (1 - pA)`
*   Both eliminated: `pA * pD`
*   No elimination: `(1 - pA) * (1 - pD)`

Representative archetype pair outcomes (before stalemate re-cycles):

| Pairing | A survives | B survives | Both eliminated | No elimination |
| :--- | :---: | :---: | :---: | :---: |
| Tank vs Standard (`p=33%` each) | 22% | 22% | 11% | 44% |
| Tank vs Striker (`p=50%` each) | 25% | 25% | 25% | 25% |
| Standard vs Striker (`p=66%` each) | 22% | 22% | 44% | 11% |

---

### 3. Tactical Implications

*   **Tank Stalling**: Place Tanks in low-priority lanes to increase unresolved/contested outcomes.
*   **Tempo Push**: Use Strikers with `Twist of Fate` in must-win lanes for early pressure.
*   **Risk Management**: Save `Second Chance` for lanes with high strategic value where a miss would swing the round.
