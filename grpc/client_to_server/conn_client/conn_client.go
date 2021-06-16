package conn_client

import (
	"context"
	"fmt"
	"github.com/aleale2121/DSP_LAB/Music_Service/constant/model"
	protos "github.com/aleale2121/DSP_LAB/Music_Service/grpc/server/services/connect"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"google.golang.org/grpc"
	"log"
	"os"
	"strconv"
	"time"
)

// liveClient holds the long lived gRPC client fields
type liveClient struct {
	client   protos.LiveConnectionClient // client is the long lived gRPC client
	conn     *grpc.ClientConn            // conn is the client gRPC connection
	connData model.Info
}

var (
	changeMain = 0
	changeSub  = -1
)

// NewLiveClient creates a new client instance
func NewLiveClient(
	client protos.LiveConnectionClient, // client is the long lived gRPC client
	conn *grpc.ClientConn, // conn is the client gRPC connection
	connData model.Info) (*liveClient, error) {

	return &liveClient{
		client:   client,
		conn:     conn,
		connData: connData,
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
		Port:  int32(c.connData.Port),
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
		changeMain += 1
		DisplayLiveClientsAndSons(response.Songs)

		//log.Printf("Client ID %s got response: %q", c.connData.Id, len(response.Songs))
	}
}

// sleep is used to give the server time to unsubscribe the client and reset the stream
func (c *liveClient) sleep() {
	time.Sleep(time.Second * 5)
}

func DisplayLiveClientsAndSons(songs []*protos.SongData) {
	changeSub = changeMain
	for changeSub == changeMain {
		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.SetColumnConfigs([]table.ColumnConfig{
			{
				Name:         "Online Clients",
				Align:        text.AlignLeft,
				AlignFooter:  text.AlignLeft,
				AlignHeader:  text.AlignLeft,
				Colors:       text.Colors{text.BgBlack, text.FgRed},
				ColorsHeader: text.Colors{text.BgRed, text.FgBlack, text.Bold},
				ColorsFooter: text.Colors{text.BgRed, text.FgBlack},
				Hidden:       false,
				VAlign:       text.VAlignMiddle,
				VAlignFooter: text.VAlignTop,
				VAlignHeader: text.VAlignBottom,
				WidthMin:     20,
				WidthMax:     64,
			},
		})

		//fmt.Println("Online Clients")
		t.AppendHeader(table.Row{"CLIENT-ID", "IP", "MUSICS"})
		for _, s := range songs {

			songText := ""
			for _, song := range s.Songs {
				songText += song + "\n"
			}
			songText += ""
			ip := "127.0.0.1:" + strconv.Itoa(int(s.Port))
			t.AppendRow([]interface{}{s.Id, ip, songText})
			t.AppendSeparator()
		}
		//t.SetStyle(table.StyleColoredBright)
		t.SetAllowedRowLength(100)
		t.Render()

		fmt.Println("Enter")
		fmt.Println("1--------To Download Music")
		fmt.Println("2--------To Send Music")
		fmt.Println("3--------To Continue")

		var choice = ""
		_, _ = fmt.Scanln(&choice)
		if choice == "1" {
			fmt.Println("-------------1")
			//Download()
		} else if choice == "2" {
			fmt.Println("------------2")
			//Send()
		} else {
			fmt.Println("-------------default")

		}

		fmt.Printf("old:%d--------new:%d\n", changeMain, changeSub)
	}

}
func Download() {
	//reader := bufio.NewReader(os.Stdin)
	//choice, _ := reader.ReadString('\n')
	//fmt.Println(choice)
	fmt.Println("download")
}

func Send() {
	fmt.Println("send")
}
