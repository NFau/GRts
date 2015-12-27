import GameManager  from "./game"
import ko           from "knockout"
import Enum         from "es6-enum"
import pb           from "protobufjs"

const STEPS = Enum(
    "CONNECT",
    "LOBBY",
    "GAME",
);

class MainVM {
    constructor() {
        this.ws = null;
        this.gameManager = null;

        this.currentStep = ko.observable(STEPS.CONNECT);

        this.isConnecting = ko.computed(() => {
            return (this.currentStep() == STEPS.CONNECT);
        });

        this.isInLobby = ko.computed(() => {
            return (this.currentStep() == STEPS.LOBBY);
        });

        this.isInGame = ko.computed(() => {
            return (this.currentStep() == STEPS.GAME);
        });
    }

    start() {
//        this.builder = pb.loadProtoFile("../protocol/lobby.proto")
//        console.log(this.builder);
//        this.grtsproto = this.builder.build('grtsproto')
        this.ws = new WebSocket("ws://127.0.0.1:8080/socket");
        this.ws.onopen = (event) => {
            console.log("Connected");
            this.ws.send("Hello boy");
            setTimeout(() => {
                // FIRST INSTRUCTION!!
                this.ws.send("QUIT GAME");
            }, 4000)
        };

        this.ws.onmessage = (event) => {
            console.log("Received message", event.data);

            if (this.isInGame()) {
                // TODO: Send info to GM
            } else if (this.isInLobby()) {
                // TODO on start game:
                // this.gm = new GameManager;
                // this.gm.start()
            } else if (this.isConnecting()){
                console.log(event)
//                this.currentStep(STEPS.LOBBY);
            }
        };

        this.ws.onclose = (event) => {
            console.log("Connection closed", event.data);
            this.currentStep(STEPS.CONNECT)
            // TODO, resend game id on reconnect to be able to resume the game
            this.ws.connect()
            this.gm = null;
        }
    }
}

let vm = new MainVM;
ko.applyBindings(vm);
vm.start()
