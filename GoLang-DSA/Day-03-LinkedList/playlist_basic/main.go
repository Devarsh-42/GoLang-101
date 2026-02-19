package main

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)


//  I/O Helpers -> bcoz fmt.scan can't read song names with spaces


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
		reader.ReadString('\n') // flush rest of line
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
		reader.ReadString('\n') // flush rest of line
		if err == nil {
			return val
		}
		fmt.Println(" Invalid input. Please enter a valid number.")
	}
}

type Song struct {
	SongName     string
	SingerName   string
	SongDuration float32 
	Next         *Song
	Prev         *Song
}

func NewSong(songName, singerName string, songDuration float32) (*Song, error) {
	songName = strings.TrimSpace(songName)
	singerName = strings.TrimSpace(singerName)
	if songName == "" || singerName == "" || songDuration <= 0 {
		return nil, errors.New("name/singer must be non-empty and duration must be positive")
	}
	return &Song{
		SongName:     songName,
		SingerName:   singerName,
		SongDuration: songDuration,
	}, nil
}

type Playlist struct {
	head *Song
	tail *Song
	Size int
}



func fmtDuration(d float32) string {
	t := int(d)
	return fmt.Sprintf("%d:%02d", t/60, t%60)
}

func (pl *Playlist) Display() {
	if pl.head == nil {
		fmt.Println("  (playlist is empty)")
		return
	}
	fmt.Printf("  %-4s  %-30s %-24s %s\n", "#", "Song", "Singer", "Duration")
	fmt.Println("\n--------------------\n")
	cur := pl.head
	for i := 1; cur != nil; i++ {
		fmt.Printf("  %-4d  %-30s %-24s %s\n",
			i, cur.SongName, cur.SingerName, fmtDuration(cur.SongDuration))
		cur = cur.Next
	}
	fmt.Printf("\n  Total: %d song(s)\n", pl.Size)
}


func (pl *Playlist) AddSong() { // O(1) tail insertion
	fmt.Println("\n Add New Song")
	fmt.Println("\n--------------------\n")
	songName := readLine("    Song Name    : ")
	singerName := readLine("    Singer Name  : ")
	dur := readFloat("    Duration (s) : ")

	newSong, err := NewSong(songName, singerName, dur)
	if err != nil {
		fmt.Println(" Error:", err)
		return
	}

	if pl.head == nil {
		pl.head = newSong
		pl.tail = newSong
	} else {
		newSong.Prev = pl.tail
		pl.tail.Next = newSong
		pl.tail = newSong
	}
	pl.Size++
	fmt.Printf(" \"%s\" by %s added to playlist.\n", newSong.SongName, newSong.SingerName)
}


func (pl *Playlist) DeleteSong() {  // O(n) search, O(1) unlink
	if pl.head == nil {
		fmt.Println("  Playlist is empty.")
		return
	}
	fmt.Println("\n  Delete Song")
	fmt.Println("\n--------------------\n")
	query := readLine("    Enter Song Name to delete: ")

	cur := pl.head
	for cur != nil {
		if strings.EqualFold(cur.SongName, query) {
			// Unlink node
			if cur.Prev != nil {
				cur.Prev.Next = cur.Next
			} else {
				pl.head = cur.Next // deleted head
			}
			if cur.Next != nil {
				cur.Next.Prev = cur.Prev
			} else {
				pl.tail = cur.Prev // deleted tail
			}
			deleted := cur.SongName
			cur.Next = nil
			cur.Prev = nil
			pl.Size--
			fmt.Printf("  \"%s\" deleted from playlist.\n", deleted)
			return
		}
		cur = cur.Next
	}
	fmt.Printf("  Song \"%s\" not found.\n", query)
}


func (pl *Playlist) SearchSong() { // case-insensitive substring match on name or singer
	if pl.head == nil {
		fmt.Println("  Playlist is empty.")
		return
	}
	fmt.Println("\n Search Song")
	fmt.Println("\n--------------------\n")
	query := strings.ToLower(readLine("    Search by Song Name or Singer: "))

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
		fmt.Printf("  No songs found matching \"%s\".\n", query)
		return
	}

	fmt.Printf("\n  Found %d result(s):\n\n", len(results))
	fmt.Printf("  %-4s  %-30s %-24s %s\n", "#", "Song", "Singer", "Duration")
	fmt.Println("\n--------------------\n")
	for i, s := range results {
		fmt.Printf("  %-4d  %-30s %-24s %s\n",
			i+1, s.SongName, s.SingerName, fmtDuration(s.SongDuration))
	}
}


func (pl *Playlist) ShufflePlaylist() { // Fisher-Yates on a slice, then rebuild list -> O(n) time, O(n) space
	if pl.Size < 2 {
		fmt.Println("  Need at least 2 songs to shuffle.")
		return
	}

	// Flatten to slice
	songs := make([]*Song, 0, pl.Size)
	cur := pl.head
	for cur != nil {
		songs = append(songs, cur)
		cur = cur.Next
	}

	// Fisher-Yates shuffle (cryptographically independent, seeded per call)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	rng.Shuffle(len(songs), func(i, j int) {
		songs[i], songs[j] = songs[j], songs[i]
	})

	// Rebuild doubly linked list from shuffled slice
	for i, s := range songs { // O(n) relinking after O(n) shuffle
		if i == 0 {
			s.Prev = nil
		} else {
			s.Prev = songs[i-1]
			songs[i-1].Next = s
		}
		s.Next = nil
	}
	pl.head = songs[0]
	pl.tail = songs[len(songs)-1]

	fmt.Println("  Playlist shuffled! New order:\n")
	pl.Display()
}

func main() {
	playlist := &Playlist{}

	for {
		fmt.Println()
		fmt.Println("\n--------------------\n")
		fmt.Println(" Music Playlist Manager")
		fmt.Println("\n--------------------\n")
		fmt.Println("  1. Add Song")
		fmt.Println("  2. Delete Song")
		fmt.Println("  3. Search Song")
		fmt.Println("  4. Shuffle Playlist")
		fmt.Println("  5. Display Playlist")
		fmt.Println("  0. Exit")
		fmt.Println("\n--------------------\n")
		choice := readInt("  Choose an option: ")
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
			fmt.Println("  Current Playlist")
			playlist.Display()
		case 0:
			fmt.Println("  Goodbye! ")
			return
		default:
			fmt.Println("  Invalid option. Please try again.")
		}
	}
}
