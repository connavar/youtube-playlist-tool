package client

import (
	"fmt"
	"github.com/connavar/youtube-playlist-tool/model"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/youtube/v3"
	"io/ioutil"
	"log"
)

type Driver string

var (
	DriverYoutube Driver = "YouTubeDriver"
	DriverSpotify Driver = "SpotifyDriver"
	DriverApple   Driver = "SpotifyDriver"
)

var Drivers = map[Driver]PlaylistReader{
	DriverYoutube: NewYoutubeClient(),
}

type PlaylistReader interface {
	Playlists() chan model.Playlist
	PlaylistItems(id string) chan model.PlaylistItem
}

type YoutubeClient struct {
	service *youtube.Service
	part    string
}

func NewYoutubeClient() (c *YoutubeClient) {
	fmt.Println("> Building YT Client")
	return &YoutubeClient{
		service: createService(),
		part:    "snippet",
	}
}

func (c *YoutubeClient) Playlists() chan model.Playlist {
	fmt.Println("> Playlists()")
	playlists := make(chan model.Playlist)
	go func() {
		items := c.playlistsList("").Items
		defer close(playlists)
		for _, item := range items {
			playlists <- model.Playlist{
				Title: item.Snippet.Title,
				Id:    item.Id,
			}
		}
	}()
	return playlists
}

func (c *YoutubeClient) PlaylistItems(id string) chan model.PlaylistItem {
	playlistItems := make(chan model.PlaylistItem)
	items := c.playlistsItems(id, 100, "").Items
	go func() {
		defer close(playlistItems)
		for _, item := range items {
			playlistItems <- model.PlaylistItem{
				Title:             item.Snippet.Title,
				VideoChannelOwner: item.Snippet.VideoOwnerChannelTitle,
			}
		}
	}()

	return playlistItems
}

func (c *YoutubeClient) playlistsList(pageToken string) *youtube.PlaylistListResponse {
	call := c.service.Playlists.List([]string{c.part})
	call = call.MaxResults(40)
	call = call.Mine(true)
	if pageToken != "" {
		call = call.PageToken(pageToken)
	}
	response, err := call.Do()
	HandleError(err, "")
	return response
}

func (c *YoutubeClient) playlistsItems(playlistId string, maxResults int64, pageToken string) *youtube.PlaylistItemListResponse {
	//call := service.Playlists.List([]string{part})
	call := c.service.PlaylistItems.List([]string{c.part})

	if playlistId != "" {
		call = call.PlaylistId(playlistId)
	}

	if pageToken != "" {
		call = call.PageToken(pageToken)
	}

	call = call.MaxResults(maxResults)

	response, err := call.Do()
	HandleError(err, "")
	return response
}

func createService() *youtube.Service {
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
