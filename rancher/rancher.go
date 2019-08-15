package rancher

import (
	"strconv"
	"strings"
	"time"

	"github.com/rancher/norman/clientbase"
	client "github.com/rancher/types/client/project/v3"
)

//ServerConfig holds the config for each server the user has setup
type ServerConfig struct {
	AccessKey string `json:"accessKey"`
	SecretKey string `json:"secretKey"`
	TokenKey  string `json:"tokenKey"`
	URL       string `json:"url"`
	Project   string `json:"project"`
	CACerts   string `json:"cacert"`
}

type Rancher struct {
	config ServerConfig
	client *client.Client
}

type KubeConfigApiResponse struct {
	Config string `json:"config"`
}

func New(token string, project string, serverUrl string) *Rancher {
	r := new(Rancher)
	splitToken := strings.Split(token, ":")

	projectUrl := serverUrl + "/v3/project/" + project

	rancherClient, err := client.NewClient(&clientbase.ClientOpts{
		AccessKey: splitToken[0],
		SecretKey: splitToken[1],
		TokenKey:  token,
		URL:       projectUrl,
		CACerts:   "",
	})

	if err != nil {
		panic(err)
	}
	r.client = rancherClient
	return r
}

func (r *Rancher) Redeploy(id string) error {
	workload, err := r.client.Workload.ByID(id)
	if err != nil {
		return err
	}
	newWorkload := &client.Workload{}
	newWorkload.Labels = workload.Labels
	newWorkload.Labels["updated-at"] = strconv.FormatInt(time.Now().Unix(), 10)
	workload, err = r.client.Workload.Update(workload, newWorkload)
	if err != nil {
		return err
	}
	return nil
}

func (r *Rancher) UpdateImage(id string, image string, deployment string) error {
	workload, err := r.client.Workload.ByID(id)
	if err != nil {
		return err
	}
	newWorkload := &client.Workload{}
	workload.Containers[0].Image = image

	newWorkload.Containers = append(newWorkload.Containers, workload.Containers[0])
	newWorkload.Labels = workload.Labels
	newWorkload.Labels["updated-at"] = strconv.FormatInt(time.Now().Unix(), 10)
	workload, err = r.client.Workload.Update(workload, newWorkload)
	if err != nil {
		return err
	}
	return nil
}
