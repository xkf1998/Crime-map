import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { AgmCoreModule } from '@agm/core';
import { FormsModule } from '@angular/forms';
import { CrimeService } from './crime.service';



import { AppComponent } from './app.component';

@NgModule({
  imports: [
    BrowserModule,
    FormsModule,
    AgmCoreModule.forRoot({
      apiKey: 'AIzaSyBMyIUzztDz_jKItOcYChVYOxAdbm7rIFI'
    })
  ],
  providers: [CrimeService],
  declarations: [AppComponent],
  bootstrap: [AppComponent]
})
export class AppModule { }