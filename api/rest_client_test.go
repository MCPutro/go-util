package api

import (
	"context"
	"fmt"
	"testing"
	"time"
)

type Data struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func TestGet(t *testing.T) {
	restClient := NewRestClient(
		"http://localhost:8888",
		100*time.Second,
	)

	var temp Data

	header := map[string]string{
		"ad": "aaaaaaaaaaaaaaaaa",
	}

	err := restClient.GetRequest(context.Background(), "/coba2", header, &temp)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(temp)
	fmt.Println(temp.Id)
	fmt.Println(temp.Name)
}
