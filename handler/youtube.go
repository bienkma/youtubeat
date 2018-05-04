package handler

import (
"flag"
"log"

"google.golang.org/api/youtube/v3"
)

var (
	maxResults = flag.Int64("max-results", 25, "Max YouTube results")
)

// SearchKeyWord function extend from youtube api v3.
func SearchKeyWord(service *youtube.Service, query string) map[string]string {
	flag.Parse()
	// Make the API call to YouTube.
	call := service.Search.List("id,snippet").
		Q(query).
		MaxResults(*maxResults)
	response, err := call.Do()
	log.Printf("%v", err)

	// Group video, channel, and playlist results in separate lists.
	videos := make(map[string]string)
	channels := make(map[string]string)
	playlists := make(map[string]string)

	// Iterate through each item and add it to the correct list.
	for _, item := range response.Items {
		switch item.Id.Kind {
		case "youtube#video":
			videos[item.Id.VideoId] = item.Snippet.Title
		case "youtube#channel":
			channels[item.Id.ChannelId] = item.Snippet.Title
		case "youtube#playlist":
			playlists[item.Id.PlaylistId] = item.Snippet.Title
		}
	}
	// Return only videos results in list
	return videos
}

