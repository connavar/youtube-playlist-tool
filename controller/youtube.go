package controller

import (
	"fmt"
	"github.com/connavar/youtube-playlist-tool/client"
	"github.com/connavar/youtube-playlist-tool/model"
	"sync"
)

func GetAll() {

	fmt.Println("Starting...")
	c := client.Drivers[client.DriverYoutube]

	tp := NewThreadPool(4, func(playlist model.Playlist) {
		for playListItem := range c.PlaylistItems(playlist.Id) {
			fmt.Println(" - ", playListItem.Title, playListItem.VideoChannelOwner)
		}
	})

	for playlist := range c.Playlists() {
		// Print the playlist ID and title for the playlist resource.
		fmt.Println(playlist.Title)
		tp.in <- playlist
	}
	tp.Close()
}

type ThreadPool struct {
	in chan model.Playlist
	wg sync.WaitGroup
}

func (tp *ThreadPool) Close() {
	close(tp.in)
	tp.wg.Wait()

}

func NewThreadPool(threads int, f func(playlist model.Playlist)) *ThreadPool {
	result := &ThreadPool{
		in: make(chan model.Playlist, 32),
	}
	for i := 0; i < threads; i++ {
		result.wg.Add(1)
		go func() {
			for next := range result.in {
				f(next)
			}
			result.wg.Done()
		}()
	}
	return result
}
