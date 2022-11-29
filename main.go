package main

import (
	"fmt"
	"github.com/connavar/youtube-playlist-tool/controller"
)

func main() {

	youtubePlaylists := controller.GetAllYoutubePlaylists()
	for _, youtubePlaylist := range youtubePlaylists {
		if controller.IsYoutubeMusic(&youtubePlaylist) {
			fmt.Println(youtubePlaylist.Title)
			for _, item := range youtubePlaylist.PlaylistItems {
				fmt.Println(" - ", item.Title, " | ", item.VideoChannelOwner)
			}
		}
	}

}
