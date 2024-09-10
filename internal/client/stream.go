package client

import (
	"errors"
	"fmt"
	"net/http"

	resty "github.com/go-resty/resty/v2"
	"github.com/lindb/common/pkg/encoding"

	"github.com/lindb/lindb/sql/execution/model"
)

type StreamCli interface {
	SendTaskResultSet(result *model.TaskResultSet) error
}

type streamCli struct {
	Base
}

func NewStreamCli(endpoint string) StreamCli {
	cli := resty.New()
	cli.SetBaseURL(endpoint)
	return &streamCli{
		Base{
			cli: cli,
		},
	}
}

func (cli *streamCli) SendTaskResultSet(result *model.TaskResultSet) error {
	fmt.Printf("task result=%v,%v\n", result, string(encoding.JSONMarshal(result)))
	resp, err := cli.cli.R().
		SetBody(encoding.JSONMarshal(result)).
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		Put("/task/result")
	if err != nil {
		return err
	}
	if resp.StatusCode() == http.StatusOK {
		return nil
	}
	return errors.New(string(resp.Body()))
}
