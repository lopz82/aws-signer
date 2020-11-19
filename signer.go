package signature

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	v4 "github.com/aws/aws-sdk-go/aws/signer/v4"
	"io/ioutil"
	"net/http"
	"text/template"
	"time"
)

// Config the plugin configuration.
type Config struct {
	Service string `json:"service"`
	Region  string `json:"region"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		Service: "",
		Region:  "",
	}
}

// AWSSigner a Traefik plugin.
type AWSSigner struct {
	next     http.Handler
	service  string
	region   string
	name     string
	template *template.Template
}

// New created a new AWSSigner plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if config.Service == "" {
		return nil, fmt.Errorf("service cannot be empty")
	}
	if config.Region == "" {
		return nil, fmt.Errorf("region cannot be empty")
	}
	if err := IsExistingServiceInRegion(config.Region, config.Service); err != nil {
		return nil, fmt.Errorf(err.Error())
	}

	return &AWSSigner{
		service:  config.Service,
		region:   config.Region,
		next:     next,
		name:     name,
		template: template.New("aws-signer").Delims("[[", "]]"),
	}, nil
}

func (a *AWSSigner) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	creds := credentials.NewEnvCredentials()
	err := SignRequest(req, a, creds)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	a.next.ServeHTTP(rw, req)
}

func SignRequest(req *http.Request, a *AWSSigner, creds *credentials.Credentials) error {
	signer := v4.NewSigner(creds)
	body, err := ReadRequestsBody(req)
	if err != nil {
		return fmt.Errorf("error reading request's body: %s", err)
	}

	_, err = signer.Sign(req, bytes.NewReader(body), a.service, a.region, time.Now())
	if err != nil {
		return fmt.Errorf("error signing the request: %s", err)
	}
	return nil
}

func ReadRequestsBody(req *http.Request) ([]byte, error) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func IsExistingServiceInRegion(r, s string) error {
	regions, exists := endpoints.RegionsForService(endpoints.DefaultPartitions(), endpoints.AwsPartitionID, s)
	if !exists {
		return fmt.Errorf("service %s does not exist", s)
	}
	for _, region := range regions {
		if region.ID() == r {
			return nil
		}
	}
	return fmt.Errorf("region %s does not exist", r)
}
