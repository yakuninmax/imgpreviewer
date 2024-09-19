package cache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestQueue(t *testing.T) {
	testFiles := []file{
		{
			url:  "../../examples/_gopher_original_1024x504.jpg",
			size: 64212,
			name: "_gopher_original_1024x504",
		},
		{
			url:  "../../examples/gopher_50x50.jpg",
			size: 1956,
			name: "gopher_50x50",
		},
		{
			url:  "../../examples/gopher_200x700.jpg",
			size: 30146,
			name: "gopher_200x700",
		},
		{
			url:  "../../examples/gopher_256x126.jpg",
			size: 10121,
			name: "gopher_256x126",
		},
		{
			url:  "../../examples/gopher_333x666.jpg",
			size: 41562,
			name: "gopher_333x666",
		},
		{
			url:  "../../examples/gopher_500x500.jpg",
			size: 47656,
			name: "gopher_500x500",
		},
		{
			url:  "../../examples/gopher_1024x252.jpg",
			size: 41771,
			name: "gopher_1024x252",
		},
		{
			url:  "../../examples/gopher_2000x1000.jpg",
			size: 226943,
			name: "gopher_2000x1000",
		},
	}

	t.Run("empty queue", func(t *testing.T) {
		q := newQueue()

		require.Nil(t, q.getFront())
		require.Nil(t, q.getBack())
	})

	t.Run("push item to empty queue", func(t *testing.T) {
		q := newQueue()

		q.pushFront(testFiles[0])
		require.Equal(t, testFiles[0].size, q.size)
		require.Equal(t, testFiles[0], q.getFront().file)
		require.Equal(t, testFiles[0], q.getBack().file)
	})

	t.Run("push multiple items", func(t *testing.T) {
		q := newQueue()

		for _, file := range testFiles {
			q.pushFront(file)
			require.Equal(t, file, q.getFront().file)
		}

		require.Equal(t, int64(464367), q.size)
		require.Equal(t, testFiles[7], q.getFront().file)
		require.Equal(t, testFiles[0], q.getBack().file)
	})

	t.Run("move items to front", func(t *testing.T) {
		q := newQueue()

		for _, file := range testFiles {
			q.pushFront(file)
		}

		for i := 0; i < len(testFiles)-1; i++ {
			q.moveToFront(q.getBack())

			require.Equal(t, int64(464367), q.size)
			require.Equal(t, testFiles[i], q.getFront().file)
			require.Equal(t, testFiles[i+1], q.getBack().file)
		}
	})

	t.Run("remove items", func(t *testing.T) {
		q := newQueue()

		for _, file := range testFiles {
			q.pushFront(file)
		}

		q.remove(q.getFront())

		require.Equal(t, 464367-testFiles[7].size, q.size)
		require.Equal(t, testFiles[6], q.getFront().file)
		require.Equal(t, testFiles[0], q.getBack().file)

		q.remove(q.getBack())

		require.Equal(t, 464367-testFiles[7].size-testFiles[0].size, q.size)
		require.Equal(t, testFiles[6], q.getFront().file)
		require.Equal(t, testFiles[1], q.getBack().file)
	})
}
