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

const proxyAddr = "http://localhost:8888/"

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
	response, err := p.client.Get(proxyAddr + "fill/100/200/nginx/fakeimage.jpg")
	p.Require().NoError(err)
	p.Require().Equal(502, response.StatusCode)

	bodyString, err := getResponseBodyString(*response)
	p.Require().NoError(err)
	p.Require().Equal("remote server return: 404 Not Found\n", bodyString)

	response.Body.Close()
}

func (p *ProxySuite) TestServerConnectionError() {
	response, err := p.client.Get(proxyAddr + "fill/100/200/fake.serv/fakeimage.jpg")
	p.Require().NoError(err)
	p.Require().Equal(502, response.StatusCode)

	bodyString, err := getResponseBodyString(*response)
	p.Require().NoError(err)
	p.Require().Contains(bodyString, "no such host")

	response.Body.Close()
}

func (p *ProxySuite) TestRemoteServerError() {
	response, err := p.client.Get(proxyAddr + "fill/100/200/nginx/error")
	p.Require().NoError(err)
	p.Require().Equal(502, response.StatusCode)

	bodyString, err := getResponseBodyString(*response)
	p.Require().NoError(err)
	p.Require().Equal(bodyString, "remote server return: 503 Service Temporarily Unavailable\n")

	response.Body.Close()
}

func (p *ProxySuite) TestInvalidFileType() {
	response, err := p.client.Get(proxyAddr + "fill/100/200/nginx/text.file")
	p.Require().NoError(err)
	p.Require().Equal(502, response.StatusCode)

	bodyString, err := getResponseBodyString(*response)
	p.Require().NoError(err)
	p.Require().Equal(bodyString, "invalid file type\n")

	response.Body.Close()
}

func (p *ProxySuite) TestInvalidImageSize() {
	response, err := p.client.Get(proxyAddr + "fill/3000/5000/nginx/_gopher_original_1024x504.jpg")
	p.Require().NoError(err)
	p.Require().Equal(502, response.StatusCode)

	bodyString, err := getResponseBodyString(*response)
	p.Require().NoError(err)
	p.Require().Contains(bodyString, "target size is larger than original")

	response.Body.Close()
}

func (p *ProxySuite) TestCustomHeaderPass() {
	ctx := context.Background()
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, proxyAddr+"fill/100/200/nginx/protected/_gopher_original_1024x504.jpg", nil)
	p.Require().NoError(err)

	request.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("testuser:testpass")))

	response, err := p.client.Do(request)
	p.Require().NoError(err)
	p.Require().Equal(200, response.StatusCode)
	p.Require().NotNil(response.Body)

	response.Body.Close()
}

func (p *ProxySuite) TestSuccessfullProcessing() {
	response, err := p.client.Get(proxyAddr + "fill/100/200/nginx/_gopher_original_1024x504.jpg")
	p.Require().NoError(err)
	p.Require().Equal(200, response.StatusCode)
	p.Require().NoError(err)
	p.Require().NotNil(response.Body)

	response.Body.Close()
}

func getResponseBodyString(response http.Response) (string, error) {
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return string(bodyBytes), nil
}

func TestProxySuite(t *testing.T) {
	suite.Run(t, new(ProxySuite))
}
