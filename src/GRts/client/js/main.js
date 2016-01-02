import GameManager  from "./game"
import NetworkManager  from "./network"
import ko           from "knockout"
import Enum         from "es6-enum"

const STEPS = Enum(
    "CONNECTING",
    "CONNECTED",
    "LOBBY",
    "GAME",
    "DISCONNECTED",
);

class MainVM {
    constructor() {
        this.networkManager = new NetworkManager();
        this.gameManager = null;

        this.currentStep = ko.observable(STEPS.CONNECTING);

        this.isConnecting = ko.computed(() => {
            return (this.currentStep() == STEPS.CONNECTING);
        });

        this.isConnected = ko.computed(() => {
            return (this.currentStep() == STEPS.CONNECTED);
        });

        this.isDisconnected = ko.computed(() => {
            return (this.currentStep() == STEPS.DISCONNECTED);
        });

        this.isInLobby = ko.computed(() => {
            return (this.currentStep() == STEPS.LOBBY);
        });

        this.isInGame = ko.computed(() => {
            return (this.currentStep() == STEPS.GAME);
        });

        this.currentPlayerId = ko.observable();
        this.currentPlayerLogin = ko.observable();

        this.loginEditMode = ko.observable(false);

        this.lobbyName = ko.observable();
        this.lobbyPlayers = ko.observableArray([]);
    }

    editLogin() {
        this.loginEditMode(true);
    }

    changeLogin() {
        this.loginEditMode(false);
        this.networkManager.send("CHANGE_LOGIN", { login: this.currentPlayerLogin() });
        // TODO Send login changed
    }

    resetLobby() {
        this.loginEditMode = ko.observable(false);
        this.lobbyName = ko.observable();
        this.lobbyPlayers = ko.observableArray();
    }

    start() {
        this.networkManager.registerOnWSOpenHandler(this.onWSOpen.bind(this));
        this.networkManager.registerOnWSCloseHandler(this.onWSClose.bind(this));
        this.networkManager.registerOnMessageHandler(this.onWSMessage.bind(this));

        this.networkManager.connect();
    }

    onWSOpen() {
        console.log("WS Opened :)");
        this.currentStep(STEPS.CONNECTED);
    }

    onWSClose() {
        console.log("WS Closed :(");
        this.currentStep(STEPS.DISCONNECTED);
    }

    onWSMessage(message) {
        console.log("WS Message", message);

        if (message.type == "CONNECTED") {
            this.currentPlayerId(message.data.playerId);
            this.currentPlayerLogin(message.data.playerLogin);
        }
        else if (message.type == "JOINED_GAME") {
            this.lobbyName(message.data.gameName);
            message.data.players.forEach((player) => {
                this.lobbyPlayers.push({
                    id: player.id,
                    login: ko.observable(player.login)
                })
            });
            this.currentStep(STEPS.LOBBY);
        }
        else if (message.type == "PLAYER_JOINED_GAME") {
            this.lobbyPlayers.push({
                id: message.data.player.id,
                login: ko.observable(message.data.player.login)
            });
        }
        else if (message.type == "PLAYER_CHANGED_LOGIN") {
            ko.utils.arrayForEach(this.lobbyPlayers(), (player) => {
                if (player.id == message.data.player.id)
                    player.login(message.data.player.login);
            });
            if (message.data.player.id == this.currentPlayerId())
                this.currentPlayerLogin(message.data.player.login);
        }
    }
}

let vm = new MainVM;
ko.applyBindings(vm);
vm.start()
