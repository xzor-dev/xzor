package messenger_test

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/xzor-dev/xzor/internal/module/messenger"
	"github.com/xzor-dev/xzor/internal/module/messenger/command"
	"github.com/xzor-dev/xzor/internal/xzor/action"
	"github.com/xzor-dev/xzor/internal/xzor/module"
	"github.com/xzor-dev/xzor/internal/xzor/storage"
	"github.com/xzor-dev/xzor/internal/xzor/storage/file"
	"github.com/xzor-dev/xzor/internal/xzor/storage/json"
)

func TestModule(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("%v", err)
	}
	store := &storage.Service{
		EncodeDecoder: &json.EncodeDecoder{},
		Store: &file.RecordStore{
			RootDir: dir + "/testdata",
		},
	}
	srv := &messenger.Service{
		Storage: store,
	}
	mod := &messenger.Module{
		Commander: command.NewCommander(srv),
		Service:   srv,
	}
	actions := &action.Service{
		Modules: map[module.Name]module.Module{
			mod.Name(): mod,
		},
	}

	type testRun struct {
		Name     string
		Action   func() *action.Action
		Validate func(*action.Response) error
		TearDown func(*action.Response) error
	}

	var boardHash messenger.BoardHash
	var messageHash messenger.MessageHash
	var threadHash messenger.ThreadHash

	testRuns := []testRun{
		testRun{
			Name: "Create Board",
			Action: func() *action.Action {
				return &action.Action{
					Module:    mod.Name(),
					Command:   command.CreateBoardName,
					Arguments: []interface{}{"foo"},
				}
			},
			Validate: func(res *action.Response) error {
				board, ok := res.Value.(*messenger.Board)
				if !ok {
					return errors.New("could not convert response value to board struct")
				}
				if board.Title != "foo" {
					return fmt.Errorf("invalid board title: wanted %s, got %s", "foo", board.Title)
				}
				if board.Hash == "" {
					return errors.New("expected newly created board to have a hash")
				}
				_, err := srv.Board(board.Hash)
				if err != nil {
					return err
				}
				boardHash = board.Hash
				return nil
			},
		},
		testRun{
			Name: "Create Thread",
			Action: func() *action.Action {
				return &action.Action{
					Module:    mod.Name(),
					Command:   command.CreateThreadName,
					Arguments: []interface{}{boardHash, "Test Thread"},
				}
			},
			Validate: func(res *action.Response) error {
				thread, ok := res.Value.(*messenger.Thread)
				if !ok {
					return errors.New("could not convert response value to a Thread")
				}
				if thread.Title != "Test Thread" {
					return fmt.Errorf("unexpected thread title, wanted %s, got %s", "Test Thread", thread.Title)
				}
				if thread.Hash == "" {
					return errors.New("expected thread to have a hash")
				}
				_, err := srv.Thread(thread.Hash)
				if err != nil {
					return err
				}
				board, err := srv.Board(boardHash)
				if err != nil {
					return err
				}
				if !board.HasThread(thread.Hash) {
					return errors.New("expected board to have thread hash")
				}
				threadHash = thread.Hash
				return nil
			},
		},
		testRun{
			Name: "Create Message",
			Action: func() *action.Action {
				return &action.Action{
					Module:    mod.Name(),
					Command:   command.CreateMessageName,
					Arguments: []interface{}{threadHash, "hello world"},
				}
			},
			Validate: func(res *action.Response) error {
				message, ok := res.Value.(*messenger.Message)
				if !ok {
					return errors.New("could not convert response value to message")
				}
				if message.Body != "hello world" {
					return fmt.Errorf("unexpected message body: wanted %s, got %s", "hello world", message.Body)
				}
				if message.Hash == "" {
					return errors.New("message was not assigned a hash")
				}
				_, err := srv.Message(message.Hash)
				if err != nil {
					return err
				}
				thread, err := srv.Thread(threadHash)
				if err != nil {
					return err
				}
				if !thread.HasMessage(message.Hash) {
					return errors.New("expected thread to have message hash")
				}
				messageHash = message.Hash
				return nil
			},
		},
	}

	for _, r := range testRuns {
		t.Run(r.Name, func(t *testing.T) {
			act := r.Action()
			res, err := actions.Execute(act)
			if err != nil {
				t.Fatalf("%v", err)
			}

			if r.Validate != nil {
				err = r.Validate(res)
				if err != nil {
					t.Fatalf("%v", err)
				}
			}

			if r.TearDown != nil {
				err = r.TearDown(res)
				if err != nil {
					t.Fatalf("%v", err)
				}
			}
		})
	}

	t.Run("Cleanup", func(t *testing.T) {
		err := srv.DeleteBoard(boardHash)
		if err != nil {
			t.Fatalf("%v", err)
		}

		_, err = srv.Thread(threadHash)
		if err == nil {
			t.Fatal("expected thread to be deleted from board deletion")
		}

		_, err = srv.Message(messageHash)
		if err == nil {
			t.Fatal("expected message to be delete from board deletion")
		}
	})
}
