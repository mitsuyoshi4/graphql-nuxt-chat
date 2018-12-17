package graphqlchat

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/go-redis/redis"
	"github.com/segmentio/ksuid"
)

// Resolver ...
type Resolver struct {
	RedisClient     *redis.Client
	MessageChannels map[string]chan Message
	UserChannels    map[string]chan string
	Mutex           sync.Mutex
}

// Mutation ...
func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}

// Query ...
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

// Subscription ...
func (r *Resolver) Subscription() SubscriptionResolver {
	return &subscriptionResolver{r}
}

// New ...
func New(redisURL string) Config {
	client := redis.NewClient(&redis.Options{
		Addr: redisURL,
	})
	pong, err := client.Ping().Result()
	log.Println(pong, err)
	if err != nil {
		panic(err)
	}

	resolver := Resolver{
		RedisClient:     client,
		MessageChannels: map[string]chan Message{},
		UserChannels:    map[string]chan string{},
		Mutex:           sync.Mutex{},
	}
	return Config{
		Resolvers: &resolver,
	}
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) PostMessage(ctx context.Context, user string, text string) (*Message, error) {
	err := r.createUser(user)
	if err != nil {
		return nil, err
	}

	m := Message{
		ID:        ksuid.New().String(),
		CreatedAt: time.Now().UTC(),
		Text:      text,
		User:      user,
	}
	mj, _ := json.Marshal(m)
	if err := r.RedisClient.LPush("messages", mj).Err(); err != nil {
		log.Println(err)
		return nil, err
	}

	r.Mutex.Lock()
	for _, ch := range r.MessageChannels {
		ch <- m
	}
	r.Mutex.Unlock()
	return &m, nil
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Messages(ctx context.Context) ([]Message, error) {
	//cmd := r.RedisClient.LRange("messages", 0, -1)
	cmd := r.RedisClient.LRange("messages", 0, -1)
	if cmd.Err() != nil {
		log.Println(cmd.Err())
		return nil, cmd.Err()
	}
	res, err := cmd.Result()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	messages := []Message{}
	for _, mj := range res {
	//for i := len(res)-1; i > 0; i-- {
		var m Message
	//	err = json.Unmarshal([]byte(res[i]), &m)
		err = json.Unmarshal([]byte(mj), &m)
		messages = append(messages, m)
	}
	return messages, nil
}
func (r *queryResolver) Users(ctx context.Context) ([]string, error) {
	cmd := r.RedisClient.SMembers("users")
	if cmd.Err() != nil {
		log.Println(cmd.Err())
		return nil, cmd.Err()
	}
	res, err := cmd.Result()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return res, nil
}

type subscriptionResolver struct{ *Resolver }

func (r *subscriptionResolver) MessagePosted(ctx context.Context, user string) (<-chan Message, error) {
	err := r.createUser(user)
	if err != nil {
		return nil, err
	}

	messages := make(chan Message, 1)
	r.Mutex.Lock()
	r.MessageChannels[user] = messages
	r.Mutex.Unlock()

	go func() {
		<-ctx.Done()
		r.Mutex.Lock()
		delete(r.MessageChannels, user)
		r.Mutex.Unlock()
	}()

	return messages, nil
}
func (r *subscriptionResolver) UserJoined(ctx context.Context, user string) (<-chan string, error) {
	err := r.createUser(user)
	if err != nil {
		return nil, err
	}

	users := make(chan string, 1)
	r.Mutex.Lock()
	r.UserChannels[user] = users
	r.Mutex.Unlock()

	go func() {
		<-ctx.Done()
		r.Mutex.Lock()
		delete(r.UserChannels, user)
		r.Mutex.Unlock()
	}()

	return users, nil
}

func (r *Resolver) createUser(user string) error {
	if err := r.RedisClient.SAdd("users", user).Err(); err != nil {
		return err
	}
	r.Mutex.Lock()
	for _, ch := range r.UserChannels {
		ch <- user
	}
	r.Mutex.Unlock()
	return nil
}
