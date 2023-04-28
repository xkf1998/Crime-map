import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { Crime, GetCrimesRequest, GetCrimesResponse } from './crime_generated/crime_pb';
import { CrimeMapClient } from './crime_generated/crime_grpc_web_pb';

@Injectable({
  providedIn: 'root'
})
export class CrimeService {
  private client: CrimeMapClient;

  constructor() {
    this.client = new CrimeMapClient('http://localhost:8080', null, null);
  }
  getCrimes(request: GetCrimesRequest): Observable<GetCrimesResponse> {
    return new Observable<GetCrimesResponse>(observer => {
      this.client.getCrimes(request, {}, (err, response) => {
        if (err) {
          observer.error(err);
        } else {
          observer.next(response);
        }
        observer.complete();
      });
    });
  }

}
