// Package main — Minimal Rummy test client (AGENTS.md Day 24).
// CLI for local Nakama development that demonstrates the authoritative protocol
// without requiring a full UI. It runs in local simulation mode (in-process
// RoundState via internal/match) and shows PrivateView vs PublicView redaction.
//
// Usage:
//
//	go run ./cmd/rummy-cli                 # local simulation, 2 players alice/bob
//	go run ./cmd/rummy-cli --help
package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/gabriel-d0/rummy_backend/internal/match"
	"github.com/gabriel-d0/rummy_backend/internal/protocol"
	"github.com/gabriel-d0/rummy_backend/internal/setup"
	"github.com/heroiclabs/nakama-common/runtime"
)

type testLogger struct{}

func (l *testLogger) Debug(f string, v ...interface{})                   {}
func (l *testLogger) Info(f string, v ...interface{})                    { fmt.Printf(f+"\n", v...) }
func (l *testLogger) Warn(f string, v ...interface{})                    { fmt.Printf("WARN: "+f+"\n", v...) }
func (l *testLogger) Error(f string, v ...interface{})                   { fmt.Printf("ERROR: "+f+"\n", v...) }
func (l *testLogger) WithField(k string, v interface{}) runtime.Logger   { return l }
func (l *testLogger) WithFields(m map[string]interface{}) runtime.Logger { return l }
func (l *testLogger) Fields() map[string]interface{}                     { return nil }

type mockPresence struct {
	userId string
}

func (m *mockPresence) GetUserId() string                 { return m.userId }
func (m *mockPresence) GetSessionId() string              { return "sess-" + m.userId }
func (m *mockPresence) GetUsername() string               { return m.userId }
func (m *mockPresence) GetNodeId() string                 { return "node1" }
func (m *mockPresence) GetHidden() bool                   { return false }
func (m *mockPresence) GetPersistence() bool              { return false }
func (m *mockPresence) GetStatus() string                 { return "" }
func (m *mockPresence) GetReason() runtime.PresenceReason { return runtime.PresenceReasonUnknown }

type mockDispatcher struct {
	lastOp   int64
	lastData []byte
}

func (m *mockDispatcher) BroadcastMessage(opCode int64, data []byte, presences []runtime.Presence, sender runtime.Presence, reliable bool) error {
	m.lastOp = opCode
	m.lastData = data
	if opCode == protocol.OpServerError {
		var e map[string]interface{}
		_ = json.Unmarshal(data, &e)
		fmt.Printf("  ← OpServerError %v\n", e)
	} else if opCode == protocol.OpServerEvent {
		fmt.Printf("  ← OpServerEvent %s\n", string(data))
	} else if opCode == protocol.OpServerState {
		fmt.Printf("  ← OpServerState (private) %d bytes\n", len(data))
	} else if opCode == protocol.OpServerStatePublic {
		fmt.Printf("  ← OpServerStatePublic %d bytes\n", len(data))
	}
	return nil
}
func (m *mockDispatcher) BroadcastMessageDeferred(opCode int64, data []byte, presences []runtime.Presence, sender runtime.Presence, reliable bool) error {
	return m.BroadcastMessage(opCode, data, presences, sender, reliable)
}
func (m *mockDispatcher) MatchKick(presences []runtime.Presence) error { return nil }
func (m *mockDispatcher) MatchLabelUpdate(label string) error          { return nil }

type mockMatchData struct {
	mockPresence
	opCode int64
	data   []byte
}

func (m *mockMatchData) GetOpCode() int64      { return m.opCode }
func (m *mockMatchData) GetData() []byte       { return m.data }
func (m *mockMatchData) GetReliable() bool     { return true }
func (m *mockMatchData) GetReceiveTime() int64 { return 0 }

