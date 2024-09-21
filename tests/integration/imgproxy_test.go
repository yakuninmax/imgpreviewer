package integration_test

import (
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"
)

const proxyAddr = "http://localhost:8888"

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
	response, err := p.client.Get(proxyAddr + "/crop/100/200/nginx/fakeimage.jpg")
	p.Require().NoError(err)
	p.Require().Equal(500, response.StatusCode)

	bodyString, err := getResponseBodyString(*response)
	p.Require().NoError(err)
	p.Require().Equal("404 Not Found\n", bodyString)
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
