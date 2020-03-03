package rpc

import (
    "net"
    "testing"
    "time"
)

import (
    log "github.com/AlexStocks/log4go"
    jerrors "github.com/juju/errors"
    "github.com/stretchr/testify/suite"
)

type ClientTestSuite struct {
    suite.Suite
    client  *Client
    server  *Server
}

func (suite *ClientTestSuite) SetupTest() {
    var err error
    initConf()
    suite.server, err = NewServer(serverConf)
    suite.Nil(err)
    err = suite.server.Register(&TestService{})
    suite.Nil(err)
    suite.server.Start()
    suite.client, err = NewClient(clientConf)
    suite.Nil(err)
}

func (suite *ClientTestSuite) TearDownTest() {
    suite.server.Stop()
    suite.client.Close()
    suite.Nil(suite.client.pool)
}

func (suite *ClientTestSuite) TestClient_Json_CallOneway() {
    var err error
    ts := TestService{}
    addr := net.JoinHostPort(serverConf.Host, serverConf.Ports[0])

    eventReq := EventReq{}
    err = suite.client.CallOneway(CodecJson, addr, ts.Service(), "Event", &eventReq,
        WithCallRequestTimeout(100e6), WithCallResponseTimeout(100e6))
    suite.Nil(err)
}

func (suite *ClientTestSuite) TestClient_Json_Call() {
    var err error
    ts := TestService{}
    addr := net.JoinHostPort(serverConf.Host, serverConf.Ports[0])

    testReq := TestReq{}
    testRsp := TestRsp{}
    err = suite.client.Call(CodecJson, addr, ts.Service(), "Test", &testReq, &testRsp,
        WithCallRequestTimeout(100e6), WithCallResponseTimeout(100e6))
    suite.Nil(err)
}

func (suite *ClientTestSuite) TestClient_Json_AsyncCall() {
    var err error
    ts := TestService{}
    addr := net.JoinHostPort(serverConf.Host, serverConf.Ports[0])

    testReq := TestReq{}
    testRsp := TestRsp{}
    err = suite.client.AsyncCall(CodecJson, addr,
        ts.Service(), "Add", &testReq, Callback, &testRsp,
        WithCallRequestTimeout(100e6), WithCallResponseTimeout(100e6),
        WithCallMeta("hello", "Service::Add::Json"))
    suite.Nil(err)
}

func TestClientTestSuite(t *testing.T) {
    suite.Run(t, new(ClientTestSuite))
}

func Callback(rsp CallResponse) {
    log.Info("method:%s, cost time span:%s, error:%s, reply:%#v",
        rsp.Opts.Meta["hello"].(string),
        time.Since(rsp.Start),
        jerrors.ErrorStack(rsp.Cause),
        rsp.Reply)
}
