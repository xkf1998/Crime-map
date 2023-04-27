import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { AgmCoreModule } from '@agm/core';
import { FormsModule } from '@angular/forms';
import { GrpcService } from './grpc.service';



import { AppComponent } from './app.component';

@NgModule({
  imports: [
    BrowserModule,
    FormsModule,
    AgmCoreModule.forRoot({
      apiKey: 'AIzaSyBMyIUzztDz_jKItOcYChVYOxAdbm7rIFI'
    })
  ],
  providers: [GrpcService],
  declarations: [AppComponent],
  bootstrap: [AppComponent]
})
export class AppModule { }