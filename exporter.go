package bogo

// PingExporter that needs to be implemented to be a PingMessage exporter
type PingExproter interface {
	Initialize(in chan PingMessage, wait chan int)
}
