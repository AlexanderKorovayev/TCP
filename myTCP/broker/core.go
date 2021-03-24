/*
можуль описывает базовую структуру для работы с брокерами сообщений
*/
package broker

type IBaseAMQP interface {
	Publish(task []byte)
	Consume()
}
