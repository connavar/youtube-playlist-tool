package model

type Playlist struct {
	Title         string
	Id            string
	PlaylistItems []PlaylistItem
}

type PlaylistItem struct {
	Title             string
	VideoChannelOwner string
}

func (playlist *Playlist) AddItem(item PlaylistItem) {
	playlist.PlaylistItems = append(playlist.PlaylistItems, item)
}
