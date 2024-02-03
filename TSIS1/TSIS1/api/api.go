package api

type Team struct {
    ID             string `json:"id"`
    TeamName       string `json:"team"`
    CountOfMembers int    `json:"count"`
    Leader         string `json:"leader"`
}

var Teams = []Team{
    {ID: "1", TeamName: "Hamming Bird", CountOfMembers: 6, Leader: "Jay Jo"},
    {ID: "2", TeamName: "Zephyrus", CountOfMembers: 3, Leader: "TJ"},
    {ID: "3", TeamName: "Sabbath", CountOfMembers: 3, Leader: "Wooin"},
    {ID: "4", TeamName: "Monster Bull", CountOfMembers: 5, Leader: "Monster"},
    {ID: "5", TeamName: "Tridents", CountOfMembers: 4, Leader: "Juhwan"},
    {ID: "6", TeamName: "Ghost", CountOfMembers: 7, Leader: "Hwangyeon Choi"},
    {ID: "7", TeamName: "Light Cavalry", CountOfMembers: 3, Leader: "Owen Knight   "},
}

func GetTeam(id string) *Team {
	for i := range Teams {
		if Teams[i].ID == id {
			return &Teams[i]
		}
	}
	return nil
}