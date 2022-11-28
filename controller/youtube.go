package controller

import (
	"fmt"
	"github.com/connavar/youtube-playlist-tool/model"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/youtube/v3"
	"io/ioutil"
	"log"
)

var (
	service *youtube.Service
	part    string
)

func init() {
	service = CreateService()
	part = "snippet"
}

func GetAll() {

	for _, playlist := range Playlists() {

		// Print the playlist ID and title for the playlist resource.
		fmt.Println(playlist.Title)

		for _, playListItem := range PlaylistItems(&playlist) {
			fmt.Println(" - ", playListItem.Title, playListItem.VideoChannelOwner)
		}
	}

}

func Playlists() map[string]model.Playlist {
	playlists := make(map[string]model.Playlist)

	items := playlistsList(service, part, "", "", 40, true, "", "", "").Items

	for _, item := range items {

		playlists[item.Id] = model.Playlist{
			Title: item.Snippet.Title,
			Id:    item.Id,
		}

	}

	return playlists
}

func PlaylistItems(playlist *model.Playlist) []model.PlaylistItem {
	playlistItems := make([]model.PlaylistItem, 0)

	items := playlistsItems(service, part, playlist.Id, 100, "", "").Items

	for _, item := range items {
		playlistItems = append(playlistItems, model.PlaylistItem{
			Title:             item.Snippet.Title,
			VideoChannelOwner: item.Snippet.VideoOwnerChannelTitle,
		})
	}

	return playlistItems
}

func playlistsList(service *youtube.Service, part string, channelId string, hl string, maxResults int64, mine bool, onBehalfOfContentOwner string, pageToken string, playlistId string) *youtube.PlaylistListResponse {
	call := service.Playlists.List([]string{part})
	if channelId != "" {
		call = call.ChannelId(channelId)
	}
	if hl != "" {
		call = call.Hl(hl)
	}
	call = call.MaxResults(maxResults)
	if mine != false {
		call = call.Mine(true)
	}
	if onBehalfOfContentOwner != "" {
		call = call.OnBehalfOfContentOwner(onBehalfOfContentOwner)
	}
	if pageToken != "" {
		call = call.PageToken(pageToken)
	}
	if playlistId != "" {
		call = call.Id(playlistId)
	}
	response, err := call.Do()
	HandleError(err, "")
	return response
}

func playlistsItems(service *youtube.Service, part string, playlistId string, maxResults int64, onBehalfOfContentOwner string, pageToken string) *youtube.PlaylistItemListResponse {
	//call := service.Playlists.List([]string{part})
	call := service.PlaylistItems.List([]string{part})

	if playlistId != "" {
		call = call.PlaylistId(playlistId)
	}

	if pageToken != "" {
		call = call.PageToken(pageToken)
	}

	call = call.MaxResults(maxResults)

	if onBehalfOfContentOwner != "" {
		call = call.OnBehalfOfContentOwner(onBehalfOfContentOwner)
	}

	response, err := call.Do()
	HandleError(err, "")
	return response
}

func CreateService() *youtube.Service {
	ctx := context.Background()

	b, err := ioutil.ReadFile("client_secret.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved credentials
	// at ~/.credentials/youtube-go-quickstart.json
	config, err := google.ConfigFromJSON(b, youtube.YoutubeReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(ctx, config)
	service, err := youtube.New(client)

	HandleError(err, "Error creating YouTube client")

	return service
}

func HandleError(err error, message string) {
	if message == "" {
		message = "Error making API call"
	}
	if err != nil {
		log.Fatalf(message+": %v", err.Error())
	}
}
