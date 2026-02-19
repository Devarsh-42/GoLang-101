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


type Node struct {
	SongID       string // unique identifier  e.g. "song-000042"
	SongName     string
	SingerName   string
	SongDuration float32 // seconds
	Next         *Node
	Prev         *Node
}

//  OptimizedPlaylist  ──  DLL + two HashMaps
//
//   idMap   : songID  (string)  → *Node   →  O(1) lookup / delete by ID
//   nameMap : lower(songName)   → *Node   →  O(1) lookup / delete by name
//
//  DLL maintains insertion order and enables O(1) prev/next navigation.
//  HashMaps eliminate the O(n) traversal that a plain DLL requires for
//  lookup and deletion by key

type OptimizedPlaylist struct {
	head    *Node
	tail    *Node
	idMap   map[string]*Node // songID → *Node
	nameMap map[string]*Node // lower(songName) → *Node
	Size    int
	nextID  int
}

func NewOptimizedPlaylist() *OptimizedPlaylist {
	return &OptimizedPlaylist{
		idMap:   make(map[string]*Node),
		nameMap: make(map[string]*Node),
	}
}

func (pl *OptimizedPlaylist) genID() string {
	pl.nextID++
	return fmt.Sprintf("song-%06d", pl.nextID)
}

// Display 

func printHeader() {
	fmt.Printf("  %-12s  %-28s %-22s %s\n", "ID", "Song", "Singer", "Duration")
	divider()
}

func printNodeRow(n *Node) {
	fmt.Printf("  %-12s  %-28s %-22s %s\n",
		n.SongID, n.SongName, n.SingerName, fmtDuration(n.SongDuration))
}

func (pl *OptimizedPlaylist) Display() {
	if pl.head == nil {
		fmt.Println("  (playlist is empty)")
		return
	}
	printHeader()
	cur := pl.head
	for cur != nil {
		printNodeRow(cur)
		cur = cur.Next
	}
	fmt.Printf("\n  Total: %d song(s)\n", pl.Size)
}

//  AddSong  ──  O(1)  tail-insert into DLL  +  two O(1) map inserts
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

	node := &Node{
		SongID:       pl.genID(),
		SongName:     songName,
		SingerName:   singerName,
		SongDuration: dur,
	}
	pl.tailInsert(node)
	fmt.Printf("  [%s] \"%s\" by %s added.\n", node.SongID, node.SongName, node.SingerName)
}

// tailInsert links node at the end of the DLL and registers it in both maps.
func (pl *OptimizedPlaylist) tailInsert(node *Node) {
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

	var node *Node
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

// unlinkNode removes a node from the DLL and cleans both maps.  O(1).
func (pl *OptimizedPlaylist) unlinkNode(node *Node) {
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
		printNodeRow(node)

	case 2:
		id := readLine("  Song ID: ")
		node := pl.idMap[id]
		if node == nil {
			fmt.Printf("  ID \"%s\" not found.\n", id)
			return
		}
		printHeader()
		printNodeRow(node)

	case 3:
		query := strings.ToLower(readLine("  Substring (name or singer): "))
		var results []*Node
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
			printNodeRow(n)
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
	nodes := make([]*Node, 0, pl.Size)
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


func main() {
	playlist := NewOptimizedPlaylist()

	for {
		fmt.Println()
		divider()
		fmt.Println("  Optimized Playlist Manager  [DLL + HashMap]")
		divider()
		fmt.Println("  1. Add Song")
		fmt.Println("  2. Delete Song")
		fmt.Println("  3. Search Song")
		fmt.Println("  4. Shuffle Playlist")
		fmt.Println("  5. Display Playlist")
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
		case 0:
			fmt.Println("  Goodbye! ")
			return
		default:
			fmt.Println("  Invalid option.")
		}
	}
}
