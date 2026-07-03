### Proposal: Token-Only Tactical Playtest (Rules-Ready Draft)

This revision removes `Attack/Defense` and `Rush Lane` mechanics for now.
The playtest focus is strictly on high-impact token decisions and coach identity through token access.

---

### 1. Playtest Scope and Core Rules

#### Scope
*   This phase tests only token-driven decision making.
*   No `ATK/DEF` stats.
*   No `Rush Lane` declarations.

#### Duel Flow (Token Windows)
1. **Pre-roll Window**: Coaches may declare one Pre-roll token.
2. **Roll**: Both players roll simultaneously.
3. **Reveal**: Results are revealed.
4. **Reaction Window**: Eligible Reaction tokens may be declared.
5. **Resolution Window**: Eligible Resolution tokens may be declared.
6. **Lane Update**: Apply final outcome for the duel.

#### Global Token Constraints
*   Each coach equips exactly **3 tokens** per match.
*   Each equipped token is **single-use**.
*   Max **1 token per coach per duel**.
*   Both coaches may use a token in the same duel.
*   Tokens do not stack for the same coach in a single duel.

---

### 2. Token Library (Expanded for User Testing)

#### Twist of Fate (Advantage)
*   **Timing**: Pre-roll Window.
*   **Effect**: Roll `2d6` and keep the highest.
*   **Intent**: Front-load pressure in must-win lanes.

#### Second Chance (Reactive Re-roll)
*   **Timing**: Reaction Window, only if your first roll did not win.
*   **Effect**: Re-roll your own die once; second result replaces the first.
*   **Intent**: Stabilize critical moments after a miss or tie.

#### Power Play (Flat Boost)
*   **Timing**: Pre-roll Window.
*   **Effect**: Gain `+1` to your roll total this duel.
*   **Intent**: Reliable low-variance push.

#### Brace (Damage Control)
*   **Timing**: Pre-roll Window.
*   **Effect**: Opponent gets `-1` to their roll total this duel.
*   **Intent**: Defensive denial and tempo slowdown.

#### Precision Strike (Near-Miss Conversion)
*   **Timing**: Reaction Window, only if you lost by exactly `1`.
*   **Effect**: Add `+1` to your revealed total.
*   **Intent**: Skill-expression token for close reads.

#### Jamming Signal (Pre-roll Counter)
*   **Timing**: Pre-roll Window.
*   **Effect**: Cancel the opponent's declared Pre-roll token.
*   **Intent**: Anti-pattern counterplay.
*   **Constraint**: Cannot cancel Reaction or Resolution tokens.

#### Last Stand (Resolution Save)
*   **Timing**: Resolution Window, only if you would lose the duel.
*   **Effect**: Prevent your elimination this duel; lane remains unresolved.
*   **Intent**: Comeback insurance for high-value lanes.

#### Momentum Surge (Win Streak Convert)
*   **Timing**: Pre-roll Window, only if you won your previous duel.
*   **Effect**: Gain `+2` this duel.
*   **Intent**: Snowball option with explicit condition gate.

#### Ice in Veins (Tie Breaker)
*   **Timing**: Resolution Window, only on a tie.
*   **Effect**: Convert tie into a win for your side.
*   **Intent**: Tie-state control and clutch finish potential.

#### Smoke Screen (Information Denial)
*   **Timing**: Pre-roll Window.
*   **Effect**: Your token declaration remains hidden until after reveal.
*   **Intent**: Mind-game tool to punish reactive opponents.

---

### 3. Coach Personas (Token Access Limits)

At match start, the player selects one coach persona.
Each persona restricts which tokens can be equipped for that match.

#### Coach Loadout Rule
*   Equip exactly **3 single-use tokens** from your persona's allowed list.
*   No duplicate token names in the same loadout.

#### Persona A: Vanguard Coach (Aggressive)
*   **Allowed Tokens**: `Twist of Fate`, `Power Play`, `Momentum Surge`, `Precision Strike`, `Ice in Veins`
*   **Playstyle**: Tempo-first and lane conversion pressure.
*   **Restriction Theme**: No hard defensive save tokens.

#### Persona B: Bastion Coach (Defensive)
*   **Allowed Tokens**: `Second Chance`, `Brace`, `Last Stand`, `Jamming Signal`, `Ice in Veins`
*   **Playstyle**: Attrition, denial, and mistake recovery.
*   **Restriction Theme**: No high-spike momentum token.

#### Persona C: Trickster Coach (Control)
*   **Allowed Tokens**: `Jamming Signal`, `Smoke Screen`, `Precision Strike`, `Second Chance`, `Power Play`
*   **Playstyle**: Information warfare and timing traps.
*   **Restriction Theme**: No direct resolution-save token.

#### Persona D: Wildcard Coach (Flexible)
*   **Allowed Tokens**: Any token except `Momentum Surge` and `Last Stand`.
*   **Playstyle**: Broad adaptability with reduced extreme effects.
*   **Restriction Theme**: No strongest snowball/safety endpoints.

---

### 4. Playtest Objectives

Track these outcomes to evaluate whether token mechanics are fun and understandable:
*   **Token usage rate** by token and by persona.
*   **Round swing rate** (how often a token changes the duel winner).
*   **Persona pick rate** and win rate spread.
*   **Perceived fairness** (post-match user rating).
*   **Clarity score** (players can explain what each used token did).
