import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { AgmCoreModule } from '@agm/core';
import { FormsModule } from '@angular/forms';
import { CrimeService } from './crime.service';



import { AppComponent } from './app.component';
import { LoadingSpinnerComponent } from './loading-spinner/loading-spinner.component';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatCardModule } from '@angular/material/card';



@NgModule({
  imports: [
    BrowserModule,
    FormsModule,
    AgmCoreModule.forRoot({
      apiKey: 'AIzaSyBMyIUzztDz_jKItOcYChVYOxAdbm7rIFI'
    }),
    BrowserAnimationsModule,
    MatProgressSpinnerModule,
    MatCardModule,
  ],
  providers: [CrimeService],
  declarations: [AppComponent, LoadingSpinnerComponent],
  bootstrap: [AppComponent]
})
export class AppModule { }