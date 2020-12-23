package parser_test

import (
	"errors"
	"testing"

	"github.com/reesilva/quake-log/pkg/parser"
	"github.com/stretchr/testify/assert"
)

type Parameters struct {
	Matchs []parser.Match
	Line   string
}

func TestParseLine(t *testing.T) {
	tests := []struct {
		name          string
		expectedError string
		want          []parser.Match
		expectError   bool
		parameters    Parameters
	}{
		// Line can not be empty
		{
			name:          "Line can not be empty",
			expectedError: "Error on Parse Line",
			expectError:   true,
			parameters: Parameters{
				Matchs: []parser.Match{
					{
						Players: []parser.Player{},
						Events:  []parser.Kill{},
					},
				},
				Line: "",
			},
			want: []parser.Match{},
		},
		// Line Init Game
		{
			name: "Line Init Game",
			want: []parser.Match{
				{
					Players: []parser.Player{},
					Events:  []parser.Kill{},
				},
			},
			expectError: false,
			parameters: Parameters{
				Matchs: []parser.Match{},
				Line:   "  0:00 InitGame: \\sv_floodProtect\\1\\sv_maxPing\\0\\sv_minPing\\0\\sv_maxRate\\10000\\sv_minRate\\0\\sv_hostname\\Code Miner Server\\g_gametype\\0\\sv_privateClients\\2\\sv_maxclients\\16\\sv_allowDownload\\0\\dmflags\\0\\fraglimit\\20\\timelimit\\15\\g_maxGameClients\\0\\capturelimit\\8\\version\\ioq3 1.36 linux-x86_64 Apr 12 2009\\protocol\\68\\mapname\\q3dm17\\gamename\\baseq3\\g_needpass\\0",
			},
		},
		// Trying to append player to an unitialized match
		{
			name:          "Trying to append player to an unitialized match",
			want:          []parser.Match{},
			expectError:   true,
			expectedError: "ClientConnect line without an initialized match",
			parameters: Parameters{
				Matchs: []parser.Match{},
				Line:   " 20:34 ClientConnect: 2",
			},
		},
		// Line Client connect
		{
			name: "Line Client connect",
			want: []parser.Match{
				{
					Players: []parser.Player{
						{
							ID:   2,
							Name: "",
						},
					},
					Events: []parser.Kill{},
				},
			},
			expectError: false,
			parameters: Parameters{
				Matchs: []parser.Match{
					{
						Players: []parser.Player{},
						Events:  []parser.Kill{},
					},
				},
				Line: " 20:34 ClientConnect: 2",
			},
		},
		// Init a new game when a previous one already exists
		{
			name: "Init a new game when a previous one already exists",
			want: []parser.Match{
				{
					Players: []parser.Player{
						{
							ID:   2,
							Name: "isgalamido",
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
							MeanOfDeath: "test",
						},
						{
							KillerID:    3,
							VictimID:    2,
							MeanOfDeath: "test",
						},
					},
				},
				{
					Players: []parser.Player{},
					Events:  []parser.Kill{},
				},
			},
			expectError: false,
			parameters: Parameters{
				Matchs: []parser.Match{
					{
						Players: []parser.Player{
							{
								ID:   2,
								Name: "isgalamido",
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
								MeanOfDeath: "test",
							},
							{
								KillerID:    3,
								VictimID:    2,
								MeanOfDeath: "test",
							},
						},
					},
				},
				Line: "  0:00 InitGame: \\sv_floodProtect\\1\\sv_maxPing\\0\\sv_minPing\\0\\sv_maxRate\\10000\\sv_minRate\\0\\sv_hostname\\Code Miner Server\\g_gametype\\0\\sv_privateClients\\2\\sv_maxclients\\16\\sv_allowDownload\\0\\dmflags\\0\\fraglimit\\20\\timelimit\\15\\g_maxGameClients\\0\\capturelimit\\8\\version\\ioq3 1.36 linux-x86_64 Apr 12 2009\\protocol\\68\\mapname\\q3dm17\\gamename\\baseq3\\g_needpass\\0",
			},
		},
		// Client connect to a new game when a previous one already exists
		{
			name: "Client connect to a new game when a previous one already exists",
			want: []parser.Match{
				{
					Players: []parser.Player{
						{
							ID:   2,
							Name: "isgalamido",
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
							MeanOfDeath: "test",
						},
						{
							KillerID:    3,
							VictimID:    2,
							MeanOfDeath: "test",
						},
					},
				},
				{
					Players: []parser.Player{
						{
							ID:   2,
							Name: "",
						},
					},
					Events: []parser.Kill{},
				},
			},
			expectError: false,
			parameters: Parameters{
				Matchs: []parser.Match{
					{
						Players: []parser.Player{
							{
								ID:   2,
								Name: "isgalamido",
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
								MeanOfDeath: "test",
							},
							{
								KillerID:    3,
								VictimID:    2,
								MeanOfDeath: "test",
							},
						},
					},
					{
						Players: []parser.Player{},
						Events:  []parser.Kill{},
					},
				},
				Line: " 20:34 ClientConnect: 2",
			},
		},
		// Client user info changed when has no matches
		{
			name:          "Client user info changed when has no matches",
			want:          []parser.Match{},
			expectError:   true,
			expectedError: "Updating player with no matches running",
			parameters: Parameters{
				Matchs: []parser.Match{},
				Line:   ` 20:38 ClientUserinfoChanged: 2 n\Isgalamido\t\0\model\uriel/zael\hmodel\uriel/zael\g_redteam\\g_blueteam\\c1\5\c2\5\hc\100\w\0\l\0\tt\0\tl\0`,
			},
		},
		// Client user info changed when has no player on match
		{
			name:          "Client user info changed when has no player on match",
			want:          []parser.Match{},
			expectError:   true,
			expectedError: "Updating player with no players on match",
			parameters: Parameters{
				Matchs: []parser.Match{
					{
						Players: []parser.Player{},
						Events:  []parser.Kill{},
					},
				},
				Line: ` 20:38 ClientUserinfoChanged: 2 n\Isgalamido\t\0\model\uriel/zael\hmodel\uriel/zael\g_redteam\\g_blueteam\\c1\5\c2\5\hc\100\w\0\l\0\tt\0\tl\0`,
			},
		},
		// Client user info changed
		{
			name: "Client user info changed",
			want: []parser.Match{
				{
					Players: []parser.Player{
						{
							ID:   2,
							Name: "Isgalamido",
						},
					},
					Events: []parser.Kill{},
				},
			},
			expectError: false,
			parameters: Parameters{
				Matchs: []parser.Match{
					{
						Players: []parser.Player{
							{
								ID:   2,
								Name: "",
							},
						},
						Events: []parser.Kill{},
					},
				},
				Line: ` 20:38 ClientUserinfoChanged: 2 n\Isgalamido\t\0\model\uriel/zael\hmodel\uriel/zael\g_redteam\\g_blueteam\\c1\5\c2\5\hc\100\w\0\l\0\tt\0\tl\0`,
			},
		},
		// Client user info changed for a user that doesn't exists
		{
			name:          "Client user info changed for a user that doesn't exists",
			want:          []parser.Match{},
			expectError:   true,
			expectedError: "Trying to update a user that doesn't exists",
			parameters: Parameters{
				Matchs: []parser.Match{
					{
						Players: []parser.Player{
							{
								ID:   2,
								Name: "",
							},
						},
						Events: []parser.Kill{},
					},
				},
				Line: ` 20:38 ClientUserinfoChanged: 3 n\Isgalamido\t\0\model\uriel/zael\hmodel\uriel/zael\g_redteam\\g_blueteam\\c1\5\c2\5\hc\100\w\0\l\0\tt\0\tl\0`,
			},
		},
		// Client user info changed in a new game when a previous one already exists
		{
			name: "Client user info changed in a new game when a previous one already exists",
			want: []parser.Match{
				{
					Players: []parser.Player{
						{
							ID:   2,
							Name: "test2",
						},
						{
							ID:   3,
							Name: "test3",
						},
					},
					Events: []parser.Kill{
						{
							KillerID:    2,
							VictimID:    3,
							MeanOfDeath: "test",
						},
						{
							KillerID:    3,
							VictimID:    2,
							MeanOfDeath: "test",
						},
					},
				},
				{
					Players: []parser.Player{
						{
							ID:   2,
							Name: "Isgalamido",
						},
					},
					Events: []parser.Kill{},
				},
			},
			expectError: false,
			parameters: Parameters{
				Matchs: []parser.Match{
					{
						Players: []parser.Player{
							{
								ID:   2,
								Name: "test2",
							},
							{
								ID:   3,
								Name: "test3",
							},
						},
						Events: []parser.Kill{
							{
								KillerID:    2,
								VictimID:    3,
								MeanOfDeath: "test",
							},
							{
								KillerID:    3,
								VictimID:    2,
								MeanOfDeath: "test",
							},
						},
					},
					{
						Players: []parser.Player{
							{
								ID:   2,
								Name: "",
							},
						},
						Events: []parser.Kill{},
					},
				},
				Line: ` 20:38 ClientUserinfoChanged: 2 n\Isgalamido\t\0\model\uriel/zael\hmodel\uriel/zael\g_redteam\\g_blueteam\\c1\5\c2\5\hc\100\w\0\l\0\tt\0\tl\0`,
			},
		},
		// Client user info changed in a game with more then one player
		{
			name: "Client user info changed in a game with more then one player",
			want: []parser.Match{
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
			expectError: false,
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
								Name: "",
							},
						},
						Events: []parser.Kill{},
					},
				},
				Line: ` 20:38 ClientUserinfoChanged: 3 n\Mocinha\t\0\model\uriel/zael\hmodel\uriel/zael\g_redteam\\g_blueteam\\c1\5\c2\5\hc\100\w\0\l\0\tt\0\tl\0`,
			},
		},
		// Client disconnected when has no matches
		{
			name:          "Client disconnected when has no matches",
			want:          []parser.Match{},
			expectError:   true,
			expectedError: "Disconnecting player with no matches running",
			parameters: Parameters{
				Matchs: []parser.Match{},
				Line:   ` 21:10 ClientDisconnect: 2`,
			},
		},
		// Client disconnected when has no player on match
		{
			name:          "Client disconnected when has no player on match",
			want:          []parser.Match{},
			expectError:   true,
			expectedError: "Disconnecting player with no players on match",
			parameters: Parameters{
				Matchs: []parser.Match{
					{
						Players: []parser.Player{},
						Events:  []parser.Kill{},
					},
				},
				Line: ` 21:10 ClientDisconnect: 2`,
			},
		},
		// Client disconnected
		{
			name: "Client disconnected",
			want: []parser.Match{
				{
					Players: []parser.Player{},
					Events:  []parser.Kill{},
				},
			},
			expectError: false,
			parameters: Parameters{
				Matchs: []parser.Match{
					{
						Players: []parser.Player{
							{
								ID:   2,
								Name: "",
							},
						},
						Events: []parser.Kill{},
					},
				},
				Line: ` 21:10 ClientDisconnect: 2`,
			},
		},
		// Client disconnected for a user that doesn't exists
		{
			name:          "Client disconnected for a user that doesn't exists",
			want:          []parser.Match{},
			expectError:   true,
			expectedError: "Disconnecting a user that doesn't exists",
			parameters: Parameters{
				Matchs: []parser.Match{
					{
						Players: []parser.Player{
							{
								ID:   3,
								Name: "Mocinha",
							},
						},
						Events: []parser.Kill{},
					},
				},
				Line: ` 21:10 ClientDisconnect: 2`,
			},
		},
		// Client disconnected in a new game when a previous one already exists
		{
			name: "Client disconnected in a new game when a previous one already exists",
			want: []parser.Match{
				{
					Players: []parser.Player{
						{
							ID:   2,
							Name: "test2",
						},
						{
							ID:   3,
							Name: "test3",
						},
					},
					Events: []parser.Kill{
						{
							KillerID:    2,
							VictimID:    3,
							MeanOfDeath: "test",
						},
						{
							KillerID:    3,
							VictimID:    2,
							MeanOfDeath: "test",
						},
					},
				},
				{
					Players: []parser.Player{},
					Events:  []parser.Kill{},
				},
			},
			expectError: false,
			parameters: Parameters{
				Matchs: []parser.Match{
					{
						Players: []parser.Player{
							{
								ID:   2,
								Name: "test2",
							},
							{
								ID:   3,
								Name: "test3",
							},
						},
						Events: []parser.Kill{
							{
								KillerID:    2,
								VictimID:    3,
								MeanOfDeath: "test",
							},
							{
								KillerID:    3,
								VictimID:    2,
								MeanOfDeath: "test",
							},
						},
					},
					{
						Players: []parser.Player{
							{
								ID:   2,
								Name: "Isgalamido",
							},
						},
						Events: []parser.Kill{},
					},
				},
				Line: ` 21:10 ClientDisconnect: 2`,
			},
		},
		// Client disconnected in a game with more then one player
		{
			name: "Client disconnected in a game with more then one player",
			want: []parser.Match{
				{
					Players: []parser.Player{
						{
							ID:   3,
							Name: "Mocinha",
						},
					},
					Events: []parser.Kill{},
				},
			},
			expectError: false,
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
				Line: ` 21:10 ClientDisconnect: 2`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := parser.ParseLine(len(tt.parameters.Matchs)-1, &tt.parameters.Matchs, tt.parameters.Line)
			if tt.expectError {
				if assert.Error(t, err) {
					expected := errors.New(tt.expectedError)
					assert.Equal(t, expected, err)
				}
			} else {
				assert.Equal(t, tt.want, tt.parameters.Matchs)
			}
		})
	}
}