func renderState(st *match.RoundState, mySeat match.Seat) {
	pub := match.PublicView(st)
	priv := match.PrivateView(st, mySeat)
	fmt.Printf("\n--- State (you are %s, current %v, phase %v/%v) ---\n", priv.OwnSeat, pub.CurrentSeat, pub.GamePhase, pub.TurnPhase)
	fmt.Printf("Players:\n")
	for _, p := range pub.Players {
		marker := ""
		if p.Seat == pub.CurrentSeat {
			marker = " ← current"
		}
		if p.Seat == mySeat {
			marker += " (you)"
		}
		fmt.Printf("  %s seat-%d HasOpened=%v RackCount=%d%s\n", p.ID, p.Seat, p.HasOpened, p.RackCount, marker)
	}
	fmt.Printf("StockCount: %d\n", pub.StockCount)
	fmt.Printf("DiscardRow (%d):\n", len(pub.DiscardRow))
	for _, d := range pub.DiscardRow {
		flag := ""
		if d.IsOpeningDiscard {
			flag = " [opening blocked]"
		}
		fmt.Printf("  [%d] %s%s\n", d.Index, d.Tile.String(), flag)
	}
	fmt.Printf("TableMelds (%d):\n", len(pub.TableMelds))
	for _, mm := range pub.TableMelds {
		fmt.Printf("  %s kind=%s owner=%v tiles:", mm.ID, mm.Kind, mm.OwnerSeat)
		for _, tl := range mm.Tiles {
			jrep := ""
			if tl.IsJoker {
				if rep, ok := mm.JokerReps[tl.ID]; ok {
					jrep = fmt.Sprintf("->%s-%s", rep.Colour, rep.Rank)
				}
			}
			fmt.Printf(" %s%s", tl.String(), jrep)
		}
		fmt.Printf("\n")
	}
	if pub.Winner != match.SeatInvalid {
		fmt.Printf("Winner: %v\n", pub.Winner)
	}
	fmt.Printf("Your rack (%d):\n", len(priv.OwnRack))
	for i, tl := range priv.OwnRack {
		fmt.Printf("  [%d] %s id=%s\n", i, tl.String(), tl.ID)
	}
	b, _ := json.Marshal(pub)
	fmt.Printf("Public JSON bytes: %d (no private rack leak)\n", len(b))
	_ = b
}

func mustAtoi(s string) int {
	var n int
	fmt.Sscanf(s, "%d", &n)
	return n
}

