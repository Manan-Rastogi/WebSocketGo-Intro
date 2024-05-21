To Run:
```
    -   go run .
    -   Go to Web Browser > Press F12
        -   let socket= new WebSocket("ws://localhost:3100/ws")
        -   socket.onmessage = (event) => {console.log("Received from server - ", event.data)}
        -   socket.send("Hi, How are you?")


```