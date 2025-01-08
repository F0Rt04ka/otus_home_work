package hw04lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("simple purge logic", func(t *testing.T) {
		c := NewCache(3)
		c.Set("1", 10)
		c.Set("2", 20)
		c.Set("3", 30)

		c.Set("4", 40) // [40 30 20]

		result := false
		_, result = c.Get("1")
		require.False(t, result)
		_, result = c.Get("2")
		require.True(t, result)
		_, result = c.Get("3")
		require.True(t, result)
		_, result = c.Get("4")
		require.True(t, result)
	})

	t.Run("purge logic with get and set", func(t *testing.T) {
		c := NewCache(3)
		c.Set("1", 10)
		c.Set("2", 20)
		c.Set("3", 30)

		c.Set("1", 11) // [11 30 20]
		c.Set("4", 40) // [40 11 30]

		result := false
		_, result = c.Get("3")
		require.True(t, result)
		_, result = c.Get("1")
		require.True(t, result)
		_, result = c.Get("4")
		require.True(t, result)
		_, result = c.Get("2")
		require.False(t, result)

		c.Get("3")     // [30 40 11]
		c.Set("5", 50) // [50 30 40]

		_, result = c.Get("4")
		require.True(t, result)
		_, result = c.Get("3")
		require.True(t, result)
		_, result = c.Get("5")
		require.True(t, result)
		_, result = c.Get("1")
		require.False(t, result)
	})
}

func TestCacheMultithreading(t *testing.T) {

	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}
