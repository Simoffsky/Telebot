package opendota

import (
	"strconv"
	"time"
	"vladOS/pkg/telebot"
)

var (
	danilaLastMatchId int64
)

func init() {
	danilaLastMatchId = GetLastMatch().MatchId
}

// send message to the group if new match was played
func CheckPlayedMatch(bot *telebot.Telebot, telegramGroupId int64) {
	for {
		lastMatch := GetLastMatch()
		if lastMatch.MatchId != danilaLastMatchId {
			danilaLastMatchId = lastMatch.MatchId
			heroName := GetHeroName(lastMatch.HeroId)
			message := getMessage(lastMatch, heroName)
			bot.SendMessage(telegramGroupId, message)
		}
		time.Sleep(1 * time.Minute)
	}

}

func getMessage(match Match, heroName string) string {
	return `
Поздравляем @bozhedoms! Он сыграл великолепную игру на ` + heroName + `!
Он совершил ` + strconv.Itoa(match.Kills) + ` убийств, ` + `умер ` + strconv.Itoa(match.Deaths) + ` раз и нанес ` + strconv.Itoa(match.HeroDmg) + ` урона героям!
Игра длилась ` + strconv.Itoa(match.Duration/60) + ` минут и ` + strconv.Itoa(match.Duration%60) + ` секунд!
	`
}
