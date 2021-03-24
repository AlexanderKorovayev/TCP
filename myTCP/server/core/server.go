package core

import (
	"bufio"
	"io"
	"log"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/AlexanderKorovaev/TCP/myTCP/server/broker"
)

//Server объект сервера
type Server struct {
	Port        string
	IdleTimeout time.Duration
	Broker      broker.IBaseAMQP

	listener   net.Listener
	conns      map[*conn]struct{}
	mu         sync.Mutex
	inShutdown bool
}

//ListenAndServe функция прослушивания подключений
func (srv *Server) ListenAndServe(wg *sync.WaitGroup, maxWorkers int) error {
	tasksCh := make(chan *conn)
	// запустим воркеры для ожидания
	// задания будут поступать в очередь, воркеры разбирают задачи из очереди
	for i := 0; i < maxWorkers; i++ {
		go srv.worker(tasksCh, wg)
	}

	port := srv.Port
	listener, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}
	defer listener.Close()
	srv.listener = listener
	for {
		if srv.inShutdown {
			close(tasksCh)
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
		// надо начитать данные и поместить их в очередь
		tasksCh <- connObj
	}
	return nil
}

func (srv *Server) worker(tasksCh <-chan *conn, wg *sync.WaitGroup) {
	//return errors.New("Not implemented handler")
	for {
		conn, ok := <-tasksCh
		if !ok {
			wg.Done()
			return
		}
		for {
			data, err := bufio.NewReader(conn).ReadString('\n')

			if err == io.EOF {
				log.Printf("конец строчки")
				continue
			} else if err != nil {
				log.Printf("ошибка: %v", err)
				break
			}

			if strings.TrimSpace(string(data)) == "STOP" {
				log.Printf("закрыто соединение: %v", conn.RemoteAddr())
				break
			}
			res := strings.ToUpper(string(data))
			log.Printf("-> %v", res)
			conn.Write([]byte(res))
		}
		conn.Close()
		srv.deleteConn(conn)
		srv.Broker.Publish([]byte(conn.RemoteAddr().String()))
	}
}

func (srv *Server) trackConn(c *conn) {
	defer srv.mu.Unlock()
	srv.mu.Lock()
	if srv.conns == nil {
		srv.conns = make(map[*conn]struct{})
	}
	srv.conns[c] = struct{}{}
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
			return
		}
	}
}
