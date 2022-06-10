import {Kafka} from 'kafkajs';

const kafka = new Kafka({
    clientId: 'chat-app',
    brokers: ['one-node-cluster-0.one-node-cluster.panda-chat.svc.cluster.local:9092']
});
// const kafka = new Kafka({
//     clientId: 'chat-app',
//     brokers: ['one-node-cluster-0.one-node-cluster.panda-chat.svc.cluster.local:9092']
// });

const producer = kafka.producer();

export function getConnection(user: string){
  return producer.connect().then(() => {
    return (message: string) => {
      return producer.send({
        topic: 'chat-rooms', // the topic created before
        messages: [//we send the message and the user who sent it
          {value: JSON.stringify({message, user})},
        ],
      })
    }
  })
}

export function disconnect(){
  return producer.disconnect()
}

getConnection('test').then(sendMessage => {
  sendMessage('test message')
})
