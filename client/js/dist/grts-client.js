(function e(t,n,r){function s(o,u){if(!n[o]){if(!t[o]){var a=typeof require=="function"&&require;if(!u&&a)return a(o,!0);if(i)return i(o,!0);var f=new Error("Cannot find module '"+o+"'");throw f.code="MODULE_NOT_FOUND",f}var l=n[o]={exports:{}};t[o][0].call(l.exports,function(e){var n=t[o][1][e];return s(n?n:e)},l,l.exports,e,t,n,r)}return n[o].exports}var i=typeof require=="function"&&require;for(var o=0;o<r.length;o++)s(r[o]);return s})({1:[function(require,module,exports){
'use strict';

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { 'default': obj }; }

function _classCallCheck(instance, Constructor) { if (!(instance instanceof Constructor)) { throw new TypeError('Cannot call a class as a function'); } }

var _websocket = require('./websocket');

var _websocket2 = _interopRequireDefault(_websocket);

var GameManager = (function () {
    function GameManager() {
        _classCallCheck(this, GameManager);

        this.webSocketHandler = new _websocket2['default']();
    }

    GameManager.prototype.start = function start() {
        this.webSocketHandler.connect();
    };

    return GameManager;
})();

var gm = new GameManager();
gm.start();

},{"./websocket":2}],2:[function(require,module,exports){
"use strict";

exports.__esModule = true;

function _classCallCheck(instance, Constructor) { if (!(instance instanceof Constructor)) { throw new TypeError("Cannot call a class as a function"); } }

var WebSocketHandler = (function () {
    function WebSocketHandler() {
        _classCallCheck(this, WebSocketHandler);

        this.webSocket = null;
    }

    WebSocketHandler.prototype.connect = function connect() {
        this.webSocket = new WebSocket("ws://127.0.0.1:8080/socket");
        this.webSocket.onopen = function (event) {
            console.log("Connection opened");
        };
        this.webSocket.onmessage = function (event) {
            console.log("Received message", event.data);
        };
        this.webSocket.onclose = function (event) {
            console.log("Connection closed", event.data);
        };
    };

    return WebSocketHandler;
})();

exports["default"] = WebSocketHandler;
module.exports = exports["default"];

},{}]},{},[1,2]);
