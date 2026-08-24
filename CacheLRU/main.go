package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type Node struct {
	Key    string
	Value  string
	Next   *Node
	Prev   *Node
	Expiry time.Time
}

type LRUCache struct {
	Capacity     int
	ItemsHashMap map[string]*Node
	mutex        sync.Mutex //it guards the map
	Head         *Node
	Tail         *Node
}

func newNode(value string, key string, ttl time.Duration) *Node {

	return &Node{
		Key:    key,
		Value:  value,
		Next:   nil,
		Prev:   nil,
		Expiry: time.Now().Add(ttl),
	}
}

func NewLRUCache(capacity int) *LRUCache {

	//the sentinel
	head := &Node{}
	tail := &Node{}
	head.Next = tail
	tail.Prev = head

	return &LRUCache{Capacity: capacity,
		ItemsHashMap: make(map[string]*Node),
		mutex:        sync.Mutex{},
		Head:         head,
		Tail:         tail}
}

func (c *LRUCache) UnlinkNode(node *Node) {
	node.Prev.Next = node.Next
	node.Next.Prev = node.Prev
}

func (c *LRUCache) PutNodeAtFront(node *Node) {
	node.Next = c.Head.Next
	node.Prev = c.Head
	c.Head.Next.Prev = node
	c.Head.Next = node

}

func (c *LRUCache) evictLastNode() {

	lastNode := c.Tail.Prev
	lastNode.Prev.Next = c.Tail
	c.Tail.Prev = lastNode.Prev
	delete(c.ItemsHashMap, lastNode.Key)
}

/*Set(key, value string, ttl time.Duration) —
if the key already exists, update it and move to front;
if not, create a new node, add it at the front,
and if the map's now over Capacity, evict whoever's at the tail*/

func (c *LRUCache) isOverCapacity() bool {

	return len(c.ItemsHashMap) > c.Capacity
}

func (c *LRUCache) Set(key, value string, ttl time.Duration) {

	c.mutex.Lock()
	defer c.mutex.Unlock()
	node, ok := c.ItemsHashMap[key] // we check if key already exists
	if ok {                         // if they exist we update it and put it in front
		node.Value = value
		node.Expiry = time.Now().Add(ttl)
		c.UnlinkNode(node)
		c.PutNodeAtFront(node)

	} else { //if not we make a new node, add it to the front.

		auxNode := newNode(value, key, ttl)
		c.PutNodeAtFront(auxNode)
		c.ItemsHashMap[key] = auxNode
	}

	if c.isOverCapacity() { // here the c is over capacity.
		c.evictLastNode()
	}
}

func (c *LRUCache) Get(key string) (string, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	node, ok := c.ItemsHashMap[key]

	if !ok {
		return "", false
	}

	if time.Now().After(node.Expiry) {
		c.UnlinkNode(node)
		delete(c.ItemsHashMap, key)
		return "", false
	}

	c.UnlinkNode(node)
	c.PutNodeAtFront(node)

	return node.Value, true

}

func (c *LRUCache) GetAll() []DashboardItem {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	items := make([]DashboardItem, 0)
	now := time.Now()

	aux := c.Head.Next
	for aux != c.Tail {
		next := aux.Next

		if now.Before(aux.Expiry) {
			items = append(items, DashboardItem{
				Key:                 aux.Key,
				Value:               aux.Value,
				TTLRemainingSeconds: int(aux.Expiry.Sub(now).Seconds()),
			})
		} else {
			c.UnlinkNode(aux)
			delete(c.ItemsHashMap, aux.Key)
		}

		aux = next
	}
	return items
}

// --- JSON STRUCTS ---

type SetRequest struct {
	Key        string `json:"key"`
	Value      string `json:"value"`
	TTLSeconds int    `json:"ttl_seconds"`
}

type GetResponse struct {
	Value string `json:"value"`
	Found bool   `json:"found"`
}

type DashboardItem struct {
	Key                 string `json:"key"`
	Value               string `json:"value"`
	TTLRemainingSeconds int    `json:"ttl_remaining_seconds"`
}

// --- HTTP SERVER ---

type CacheServer struct {
	cache *LRUCache
}

func (s *CacheServer) handleSet(w http.ResponseWriter, r *http.Request) {
	var req SetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	s.cache.Set(req.Key, req.Value, time.Duration(req.TTLSeconds)*time.Second)
	w.WriteHeader(http.StatusOK)
}

func (s *CacheServer) handleGet(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	val, found := s.cache.Get(key)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(GetResponse{Value: val, Found: found})
}

func (s *CacheServer) handleAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s.cache.GetAll())
}

func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}


		next(w, r)
	}
}

func main() {
	server := &CacheServer{cache: NewLRUCache(5)} // Capacity of 5 for easy testing

	http.HandleFunc("/set", enableCORS(server.handleSet))
	http.HandleFunc("/get", enableCORS(server.handleGet))
	http.HandleFunc("/all", enableCORS(server.handleAll))
	/*http.HandleFunc("/reset", enableCORS(server.handleReset)) maybe if in the future we add Reset*/

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server running at port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
