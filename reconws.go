package reconws

import (
	"github.com/gorilla/websocket"
	"log"
)

type Client struct {
	conn             *websocket.Conn
	chans            channeler
	isShutdown       bool
	url              string
	cbs              callbacks
	subscriptionData []byte
}

type channeler struct {
	read  chan []byte
	recon chan bool
	done  chan bool
	quit  chan bool
}

type callbacks struct {
	onConnect    func()
	onDisconnect func()
	onReconnect  func()
}

func NewClient() *Client {
	c := &Client{
		conn: nil,
		url:  "",
		chans: channeler{
			read:  make(chan []byte),
			done:  make(chan bool),
			quit:  make(chan bool),
			recon: make(chan bool),
		},
		cbs: callbacks{
			onConnect:    func() {},
			onDisconnect: func() {},
			onReconnect:  func() {},
		},
	}

	go c.fireUpReconChannel()
	return c
}
func (c *Client) SetChannels(read chan []byte, done chan bool) *Client {
	c.chans.read = read
	c.chans.done = done
	return c
}
func (c *Client) OnConnect(fn func()) *Client {
	c.cbs.onConnect = fn
	return c
}
func (c *Client) OnDisconnect(fn func()) *Client {
	c.cbs.onDisconnect = fn
	return c
}
func (c *Client) OnReconnect(fn func()) *Client {
	c.cbs.onReconnect = fn
	return c
}
func (c *Client) SetSubscriptionData(data []byte) *Client {
	c.subscriptionData = data
	return c
}
func (c *Client) fireUpReconChannel() {
	for {
		select {
		case <-c.chans.done:
			return
		case <-c.chans.quit:
			return
		case <-c.chans.recon:
			_ = c.conn.Close()
			_, err := c.Connect(c.url)
			if err != nil {
				log.Fatal("reconnection channel err:", err)
				return
			}
			c.cbs.onReconnect()
		}
	}
}
func (c *Client) Connect(url string) (*websocket.Conn, error) {
	if c.url == "" {
		c.url = url
	}

	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, err
	}
	c.conn = conn

	if len(c.subscriptionData) > 0 {
		err = c.Write(websocket.TextMessage, c.subscriptionData)
		if err != nil {
			_ = c.conn.Close()
			return nil, err
		}
	}

	c.cbs.onConnect()
	go c.read()
	return conn, nil
}

func (c *Client) Write(msgType int, data []byte) error {
	if msgType == 0 {
		msgType = websocket.TextMessage
	}

	err := c.conn.WriteMessage(msgType, data)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) Close() {
	c.isShutdown = true
	c.chans.quit <- true
	err := c.conn.Close()
	if err != nil {
		return
	}
}
func (c *Client) ReadChan() chan []byte {
	return c.chans.read
}
func (c *Client) read() {
	defer func() {
		err := c.conn.Close()
		if err != nil {
			return
		}
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		if c.isShutdown {
			_ = c.conn.Close()
			return
		}

		c.chans.read <- message
	}

	c.cbs.onDisconnect()
	if !c.isShutdown {
		c.chans.recon <- true
	}
}
