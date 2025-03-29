package handlers

import (
	"context"
	"encoding/json"
	"github.com/mattermost/mattermost/server/public/model"
	"golang.org/x/exp/slog"
	"regexp"
	"strings"
	"votty/internal/service"
	"votty/internal/storage/tarantool"
)

var (
	createPollRegex     = regexp.MustCompile(`^/create\s+([^|]+)\s*\|\s*([^|]+(?:\s*\|\s*[^|]+)*)$`)
	deleteCommandRegex  = regexp.MustCompile(`^/delete\s+([a-zA-Z0-9_-]+)$`)
	endCommandRegex     = regexp.MustCompile(`^/end\s+([a-zA-Z0-9_-]+)$`)
	resultsCommandRegex = regexp.MustCompile(`^/results\s+([a-zA-Z0-9_-]+)$`)
	voteCommandRegex    = regexp.MustCompile(`^/vote\s+([a-zA-Z0-9_-]+)\s+([1-9][0-9]*)$`)
)

func PostHandler(ctx context.Context, storage *tarantool.Storage, log *slog.Logger, client *model.Client4, event *model.WebSocketEvent, botID string) {
	postData, ok := event.GetData()["post"].(string)
	if !ok {
		log.Warn("Failed to process event: invalid data type for 'post'.")
		return
	}
	post := &model.Post{}

	err := json.Unmarshal([]byte(postData), &post)
	if err != nil {
		log.Warn("Failed to parse 'post' data: ", err.Error())
		return
	}
	if post.UserId == botID {
		return
	}

	var r *model.Post
	switch {
	case strings.HasPrefix(post.Message, "/create"):
		matches := createPollRegex.FindStringSubmatch(post.Message)

		r = service.CreatePoll(log, storage, post, matches)

	case strings.HasPrefix(post.Message, "/vote"):
		matches := voteCommandRegex.FindStringSubmatch(post.Message)

		r = service.Vote(storage, log, post, matches)

	case strings.HasPrefix(post.Message, "/end"):
		matches := endCommandRegex.FindStringSubmatch(post.Message)

		r = service.EndPoll(storage, log, post, matches)

	case strings.HasPrefix(post.Message, "/delete"):
		matches := deleteCommandRegex.FindStringSubmatch(post.Message)

		r = service.DeletePoll(storage, log, post, matches)

	case strings.HasPrefix(post.Message, "/results"):
		matches := resultsCommandRegex.FindStringSubmatch(post.Message)

		r = service.PullResults(log, storage, post, matches)

	case strings.HasPrefix(post.Message, "/guide"):
		r = &model.Post{
			Message: "Привет! Меня зовут Вотти, я помогу тебе проводить опросы быстро и эффективно." +
				"\nВот основные команды:" +
				"\nЧтобы создать новый опрос нужно ввести ```/create Ok? | var1 | var2 | var3```, где ```Ok?``` – любой вопрос по твоему усмотрению, ```var1, var2...``` – варианты ответов" +
				"\nПример: ```/create Это понятный пример? | Да | Нет``` этот запрос вернет тебе пронумерованные варианты ответов и ID опроса " +
				"\nВсе участники (в том числе и ты), которые получат доступ к ID опроса (PollID) могут проголосовать с помощью команды ```/vote PollID 1```, где ```PollID``` – полученный ID в /create (в след. примерах тоже)" +
				"\nЕще все могут посмотреть результаты опроса с помощью команды ```/results PollID```" +
				"\nЕсли ты собрал достаточно голосов, то можно завершить опрос командой ```/end PollID``` и тогда можно будет по прежнему смотреть результаты командой ```/results```, но ```vote``` перестанет быть доступным" +
				"\nЕсли результат опроса больше не интересен, то можно удалить опрос командой ```/delete PollID```",
		}

	}

	if r != nil {
		r.ChannelId = post.ChannelId
		_, _, err = client.CreatePost(ctx, r)

		if err != nil {
			slog.Error("Failed to send the message",
				slog.String("user_id", post.UserId),
				slog.String("message", post.Message),
			)
		}
	}

}
