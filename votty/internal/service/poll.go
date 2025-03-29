package service

import (
	"errors"
	"fmt"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/mattermost/mattermost/server/public/model"
	"golang.org/x/exp/slog"
	"strings"
	"votty/internal/storage/tarantool"
)

func CreatePoll(log *slog.Logger, storage *tarantool.Storage, post *model.Post, parts []string) (r *model.Post) {
	if len(parts) < 3 {
		r = &model.Post{
			Message: "Произошла ошибка при обработки команды, запрос на создание должен быть в формате ```/create Вопрос? | Вариант1 | Вариант2 | Вариант3```",
		}

		log.Warn("Failed to parse the /create command",
			slog.String("user_id", post.UserId),
			slog.String("message", post.Message),
		)
		return
	}

	question := parts[1]
	options := strings.Split(parts[2], "|")

	for i := range options {
		options[i] = strings.TrimSpace(options[i])
	}

	id, err := gonanoid.New(10)
	if err != nil {
		r = &model.Post{
			Message: "Произошла ошибка при создании опроса :(",
		}

		log.Error("Failed to create the id for poll",
			slog.String("user_id", post.UserId),
			slog.String("pollQuestion", question),
			slog.String("error", err.Error()),
		)

		return
	}
	err = storage.CreatePoll(id, post.UserId, question, options)

	if err != nil {
		r = &model.Post{
			Message: "Произошла ошибка при создании опроса :(",
		}

		log.Error("Failed to create the poll",
			slog.String("user_id", post.UserId),
			slog.String("pollQuestion", question),
			slog.String("error", err.Error()),
		)

		return
	}

	log.Info("Create",
		slog.String("user_id", post.UserId),
		slog.String("pollQuestion", question),
	)

	message := fmt.Sprintf("Голосование \"%s\" было создано!\nID: ```%s```\nВарианты ответов:\n", question, id)

	for i, option := range options {
		message += fmt.Sprintf("\t%v. %s\n", i+1, option)
	}
	message += fmt.Sprintf("Перешли это сообщения всем участникам\nНапример, для того чтобы проголосовать за 1 вариант (%v) нужно отправить команду ```/vote %v 1```", options[0], id)
	r = &model.Post{
		Message: message,
	}
	return

}

func DeletePoll(storage *tarantool.Storage, log *slog.Logger, post *model.Post, parts []string) (r *model.Post) {
	if len(parts) < 2 {
		r = &model.Post{
			Message: "Произошла ошибка при обработки команды, запрос на удаление должен быть в формате ```/delete pollID```, где pollID – id опроса",
		}

		log.Warn("Failed to parse the /delete command",
			slog.String("user_id", post.UserId),
			slog.String("message", post.Message),
		)
		return
	}

	pollID := parts[1]

	poll, err := storage.GetPoll(pollID)

	if errors.Is(tarantool.ErrNotFound, err) {
		r = &model.Post{
			Message: "Такого опроса не существует :(",
		}

		log.Warn("Failed to find the poll",
			slog.String("user_id", post.UserId),
			slog.String("message", post.Message),
			slog.String("pollID", pollID),
		)

		return
	}
	if err != nil {
		r = &model.Post{
			Message: "Произошла какая то ошибка",
		}

		log.Error("error on get poll",
			slog.String("user_id", post.UserId),
			slog.String("message", post.Message),
			slog.String("pollID", pollID),
		)
		return
	}

	if poll.OwnerID != post.UserId {
		r = &model.Post{
			Message: "Ты не можешь удалить этот опрос, потому что ты не являешься его владельцем",
		}
		log.Warn("unauthorized: user are not the owner of this poll",
			slog.String("user_id", post.UserId),
			slog.String("message", post.Message),
			slog.String("pollID", pollID),
		)
		return
	}
	err = storage.DeletePoll(pollID)
	if err != nil {
		r = &model.Post{
			Message: "Произошла какая то ошибка :(",
		}
		log.Error("error on delete poll",
			slog.String("user_id", post.UserId),
			slog.String("message", post.Message),
			slog.String("pollID", pollID),
		)
		return
	}
	r = &model.Post{
		Message: fmt.Sprintf("Голосование ```%s``` было удалено", pollID),
	}

	log.Info("poll has been delete",
		slog.String("user_id", post.UserId),
		slog.String("message", post.Message),
		slog.String("pollID", pollID),
	)
	return
}

