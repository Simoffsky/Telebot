package opendota

import (
	"encoding/json"
	"log"
	"net/http"
)

const (
	danilaId = "216527063"
)

type Match struct {
	MatchId  int64 `json:"match_id"`
	Duration int   `json:"duration"`
	HeroId   int   `json:"hero_id"`
	Kills    int   `json:"kills"`
	Deaths   int   `json:"deaths"`
	HeroDmg  int   `json:"hero_damage"`
}

type Hero struct {
	Id   int    `json:"id"`
	Name string `json:"localized_name"`
}

func GetLastMatch() Match {
	resp, err := http.Get("https://api.opendota.com/api/players/" + danilaId + "/recentMatches")
	if err != nil {
		log.Fatal(err)
	}

	matchArray := make([]Match, 1)

	err = json.NewDecoder(resp.Body).Decode(&matchArray)

	if err != nil {
		log.Fatal(err)
	}
	return matchArray[0]
}

func GetHeroName(id int) string {
	resp, err := http.Get("https://api.opendota.com/api/heroes/")
	if err != nil {
		log.Fatal(err)
	}

	heroes := make([]Hero, 1)

	err = json.NewDecoder(resp.Body).Decode(&heroes)

	if err != nil {
		log.Fatal(err)
	}

	for _, hero := range heroes {
		if hero.Id == id {
			return hero.Name
		}
	}
	return "Hero not found"
}
