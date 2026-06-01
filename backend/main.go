package main

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Contact struct {
	ID    string `json:"id"`
	Name  string `json:"name"  binding:"required"`
	Email string `json:"email" binding:"required,email"`
	Phone string `json:"phone"`
}

type Store struct {
	mu       sync.RWMutex
	contacts map[string]Contact
}

func NewStore() *Store {
	return &Store{contacts: make(map[string]Contact)}
}

func (s *Store) Create(c Contact) Contact {
	c.ID = uuid.NewString()
	s.mu.Lock()
	s.contacts[c.ID] = c
	s.mu.Unlock()
	return c
}

func (s *Store) List() []Contact {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Contact, 0, len(s.contacts))
	for _, c := range s.contacts {
		out = append(out, c)
	}
	return out
}

func (s *Store) Get(id string) (Contact, bool) {
	s.mu.RLock()
	c, ok := s.contacts[id]
	s.mu.RUnlock()
	return c, ok
}

func (s *Store) Update(id string, c Contact) (Contact, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.contacts[id]; !ok {
		return Contact{}, false
	}
	c.ID = id
	s.contacts[id] = c
	return c, true
}

func (s *Store) Delete(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.contacts[id]; !ok {
		return false
	}
	delete(s.contacts, id)
	return true
}

const (
	contactByIDPath  = "/contacts/:id"
	errNotFound      = "contact not found"
)

func setupRouter(store *Store) *gin.Engine {
	r := gin.New()

	r.GET("/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	r.POST("/contacts", func(c *gin.Context) {
		var input Contact
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, store.Create(input))
	})

	r.GET("/contacts", func(c *gin.Context) {
		c.JSON(http.StatusOK, store.List())
	})

	r.GET(contactByIDPath, func(c *gin.Context) {
		contact, ok := store.Get(c.Param("id"))
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{"error": errNotFound})
			return
		}
		c.JSON(http.StatusOK, contact)
	})

	r.PUT(contactByIDPath, func(c *gin.Context) {
		var input Contact
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		contact, ok := store.Update(c.Param("id"), input)
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{"error": errNotFound})
			return
		}
		c.JSON(http.StatusOK, contact)
	})

	r.DELETE(contactByIDPath, func(c *gin.Context) {
		if !store.Delete(c.Param("id")) {
			c.JSON(http.StatusNotFound, gin.H{"error": errNotFound})
			return
		}
		c.Status(http.StatusNoContent)
	})

	return r
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	store := NewStore()
	r := setupRouter(store)
	r.Run(":8080")
}
