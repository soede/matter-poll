package tarantool

import (
	"context"
	"errors"
	"fmt"
	"github.com/tarantool/go-tarantool/v2"
	"golang.org/x/exp/slog"
	"time"
	"votty/internal/config"
	"votty/internal/models"
)

var (
	ErrNotFound = errors.New("data not found")
)

type Storage struct {
	Conn *tarantool.Connection
}

func New(log *slog.Logger, cfg *config.Config) *Storage {
	var conn *tarantool.Connection
	dialer := tarantool.NetDialer{
		Address:  cfg.TarantoolHost,
		User:     cfg.TarantoolUser,
		Password: cfg.TarantoolPassword,
	}

	opts := tarantool.Opts{
		Timeout: time.Second,
	}

	var err error

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	conn, err = tarantool.Connect(ctx, dialer, opts)
	if err != nil {
		log.Error("Tarantool connection error", err.Error())
		return nil
	}
	log.Info("Successfully connected to Tarantool!")
	return &Storage{Conn: conn}
}

func (s *Storage) CreatePoll(id, owner, question string, options []string) error {
	request := tarantool.NewInsertRequest("polls").Tuple([]interface{}{id, owner, question, options, true})

	future := s.Conn.Do(request)

	_, err := future.Get()
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) GetPoll(id string) (*models.Poll, error) {
	data, err := s.Conn.Do(
		tarantool.NewSelectRequest("polls").
			Limit(1).
			Iterator(tarantool.IterEq).
			Key([]interface{}{id}),
	).Get()

	if err != nil {
		return nil, err
	}

	if len(data) > 0 {
		tuple := data[0].([]interface{})
		poll := models.Poll{
			ID:       tuple[0].(string),
			OwnerID:  tuple[1].(string),
			Question: tuple[2].(string),
			Options:  toStringSlice(tuple[3].([]interface{})),
			IsActive: tuple[4].(bool),
		}
		return &poll, nil
	} else {
		return nil, ErrNotFound
	}
}

func (s *Storage) DeletePoll(id string) error {
	_, err := s.Conn.Do(
		tarantool.NewDeleteRequest("polls").
			Key([]interface{}{id}),
	).Get()

	return err
}

func (s *Storage) SelectVotes(pollID, userID string) (*models.Vote, error) {
	data, err := s.Conn.Do(
		tarantool.NewSelectRequest("votes").
			Limit(1).
			Iterator(tarantool.IterEq).
			Key([]interface{}{pollID, userID}),
	).Get()

	if err != nil {
		return nil, err
	}

	if len(data) > 0 {
		tuple := data[0].([]interface{})
		poll := models.Vote{
			PollID: tuple[0].(string),
			UserID: tuple[1].(string),
			Choice: tuple[2].(uint64),
		}
		return &poll, nil
	}
	return nil, ErrNotFound
}

func (s *Storage) UpsertVote(pollID, userID string, choice uint64) error {
	_, err := s.Conn.Do(
		tarantool.NewUpsertRequest("votes").
			Tuple([]interface{}{pollID, userID, choice}).
			Operations(tarantool.NewOperations().Assign(2, choice)),
	).Get()
	if err != nil {
		return err
	}
	return nil

}

func (s *Storage) PollResults(pollID string, optionsSize int) ([]int, error) {

	data, err := s.Conn.Do(
		tarantool.NewSelectRequest("votes").
			Iterator(tarantool.IterEq).
			Key([]interface{}{pollID}),
	).Get()

	if err != nil {
		return nil, err
	}

	voteCounts := make([]int, optionsSize)
	for _, record := range data {
		tuple := record.([]interface{})
		choice := tuple[2].(uint64)

		voteCounts[int(choice)]++
	}
	return voteCounts, nil
}

func (s *Storage) EndPoll(pollID string) error {
	request := tarantool.NewUpdateRequest("polls").
		Key([]interface{}{pollID}).
		Operations(tarantool.NewOperations().Assign(4, false))

	_, err := s.Conn.Do(request).Get()
	if err != nil {
		return fmt.Errorf("failed to end poll: %w", err)
	}

	return nil
}
func toStringSlice(data []interface{}) []string {
	result := make([]string, len(data))
	for i, v := range data {
		result[i] = v.(string)
	}
	return result
}
