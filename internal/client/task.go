package client

import (
	"errors"
	"net/http"

	resty "github.com/go-resty/resty/v2"
	"github.com/lindb/common/pkg/encoding"

	"github.com/lindb/lindb/execution/model"
)

type TaskCli interface {
	Submit(task *model.TaskRequest) error
}

type taskCli struct {
	Base
}

func NewTaskCli(endpoint string) TaskCli {
	cli := resty.New()
	cli.SetBaseURL(endpoint)
	return &taskCli{
		Base{
			cli: cli,
		}}
}

func (cli *taskCli) Submit(task *model.TaskRequest) error {
	resp, err := cli.cli.R().
		SetBody(encoding.JSONMarshal(task)).
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		Post("/task")
	if err != nil {
		return err
	}
	if resp.StatusCode() == http.StatusOK {
		return nil
	}
	return errors.New(string(resp.Body()))
}
