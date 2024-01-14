package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const SIZE = 5

type Node struct {
	Val   string
	Left  *Node
	Right *Node
}

type Queue struct {
	Head   *Node
	Tail   *Node
	length int
}

type Hash map[string]*Node

type Cache struct {
	Queue Queue
	Hash  Hash
}

func NewCache() Cache {
	return Cache{Queue: NewQueue(), Hash: Hash{}}
}

func NewQueue() Queue {
	head := &Node{}
	tail := &Node{}

	head.Right = tail
	tail.Left = head

	return Queue{Head: head, Tail: tail, length: 0}
}

func (c *Cache) Check(word string) {
	node := &Node{}

	if val, ok := c.Hash[word]; ok {
		node = c.Remove(val)
	} else {
		node = &Node{Val: word}
	}
	c.Add(node)
	c.Hash[word] = node
}

func (c *Cache) Remove(node *Node) *Node {
	left := node.Left
	right := node.Right

	right.Left = left
	left.Right = right

	c.Queue.length--

	delete(c.Hash, node.Val)

	return node
}

func (c *Cache) Add(node *Node) {
	tmp := c.Queue.Head.Right
	c.Queue.Head.Right = node
	node.Left = c.Queue.Head
	node.Right = tmp
	tmp.Left = node

	c.Queue.length++

	if c.Queue.length > SIZE {
		c.Remove(c.Queue.Tail.Left)
	}
}

func (c *Cache) Display() {
	node := c.Queue.Head.Right
	fmt.Printf("%d - [", c.Queue.length)
	for i := 0; i < c.Queue.length; i++ {
		if node != c.Queue.Head && node != c.Queue.Tail {
			fmt.Printf("%s ", node.Val)
		}
		node = node.Right
	}

	fmt.Printf("]\n")
}

func handleCommand(cache *Cache, command string) {
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return
	}

	switch strings.ToLower(parts[0]) {
	case "add":
		if len(parts) > 1 {
			cache.Check(parts[1])
		} else {
			fmt.Println("Usage: add <word>")
		}
	case "remove":
		if len(parts) > 1 {
			if node, ok := cache.Hash[parts[1]]; ok {
				cache.Remove(node)
			} else {
				fmt.Printf("Word '%s' not found in the cache.\n", parts[1])
			}
		} else {
			fmt.Println("Usage: remove <word>")
		}
	case "display":
		cache.Display()
	default:
		fmt.Printf("Unknown command: %s\n", parts[0])
	}
}

func main() {
	cache := NewCache()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print(">> ")
		if !scanner.Scan() {
			break
		}

		command := scanner.Text()
		if strings.ToLower(command) == "exit" {
			break
		}
		handleCommand(&cache, command)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}
