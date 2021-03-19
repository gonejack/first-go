package main

func main() {
	var c chan int = make(chan int)

	close(c)
	close(c)
}
