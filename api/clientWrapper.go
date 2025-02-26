package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/goph/emperror"
	"github.com/hashicorp/yamux"
	"github.com/mintance/go-uniqid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/metadata"
	"net"
	"net/url"
	"reflect"
	"time"
)

type ClientWrapper struct {
	instanceName        string
	session             **yamux.Session
	clientServiceClient *ClientServiceClient
	conn                *grpc.ClientConn
}

func NewClientWrapper(instanceName string, session **yamux.Session) *ClientWrapper {
	cw := &ClientWrapper{instanceName: instanceName,
		session:             session,
		clientServiceClient: nil,
		conn:                nil,
	}
	return cw
}

func (cw *ClientWrapper) connect() (err error) {
	if *cw.session == nil {
		cw.clientServiceClient = nil
		return errors.New(fmt.Sprintf("session closed"))
	}

	// it's a singleton
	if cw.clientServiceClient != nil {
		return nil
	}
	// gRPC dial over incoming net.Conn
	// singleton!!!
	doDial := cw.conn == nil
	if cw.conn != nil {
		if cw.conn.GetState() == connectivity.TransientFailure {
			cw.conn.Close()
			doDial = true
		}
		if cw.conn.GetState() == connectivity.Shutdown {
			doDial = true
		}
	}
	if doDial {
		cw.conn, err = grpc.Dial(":7777", grpc.WithInsecure(),
			grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
				if *cw.session == nil {
					return nil, errors.New(fmt.Sprintf("session %s closed", s))
				}
				return (*cw.session).Open()
			}),
		)
		if err != nil {
			return errors.New("cannot dial grpc connection to :7777")
		}
	}
	c := NewClientServiceClient(cw.conn)
	cw.clientServiceClient = &c
	return nil
}

func (cw *ClientWrapper) Ping(traceId string, targetInstance string) (string, error) {
	if traceId == "" {
		traceId = uniqid.New(uniqid.Params{"traceid_", false})
	}
	if err := cw.connect(); err != nil {
		return "", emperror.Wrapf(err, "cannot connect to %v", targetInstance)
	}

	ctx := metadata.AppendToOutgoingContext(context.Background(), "sourceInstance", cw.instanceName, "targetInstance", targetInstance, "traceId", traceId)
	pingResult, err := (*cw.clientServiceClient).Ping(ctx, &String{Value: "ping"})
	if err != nil {
		return "", emperror.Wrapf(err, "error pinging %v", targetInstance)
	}
	return pingResult.GetValue(), nil
}

func (cw *ClientWrapper) Click(traceId string, targetInstance string, waitFor string, timeout time.Duration, position interface{}) error {
	if traceId == "" {
		traceId = uniqid.New(uniqid.Params{"traceid_", false})
	}
	if err := cw.connect(); err != nil {
		return emperror.Wrapf(err, "cannot connect to %v", targetInstance)
	}

	ctx := metadata.AppendToOutgoingContext(context.Background(), "sourceInstance", cw.instanceName, "targetInstance", targetInstance, "traceId", traceId)

	clickMsg := &ClickMessage{
		Target:      nil,
		Timeout:     int64(timeout),
		Waitvisible: waitFor,
	}
	switch pos := position.(type) {
	case string:
		clickMsg.Target = &ClickMessage_Element{Element: pos}
	case struct{ x, y int64 }:
		clickMsg.Target = &ClickMessage_Coord{Coord: &MouseCoord{
			X: pos.x,
			Y: pos.y,
		}}
	default:
		return errors.New("invalid position type")
	}

	_, err := (*cw.clientServiceClient).MouseClick(ctx, clickMsg)
	if err != nil {
		return emperror.Wrapf(err, "MouseClick %v on %v failed", position, targetInstance)
	}
	return nil
}

func (cw *ClientWrapper) Navigate(traceId string, targetInstance string, url *url.URL, nextStatus string) error {
	if traceId == "" {
		traceId = uniqid.New(uniqid.Params{"traceid_", false})
	}
	if err := cw.connect(); err != nil {
		return emperror.Wrapf(err, "cannot connect to %v", targetInstance)
	}

	ctx := metadata.AppendToOutgoingContext(context.Background(), "sourceInstance", cw.instanceName, "targetInstance", targetInstance, "traceId", traceId)
	param := &NavigateParam{Url: url.String(), NextStatus: nextStatus}
	_, err := (*cw.clientServiceClient).Navigate(ctx, param)
	if err != nil {
		return emperror.Wrap(err, "error navigating client")
	}
	return nil
}

func (cw *ClientWrapper) StartBrowser(traceId string, targetInstance string, execOptions *map[string]interface{}) error {
	if traceId == "" {
		traceId = uniqid.New(uniqid.Params{"traceid_", false})
	}
	if err := cw.connect(); err != nil {
		return emperror.Wrapf(err, "cannot connect to %v", targetInstance)
	}

	ctx := metadata.AppendToOutgoingContext(context.Background(), "sourceInstance", cw.instanceName, "targetInstance", targetInstance, "traceId", traceId)

	flags := []*BrowserInitFlag{}
	for name, val := range *execOptions {
		bif := &BrowserInitFlag{}
		bif.Name = name
		switch val.(type) {
		case bool:
			bif.Value = &BrowserInitFlag_Bval{val.(bool)}
		case string:
			bif.Value = &BrowserInitFlag_Strval{val.(string)}
		default:
			return errors.New(fmt.Sprintf("invalid value type %v for %v=%v", reflect.TypeOf(val), name, val))
		}
		flags = append(flags, bif)
	}
	browserInitFlags := &BrowserInitFlags{Flags: flags}
	_, err := (*cw.clientServiceClient).StartBrowser(ctx, browserInitFlags)
	if err != nil {
		return emperror.Wrap(err, "error starting browser")
	}
	return nil
}

