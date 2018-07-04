package main

import (
	"bytes"
	"context"
	"io"
	"os"

	fdk "github.com/fnproject/fdk-go"
	"github.com/go-redis/redis"
)

var redisHost string
var redisPort string

func main() {

	redisHost = os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "localhost"
	}

	redisPort = os.Getenv("REDIS_PORT")
	if redisPort == "" {
		redisPort = "6379"
	}
	fdk.Handle(fdk.HandlerFunc(myHandler))

}

func myHandler(ctx context.Context, in io.Reader, out io.Writer) {

	opts := redis.Options{Addr: redisHost + ":" + redisPort}
	client := redis.NewClient(&opts)
	_, conErr := client.Ping().Result()

	if conErr != nil {
		out.Write([]byte("connection error - " + conErr.Error()))
		return
	}

	defer client.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(in)
	key := buf.String()

	val, err := client.Get(key).Result()

	if err != nil {
		out.Write([]byte("GET error - " + err.Error()))
		return
	}

	out.Write([]byte("value for key " + key + " is " + val))

}
