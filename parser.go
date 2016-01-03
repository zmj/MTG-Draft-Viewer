package main

import (
	"bufio"
	"errors"
	"io"
	"regexp"
	"strings"
)

type Draft struct {
	Event    string
	Player   string
	Players  []string
	Date     string
	Sets     []string
	Packs    []*Pack
	Deck     []string
	Comments []string
	Image    string
	HasDeck  bool
}

type Pack struct {
	Set   string
	Num   int
	Picks []*Pick
}

type Pick struct {
	Num    int
	Cards  []string
	Choice string
	Deck   []string
}

const (
	LineSeparator = "\n"
	EventId       = iota
	EventTime
	StartPlayerList
	PlayerList
	BlankLine
	NewPack
	NewPick
	CardList
)

var (
	emptyLineReg, _ = regexp.Compile(`^\s*$`)
	eventIdReg, _   = regexp.Compile(`^Event #:\s+(\d+)`)
	eventTimeReg, _ = regexp.Compile(`^Time:\s+([0-9/]+)`)

	currentItemReg, _ = regexp.Compile(`^-->\s+(\S+.*)$`)
	otherItemReg, _   = regexp.Compile(`^\s+(\S+.*)$`)

	playersReg, _ = regexp.Compile(`Players:`)
	packReg, _    = regexp.Compile(`-+ (\w+) -+`)
	pickReg, _    = regexp.Compile(`Pack (\d+) pick (\d+):`)
)

func NewDraft(log io.Reader) (*Draft, error) {
	state := EventId
	draft := new(Draft)
	var pack *Pack
	var pick *Pick
	var m []string
	buf := bufio.NewReader(log)
	for {
		lineBytes, _, readErr := buf.ReadLine()
		if readErr == io.EOF {
			break
		} else if readErr != nil {
			return nil, readErr
		}
		line := string(lineBytes)
		if strings.HasPrefix(line, "--:") {
			_, _, readErr := buf.ReadLine() // skip empty
			if readErr == io.EOF {
				break
			} else if readErr != nil {
				return nil, readErr
			}
			continue
		}
		switch state {
		case EventId:
			m = eventIdReg.FindStringSubmatch(line)
			if m == nil {
				return nil, errors.New("No eventID match")
			}
			draft.Event = m[1]
			state = EventTime
		case EventTime:
			m = eventTimeReg.FindStringSubmatch(line)
			if m == nil {
				return nil, errors.New("No eventTime match")
			}
			draft.Date = m[1]
			state = StartPlayerList
		case StartPlayerList:
			state = PlayerList
		case PlayerList:
			if emptyLineReg.MatchString(line) {
				state = NewPack
				continue
			}
			m = currentItemReg.FindStringSubmatch(line)
			if m == nil {
				m = otherItemReg.FindStringSubmatch(line)
				if m == nil {
					return nil, errors.New("No Player match")
				}
			} else {
				draft.Player = m[1]
			}
			draft.Players = append(draft.Players, m[1])
		case NewPack:
			if emptyLineReg.MatchString(line) {
				state = NewPick
				continue
			}
			m = packReg.FindStringSubmatch(line)
			if m == nil {
				return nil, errors.New("No Pack match")
			}
			pack = new(Pack)
			draft.Packs = append(draft.Packs, pack)
			pack.Num = len(draft.Packs)
			pack.Set = m[1]
			draft.Sets = append(draft.Sets, m[1])
		case NewPick:
			pick = new(Pick)
			pack.Picks = append(pack.Picks, pick)
			pick.Num = len(pack.Picks)
			state = CardList
		case CardList:
			if emptyLineReg.MatchString(line) {
				if len(pick.Cards) == 1 {
					state = NewPack
				} else {
					state = NewPick
				}
				continue
			}
			m = currentItemReg.FindStringSubmatch(line)
			if m == nil {
				m = otherItemReg.FindStringSubmatch(line)
				if m == nil {
					return nil, errors.New("No Card match")
				}
			} else {
				pick.Choice = m[1]
				pick.Deck = draft.Deck[:]
				draft.Deck = append(draft.Deck, m[1])
			}
			pick.Cards = append(pick.Cards, m[1])
		}
	}
	return draft, nil
}