func (cw *ClientWrapper) ShutdownBrowser(traceId string, targetInstance string) error {
	if traceId == "" {
		traceId = uniqid.New(uniqid.Params{"traceid_", false})
	}
	if err := cw.connect(); err != nil {
		return emperror.Wrapf(err, "cannot connect to %v", targetInstance)
	}

	ctx := metadata.AppendToOutgoingContext(context.Background(), "sourceInstance", cw.instanceName, "targetInstance", targetInstance, "traceId", traceId)
	_, err := (*cw.clientServiceClient).ShutdownBrowser(ctx, &empty.Empty{})
	if err != nil {
		return emperror.Wrapf(err, "error shutting down browser of %v", targetInstance)
	}
	return nil
}

func (cw *ClientWrapper) GetBrowserLog(traceId string, targetInstance string) ([]string, error) {
	if traceId == "" {
		traceId = uniqid.New(uniqid.Params{"traceid_", false})
	}
	if err := cw.connect(); err != nil {
		return nil, emperror.Wrapf(err, "cannot connect to %v", targetInstance)
	}

	ctx := metadata.AppendToOutgoingContext(context.Background(), "sourceInstance", cw.instanceName, "targetInstance", targetInstance, "traceId", traceId)
	ret, err := (*cw.clientServiceClient).GetBrowserLog(ctx, &empty.Empty{})
	if err != nil {
		return nil, emperror.Wrapf(err, "error getting status of %v", targetInstance)
	}
	return ret.GetEntry(), nil
}

func (cw *ClientWrapper) GetHTTPSAddr(traceId string, targetInstance string) (string, error) {
	if traceId == "" {
		traceId = uniqid.New(uniqid.Params{"traceid_", false})
	}
	if err := cw.connect(); err != nil {
		return "", emperror.Wrapf(err, "cannot connect to %v", targetInstance)
	}

	ctx := metadata.AppendToOutgoingContext(context.Background(), "sourceInstance", cw.instanceName, "targetInstance", targetInstance, "traceId", traceId)
	ret, err := (*cw.clientServiceClient).GetHTTPSAddr(ctx, &empty.Empty{})
	if err != nil {
		return "", emperror.Wrapf(err, "error getting status of %v", targetInstance)
	}
	return ret.GetValue(), nil
}

func (cw *ClientWrapper) GetStatus(traceId string, targetInstance string) (string, error) {
	if traceId == "" {
		traceId = uniqid.New(uniqid.Params{"traceid_", false})
	}
	if err := cw.connect(); err != nil {
		return "", emperror.Wrapf(err, "cannot connect to %v", targetInstance)
	}

	ctx := metadata.AppendToOutgoingContext(context.Background(), "sourceInstance", cw.instanceName, "targetInstance", targetInstance, "traceId", traceId)
	ret, err := (*cw.clientServiceClient).GetStatus(ctx, &empty.Empty{})
	if err != nil {
		return "", emperror.Wrapf(err, "error getting status of %v", targetInstance)
	}
	return ret.GetValue(), nil
}

func (cw *ClientWrapper) SetStatus(traceId string, targetInstance string, stat string) error {
	if traceId == "" {
		traceId = uniqid.New(uniqid.Params{"traceid_", false})
	}
	if err := cw.connect(); err != nil {
		return emperror.Wrapf(err, "cannot connect to %v", targetInstance)
	}

	ctx := metadata.AppendToOutgoingContext(context.Background(), "sourceInstance", cw.instanceName, "targetInstance", targetInstance, "traceId", traceId)
	_, err := (*cw.clientServiceClient).SetStatus(ctx, &String{Value: stat})
	if err != nil {
		return emperror.Wrapf(err, "error setting status of %v", targetInstance)
	}
	return nil
}

func (cw *ClientWrapper) WebsocketMessage(traceId string, sourceInstance string, targetInstance string, data []byte) error {
	if traceId == "" {
		traceId = uniqid.New(uniqid.Params{"traceid_", false})
	}
	if err := cw.connect(); err != nil {
		return emperror.Wrapf(err, "cannot connect to %v", targetInstance)
	}

	if sourceInstance == "" {
		sourceInstance = cw.instanceName
	}
	ctx := metadata.AppendToOutgoingContext(context.Background(),
		"sourceInstance", sourceInstance,
		"targetInstance", targetInstance,
		"traceId", traceId)
	_, err := (*cw.clientServiceClient).WebsocketMessage(ctx, &Bytes{Value: data})
	if err != nil {
		return emperror.Wrapf(err, "error sending websocket message to %v", targetInstance)
	}

	return nil
}

func (cw *ClientWrapper) MouseClick(traceId string, sourceInstance string, targetInstance string, x, y int64) error {
	if traceId == "" {
		traceId = uniqid.New(uniqid.Params{"traceid_", false})
	}
	if err := cw.connect(); err != nil {
		return emperror.Wrapf(err, "cannot connect to %v", targetInstance)
	}

	if sourceInstance == "" {
		sourceInstance = cw.instanceName
	}
	ctx := metadata.AppendToOutgoingContext(context.Background(),
		"sourceInstance", sourceInstance,
		"targetInstance", targetInstance,
		"traceId", traceId)
	_, err := (*cw.clientServiceClient).MouseClick(ctx, &ClickMessage{
		Target: &ClickMessage_Coord{
			Coord: &MouseCoord{
				X: x,
				Y: y,
			},
		},
	})
	if err != nil {
		return emperror.Wrapf(err, "error clicking mouse on %v", targetInstance)
	}

	return nil
}
