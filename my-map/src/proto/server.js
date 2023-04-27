const grpc = require('grpc');
const { GreeterService } = require('./helloworld_grpc_pb');

function sayHello(call, callback) {
    const name = call.request.getName();
    const response = new HelloResponse();
    response.setMessage(`Hello, ${name}!`);
    callback(null, response);
}

function main() {
    const server = new grpc.Server();
    server.addService(GreeterService, { sayHello });
    server.bind('localhost:8080', grpc.ServerCredentials.createInsecure());
    server.start();
    console.log('gRPC server started on port 8080');
}

main();
