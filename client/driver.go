package client

type Driver string

var (
	DriverYoutube Driver = "YouTubeDriver"
	DriverSpotify Driver = "SpotifyDriver"
	DriverApple   Driver = "SpotifyDriver"
)

var Drivers = map[Driver]PlaylistReader{
	DriverYoutube: NewYoutubeClient(),
	DriverSpotify: NewSpotifyClient(),
}
