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
	subscribers *sync.Map // subscribers is a concurrent map that holds mapping from a client ID to it's subscriber
}

func NewLiveServer(subscribers *sync.Map) *liveServer {
	return &liveServer{subscribers: subscribers}
}

type sub struct {
	stream   protos.LiveConnection_SubscribeServer // stream is the server side of the RPC stream
	finished chan<- bool                      // finished is used to signal closure of a client subscribing goroutine
	connData model.Info
}

// Subscribe handles a subscribe request from a client
func (s *liveServer) Subscribe(request *protos.ConnectRequest, stream protos.LiveConnection_SubscribeServer) error {
	// Handle subscribe request
	log.Printf("Received subscribe request from ID: %s", request.ConnectionData.Id)

	fin := make(chan bool)
	// Save the subscriber stream according to the given client ID

	s.subscribers.Store(request.ConnectionData.Id, sub{stream: stream, finished: fin,connData: model.Info{
		Id:    request.ConnectionData.Id,
		Port:  int(request.ConnectionData.Port),
		Files: request.ConnectionData.Songs,
	}})

	ctx := stream.Context()
	// Keep this scope alive because once this scope exits - the stream is closed
	for {
		select {
		case <-fin:
			log.Printf("Closing stream for client ID: %s", request.ConnectionData.Id)
			return nil
		case <- ctx.Done():
			log.Printf("Client ID %s has disconnected", request.ConnectionData.Id)
			return nil
		}
	}
}

// Unsubscribe handles a unsubscribe request from a client
// Note: this function is not called but it here as an example of an unary RPC for unsubscribing clients
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


func (s *liveServer) getClientInfo ()[]model.Info {
	connections:=make( []model.Info,0)
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
	isSentBefore:=false
	prev:=make([]model.Info,0)
	for {
		time.Sleep(time.Second)

		// A list of clients to unsubscribe in case of error
		var unsubscribe []string

		connections:=s.getClientInfo()
		if isSentBefore&&reflect.DeepEqual(prev,connections) {
			continue
		}
		protoConInfo:=make([]*protos.SongData,0)
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
			err := sub.stream.Send(&protos.Response{Songs:protoConInfo})
			if  err != nil {
				log.Printf("Failed to send data to client: %v", err)
				select {
				case sub.finished <- true:
					log.Printf("Unsubscribed client: %d", id)
				default:
					// Default case is to avoid blocking in case client has already unsubscribed
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
		prev=connections
		isSentBefore=true
	}
}

func ConvertToProtoConnInfo(data *model.Info) *protos.SongData {
	return &protos.SongData{
		Id:    data.Id,
		Port:  int32(data.Port),
		Songs: data.Files,
	}
}



//func send(s *liveServer) {
//
//	// A list of clients to unsubscribe in case of error
//	var unsubscribe []int32
//
//	// Iterate over all subscribers and send data to each client
//	s.subscribers.Range(func(k, v interface{}) bool {
//		id, ok := k.(int32)
//		if !ok {
//			log.Printf("Failed to cast subscriber key: %T", k)
//			return false
//		}
//		sub, ok := v.(sub)
//		if !ok {
//			log.Printf("Failed to cast subscriber value: %T", v)
//			return false
//		}
//		// Send data over the gRPC stream to the client
//		err := sub.stream.Send(&protos.Response{
//			Songs:
//			fmt.Sprintf("data mock for: %d", id)})
//		if  err != nil {
//			log.Printf("Failed to send data to client: %v", err)
//			select {
//			case sub.finished <- true:
//				log.Printf("Unsubscribed client: %d", id)
//			default:
//				// Default case is to avoid blocking in case client has already unsubscribed
//			}
//			// In case of error the client would re-subscribe so close the subscriber stream
//			unsubscribe = append(unsubscribe, id)
//		}
//		return true
//	})
//
//	// Unsubscribe erroneous client streams
//	for _, id := range unsubscribe {
//		s.subscribers.Delete(id)
//	}
//
//
//}