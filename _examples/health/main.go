package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/CodyGuo/go-pkg/health"
)

type Mysql struct {
}

func (m *Mysql) Init() {
	health.Register("mysql", m)
}

func (m *Mysql) Ping(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return errors.New("mysql connect timeout")
}

type Redis struct {
}

func (r *Redis) Init() {
	health.Register("redis", r)
}

func (r *Redis) Ping(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return errors.New("redis connect tiemout")
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	mysql := Mysql{}
	mysql.Init()

	redis := Redis{}
	redis.Init()

	ping := health.Ping(ctx)
	data, _ := json.Marshal(ping)
	fmt.Printf("%s\n", data)
	// Output:
	// {"status":"down","details":[{"name":"mysql","status":"down","error":"mysql connect timeout"},{"name":"redis","status":"down","error":"redis connect tiemout"}]}
}
