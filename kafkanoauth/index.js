"use strict";
exports.__esModule = true;
exports.disconnect = exports.getConnection = void 0;
var kafkajs_1 = require("kafkajs");
var kafka = new kafkajs_1.Kafka({
    clientId: 'cxnpl-poc-sa',
    brokers: ['0.rp-4260ba7.e539449.byoc.vectorized.cloud:30684']
});
// const kafka = new Kafka({
//     clientId: 'chat-app',
//     brokers: ['one-node-cluster-0.one-node-cluster.panda-chat.svc.cluster.local:9092']
// });
var producer = kafka.producer();
function getConnection(user) {
    return producer.connect().then(function () {
        return function (message) {
            return producer.send({
                topic: 'test.topic',
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
getConnection('test').then(function (sendMessage) {
    sendMessage('test message');
});
