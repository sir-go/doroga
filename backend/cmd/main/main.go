package main

func main() {
	if err := DBC.Connect(); err != nil {
		LOG.Panic(err)
	}
	defer DBC.Disconnect()
	SRV.Run()
}
