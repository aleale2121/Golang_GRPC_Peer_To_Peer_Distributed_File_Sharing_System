package conn_client

import (
	"context"
	"fmt"
	"github.com/aleale2121/DSP_LAB/Music_Service/constant/model"
	protos "github.com/aleale2121/DSP_LAB/Music_Service/grpc/server/services/connect"
	"github.com/jedib0t/go-pretty/v6/table"
	"google.golang.org/grpc"
	"log"
	"os"
	"time"
)

// liveClient holds the long lived gRPC client fields
type liveClient struct {
	client   protos.LiveConnectionClient // client is the long lived gRPC client
	conn     *grpc.ClientConn       // conn is the client gRPC connection
	connData model.Info
}

// NewLiveClient creates a new client instance
func NewLiveClient(
	client   protos.LiveConnectionClient ,// client is the long lived gRPC client
	conn     *grpc.ClientConn ,      // conn is the client gRPC connection
	connData model.Info) (*liveClient, error) {


	return &liveClient{
		client: client,
		conn:   conn,
		connData:connData,
	}, nil
}

// close is not used but is here as an example of how to close the gRPC client connection
func (c *liveClient) close() {
	if err := c.conn.Close(); err != nil {
		log.Fatal(err)
	}
}

// subscribe subscribes to messages from the gRPC server
func (c *liveClient) subscribe() (protos.LiveConnection_SubscribeClient, error) {
	log.Printf("Subscribing client ID: %s", c.connData.Id)
	return c.client.Subscribe(context.Background(), &protos.ConnectRequest{ConnectionData: &protos.SongData{
		Id:    c.connData.Id,
		Port: int32(c.connData.Port),
		Songs: c.connData.Files,
	}})
}

// unsubscribe unsubscribes to messages from the gRPC server
func (c *liveClient) unsubscribe() error {
	log.Printf("Unsubscribing client ID %s", c.connData.Id)
	_, err := c.client.Unsubscribe(context.Background(), &protos.UnSubscribeRequest{Id: c.connData.Id})
	return err
}

func (c *liveClient) Start() {
	var err error
	// stream is the client side of the RPC stream
	var stream protos.LiveConnection_SubscribeClient
	for {
		if stream == nil {
			if stream, err = c.subscribe(); err != nil {
				log.Printf("Failed to subscribe: %v", err)
				c.sleep()
				// Retry on failure
				continue
			}
		}
		response, err := stream.Recv()
		if err != nil {
			log.Printf("Failed to receive message: %v", err)
			// Clearing the stream will force the client to resubscribe on next iteration
			stream = nil
			c.sleep()
			// Retry on failure
			continue
		}

		DisplayLiveClientsAndSons(response.Songs)

		log.Printf("Client ID %s got response: %q", c.connData.Id, len(response.Songs))
	}
}

// sleep is used to give the server time to unsubscribe the client and reset the stream
func (c *liveClient) sleep() {
	time.Sleep(time.Second * 5)
}

func DisplayLiveClientsAndSons(songs []*protos.SongData)  {

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	//fmt.Println("Online Clients")
	t.AppendHeader(table.Row{"Online Clients","",""})
	t.AppendHeader(table.Row{ "CLIENT-ID", "IP", "MUSICS"})
	for _, s := range songs {
		fmt.Printf("Songs For Client %s AND IP:%d\n",s.Id,s.Port)
		fmt.Println("-------------------------------------------")
		songText:=""
		for _, song := range s.Songs {
			fmt.Println(song)
			songText+=song+"\n"
		}
		t.AppendRow([]interface{}{ s.Id, "127.0.0.1"+string(s.Port), songText})
		t.AppendSeparator()
	}
	t.SetStyle(table.StyleColoredBright)
	t.Render()


}


