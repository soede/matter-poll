package service

import (
	"errors"
	"fmt"
	"github.com/mattermost/mattermost/server/public/model"
	"golang.org/x/exp/slog"
	"strconv"
	"votty/internal/storage/tarantool"
)

func Vote(storage *tarantool.Storage, log *slog.Logger, post *model.Post, parts []string) (r *model.Post) {
	if len(parts) < 3 {
		r = &model.Post{
			Message: "Произошла ошибка при обработке, команда голосования должна быть в формате ```/vote pollID 1```",
		}

		log.Warn("Failed to parse the /vote command",
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

	if !poll.IsActive {
		r = &model.Post{
			Message: "Опрос уже не актуален",
		}

		log.Warn("Failed to vote the poll",
			slog.String("user_id", post.UserId),
			slog.String("message", post.Message),
			slog.String("pollID", pollID),
		)

		return
	}

	if err != nil || poll == nil {
		r = &model.Post{
			Message: "Произошла какая то ошибка",
		}

		log.Error("Failed to find the poll",
			slog.String("user_id", post.UserId),
			slog.String("message", post.Message),
			slog.String("pollID", pollID),
		)
		return
	}

	choice, err := strconv.Atoi(parts[2])
	if err != nil || choice > len(poll.Options) || choice < 1 {
		r = &model.Post{
			Message: "Такого варианта не существует в опросе",
		}
		log.Warn("Failed to parse the /vote",
			slog.String("user_id", post.UserId),
			slog.String("message", post.Message),
		)
		return
	}

	v, err := storage.SelectVotes(pollID, post.UserId)

	if err != nil {
		r = &model.Post{
			Message: "Что-то пошло не так :(",
		}
		log.Warn("Failed to find the vote",
			slog.String("user_id", post.UserId),
			slog.String("message", post.Message),
			slog.String("pollID", pollID),
		)
	}

	if v != nil {
		err = storage.UpsertVote(pollID, post.UserId, uint64(choice-1))
		if err != nil {
			r = &model.Post{
				Message: "Что-то пошло не так :(",
			}
			log.Warn("Failed to update the vote",
				slog.String("user_id", post.UserId),
				slog.String("message", post.Message),
				slog.String("pollID", pollID),
			)
			return
		}

		r = &model.Post{
			Message: fmt.Sprintf("Ты успешно изменил свой выбор на %v (%s)", choice, poll.Options[choice-1]),
		}

		log.Warn("Choice has been edit",
			slog.String("user_id", post.UserId),
			slog.String("message", post.Message),
			slog.String("pollID", pollID),
			slog.Int("choice", choice),
		)
		return
	}

	err = storage.UpsertVote(pollID, post.UserId, uint64(choice-1))
	if err != nil {
		r = &model.Post{
			Message: "Что-то пошло не так :(",
		}
		log.Warn("Failed to update the vote",
			slog.String("user_id", post.UserId),
			slog.String("message", post.Message),
			slog.String("pollID", pollID),
		)
		return
	}

	r = &model.Post{
		Message: fmt.Sprintf("Ты успешно сделал свой голос %v (%s)", choice, poll.Options[choice-1]),
	}

	return
}
