
var app = Vue.createApp({

    data() {
        return {
            ws: null,
            serverUrl: "ws://127.0.0.1:8080/ws",
            roomInput: null,
            rooms: [],
            user: { name: "" },
            messages: [],
            newMessage: "",
            users: [],
        }
    },
    mounted: function () {
    },
    methods: {
        connect() {
            this.connectTowebsocket()
        },
        connectTowebsocket() {

            this.ws = new WebSocket(this.serverUrl + "?name=" + this.user.name);
            this.ws.addEventListener('open', (event) => { this.onWebsocketOpen(event) });
            this.ws.addEventListener('message', (event) => { this.handleNewMessage(event) });
            this.ws.addEventListener('close', (event) => { this.onWebsocketClose(event) });



        },

        onWebsocketOpen() {
            console.log("connected to ws");
        },
        onWebsocketClose() {
            console.log("connection close")
        },
        handleNewMessage(event) {

            let data = event.data;
            data = data.split(/\r?\n/);
            console.log(data)
            for (let i = 0; i < data.length; i++) {
                let msg = JSON.parse(data[i]);
                console.log(msg)

                switch (msg.action) {
                    case "send-message":
                        this.handleChatMessage(msg);
                        break;

                    case "user-join":
                        this.handleUserJoined(msg);
                        break;
                    case "user-left":
                        this.handleUserLeft(msg);
                        break;
                    case "join-room":
                        this.handleJoinRoom(msg);
                        break;

                    default:
                        break;
                }


            }

        },
        handleAllusers(msg) {

        },
        handleChatMessage(msg) {
            const room = this.findRoombyID(msg.target.id);

            if (typeof (room) !== "undefined") {
                room.message.push(msg)
            }
        },
        sendMessage(room) {

            if (room.newMessage !== "") {
                this.ws.send(JSON.stringify({
                    action: 'send-message',
                    message: room.newMessage,
                    target: {
                        id: room.id,
                        name: room.name,
                    }
                }));
                room.newMessage = "";

            }

        },
        handleJoinRoom(msg) {
            room = msg.target;
            room.name = room.private ? msg.sender.name : room.name;
            room.message = [];
            if (this.findRoombyID(room.id) === undefined) {
                this.rooms.push(room);
            }
            this.handleChatMessage(msg);
        },
        findroom(roomName) {
            for (let i = 0; i < this.rooms.length; i++) {
                if (this.rooms[i].name == roomName) {
                    return this.rooms[i]
                }
            }
        },
        findRoombyID(roomID) {
            for (let i = 0; i < this.rooms.length; i++) {
                if (this.rooms[i].id == roomID) {
                    return this.rooms[i]
                }
            }

        },
        joinRoom() {
            this.ws.send(JSON.stringify({
                action: 'join-room',
                message: this.roomInput
            }));
            this.messages = [];
            this.roomInput = ""
        },
        leaveRoom(room) {
            this.ws.send(JSON.stringify({
                action: "leave-room",
                message: room.name
            }));

            for (let i = 0; i < this.rooms.length; i++) {
                if (this.rooms[i].name === room.name) {
                    this.rooms.splice(i, 1);
                    break;
                }
            }
        },
        handleUserJoined(msg) {
            for (let i = 0; i < msg.Sender.length; i++) {
                this.users.push(msg.Sender[i])
            }

        },
        handleUserLeft(msg) {
            for (let i = 0; i < this.users.length; i++) {
                if (this.users[i].id == msg.Sender[0].id) {
                    this.users.splice(i, 1);
                }
            }
        },
        joinPrivateRoom(user) {
            this.ws.send(JSON.stringify({
                action: "join-room-private",
                message: user.id
            }))
        }


    }

})
app.mount("#chat")


