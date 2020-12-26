package parser

import (
	"errors"
	"regexp"
	"strconv"
)

const _worldID int = 1022

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
		userIndex := FindUserByID((*slc)[gameID].Players, userID)
		if userIndex == -1 {
			return errors.New("Trying to update a user that doesn't exists")
		}
		(*slc)[gameID].Players[userIndex].Name = pInfos[2]
	case "Kill":
		if len((*slc)) == 0 {
			return errors.New("Kill attempt but no match is active")
		}
		if len((*slc)[gameID].Players) == 0 {
			return errors.New("Kill attempt but no one is on the match")
		}
		re, err := regexp.Compile(`^(\d+) (\d+) (\d+): .*$`)
		if err != nil {
			return errors.New("Error on Parse Line")
		}
		pInfos := re.FindStringSubmatch(matchs[2])
		killerID, _ := strconv.Atoi(pInfos[1])
		victimID, _ := strconv.Atoi(pInfos[2])
		meanOfDeath, _ := strconv.Atoi(pInfos[3])
		killerIndex := FindUserByID((*slc)[gameID].Players, killerID)
		if killerIndex == -1 && killerID != _worldID {
			return errors.New("Kill by a non existent player")
		}
		victimIndex := FindUserByID((*slc)[gameID].Players, victimID)
		if victimIndex == -1 {
			return errors.New("Kill attempt to a non existent player")
		}
		(*slc)[gameID].Events = append((*slc)[gameID].Events, Kill{
			KillerID:    killerID,
			VictimID:    victimID,
			MeanOfDeath: meanOfDeath,
		})
	default:
		return nil
	}
	return nil
}

// FindUserByID test
func FindUserByID(players []Player, id int) int {
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
	MeanOfDeath int
}

// Match doc
type Match struct {
	Players []Player
	Events  []Kill
}
