package handlers

import (
	"fmt"
	"net/http"

	"backend-go/internal/events"

	"github.com/gin-gonic/gin"
)

func EventHandler(b *events.Broadcaster) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")

		flusher, ok := c.Writer.(http.Flusher)
		if !ok {
			c.String(http.StatusInternalServerError, "Streaming unsupported")
			return
		}

		// Subscribe client
		ch := b.Subscribe()
		defer b.Unsubscribe(ch)

		// Send messages until client disconnects
		notify := c.Writer.CloseNotify()
		for {
			select {
			case msg := <-ch:
				fmt.Fprintf(c.Writer, "data: %s\n\n", msg)
				flusher.Flush()
			case <-notify:
				return
			}
		}
	}
}
