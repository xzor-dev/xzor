package main

import (
	"github.com/xzor-dev/xzor/internal/module/messenger"
	messenger_command "github.com/xzor-dev/xzor/internal/module/messenger/command"
	"github.com/xzor-dev/xzor/internal/xzor/action"
	"github.com/xzor-dev/xzor/internal/xzor/module"
	"github.com/xzor-dev/xzor/internal/xzor/network"
	"github.com/xzor-dev/xzor/internal/xzor/storage"
	"github.com/xzor-dev/xzor/internal/xzor/storage/file"
	"github.com/xzor-dev/xzor/internal/xzor/storage/json"
)

type instance struct {
	actionService *action.Service
	node          *network.Node
}

func (i *instance) Execute(a *action.Action) (*action.Response, error) {
	return i.actionService.Execute(a)
}

func newInstance(dataDir string) (*instance, error) {
	messengerModule, err := newMessengerModule(dataDir)
	if err != nil {
		return nil, err
	}

	actionService := action.NewService([]module.Module{
		messengerModule,
	})
	return &instance{
		actionService: actionService,
		node:          network.NewNode(),
	}, nil
}

func newMessengerModule(dataDir string) (*messenger.Module, error) {
	recordEncodeDecoder := &json.EncodeDecoder{}
	recordStore := file.NewRecordStore(dataDir + "/messenger")
	storage := storage.NewService(recordEncodeDecoder, recordStore)
	service := messenger.NewService(storage)
	commander := messenger_command.NewCommander(service)

	return messenger.NewModule(service, commander), nil
}
