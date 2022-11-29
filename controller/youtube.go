package controller

import (
	"fmt"
	"github.com/connavar/youtube-playlist-tool/client"
	"github.com/connavar/youtube-playlist-tool/model"
	"strings"
	"sync"
)

func GetAllYoutubePlaylists() []model.Playlist {

	fmt.Println("Starting...")
	c := client.Drivers[client.DriverYoutube]

	youtubePlaylists := make([]model.Playlist, 0)

	tp := NewThreadPool(4, func(playlist model.Playlist) {
		for playListItem := range c.PlaylistItems(playlist.Id) {
			playlist.AddItem(playListItem)
		}
		youtubePlaylists = append(youtubePlaylists, playlist)
	})

	for playlist := range c.Playlists() {
		// Adds playlist to channel to process playlist items
		tp.in <- playlist
	}

	tp.Close()

	return youtubePlaylists
}

func IsYoutubeMusic(playlist *model.Playlist) bool {
	for _, item := range playlist.PlaylistItems {
		if strings.HasSuffix(strings.TrimSpace(item.VideoChannelOwner), "Topic") {
			return true
		}
	}
	return false
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
