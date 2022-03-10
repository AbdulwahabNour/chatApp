
var app = Vue.createApp({

    data() {
        return {
            ws: null,
            serverUrl: "ws://127.0.0.1:8080/ws",
            messages: [],
            newMessage: ""
        }
    },
    mounted: function () {

        this.connectTowebsocket()
    },
    methods: {
        connectTowebsocket() {

            this.ws = new WebSocket(this.serverUrl);
            this.ws.addEventListener('open', (event) => { this.onWebsocketOpen(event) });
            this.ws.addEventListener('message', (event) => { this.handleNewMessage(event) });
            this.ws.addEventListener('close', event => { alert(close) })

        },
        onWebsocketOpen() {
            console.log("connected to ws");
        },
        handleNewMessage(event) {
            console.log(event.data)
            let data = event.data
            data = data.split(/\r?\n/)
            for (let i = 0; i < data.length; i++) {
                let msg = JSON.parse(data[i])
                this.messages.push(msg)
            }

        },
        sendMessage() {

            if (this.newMessage !== "") {

                this.ws.send(JSON.stringify({ message: this.newMessage }))
                this.newMessage = "";
            }
        }
    }

})
app.mount("#chat")


