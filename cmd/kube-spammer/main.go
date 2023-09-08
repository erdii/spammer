package main

import (
	"context"
	"fmt"
	"math/rand"
	"runtime"
	"sync"

	"github.com/google/uuid"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sruntime "k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var scheme = k8sruntime.NewScheme()

var (
	letters          = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	longStringLength = 1024 // 1Kib
	longString       string
)

func init() {
	if err := corev1.AddToScheme(scheme); err != nil {
		panic(err)
	}

	r := make([]rune, longStringLength)
	letterLen := len(letters)
	for i := range r {
		r[i] = letters[rand.Intn(letterLen)]
	}
	longString = string(r)
}

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	restConfig := ctrl.GetConfigOrDie()
	c, err := client.New(restConfig, client.Options{})
	if err != nil {
		return err
	}

	ctx := context.Background()

	var wg sync.WaitGroup

	for i := 0; i < runtime.NumCPU(); i++ {
		fmt.Println("spawning spammer goroutine")
		wg.Add(1)
		go func() {
			if err := createForever(ctx, c); err != nil {
				panic(err)
			}
			wg.Done()
		}()
	}

	wg.Wait()
	return nil
}

func createForever(ctx context.Context, c client.Client) error {
	for {
		if err := create(ctx, c); err != nil {
			return fmt.Errorf("Creating event: %w", err)
		}
	}
}

func create(ctx context.Context, c client.Client) error {
	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-%s", "test", uuid.New().String()),
			Namespace: "default",
		},
		Data: map[string]string{
			"foo": longString,
		},
	}
	if err := c.Create(ctx, cm); err != nil {
		return err
	}
	return nil
}
