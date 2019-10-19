package instance_test

import (
	"net"
	"os"
	"testing"

	"github.com/xzor-dev/xzor/internal/module/messenger"
	messenger_command "github.com/xzor-dev/xzor/internal/module/messenger/command"
	"github.com/xzor-dev/xzor/internal/xzor/action"
	"github.com/xzor-dev/xzor/internal/xzor/instance"
	"github.com/xzor-dev/xzor/internal/xzor/module"
	"github.com/xzor-dev/xzor/internal/xzor/network"
	"github.com/xzor-dev/xzor/internal/xzor/storage"
	storage_file "github.com/xzor-dev/xzor/internal/xzor/storage/file"
	storage_json "github.com/xzor-dev/xzor/internal/xzor/storage/json"
)

func TestInstance(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("%v", err)
	}

	connA, connB := net.Pipe()
	nodeA := network.NewNode()
	nodeA.AddConnection(connA)
	nodeB := network.NewNode()
	nodeB.AddListener(&network.MockListener{
		Connections: []net.Conn{connB},
	})

	iA := newInstance(t, nodeA, dir+"/testdata/instanceA")
	iB := newInstance(t, nodeB, dir+"/testdata/instanceB")

	resA1, err := iA.ExecuteAction(&action.Action{
		Module:    messenger.ModuleName,
		Command:   messenger_command.CreateBoardName,
		Arguments: []interface{}{"My Board"},
	})
	if err != nil {
		t.Fatalf("%v", err)
	}

	boardA1, ok := resA1.Value.(*messenger.Board)
	if !ok {
		t.Fatalf("could not convert response value to board")
	}

	resourceA1, err := iA.Resource(messenger.ModuleName, boardA1.ResourceName(), boardA1.ResourceID())
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

	resourceB1, err := iB.Resource(messenger.ModuleName, boardA1.ResourceName(), boardA1.ResourceID())
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

func newInstance(t *testing.T, node *network.Node, rootDir string) *instance.Instance {
	msgMod := newMessengerModule(t, rootDir)
	modules := []module.Module{msgMod}

	i := instance.New(modules, node)
	err := i.Start()
	if err != nil {
		t.Fatalf("%v", err)
	}
	return i
}

func newMessengerModule(t *testing.T, rootDir string) *messenger.Module {
	storeED := &storage_json.EncodeDecoder{}
	fileStore := storage_file.NewRecordStore(rootDir + "/messenger")
	store := storage.NewService(storeED, fileStore)
	srv := messenger.NewService(store)
	cmdr := messenger_command.NewCommander(srv)

	return messenger.NewModule(srv, cmdr)
}
