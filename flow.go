package main

type Flow interface {
	FlowOut(size int) []chan string
	FlowIn(in chan string)
}
