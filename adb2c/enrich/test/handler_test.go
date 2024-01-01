package test

import (
	"bytes"
	json2 "encoding/json"
	api2 "github.com/miruken-go/demo.microservice/adb2c/api"
	"github.com/miruken-go/miruken/handles"
	"github.com/miruken-go/miruken/security/authorizes"
	"io"
	http2 "net/http"
	"net/http/httptest"
	"testing"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/providers/file"
	"github.com/miruken-go/demo.microservice/adb2c/enrich"
	"github.com/miruken-go/miruken/api"
	"github.com/miruken-go/miruken/api/http"
	"github.com/miruken-go/miruken/api/http/httpsrv"
	"github.com/miruken-go/miruken/api/http/httpsrv/auth"
	"github.com/miruken-go/miruken/api/json/stdjson"
	"github.com/miruken-go/miruken/config"
	koanfp "github.com/miruken-go/miruken/config/koanf"
	"github.com/miruken-go/miruken/context"
	"github.com/miruken-go/miruken/security/password"
	"github.com/miruken-go/miruken/setup"
	"github.com/stretchr/testify/suite"
)

type EnrichTestSuite struct {
	suite.Suite
	srv *httptest.Server
}

func (suite *EnrichTestSuite) Setup() *context.Context {
	ctx, _ := setup.New(
		http.Feature(), stdjson.Feature()).
		Specs(&api.GoPolymorphism{}).
		Context()
	return ctx
}

func (suite *EnrichTestSuite) SetupTest() {
	var k = koanf.New(".")
	err := k.Load(file.Provider("./login.json"), json.Parser())
	suite.Nil(err)

	ctx, _ := setup.New(
		httpsrv.Feature(), stdjson.Feature(),
		enrich.Feature(), password.Feature(),
		config.Feature(koanfp.P(k))).
		Specs(&api.GoPolymorphism{}, &EnrichTestSuite{}).
		Handlers(suite).
		Context()

	suite.srv = httptest.NewServer(
		httpsrv.Use(ctx,
			httpsrv.H[*enrich.Handler](),
			auth.WithFlowAlias("login.adb2c").Basic().Required()),
	)
}

func (suite *EnrichTestSuite) TearDownTest() {
	suite.srv.CloseClientConnections()
	suite.srv.Close()
}

func (suite *EnrichTestSuite) TestHandler() {
	suite.Run("Enrich Claims", func() {
		enrichRequest := enrich.Request{
			ObjectId: "123456789",
			Scope:    "domain1/Roles domain1/Groups domain1/Entitlements",
		}

		request, err := json2.Marshal(enrichRequest)
		suite.Nil(err, "marshal request failed")

		reqBytes := bytes.NewReader(request)

		req, err := http2.NewRequest("POST", suite.srv.URL, reqBytes)
		suite.Nil(err, "request create failed")

		req.Header.Set("Content-Type", "application/json")
		req.SetBasicAuth("user", "password")

		resp, err := http2.DefaultClient.Do(req)
		suite.Nil(err, "post failed")
		suite.True(resp.StatusCode >= 200 && resp.StatusCode < 300, "post failed")

		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(resp.Body)

		respBytes, err := io.ReadAll(resp.Body)
		suite.Nil(err, "response read failed")

		var claims map[string]any
		err = json2.Unmarshal(respBytes, &claims)
		suite.Nil(err, "unmarshal response failed")
	})

	suite.Run("Unauthorized", func() {
		request, _ := json2.Marshal(enrich.Request{})
		reqBytes := bytes.NewReader(request)
		req, _ := http2.NewRequest("POST", suite.srv.URL, reqBytes)
		req.Header.Set("Content-Type", "application/json")
		resp, _ := http2.DefaultClient.Do(req)
		suite.Equal(http2.StatusUnauthorized, resp.StatusCode)
	})
}

func (suite *EnrichTestSuite) Get(
	_ *struct {
		handles.It
		authorizes.Required
	  }, get api2.GetSubject,
) api2.Subject {
	return api2.Subject{Id: get.SubjectId}
}

func TestEnrichTestSuite(t *testing.T) {
	suite.Run(t, new(EnrichTestSuite))
}
