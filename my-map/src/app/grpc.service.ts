import { Injectable } from '@angular/core';
import { grpc } from '@improbable-eng/grpc-web';
import { GreeterClient } from './generated/HelloworldServiceClientPb';
import { HelloRequest, HelloResponse } from './generated/helloworld_pb';

@Injectable({
  providedIn: 'root'
})
export class GrpcService {
  private client: GreeterClient;

  constructor() {
    const transport = grpc.WebsocketTransport();
    this.client = new GreeterClient('http://localhost:8080', null, { transport });
  }

  public sayHello(name: string): Promise<HelloResponse> {
    return new Promise((resolve, reject) => {
      const request = new HelloRequest();
      request.setName(name);

      this.client.sayHello(request, null, (error, response) => {
        if (error) {
          reject(error);
        } else {
          resolve(response);
        }
      });
    });
  }
}
