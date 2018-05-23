package client

import (
	"fmt"
	"io"
	"net/url"
	"strconv"

	"github.com/gorilla/websocket"
)

const (
	followParam    = "follow"   // bool
	previousParam  = "previous" // bool
	tailParam      = "tail"     // int
	containerParam = "container"
)

const (
	// Maximum message size allowed from peer.
	maxMessageSize = 1024
)

type GetPodLogsParams struct {
	Namespace, Pod, Container string
	Previous, Follow          bool
	Tail                      int
}

func (client *Client) GetPodLogs(params GetPodLogsParams) (*io.PipeReader, error) {
	logUrl, err := client.podLogPath(params)
	if err != nil {
		return nil, err
	}
	conn, err := client.WSDialer.Dial(logUrl, nil)
	if err != nil {
		return nil, err
	}
	re, wr := io.Pipe()
	go client.logStream(conn, wr)
	return re, nil
}

func (client *Client) podLogPath(params GetPodLogsParams) (string, error) {
	queryUrl, err := url.Parse(fmt.Sprintf("/namespaces/%s/pods/%s/log", params.Namespace, params.Pod))
	if err != nil {
		return "", err
	}
	queryParams := queryUrl.Query()
	queryParams.Set(followParam, strconv.FormatBool(params.Follow))
	queryParams.Set(previousParam, strconv.FormatBool(params.Previous))
	queryParams.Set(tailParam, strconv.Itoa(params.Tail))
	queryParams.Set(containerParam, params.Container)
	queryUrl.RawQuery = queryParams.Encode()
	return queryUrl.String(), nil
}

func (client *Client) logStream(conn *websocket.Conn, out *io.PipeWriter) {
	defer conn.Close()
	conn.SetReadLimit(maxMessageSize)
	for {
		mtype, data, err := conn.NextReader()
		if err != nil {
			out.CloseWithError(err)
			return
		}
		if mtype != websocket.TextMessage && mtype != websocket.BinaryMessage {
			continue
		}

		_, err = io.Copy(out, data)
		switch err {
		case nil:
			// pass
		case io.ErrClosedPipe:
			return
		default:
			out.CloseWithError(err)
			return
		}
	}
}
