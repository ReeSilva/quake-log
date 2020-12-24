package parser

import (
	"errors"
	"regexp"
	"strconv"
)

// ParseLine testing
func ParseLine(gameID int, slc *[]Match, line string) error {
	if line == "" {
		return errors.New("Error on Parse Line")
	}
	r, err := regexp.Compile(`\d+:\d+ (\w+): (.*)`)
	if err != nil {
		return errors.New("Error on Parse Line")
	}
	matchs := r.FindStringSubmatch(line)
	switch matchs[1] {
	case "InitGame":
		*slc = append((*slc), Match{
			Players: []Player{},
			Events:  []Kill{},
		})
	case "ClientConnect":
		if len((*slc)) == 0 {
			return errors.New("ClientConnect line without an initialized match")
		}
		playerID, _ := strconv.Atoi(matchs[2])
		(*slc)[gameID].Players = append((*slc)[gameID].Players, Player{
			ID:   playerID,
			Name: "",
		})
	case "ClientUserinfoChanged":
		if len((*slc)) == 0 {
			return errors.New("Updating player with no matches running")
		}
		if len((*slc)[gameID].Players) == 0 {
			return errors.New("Updating player with no players on match")
		}
		re, err := regexp.Compile(`^(\d+) n\\(\w+).*$`)
		if err != nil {
			return errors.New("Error on Parse Line")
		}
		pInfos := re.FindStringSubmatch(matchs[2])
		userID, _ := strconv.Atoi(pInfos[1])
		userIndex := _findUserByID((*slc)[gameID].Players, userID)
		if userIndex == -1 {
			return errors.New("Trying to update a user that doesn't exists")
		}
		(*slc)[gameID].Players[userIndex].Name = pInfos[2]
	default:
		return nil
	}
	return nil
}

func _findUserByID(players []Player, id int) int {
	index := -1
	for key, value := range players {
		if value.ID == id {
			index = key
			break
		}
	}
	return index
}

// Player doc
type Player struct {
	ID   int
	Name string
}

// Kill doc
type Kill struct {
	KillerID    int
	VictimID    int
	MeanOfDeath string
}

// Match doc
type Match struct {
	Players []Player
	Events  []Kill
}
