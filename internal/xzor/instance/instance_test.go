package instance_test

import (
	"log"
	"net"
	"os"
	"testing"

	"github.com/xzor-dev/xzor/internal/module/messenger"
	messenger_command "github.com/xzor-dev/xzor/internal/module/messenger/command"
	"github.com/xzor-dev/xzor/internal/xzor/action"
	"github.com/xzor-dev/xzor/internal/xzor/command"
	"github.com/xzor-dev/xzor/internal/xzor/instance"
	"github.com/xzor-dev/xzor/internal/xzor/network"
	"github.com/xzor-dev/xzor/internal/xzor/resource"
	"github.com/xzor-dev/xzor/internal/xzor/storage"
	storage_file "github.com/xzor-dev/xzor/internal/xzor/storage/file"
	storage_json "github.com/xzor-dev/xzor/internal/xzor/storage/json"
)

func TestInstance(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("%v", err)
	}

	instanceA, _, nodeA := newTestSetup(t, dir+"/testdata/instanceA")
	instanceB, executorB, nodeB := newTestSetup(t, dir+"/testdata/instanceB")
	connA, connB := net.Pipe()

	nodeA.AddConnection(connA)
	nodeB.AddListener(&network.MockListener{
		Connections: []net.Conn{connB},
	})

	actionA, err := action.New(messenger.ModuleName, messenger_command.CreateBoardName, command.Params{
		"title": "My Board",
	})
	resA1, err := instanceA.ExecuteAction(actionA)
	if err != nil {
		t.Fatalf("%v", err)
	}

	boardA1, ok := resA1.Value.(*messenger.Board)
	if !ok {
		t.Fatalf("could not convert response value to board")
	}

	resourceA1, err := instanceA.Resource(messenger.ModuleName, boardA1.ResourceName(), boardA1.ResourceID())
	if err != nil {
		t.Fatalf("%v", err)
	}
	boardA2, ok := resourceA1.(*messenger.Board)
	if !ok {
		t.Fatal("could not convert resource to board")
	}
	if boardA2.Title != boardA1.Title {
		t.Fatalf("expected board title to be %s, got %s", boardA1.Title, boardA2.Title)
	}

	<-executorB.responseChan

	resourceB1, err := instanceB.Resource(messenger.ModuleName, boardA1.ResourceName(), boardA1.ResourceID())
	if err != nil {
		t.Fatalf("%v", err)
	}

	boardB1, ok := resourceB1.(*messenger.Board)
	if !ok {
		t.Fatalf("could not convert response value to board")
	}
	if boardB1.Title != boardA1.Title {
		t.Fatalf("expected board title to be %s, got %s", boardA1.Title, boardB1.Title)
	}

	err = os.RemoveAll(dir + "/testdata")
	if err != nil {
		t.Fatalf("%v", err)
	}
}

func newTestSetup(t *testing.T, rootDir string) (*instance.Instance, *actionExecutor, *network.Node) {
	msgMod := newMessengerModule(t, rootDir)
	node := network.NewNode()
	actionService := action.NewService([]command.Provider{msgMod})
	actionExecutor := newActionExecutor(actionService)
	resourceService := resource.NewService([]resource.Provider{msgMod})
	inst := instance.New(actionExecutor, node, resourceService)
	err := inst.Start()
	if err != nil {
		t.Fatalf("%v", err)
	}

	return inst, actionExecutor, node
}

func newMessengerModule(t *testing.T, rootDir string) *messenger.Module {
	storeED := &storage_json.EncodeDecoder{}
	fileStore := storage_file.NewRecordStore(rootDir + "/messenger")
	store := storage.NewService(storeED, fileStore)
	srv := messenger.NewService(store)
	commands := messenger_command.Commands(srv)

	return messenger.NewModule(srv, commands)
}

var _ action.Executor = &actionExecutor{}

type actionExecutor struct {
	actionService *action.Service
	responseChan  chan *action.Response
}

func newActionExecutor(actionService *action.Service) *actionExecutor {
	return &actionExecutor{
		actionService: actionService,
		responseChan:  make(chan *action.Response),
	}
}

func (e *actionExecutor) ExecuteAction(a *action.Action) (*action.Response, error) {
	res, err := e.actionService.ExecuteAction(a)
	if err != nil {
		log.Printf("got error from executing action: %v", err)
		return nil, err
	}
	go func() {
		e.responseChan <- res
	}()
	return res, nil
}