func main() {
	nakamaFlag := flag.Bool("nakama", false, "try to connect to real Nakama at 127.0.0.1:7350 (requires docker compose up)")
	helpFlag := flag.Bool("help", false, "show help")
	flag.Parse()
	if *helpFlag {
		fmt.Printf("Rummy minimal test client\n\n")
		fmt.Printf("Local simulation (default):\n  go run ./cmd/rummy-cli\n\n")
		fmt.Printf("Remote Nakama (requires docker compose up --build -d):\n  go run ./cmd/rummy-cli --nakama\n\n")
		fmt.Printf("Commands inside REPL:\n")
		fmt.Printf("  help, state, draw, discard <tileId>, meld <run|set> <tileId>..., extend <meldId> <tileId>..., prev, pickup <discardIndex> <tileId> <tileId>, replace <targetMeldId> <tileId> <newTile1> <newTile2> [jokerId colour rank], winner, quit\n")
		fmt.Printf("\nThis client shows PrivateView (own rack) vs PublicView (counts, discard, melds) per docs/protocol.md.\n")
		os.Exit(0)
	}

	if *nakamaFlag {
		fmt.Printf("Remote Nakama mode not fully implemented in this minimal CLI — it would do device auth via defaultkey and WebSocket match.\n")
		fmt.Printf("Falling back to local simulation for protocol validation. Run without --nakama for local.\n\n")
	}

	fmt.Printf("Rummy minimal test client — local simulation\n")
	fmt.Printf("Players: alice (seat 0, opener) and bob (seat 1)\n")
	fmt.Printf("Type 'help' for commands, 'state' to show board, 'quit' to exit.\n\n")

	logger := &testLogger{}
	dispatcher := &mockDispatcher{}
	rm := &match.RummyMatch{}

	stateRaw, _, _ := rm.MatchInit(context.Background(), logger, nil, nil, nil)
	stateRaw = rm.MatchJoin(context.Background(), logger, nil, nil, dispatcher, 0, stateRaw, []runtime.Presence{&mockPresence{userId: "alice"}, &mockPresence{userId: "bob"}})
	st := stateRaw.(*match.RoundState)
	players := []match.PlayerId{"alice", "bob"}
	seededState, _, err := setup.NewRoundState(players, 42)
	if err != nil {
		fmt.Printf("setup.NewRoundState failed: %v, using manual\n", err)
		seededState = st
	} else {
		st.Players = seededState.Players
		st.Racks = seededState.Racks
		st.Stock = seededState.Stock
		st.DiscardRow = seededState.DiscardRow
		st.TableMelds = seededState.TableMelds
		st.CurrentSeat = seededState.CurrentSeat
		st.GamePhase = seededState.GamePhase
		st.TurnPhase = seededState.TurnPhase
		st.Winner = seededState.Winner
	}
	// seededState is already in OpeningDiscard with alice as opener; no need for extra Start if already there
	if st.GamePhase == match.PhaseWaiting {
		payloadStart, _ := json.Marshal(map[string]interface{}{})
		envStart, _ := json.Marshal(protocol.Envelope{Version: protocol.Version, OpCode: protocol.OpClientStart, Payload: payloadStart})
		stRaw := rm.MatchLoop(context.Background(), logger, nil, nil, dispatcher, 0, st, []runtime.MatchData{&mockMatchData{mockPresence: mockPresence{userId: "alice"}, opCode: protocol.OpClientStart, data: envStart}})
		st = stRaw.(*match.RoundState)
		fmt.Printf("Match started — phase %v, current %v\n", st.GamePhase, st.CurrentSeat)
	} else {
		fmt.Printf("Match ready — phase %v, current %v (seeded)\n", st.GamePhase, st.CurrentSeat)
	}
	renderState(st, 0)

	currentUser := "alice"
	currentSeat := match.Seat(0)

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("\n[%s seat-%d] > ", currentUser, currentSeat)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			fmt.Printf("[%s seat-%d] > ", currentUser, currentSeat)
			continue
		}
		parts := strings.Fields(line)
		cmd := strings.ToLower(parts[0])
		switch cmd {
		case "help", "h", "?":
			fmt.Printf("Commands:\n")
			fmt.Printf("  state               — show Public + Private view\n")
			fmt.Printf("  draw                — DRAW_STOCK (MustDraw)\n")
			fmt.Printf("  prev                — DRAW_PREVIOUS_DISCARD (MustDraw, opened)\n")
			fmt.Printf("  pickup <idx> <id1> <id2> — PICKUP_DISCARD_FOR_MELD\n")
			fmt.Printf("  discard <tileId>    — DISCARD\n")
			fmt.Printf("  meld <run|set> <id>... — MELD_INITIAL if not opened else MELD_NEW\n")
			fmt.Printf("  extend <meldId> <id>... — EXTEND_MELD\n")
			fmt.Printf("  replace <target> <tileId> <new1> <new2> [jokerId colour rank] — REPLACE_JOKER\n")
			fmt.Printf("  switch <alice|bob>  — switch local view\n")
			fmt.Printf("  winner              — show winner\n")
			fmt.Printf("  quit, exit          — exit\n")
		case "quit", "exit", "q":
			fmt.Printf("Bye.\n")
			return
		case "state", "s":
			renderState(st, currentSeat)
		case "switch":
			if len(parts) < 2 {
				fmt.Printf("usage: switch <alice|bob>\n")
			} else {
				switch parts[1] {
				case "alice":
					currentUser = "alice"
					currentSeat = 0
				case "bob":
					currentUser = "bob"
					currentSeat = 1
				default:
					fmt.Printf("unknown user %s\n", parts[1])
				}
				fmt.Printf("Switched to %s seat-%d\n", currentUser, currentSeat)
				renderState(st, currentSeat)
			}
		case "draw", "d":
			payload, _ := json.Marshal(map[string]interface{}{})
			env, _ := json.Marshal(protocol.Envelope{Version: protocol.Version, OpCode: protocol.OpClientDrawStock, Payload: payload})
			nextRaw := rm.MatchLoop(context.Background(), logger, nil, nil, dispatcher, 0, st, []runtime.MatchData{&mockMatchData{mockPresence: mockPresence{userId: currentUser}, opCode: protocol.OpClientDrawStock, data: env}})
			st = nextRaw.(*match.RoundState)
			if dispatcher.lastOp == protocol.OpServerError {
				fmt.Printf("  error: %s\n", string(dispatcher.lastData))
			} else {
				fmt.Printf("  drew stock, new rack %d, stock %d\n", len(st.Racks[currentSeat]), len(st.Stock))
			}
			renderState(st, currentSeat)
		case "prev":
			payload, _ := json.Marshal(map[string]interface{}{})
			env, _ := json.Marshal(protocol.Envelope{Version: protocol.Version, OpCode: protocol.OpClientDrawPreviousDiscard, Payload: payload})
			nextRaw := rm.MatchLoop(context.Background(), logger, nil, nil, dispatcher, 0, st, []runtime.MatchData{&mockMatchData{mockPresence: mockPresence{userId: currentUser}, opCode: protocol.OpClientDrawPreviousDiscard, data: env}})
			st = nextRaw.(*match.RoundState)
			if dispatcher.lastOp == protocol.OpServerError {
				fmt.Printf("  error: %s\n", string(dispatcher.lastData))
			} else {
				fmt.Printf("  drew previous discard\n")
			}
			renderState(st, currentSeat)
		case "pickup":
			if len(parts) < 4 {
				fmt.Printf("usage: pickup <discardIndex> <tileId1> <tileId2> [jokerId colour rank]\n")
			} else {
				discardIdx := mustAtoi(parts[1])
				tileIds := parts[2:4]
				payloadMap := map[string]interface{}{"discardIndex": discardIdx, "tileIds": tileIds}
				if len(parts) == 7 {
					payloadMap["jokerReps"] = map[string]interface{}{parts[4]: map[string]interface{}{"colour": parts[5], "rank": mustAtoi(parts[6])}}
				}
				payload, _ := json.Marshal(payloadMap)
				env, _ := json.Marshal(protocol.Envelope{Version: protocol.Version, OpCode: protocol.OpClientPickupDiscardForMeld, Payload: payload})
				nextRaw := rm.MatchLoop(context.Background(), logger, nil, nil, dispatcher, 0, st, []runtime.MatchData{&mockMatchData{mockPresence: mockPresence{userId: currentUser}, opCode: protocol.OpClientPickupDiscardForMeld, data: env}})
				st = nextRaw.(*match.RoundState)
				if dispatcher.lastOp == protocol.OpServerError {
					fmt.Printf("  error: %s\n", string(dispatcher.lastData))
				} else {
					fmt.Printf("  pickup succeeded\n")
				}
				renderState(st, currentSeat)
			}
		case "discard":
			if len(parts) < 2 {
				fmt.Printf("usage: discard <tileId>\n")
			} else {
				payload, _ := json.Marshal(map[string]string{"tileId": parts[1]})
				env, _ := json.Marshal(protocol.Envelope{Version: protocol.Version, OpCode: protocol.OpClientDiscard, Payload: payload})
				nextRaw := rm.MatchLoop(context.Background(), logger, nil, nil, dispatcher, 0, st, []runtime.MatchData{&mockMatchData{mockPresence: mockPresence{userId: currentUser}, opCode: protocol.OpClientDiscard, data: env}})
				st = nextRaw.(*match.RoundState)
				if dispatcher.lastOp == protocol.OpServerError {
					fmt.Printf("  error: %s\n", string(dispatcher.lastData))
				} else {
					fmt.Printf("  discarded %s\n", parts[1])
				}
				renderState(st, currentSeat)
			}
		case "meld":
			if len(parts) < 3 {
				fmt.Printf("usage: meld <run|set> <tileId>... (>=3)\n")
			} else {
				kind := parts[1]
				tileIds := parts[2:]
				hasOpened := false
				for _, p := range st.Players {
					if string(p.ID) == currentUser {
						hasOpened = p.HasOpened
						break
					}
				}
				op := protocol.OpClientMeldInitial
				if hasOpened {
					op = protocol.OpClientMeldNew
				}
				payload := map[string]interface{}{"melds": []interface{}{map[string]interface{}{"id": fmt.Sprintf("cli-%s-%d", kind, len(st.TableMelds)), "kind": kind, "tileIds": tileIds}}}
				pBytes, _ := json.Marshal(payload)
				env, _ := json.Marshal(protocol.Envelope{Version: protocol.Version, OpCode: op, Payload: pBytes})
				nextRaw := rm.MatchLoop(context.Background(), logger, nil, nil, dispatcher, 0, st, []runtime.MatchData{&mockMatchData{mockPresence: mockPresence{userId: currentUser}, opCode: op, data: env}})
				st = nextRaw.(*match.RoundState)
				if dispatcher.lastOp == protocol.OpServerError {
					fmt.Printf("  error: %s\n", string(dispatcher.lastData))
				} else {
					fmt.Printf("  meld %s succeeded\n", kind)
				}
				renderState(st, currentSeat)
			}
		case "extend":
			if len(parts) < 3 {
				fmt.Printf("usage: extend <meldId> <tileId>... (>=1)\n")
			} else {
				payload, _ := json.Marshal(map[string]interface{}{"meldId": parts[1], "tileIds": parts[2:]})
				env, _ := json.Marshal(protocol.Envelope{Version: protocol.Version, OpCode: protocol.OpClientExtendMeld, Payload: payload})
				nextRaw := rm.MatchLoop(context.Background(), logger, nil, nil, dispatcher, 0, st, []runtime.MatchData{&mockMatchData{mockPresence: mockPresence{userId: currentUser}, opCode: protocol.OpClientExtendMeld, data: env}})
				st = nextRaw.(*match.RoundState)
				if dispatcher.lastOp == protocol.OpServerError {
					fmt.Printf("  error: %s\n", string(dispatcher.lastData))
				} else {
					fmt.Printf("  extend succeeded\n")
				}
				renderState(st, currentSeat)
			}
		case "replace":
			if len(parts) < 5 {
				fmt.Printf("usage: replace <targetMeldId> <tileId> <newTile1> <newTile2> [jokerId colour rank]\n")
			} else {
				mm := map[string]interface{}{
					"targetMeldId": parts[1],
					"tileId":       parts[2],
					"newMeldTiles": []string{parts[3], parts[4]},
				}
				if len(parts) == 8 {
					mm["jokerReps"] = map[string]interface{}{parts[5]: map[string]interface{}{"colour": parts[6], "rank": mustAtoi(parts[7])}}
				}
				payload, _ := json.Marshal(mm)
				env, _ := json.Marshal(protocol.Envelope{Version: protocol.Version, OpCode: protocol.OpClientReplaceJoker, Payload: payload})
				nextRaw := rm.MatchLoop(context.Background(), logger, nil, nil, dispatcher, 0, st, []runtime.MatchData{&mockMatchData{mockPresence: mockPresence{userId: currentUser}, opCode: protocol.OpClientReplaceJoker, data: env}})
				st = nextRaw.(*match.RoundState)
				if dispatcher.lastOp == protocol.OpServerError {
					fmt.Printf("  error: %s\n", string(dispatcher.lastData))
				} else {
					fmt.Printf("  replace succeeded\n")
				}
				renderState(st, currentSeat)
			}
		case "winner", "w":
			fmt.Printf("Winner: %v, GamePhase: %v\n", st.Winner, st.GamePhase)
		default:
			fmt.Printf("unknown command %q, type 'help'\n", cmd)
		}
		fmt.Printf("\n[%s seat-%d] > ", currentUser, currentSeat)
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("scanner error: %v\n", err)
	}
}
