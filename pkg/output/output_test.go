package output_test

import (
	"errors"
	"testing"

	"github.com/reesilva/quake-log/pkg/output"
	"github.com/reesilva/quake-log/pkg/parser"
	"github.com/stretchr/testify/assert"
)

type Parameters struct {
	Matchs       []parser.Match
	DeathByMeans bool
}

func TestCreateMatchReport(t *testing.T) {
	tests := []struct {
		name          string
		parameters    Parameters
		expectError   bool
		expectedError string
		want          map[string]output.MatchReport
	}{
		// No Matches
		{
			name: "No Matches",
			parameters: Parameters{
				Matchs:       []parser.Match{},
				DeathByMeans: false,
			},
			expectError:   false,
			expectedError: "",
			want:          map[string]output.MatchReport{},
		},
		// Match with no events
		{
			name: "Matches with no events",
			parameters: Parameters{
				Matchs: []parser.Match{
					{
						Players: []parser.Player{
							{
								ID:   2,
								Name: "Isgalamido",
							},
							{
								ID:   3,
								Name: "Mocinha",
							},
						},
						Events: []parser.Kill{},
					},
				},
				DeathByMeans: false,
			},
			expectError:   false,
			expectedError: "",
			want: map[string]output.MatchReport{
				"game_1": {
					TotalKills: 0,
					Players:    []string{"Isgalamido", "Mocinha"},
					Kills:      map[string]int{},
				},
			},
		},
		// More then one match with no events
		{
			name: "More then one match with no events",
			parameters: Parameters{
				Matchs: []parser.Match{
					{
						Players: []parser.Player{
							{
								ID:   2,
								Name: "Isgalamido",
							},
							{
								ID:   3,
								Name: "Mocinha",
							},
						},
						Events: []parser.Kill{},
					},
					{
						Players: []parser.Player{
							{
								ID:   49,
								Name: "Faker1",
							},
							{
								ID:   56,
								Name: "Faker43",
							},
						},
						Events: []parser.Kill{},
					},
				},
				DeathByMeans: false,
			},
			expectError:   false,
			expectedError: "",
			want: map[string]output.MatchReport{
				"game_1": {
					TotalKills: 0,
					Players:    []string{"Isgalamido", "Mocinha"},
					Kills:      map[string]int{},
				},
				"game_2": {
					TotalKills: 0,
					Players:    []string{"Faker1", "Faker43"},
					Kills:      map[string]int{},
				},
			},
		},
		// Match with kills only by players
		{
			name: "Match with kills only by players",
			parameters: Parameters{
				Matchs: []parser.Match{
					{
						Players: []parser.Player{
							{
								ID:   2,
								Name: "Isgalamido",
							},
							{
								ID:   3,
								Name: "Mocinha",
							},
						},
						Events: []parser.Kill{
							{
								KillerID:    2,
								VictimID:    3,
								MeanOfDeath: 22,
							},
							{
								KillerID:    3,
								VictimID:    2,
								MeanOfDeath: 25,
							},
							{
								KillerID:    3,
								VictimID:    2,
								MeanOfDeath: 25,
							},
							{
								KillerID:    2,
								VictimID:    3,
								MeanOfDeath: 22,
							},
						},
					},
				},
				DeathByMeans: false,
			},
			expectError:   false,
			expectedError: "",
			want: map[string]output.MatchReport{
				"game_1": {
					TotalKills: 4,
					Players:    []string{"Isgalamido", "Mocinha"},
					Kills: map[string]int{
						"Isgalamido": 2,
						"Mocinha":    2,
					},
				},
			},
		},
		// More then one match with kills only by players
		{
			name: "More then one match with kills only by players",
			parameters: Parameters{
				Matchs: []parser.Match{
					{
						Players: []parser.Player{
							{
								ID:   2,
								Name: "Isgalamido",
							},
							{
								ID:   3,
								Name: "Mocinha",
							},
						},
						Events: []parser.Kill{
							{
								KillerID:    2,
								VictimID:    3,
								MeanOfDeath: 22,
							},
							{
								KillerID:    3,
								VictimID:    2,
								MeanOfDeath: 25,
							},
							{
								KillerID:    3,
								VictimID:    2,
								MeanOfDeath: 25,
							},
							{
								KillerID:    2,
								VictimID:    3,
								MeanOfDeath: 22,
							},
						},
					},
					{
						Players: []parser.Player{
							{
								ID:   49,
								Name: "Faker49",
							},
							{
								ID:   57,
								Name: "Faker57",
							},
						},
						Events: []parser.Kill{
							{
								KillerID:    49,
								VictimID:    57,
								MeanOfDeath: 22,
							},
							{
								KillerID:    49,
								VictimID:    57,
								MeanOfDeath: 25,
							},
							{
								KillerID:    57,
								VictimID:    49,
								MeanOfDeath: 25,
							},
						},
					},
				},
				DeathByMeans: false,
			},
			expectError:   false,
			expectedError: "",
			want: map[string]output.MatchReport{
				"game_1": {
					TotalKills: 4,
					Players:    []string{"Isgalamido", "Mocinha"},
					Kills: map[string]int{
						"Isgalamido": 2,
						"Mocinha":    2,
					},
				},
				"game_2": {
					TotalKills: 3,
					Players:    []string{"Faker49", "Faker57"},
					Kills: map[string]int{
						"Faker49": 2,
						"Faker57": 1,
					},
				},
			},
		},
		// Match with kills by players and world
		{
			name: "Match with kills by players and world",
			parameters: Parameters{
				Matchs: []parser.Match{
					{
						Players: []parser.Player{
							{
								ID:   2,
								Name: "Isgalamido",
							},
							{
								ID:   3,
								Name: "Mocinha",
							},
						},
						Events: []parser.Kill{
							{
								KillerID:    2,
								VictimID:    3,
								MeanOfDeath: 22,
							},
							{
								KillerID:    3,
								VictimID:    2,
								MeanOfDeath: 25,
							},
							{
								KillerID:    3,
								VictimID:    2,
								MeanOfDeath: 25,
							},
							{
								KillerID:    1022,
								VictimID:    2,
								MeanOfDeath: 22,
							},
							{
								KillerID:    1022,
								VictimID:    3,
								MeanOfDeath: 22,
							},
						},
					},
				},
				DeathByMeans: false,
			},
			expectError:   false,
			expectedError: "",
			want: map[string]output.MatchReport{
				"game_1": {
					TotalKills: 5,
					Players:    []string{"Isgalamido", "Mocinha"},
					Kills: map[string]int{
						"Isgalamido": 1,
						"Mocinha":    2,
					},
				},
			},
		},
		// More then one match with kills by players and world
		{
			name: "More then one match with kills by players and world",
			parameters: Parameters{
				Matchs: []parser.Match{
					{
						Players: []parser.Player{
							{
								ID:   2,
								Name: "Isgalamido",
							},
							{
								ID:   3,
								Name: "Mocinha",
							},
						},
						Events: []parser.Kill{
							{
								KillerID:    2,
								VictimID:    3,
								MeanOfDeath: 22,
							},
							{
								KillerID:    3,
								VictimID:    2,
								MeanOfDeath: 25,
							},
							{
								KillerID:    3,
								VictimID:    2,
								MeanOfDeath: 25,
							},
							{
								KillerID:    1022,
								VictimID:    3,
								MeanOfDeath: 22,
							},
							{
								KillerID:    1022,
								VictimID:    3,
								MeanOfDeath: 22,
							},
						},
					},
					{
						Players: []parser.Player{
							{
								ID:   49,
								Name: "Faker49",
							},
							{
								ID:   57,
								Name: "Faker57",
							},
						},
						Events: []parser.Kill{
							{
								KillerID:    49,
								VictimID:    57,
								MeanOfDeath: 22,
							},
							{
								KillerID:    49,
								VictimID:    57,
								MeanOfDeath: 25,
							},
							{
								KillerID:    57,
								VictimID:    49,
								MeanOfDeath: 25,
							},
							{
								KillerID:    1022,
								VictimID:    49,
								MeanOfDeath: 22,
							},
							{
								KillerID:    1022,
								VictimID:    57,
								MeanOfDeath: 22,
							},
						},
					},
				},
				DeathByMeans: false,
			},
			expectError:   false,
			expectedError: "",
			want: map[string]output.MatchReport{
				"game_1": {
					TotalKills: 5,
					Players:    []string{"Isgalamido", "Mocinha"},
					Kills: map[string]int{
						"Isgalamido": 1,
						"Mocinha":    2,
					},
				},
				"game_2": {
					TotalKills: 5,
					Players:    []string{"Faker49", "Faker57"},
					Kills: map[string]int{
						"Faker49": 2,
						"Faker57": 1,
					},
				},
			},
		},
		// Match with kills only by players and means of death active
		{
			name: "Match with kills only by players and means of death active",
			parameters: Parameters{
				Matchs: []parser.Match{
					{
						Players: []parser.Player{
							{
								ID:   2,
								Name: "Isgalamido",
							},
							{
								ID:   3,
								Name: "Mocinha",
							},
						},
						Events: []parser.Kill{
							{
								KillerID:    2,
								VictimID:    3,
								MeanOfDeath: 22,
							},
							{
								KillerID:    3,
								VictimID:    2,
								MeanOfDeath: 25,
							},
							{
								KillerID:    3,
								VictimID:    2,
								MeanOfDeath: 25,
							},
							{
								KillerID:    2,
								VictimID:    3,
								MeanOfDeath: 22,
							},
						},
					},
				},
				DeathByMeans: true,
			},
			expectError:   false,
			expectedError: "",
			want: map[string]output.MatchReport{
				"game_1": {
					TotalKills: 4,
					Players:    []string{"Isgalamido", "Mocinha"},
					Kills: map[string]int{
						"Isgalamido": 2,
						"Mocinha":    2,
					},
					KillsByMeans: map[string]int{
						"MOD_TRIGGER_HURT":   2,
						"MOD_PROXIMITY_MINE": 2,
					},
				},
			},
		},
		// More then one match with kills only by players and death by means active
		{
			name: "More then one match with kills only by players",
			parameters: Parameters{
				Matchs: []parser.Match{
					{
						Players: []parser.Player{
							{
								ID:   2,
								Name: "Isgalamido",
							},
							{
								ID:   3,
								Name: "Mocinha",
							},
						},
						Events: []parser.Kill{
							{
								KillerID:    2,
								VictimID:    3,
								MeanOfDeath: 22,
							},
							{
								KillerID:    3,
								VictimID:    2,
								MeanOfDeath: 25,
							},
							{
								KillerID:    3,
								VictimID:    2,
								MeanOfDeath: 25,
							},
							{
								KillerID:    2,
								VictimID:    3,
								MeanOfDeath: 22,
							},
						},
					},
					{
						Players: []parser.Player{
							{
								ID:   49,
								Name: "Faker49",
							},
							{
								ID:   57,
								Name: "Faker57",
							},
						},
						Events: []parser.Kill{
							{
								KillerID:    49,
								VictimID:    57,
								MeanOfDeath: 22,
							},
							{
								KillerID:    49,
								VictimID:    57,
								MeanOfDeath: 25,
							},
							{
								KillerID:    57,
								VictimID:    49,
								MeanOfDeath: 25,
							},
						},
					},
				},
				DeathByMeans: true,
			},
			expectError:   false,
			expectedError: "",
			want: map[string]output.MatchReport{
				"game_1": {
					TotalKills: 4,
					Players:    []string{"Isgalamido", "Mocinha"},
					Kills: map[string]int{
						"Isgalamido": 2,
						"Mocinha":    2,
					},
					KillsByMeans: map[string]int{
						"MOD_TRIGGER_HURT":   2,
						"MOD_PROXIMITY_MINE": 2,
					},
				},
				"game_2": {
					TotalKills: 3,
					Players:    []string{"Faker49", "Faker57"},
					Kills: map[string]int{
						"Faker49": 2,
						"Faker57": 1,
					},
					KillsByMeans: map[string]int{
						"MOD_TRIGGER_HURT":   1,
						"MOD_PROXIMITY_MINE": 2,
					},
				},
			},
		},
		// Match with kills by players and world
		{
			name: "Match with kills by players and world",
			parameters: Parameters{
				Matchs: []parser.Match{
					{
						Players: []parser.Player{
							{
								ID:   2,
								Name: "Isgalamido",
							},
							{
								ID:   3,
								Name: "Mocinha",
							},
						},
						Events: []parser.Kill{
							{
								KillerID:    2,
								VictimID:    3,
								MeanOfDeath: 22,
							},
							{
								KillerID:    3,
								VictimID:    2,
								MeanOfDeath: 25,
							},
							{
								KillerID:    3,
								VictimID:    2,
								MeanOfDeath: 25,
							},
							{
								KillerID:    1022,
								VictimID:    2,
								MeanOfDeath: 22,
							},
							{
								KillerID:    1022,
								VictimID:    3,
								MeanOfDeath: 22,
							},
						},
					},
				},
				DeathByMeans: true,
			},
			expectError:   false,
			expectedError: "",
			want: map[string]output.MatchReport{
				"game_1": {
					TotalKills: 5,
					Players:    []string{"Isgalamido", "Mocinha"},
					Kills: map[string]int{
						"Isgalamido": 1,
						"Mocinha":    2,
					},
					KillsByMeans: map[string]int{
						"MOD_TRIGGER_HURT":   3,
						"MOD_PROXIMITY_MINE": 2,
					},
				},
			},
		},
		// More then one match with kills by players and world
		{
			name: "More then one match with kills by players and world",
			parameters: Parameters{
				Matchs: []parser.Match{
					{
						Players: []parser.Player{
							{
								ID:   2,
								Name: "Isgalamido",
							},
							{
								ID:   3,
								Name: "Mocinha",
							},
						},
						Events: []parser.Kill{
							{
								KillerID:    2,
								VictimID:    3,
								MeanOfDeath: 22,
							},
							{
								KillerID:    3,
								VictimID:    2,
								MeanOfDeath: 25,
							},
							{
								KillerID:    3,
								VictimID:    2,
								MeanOfDeath: 25,
							},
							{
								KillerID:    1022,
								VictimID:    3,
								MeanOfDeath: 22,
							},
							{
								KillerID:    1022,
								VictimID:    3,
								MeanOfDeath: 22,
							},
						},
					},
					{
						Players: []parser.Player{
							{
								ID:   49,
								Name: "Faker49",
							},
							{
								ID:   57,
								Name: "Faker57",
							},
						},
						Events: []parser.Kill{
							{
								KillerID:    49,
								VictimID:    57,
								MeanOfDeath: 22,
							},
							{
								KillerID:    49,
								VictimID:    57,
								MeanOfDeath: 25,
							},
							{
								KillerID:    57,
								VictimID:    49,
								MeanOfDeath: 25,
							},
							{
								KillerID:    1022,
								VictimID:    49,
								MeanOfDeath: 22,
							},
							{
								KillerID:    1022,
								VictimID:    57,
								MeanOfDeath: 22,
							},
						},
					},
				},
				DeathByMeans: true,
			},
			expectError:   false,
			expectedError: "",
			want: map[string]output.MatchReport{
				"game_1": {
					TotalKills: 5,
					Players:    []string{"Isgalamido", "Mocinha"},
					Kills: map[string]int{
						"Isgalamido": 1,
						"Mocinha":    2,
					},
					KillsByMeans: map[string]int{
						"MOD_TRIGGER_HURT":   3,
						"MOD_PROXIMITY_MINE": 2,
					},
				},
				"game_2": {
					TotalKills: 5,
					Players:    []string{"Faker49", "Faker57"},
					Kills: map[string]int{
						"Faker49": 2,
						"Faker57": 1,
					},
					KillsByMeans: map[string]int{
						"MOD_TRIGGER_HURT":   3,
						"MOD_PROXIMITY_MINE": 2,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := output.CreateMatchReport(tt.parameters.Matchs, tt.parameters.DeathByMeans)
			if tt.expectError {
				if assert.Error(t, err) {
					expected := errors.New(tt.expectedError)
					assert.Equal(t, expected, err)
				}
			} else {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
