//go:build integration
// +build integration

package integration_test

import (
	"context"
	"encoding/base64"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"
)

const proxy = "http://localhost:8888/"

type ProxySuite struct {
	suite.Suite
	ctx    context.Context
	client *http.Client
}

func (p *ProxySuite) SetupSuite() {
	p.ctx = context.Background()
	p.client = &http.Client{}
}

func (p *ProxySuite) TestImageNotFound() {
	resp, err := p.client.Get(proxy + "fill/100/200/nginx/fakeimage.jpg")
	p.Require().NoError(err)
	p.Require().Equal(502, resp.StatusCode)

	body, err := getResponseBodyString(*resp)
	p.Require().NoError(err)
	p.Require().Equal("remote server return: 404 Not Found\n", body)

	resp.Body.Close()
}

func (p *ProxySuite) TestServerConnectionError() {
	resp, err := p.client.Get(proxy + "fill/100/200/fake.serv/fakeimage.jpg")
	p.Require().NoError(err)
	p.Require().Equal(502, resp.StatusCode)

	bodyString, err := getResponseBodyString(*resp)
	p.Require().NoError(err)
	p.Require().Contains(bodyString, "no such host")

	resp.Body.Close()
}

func (p *ProxySuite) TestRemoteServerError() {
	resp, err := p.client.Get(proxy + "fill/100/200/nginx/error")
	p.Require().NoError(err)
	p.Require().Equal(502, resp.StatusCode)

	body, err := getResponseBodyString(*resp)
	p.Require().NoError(err)
	p.Require().Equal(body, "remote server return: 503 Service Temporarily Unavailable\n")

	resp.Body.Close()
}

func (p *ProxySuite) TestInvalidFileType() {
	resp, err := p.client.Get(proxy + "fill/100/200/nginx/text.file")
	p.Require().NoError(err)
	p.Require().Equal(502, resp.StatusCode)

	body, err := getResponseBodyString(*resp)
	p.Require().NoError(err)
	p.Require().Equal(body, "invalid file type\n")

	resp.Body.Close()
}

func (p *ProxySuite) TestInvalidImageSize() {
	resp, err := p.client.Get(proxy + "fill/3000/5000/nginx/_gopher_original_1024x504.jpg")
	p.Require().NoError(err)
	p.Require().Equal(502, resp.StatusCode)

	body, err := getResponseBodyString(*resp)
	p.Require().NoError(err)
	p.Require().Contains(body, "target size is larger than original")

	resp.Body.Close()
}

func (p *ProxySuite) TestCustomHeaderPass() {
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, proxy+"fill/100/200/nginx/protected/_gopher_original_1024x504.jpg", nil)
	p.Require().NoError(err)

	req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("testuser:testpass")))

	resp, err := p.client.Do(req)
	p.Require().NoError(err)
	p.Require().Equal(200, resp.StatusCode)
	p.Require().NotNil(resp.Body)

	resp.Body.Close()
}

func (p *ProxySuite) TestSuccessfullProcessing() {
	resp, err := p.client.Get(proxy + "fill/100/200/nginx/_gopher_original_1024x504.jpg")
	p.Require().NoError(err)
	p.Require().Equal(200, resp.StatusCode)
	p.Require().NoError(err)
	p.Require().NotNil(resp.Body)

	resp.Body.Close()
}

func getResponseBodyString(resp http.Response) (string, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func TestProxySuite(t *testing.T) {
	suite.Run(t, new(ProxySuite))
}
