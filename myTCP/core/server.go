package core

import (
	"bufio"
	"io"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

//Server объект сервера
type Server struct {
	Port         string
	IdleTimeout  time.Duration
	MaxReadBytes int64

	listener   net.Listener
	conns      map[*conn]struct{}
	mu         sync.Mutex
	inShutdown bool
}

//ListenAndServe функция прослушивания подключений
func (srv *Server) ListenAndServe() error {
	port := srv.Port
	listener, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}
	defer listener.Close()
	srv.listener = listener
	for {
		if srv.inShutdown {
			break
		}
		newConn, err := listener.Accept()
		if err != nil {
			log.Printf("ошибка получения соединения: %v", err)
			continue
		}
		connObj := &conn{
			Conn:        newConn,
			IdleTimeout: srv.IdleTimeout,
		}
		srv.trackConn(connObj)
		connObj.SetDeadline(time.Now().Add(connObj.IdleTimeout))
		go srv.handle(connObj)
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
		log.Printf("закрыто соединение: %v", conn.RemoteAddr())
		conn.Close()
		srv.deleteConn(conn)
	}()
	//return errors.New("Not implemented handler")
	for {
		data, err := bufio.NewReader(conn).ReadString('\n')

		if err == io.EOF {
			return err
		} else if err != nil {
			log.Printf("ошибка: %v", err)
			return err
		}

		if strings.TrimSpace(string(data)) == "STOP" {
			return nil
		}
		res := strings.ToUpper(string(data))
		log.Printf("-> %v", res)
		conn.Write([]byte(res))
	}
}

func (srv *Server) deleteConn(conn *conn) {
	defer srv.mu.Unlock()
	srv.mu.Lock()
	delete(srv.conns, conn)
}

//Shutdown функция для корректного завершения всех обработчиков
func (srv *Server) Shutdown(wg *sync.WaitGroup) {
	srv.mu.Lock()
	srv.inShutdown = true
	srv.mu.Unlock()
	log.Println("shutting down...")
	if srv.listener != nil {
		srv.listener.Close()
	}
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			log.Printf("waiting on %v connections", len(srv.conns))
		}
		if len(srv.conns) == 0 {
			wg.Done()
			return
		}
	}
}
