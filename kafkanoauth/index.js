"use strict";
exports.__esModule = true;
exports.disconnect = exports.getConnection = void 0;
var kafkajs_1 = require("kafkajs");
var kafka = new kafkajs_1.Kafka({
    clientId: 'chat-app',
    brokers: ['127.0.0.1:9092']
});
var producer = kafka.producer();
function getConnection(user) {
    return producer.connect().then(function () {
        return function (message) {
            return producer.send({
                topic: 'chat-room',
                messages: [
                    { value: JSON.stringify({ message: message, user: user }) },
                ]
            });
        };
    });
}
exports.getConnection = getConnection;
function disconnect() {
    return producer.disconnect();
}
exports.disconnect = disconnect;
