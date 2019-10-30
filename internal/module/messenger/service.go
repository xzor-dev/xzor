package messenger

import "github.com/xzor-dev/xzor/internal/xzor/storage"

// Service handles messaging.
type Service struct {
	Storage *storage.Service
}

// NewService creates a new messenger service.
func NewService(storage *storage.Service) *Service {
	return &Service{
		Storage: storage,
	}
}

func (s *Service) boardID(hash BoardHash) storage.RecordID {
	id := "board-" + string(hash)
	return storage.RecordID(id)
}

func (s *Service) messageID(hash MessageHash) storage.RecordID {
	id := "message-" + string(hash)
	return storage.RecordID(id)
}

func (s *Service) threadID(hash ThreadHash) storage.RecordID {
	id := "thread-" + string(hash)
	return storage.RecordID(id)
}

// Board returns a Board by its hash.
func (s *Service) Board(hash BoardHash) (*Board, error) {
	board := &Board{}
	id := s.boardID(hash)
	err := s.Storage.Read(id, board)
	if err != nil {
		return nil, err
	}
	return board, nil
}

// DeleteBoard removes a board and all its threads from the storage.
func (s *Service) DeleteBoard(hash BoardHash) error {
	board, err := s.Board(hash)
	if err != nil {
		return err
	}
	if board.Threads != nil {
		for _, threadHash := range board.Threads {
			err := s.DeleteThread(threadHash)
			if err != nil {
				return err
			}
		}
	}
	id := s.boardID(hash)
	return s.Storage.Delete(id)
}

// DeleteMessage removes a message from the storage.
func (s *Service) DeleteMessage(hash MessageHash) error {
	id := s.messageID(hash)
	return s.Storage.Delete(id)
}

// DeleteThread removes a thread and all its messages from the storage.
func (s *Service) DeleteThread(hash ThreadHash) error {
	thread, err := s.Thread(hash)
	if err != nil {
		return err
	}
	if thread.Messages != nil {
		for _, messageHash := range thread.Messages {
			err := s.DeleteMessage(messageHash)
			if err != nil {
				return err
			}
		}
	}
	id := s.threadID(hash)
	return s.Storage.Delete(id)
}

// Message returns a message by its hash.
func (s *Service) Message(hash MessageHash) (*Message, error) {
	message := &Message{}
	id := s.messageID(hash)
	err := s.Storage.Read(id, message)
	if err != nil {
		return nil, err
	}
	return message, nil
}

// NewBoard creates a new board with the provided title.
func (s *Service) NewBoard(title string) (*Board, error) {
	hash, err := NewBoardHash(title)
	if err != nil {
		return nil, err
	}
	board := &Board{
		Hash:  hash,
		Title: title,
	}
	err = s.SetBoard(board)
	if err != nil {
		return nil, err
	}
	return board, nil
}

// NewMessage creates a new message in a thread.
func (s *Service) NewMessage(threadHash ThreadHash, body string) (*Message, error) {
	thread, err := s.Thread(threadHash)
	if err != nil {
		return nil, err
	}
	message, err := thread.NewMessage(body)
	if err != nil {
		return nil, err
	}
	err = s.SetMessage(message)
	if err != nil {
		return nil, err
	}
	err = s.SetThread(thread)
	if err != nil {
		return nil, err
	}
	return message, nil
}

// NewThread creates a new thread within a board.
func (s *Service) NewThread(boardHash BoardHash, title string) (*Thread, error) {
	board, err := s.Board(boardHash)
	if err != nil {
		return nil, err
	}
	thread, err := board.NewThread(title)
	if err != nil {
		return nil, err
	}
	err = s.SetThread(thread)
	if err != nil {
		return nil, err
	}
	err = s.SetBoard(board)
	if err != nil {
		return nil, err
	}
	return thread, nil
}

// SetBoard saves a board to the storage.
func (s *Service) SetBoard(board *Board) error {
	id := s.boardID(board.Hash)
	return s.Storage.Write(id, board)
}

// SetMessage saves a message to the storage.
func (s *Service) SetMessage(message *Message) error {
	id := s.messageID(message.Hash)
	return s.Storage.Write(id, message)
}

// SetThread saves a thread to the storage.
func (s *Service) SetThread(thread *Thread) error {
	id := s.threadID(thread.Hash)
	return s.Storage.Write(id, thread)
}

// Thread gets a thread by its hash.
func (s *Service) Thread(hash ThreadHash) (*Thread, error) {
	thread := &Thread{}
	id := s.threadID(hash)
	err := s.Storage.Read(id, thread)
	if err != nil {
		return nil, err
	}
	return thread, nil
}
