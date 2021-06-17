package connection

import (
	"context"
	"fmt"
	"github.com/aleale2121/DSP_LAB/Music_Service/constant/model"
	protos "github.com/aleale2121/DSP_LAB/Music_Service/grpc/server/services/connect"
	"log"
	"reflect"
	"sync"
	"time"
	//"time"
)

type liveServer struct {
	protos.UnimplementedLiveConnectionServer
	subscribers *sync.Map
}

func NewLiveServer(subscribers *sync.Map) *liveServer {
	return &liveServer{subscribers: subscribers}
}

type sub struct {
	stream   protos.LiveConnection_SubscribeServer // stream is the server side of the RPC stream
	finished chan<- bool
	connData model.Info
}

func (s *liveServer) Subscribe(request *protos.ConnectRequest, stream protos.LiveConnection_SubscribeServer) error {
	log.Printf("Received subscribe request from ID: %s", request.ConnectionData.Id)

	finished := make(chan bool)

	s.subscribers.Store(request.ConnectionData.Id, sub{stream: stream, finished: finished, connData: model.Info{
		Id:    request.ConnectionData.Id,
		Port:  int(request.ConnectionData.Port),
		Files: request.ConnectionData.Songs,
	}})

	ctx := stream.Context()
	// Keep this scope alive because once this scope exits - the stream is closed
	for {
		select {
		case <-finished:
			log.Printf("Closing stream for client ID: %s", request.ConnectionData.Id)
			return nil
		case <-ctx.Done():
			log.Printf("Client ID %s has disconnected", request.ConnectionData.Id)
			return nil
		}
	}
}

func (s *liveServer) Unsubscribe(ctx context.Context, request *protos.UnSubscribeRequest) (*protos.Response, error) {
	v, ok := s.subscribers.Load(request.Id)
	if !ok {
		return nil, fmt.Errorf("failed to load subscriber key: %d", request.Id)
	}
	sub, ok := v.(sub)
	if !ok {
		return nil, fmt.Errorf("failed to cast subscriber value: %T", v)
	}
	select {
	case sub.finished <- true:
		log.Printf("Unsubscribed client: %s", request.Id)
	default:
		// Default case is to avoid blocking in case client has already unsubscribed
	}
	s.subscribers.Delete(request.Id)
	return &protos.Response{}, nil
}

func (s *liveServer) getClientInfo() []model.Info {
	connections := make([]model.Info, 0)
	s.subscribers.Range(func(k, v interface{}) bool {
		_, ok := k.(string)
		if !ok {
			log.Printf("Failed to cast subscriber key: %T", k)
			return false
		}
		sub, ok := v.(sub)
		if !ok {
			log.Printf("Failed to cast subscriber value: %T", v)
		}
		connections = append(connections, sub.connData)
		return true
	},
	)
	return connections
}
func (s *liveServer) StartStreamLiveClients() {
	log.Println("Starting data generation")
	isSentBefore := false
	prev := make([]model.Info, 0)
	for {
		time.Sleep(time.Second)

		// A list of clients to unsubscribe in case of error
		var unsubscribe []string

		connections := s.getClientInfo()
		if isSentBefore && reflect.DeepEqual(prev, connections) {
			continue
		}
		protoConInfo := make([]*protos.SongData, 0)
		for _, c := range connections {
			protoConInfo = append(protoConInfo, ConvertToProtoConnInfo(&c))
		}
		// Iterate over all subscribers and send data to each client
		s.subscribers.Range(func(k, v interface{}) bool {
			id, ok := k.(string)
			if !ok {
				log.Printf("Failed to cast subscriber key: %T", k)
				return false
			}
			sub, ok := v.(sub)
			if !ok {
				log.Printf("Failed to cast subscriber value: %T", v)
				return false
			}
			// Send data over the gRPC stream to the client
			err := sub.stream.Send(&protos.Response{Songs: protoConInfo})
			if err != nil {
				log.Printf("Failed to send data to client: %v", err)
				select {
				case sub.finished <- true:
					log.Printf("Unsubscribed client: %s", id)
				default:
				}
				// In case of error the client would re-subscribe so close the subscriber stream
				unsubscribe = append(unsubscribe, id)
			}
			return true
		})

		// Unsubscribe erroneous client streams
		for _, id := range unsubscribe {
			s.subscribers.Delete(id)
		}
		prev = connections
		isSentBefore = true
	}
}

func ConvertToProtoConnInfo(data *model.Info) *protos.SongData {
	return &protos.SongData{
		Id:    data.Id,
		Port:  int32(data.Port),
		Songs: data.Files,
	}
}