func PullResults(log *slog.Logger, storage *tarantool.Storage, post *model.Post, parts []string) (r *model.Post) {
	if len(parts) < 2 {
		r = &model.Post{
			Message: "Произошла ошибка при обработке, запрос на результаты должен быть в формате ```/results pollID```",
		}

		log.Warn("Failed to parse the /results command",
			slog.String("user_id", post.UserId),
			slog.String("message", post.Message),
		)
		return
	}

	pollID := parts[1]

	poll, err := storage.GetPoll(pollID)

	if errors.Is(tarantool.ErrNotFound, err) {
		r = &model.Post{
			Message: "Такого опроса не существует :(",
		}

		log.Warn("Failed to find the poll",
			slog.String("user_id", post.UserId),
			slog.String("message", post.Message),
			slog.String("pollID", pollID),
		)

		return
	}
	if err != nil {
		r = &model.Post{
			Message: "Произошла какая то ошибка",
		}

		log.Error("error on get poll",
			slog.String("user_id", post.UserId),
			slog.String("message", post.Message),
			slog.String("pollID", pollID),
		)
		return
	}

	votes, err := storage.PollResults(pollID, len(poll.Options))

	message := fmt.Sprintf("Результаты опроса для ```%s```\nВопрос: %s \nСоздатель: ```%s```\n", poll.ID, poll.Question, poll.OwnerID)

	if poll.IsActive {
		message += "Статус: активен\n"
	}
	if !poll.IsActive {
		message += "Статус: завершен\n"
	}
	for i, option := range poll.Options {
		message += fmt.Sprintf("\t%v. %s: %v\n", i+1, option, votes[i])
	}

	log.Info("Send poll results",
		slog.String("user_id", post.UserId),
		slog.String("pollID", pollID),
	)
	r = &model.Post{
		Message: message,
	}
	return

}

func EndPoll(storage *tarantool.Storage, log *slog.Logger, post *model.Post, parts []string) (r *model.Post) {
	if len(parts) < 2 {
		r = &model.Post{
			Message: "Произошла ошибка при обработки команды, запрос на завершение должен быть в формате ```/end pollID```, где pollID – id опроса",
		}

		log.Warn("Failed to parse the /end command",
			slog.String("user_id", post.UserId),
			slog.String("message", post.Message),
		)
		return
	}

	pollID := parts[1]

	poll, err := storage.GetPoll(pollID)

	if errors.Is(tarantool.ErrNotFound, err) {
		r = &model.Post{
			Message: "Такого опроса не существует :(",
		}

		log.Warn("Failed to find the poll",
			slog.String("user_id", post.UserId),
			slog.String("message", post.Message),
			slog.String("pollID", pollID),
		)

		return
	}
	if err != nil {
		r = &model.Post{
			Message: "Произошла какая то ошибка",
		}

		log.Error("error on get poll",
			slog.String("user_id", post.UserId),
			slog.String("message", post.Message),
			slog.String("pollID", pollID),
		)
		return
	}

	if poll.OwnerID != post.UserId {
		r = &model.Post{
			Message: "Ты не можешь завершить этот опрос, потому что ты не являешься его владельцем",
		}
		log.Warn("unauthorized: user are not the owner of this poll",
			slog.String("user_id", post.UserId),
			slog.String("message", post.Message),
			slog.String("pollID", pollID),
		)
		return
	}
	err = storage.EndPoll(pollID)
	if err != nil {
		r = &model.Post{
			Message: "Произошла какая то ошибка :(",
		}
		log.Error("error on end poll",
			slog.String("user_id", post.UserId),
			slog.String("message", post.Message),
			slog.String("pollID", pollID),
		)
		return
	}
	r = &model.Post{
		Message: fmt.Sprintf("Голосование ```%s``` было завершено. Результаты можно получить отправив ```/results %s```", pollID, pollID),
	}

	log.Info("poll has been delete",
		slog.String("user_id", post.UserId),
		slog.String("message", post.Message),
		slog.String("pollID", pollID),
	)
	return
}
