package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

//  I/O Helpers

var reader = bufio.NewReader(os.Stdin)

func readLine(prompt string) string {
	fmt.Print(prompt)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

func readFloat(prompt string) float32 {
	for {
		fmt.Print(prompt)
		var val float32
		_, err := fmt.Fscan(reader, &val)
		reader.ReadString('\n')
		if err == nil && val > 0 {
			return val
		}
		fmt.Println(" Invalid input. Please enter a positive number.")
	}
}

func readInt(prompt string) int {
	for {
		fmt.Print(prompt)
		var val int
		_, err := fmt.Fscan(reader, &val)
		reader.ReadString('\n')
		if err == nil {
			return val
		}
		fmt.Println(" Invalid input. Please enter a valid number.")
	}
}

func divider() { fmt.Println(strings.Repeat("─", 74)) }

func fmtDuration(d float32) string {
	t := int(d)
	return fmt.Sprintf("%d:%02d", t/60, t%60)
}

type Song struct {
	SongID       string // unique identifier  e.g. "song-000042"
	SongName     string
	SingerName   string
	SongDuration float32 // seconds
	PlayCount    int     // number of times played
	Next         *Song
	Prev         *Song
}

//  OptimizedPlaylist  ──  DLL + two HashMaps + LFU Cache
//
//   idMap   : songID  (string)  → *Song   →  O(1) lookup / delete by ID
//   nameMap : lower(songName)   → *Song   →  O(1) lookup / delete by name
//   cache   : songID            → *Song   →  hot songs (PlayCount > 5)
//
//  DLL maintains insertion order and enables O(1) prev/next navigation.
//  HashMaps eliminate the O(n) traversal that a plain DLL requires for
//  lookup and deletion by key

type OptimizedPlaylist struct {
	head         *Song
	tail         *Song
	idMap        map[string]*Song // songID → *Song
	nameMap      map[string]*Song // lower(songName) → *Song
	cache        map[string]*Song // hot songs (PlayCount > 5)
	cacheMaxSize int
	Size         int
	nextID       int
}

func NewOptimizedPlaylist() *OptimizedPlaylist {
	return &OptimizedPlaylist{
		idMap:        make(map[string]*Song),
		nameMap:      make(map[string]*Song),
		cache:        make(map[string]*Song),
		cacheMaxSize: 10,
	}
}

func (pl *OptimizedPlaylist) genID() string {
	pl.nextID++
	return fmt.Sprintf("song-%06d", pl.nextID)
}

// Display

func printHeader() {
	fmt.Printf("  %-12s  %-28s %-22s %-10s %s\n", "ID", "Song", "Singer", "Plays", "Duration")
	divider()
}

func printSongRow(s *Song) {
	fmt.Printf("  %-12s  %-28s %-22s %-10d %s\n",
		s.SongID, s.SongName, s.SingerName, s.PlayCount, fmtDuration(s.SongDuration))
}

func (pl *OptimizedPlaylist) Display() {
	if pl.head == nil {
		fmt.Println("  (playlist is empty)")
		return
	}
	printHeader()
	cur := pl.head
	for cur != nil {
		printSongRow(cur)
		cur = cur.Next
	}
	fmt.Printf("\n  Total: %d song(s)\n", pl.Size)
}

// AddSong  ──  O(1)  tail-insert into DLL  +  two O(1) map inserts
func (pl *OptimizedPlaylist) AddSong() {
	fmt.Println("\n Add New Song")
	divider()
	songName := readLine("    Song Name    : ")
	singerName := readLine("    Singer Name  : ")
	dur := readFloat("    Duration (s) : ")

	songName = strings.TrimSpace(songName)
	singerName = strings.TrimSpace(singerName)
	if songName == "" || singerName == "" || dur <= 0 {
		fmt.Println("  Invalid song details.")
		return
	}
	// Prevent duplicate name (case-insensitive)
	if _, dup := pl.nameMap[strings.ToLower(songName)]; dup {
		fmt.Printf(" \"%s\" already exists in the playlist.\n", songName)
		return
	}

	node := &Song{
		SongID:       pl.genID(),
		SongName:     songName,
		SingerName:   singerName,
		SongDuration: dur,
	}
	pl.tailInsert(node)
	fmt.Printf("  [%s] \"%s\" by %s added.\n", node.SongID, node.SongName, node.SingerName)
}

// tailInsert links song at the end of the DLL and registers it in both maps.
func (pl *OptimizedPlaylist) tailInsert(node *Song) {
	if pl.head == nil {
		pl.head = node
		pl.tail = node
	} else {
		node.Prev = pl.tail
		pl.tail.Next = node
		pl.tail = node
	}
	pl.idMap[node.SongID] = node
	pl.nameMap[strings.ToLower(node.SongName)] = node
	pl.Size++
}

//  DeleteSong  ──  O(1) map lookup  →  O(1) DLL unlink  →  O(1) map delete
//  (vs. O(n) linear scan in a plain DLL)

func (pl *OptimizedPlaylist) DeleteSong() {
	if pl.head == nil {
		fmt.Println("  Playlist is empty.")
		return
	}
	fmt.Println("\n  Delete Song")
	divider()
	fmt.Println("  Delete by:  1) Song Name   2) Song ID")
	choice := readInt("  Choice: ")

	var node *Song
	switch choice {
	case 1:
		name := readLine("  Song Name: ")
		node = pl.nameMap[strings.ToLower(name)]
	case 2:
		id := readLine("  Song ID: ")
		node = pl.idMap[id]
	default:
		fmt.Println("  Invalid choice.")
		return
	}

	if node == nil {
		fmt.Println(" Song not found.")
		return
	}
	deleted := node.SongName
	id := node.SongID
	pl.unlinkNode(node)
	fmt.Printf(" [%s] \"%s\" deleted.\n", id, deleted)
}

// unlinkNode removes a song from the DLL, both maps, and cache.  O(1).
func (pl *OptimizedPlaylist) unlinkNode(node *Song) {
	if node.Prev != nil {
		node.Prev.Next = node.Next
	} else {
		pl.head = node.Next
	}
	if node.Next != nil {
		node.Next.Prev = node.Prev
	} else {
		pl.tail = node.Prev
	}
	delete(pl.idMap, node.SongID)
	delete(pl.nameMap, strings.ToLower(node.SongName))
	delete(pl.cache, node.SongID)
	node.Next = nil
	node.Prev = nil
	pl.Size--
}

//  SearchSong
//    • Exact by Song Name : O(1)  ← nameMap lookup
//    • Exact by Song ID   : O(1)  ← idMap lookup
//    • Substring          : O(n)  ← unavoidable without an inverted index

func (pl *OptimizedPlaylist) SearchSong() {
	if pl.head == nil {
		fmt.Println(" Playlist is empty.")
		return
	}
	fmt.Println("\n  🔍  Search Song")
	divider()
	fmt.Println("  1) Exact Name  O(1)   2) Song ID  O(1)   3) Substring  O(n)")
	choice := readInt("  Choice: ")

	switch choice {
	case 1:
		name := readLine("  Song Name: ")
		node := pl.nameMap[strings.ToLower(name)]
		if node == nil {
			fmt.Printf(" \"%s\" not found.\n", name)
			return
		}
		printHeader()
		printSongRow(node)

	case 2:
		id := readLine("  Song ID: ")
		node := pl.idMap[id]
		if node == nil {
			fmt.Printf("  ID \"%s\" not found.\n", id)
			return
		}
		printHeader()
		printSongRow(node)

	case 3:
		query := strings.ToLower(readLine("  Substring (name or singer): "))
		var results []*Song
		cur := pl.head
		for cur != nil {
			if strings.Contains(strings.ToLower(cur.SongName), query) ||
				strings.Contains(strings.ToLower(cur.SingerName), query) {
				results = append(results, cur)
			}
			cur = cur.Next
		}
		if len(results) == 0 {
			fmt.Printf(" No results for \"%s\".\n", query)
			return
		}
		fmt.Printf("\n  %d result(s):\n\n", len(results))
		printHeader()
		for _, n := range results {
			printSongRow(n)
		}

	default:
		fmt.Println(" Invalid choice.")
	}
}

//  ShufflePlaylist  ──  Fisher-Yates  O(n)
//  HashMaps stay 100% valid: only DLL prev/next links are rewired;
//  the *Node pointers stored in the maps never change.

func (pl *OptimizedPlaylist) ShufflePlaylist() {
	if pl.Size < 2 {
		fmt.Println("  Need at least 2 songs to shuffle.")
		return
	}
	nodes := make([]*Song, 0, pl.Size)
	cur := pl.head
	for cur != nil {
		nodes = append(nodes, cur)
		cur = cur.Next
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	rng.Shuffle(len(nodes), func(i, j int) { nodes[i], nodes[j] = nodes[j], nodes[i] })

	// Rewire DLL links only; idMap/nameMap remain untouched
	for i, n := range nodes {
		if i == 0 {
			n.Prev = nil
		} else {
			n.Prev = nodes[i-1]
			nodes[i-1].Next = n
		}
		n.Next = nil
	}
	pl.head = nodes[0]
	pl.tail = nodes[len(nodes)-1]

	fmt.Println("  Playlist shuffled! New order:\n")
	pl.Display()
}

// PlaySong  ──  O(1) lookup via maps, increments PlayCount, caches if PlayCount > 5
func (pl *OptimizedPlaylist) PlaySong() {
	if pl.head == nil {
		fmt.Println("  Playlist is empty.")
		return
	}
	fmt.Println("\n  ▶  Play Song")
	divider()
	fmt.Println("  Find by:  1) Song Name   2) Song ID")
	choice := readInt("  Choice: ")

	var song *Song
	switch choice {
	case 1:
		name := readLine("  Song Name: ")
		song = pl.nameMap[strings.ToLower(name)]
	case 2:
		id := readLine("  Song ID: ")
		song = pl.idMap[id]
	default:
		fmt.Println("  Invalid choice.")
		return
	}

	if song == nil {
		fmt.Println("  Song not found.")
		return
	}

	song.PlayCount++
	fmt.Printf("\n  ▶  Now playing: \"%s\" by %s  [%s]\n",
		song.SongName, song.SingerName, fmtDuration(song.SongDuration))
	fmt.Printf("     Play count: %d\n", song.PlayCount)

	if song.PlayCount > 5 {
		pl.addToCache(song)
	}
}

// addToCache adds a song to the hot cache.
// If the cache is full, the song with the lowest PlayCount is evicted (LFU).  O(cache size).
func (pl *OptimizedPlaylist) addToCache(song *Song) {
	if _, exists := pl.cache[song.SongID]; exists {
		return // already cached
	}
	if len(pl.cache) >= pl.cacheMaxSize {
		// Evict the least-frequently-played cached song
		var evictID string
		minPlays := int(^uint(0) >> 1) // MaxInt
		for id, s := range pl.cache {
			if s.PlayCount < minPlays {
				minPlays = s.PlayCount
				evictID = id
			}
		}
		delete(pl.cache, evictID)
	}
	pl.cache[song.SongID] = song
	fmt.Printf("  \"%s\" added to hot cache (played %d times).\n", song.SongName, song.PlayCount)
}

// ShowMostPlayed displays all songs in the hot cache, sorted by PlayCount descending.
func (pl *OptimizedPlaylist) ShowMostPlayed() {
	if len(pl.cache) == 0 {
		fmt.Println("  No songs in the hot cache yet. Play a song more than 5 times to cache it.")
		return
	}
	songs := make([]*Song, 0, len(pl.cache))
	for _, s := range pl.cache {
		songs = append(songs, s)
	}
	// Sort descending by PlayCount (cache is small, O(k²) is fine)
	for i := 0; i < len(songs)-1; i++ {
		for j := 0; j < len(songs)-1-i; j++ {
			if songs[j].PlayCount < songs[j+1].PlayCount {
				songs[j], songs[j+1] = songs[j+1], songs[j]
			}
		}
	}
	fmt.Printf("\n Hot Cache  (%d / %d slots used)\n\n", len(pl.cache), pl.cacheMaxSize)
	printHeader()
	for _, s := range songs {
		printSongRow(s)
	}
}

func main() {
	playlist := NewOptimizedPlaylist()

	for {
		fmt.Println()
		divider()
		fmt.Println("  Optimized Playlist Manager  [DLL + HashMap + LFU Cache]")
		divider()
		fmt.Println("  1. Add Song")
		fmt.Println("  2. Delete Song")
		fmt.Println("  3. Search Song")
		fmt.Println("  4. Shuffle Playlist")
		fmt.Println("  5. Display Playlist")
		fmt.Println("  6. Play Song")
		fmt.Println("  7. Most Played (Hot Cache)")
		fmt.Println("  0. Exit")
		divider()

		choice := readInt("  Choose: ")
		fmt.Println()

		switch choice {
		case 1:
			playlist.AddSong()
		case 2:
			playlist.DeleteSong()
		case 3:
			playlist.SearchSong()
		case 4:
			playlist.ShufflePlaylist()
		case 5:
			fmt.Println("  Current Playlist\n")
			playlist.Display()
		case 6:
			playlist.PlaySong()
		case 7:
			playlist.ShowMostPlayed()
		case 0:
			fmt.Println("  Goodbye! ")
			return
		default:
			fmt.Println("  Invalid option.")
		}
	}
}
