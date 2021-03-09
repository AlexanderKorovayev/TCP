package shouter

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

type Server struct {
	Port         string
	IdleTimeout  time.Duration
	MaxReadBytes int64

	listener   net.Listener
	conns      map[*conn]struct{}
	mu         sync.Mutex
	inShutdown bool
}

func (srv *Server) ListenAndServe() error {
	port := srv.Port
	if port == "" {
		port = ":2000"
	}
	log.Printf("starting server on %v\n", port)
	listener, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}
	defer listener.Close()
	srv.listener = listener
	for {
		// should be guarded by mu
		if srv.inShutdown {
			break
		}
		newConn, err := listener.Accept()
		if err != nil {
			log.Printf("error accepting connection %v", err)
			continue
		}
		log.Printf("accepted connection from %v", newConn.RemoteAddr())
		conn1 := &conn{
			Conn:        newConn,
			IdleTimeout: srv.IdleTimeout,
		}
		srv.trackConn(conn1)
		conn1.SetDeadline(time.Now().Add(conn1.IdleTimeout))
		go srv.handle(conn1)
	}
	return nil
}

func (srv *Server) trackConn(c *conn) {
	defer srv.mu.Unlock()
	srv.mu.Lock()
	if srv.conns == nil {
		srv.conns = make(map[*conn]struct{})
	}
	srv.conns[c] = struct{}{}
}

func (srv *Server) handle(conn *conn) error {
	defer func() {
		log.Printf("closing connection from %v", conn.RemoteAddr())
		conn.Close()
		srv.deleteConn(conn)
	}()

	for {
		data, err := bufio.NewReader(conn).ReadString('\n')

		if err == io.EOF {
			fmt.Println("--end-of-file--")
			return err
		} else if err != nil {
			fmt.Println("Oops! Some error occured!", err)
			return err
		}

		if strings.TrimSpace(string(data)) == "STOP" {
			fmt.Println("Exiting TCP server!")
			return nil
		}
		res := strings.ToUpper(string(data))
		fmt.Print("-> ", res)
		conn.Write([]byte(res))
	}
}

func (srv *Server) deleteConn(conn *conn) {
	defer srv.mu.Unlock()
	srv.mu.Lock()
	delete(srv.conns, conn)
}

func (srv *Server) Shutdown() {
	// should be guarded by mu
	srv.inShutdown = true
	log.Println("shutting down...")
	srv.listener.Close()
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			log.Printf("waiting on %v connections", len(srv.conns))
		}
		if len(srv.conns) == 0 {
			return
		}
	}
}
