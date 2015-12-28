import GameManager  from "./game"
import NetworkManager  from "./network"
import ko           from "knockout"
import Enum         from "es6-enum"

const STEPS = Enum(
    "CONNECT",
    "LOBBY",
    "GAME",
);

class MainVM {
    constructor() {
        this.networkManager = new NetworkManager();
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
        this.networkManager.registerOnWSOpenHandler(this.onWSOpen.bind(this));
        this.networkManager.registerOnWSCloseHandler(this.onWSClose.bind(this));
        this.networkManager.registerOnMessageHandler(this.onWSMessage.bind(this));

        this.networkManager.connect();
    }

    onWSOpen() {
        console.log("WS Opened :)");
        this.currentStep(STEPS.LOBBY);
    }

    onWSClose() {
        console.log("WS Closed :(");

    }

    onWSMessage(message) {
        console.log("WS Message", message);
    }
}

let vm = new MainVM;
ko.applyBindings(vm);
vm.start()
