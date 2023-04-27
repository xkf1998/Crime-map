import { Component } from '@angular/core';
import { MapsAPILoader } from '@agm/core';
import { GrpcService } from './grpc.service';


@Component({
  selector: 'app-root',
  templateUrl: 'app.component.html',
  styleUrls: ['app.component.css'],
})
export class AppComponent {

  title = 'Crime Heat Map';
  lat = 51.678418;
  lng = 7.809007;
  /*  lat = 40.712776;
  lng = -74.005974;*/
  zoom = 12;
  rectangles = [];
  map: any;

  topLatBound: number;
  botLatBound: number;
  leftLongBound: number;
  rightLongBound: number;

  crimes = [
    { lat: 51.678418, lng: 7.809007, crimeNumber: 5 },
    { lat: 51.679418, lng: 7.810007, crimeNumber: 10 },
    { lat: 51.680418, lng: 7.811007, crimeNumber: 3 },
    { lat: 51.681418, lng: 7.812007, crimeNumber: 8 },
    // Add more crime data here
  ];

  gridRows = 50;
  gridCols = 100;
  grid: any[][] = Array.from({ length: this.gridRows }, () => Array(this.gridCols).fill(0));

  cities = ['New York', 'Chicago', 'Austin', 'Los Angeles', 'San Francisco', 'Seattle'];
  selectedCity: string = 'New York';



  constructor(private mapsAPILoader: MapsAPILoader, private grpcService: GrpcService) {
    this.rectangles = new Array(this.gridRows);
    for (let i = 0; i < this.gridRows; i++) {
      this.rectangles[i] = new Array(this.gridCols).fill(null);
    }

  }
  responseMessage: string;
  sendGrpcRequest(): void {
    this.grpcService.sayHello('John Doe')
      .then(response => {
        this.responseMessage = response.getMessage();
        console.log(this.responseMessage)
      })
      .catch(error => {
        console.error('Error:', error);
      });
  }

  mapCrimeToGridCell(crimeLat: number, crimeLng: number): { row: number; col: number } {
    const latRange = this.topLatBound - this.botLatBound;
    const lngRange = this.rightLongBound - this.leftLongBound;

    const row = Math.floor(((this.topLatBound - crimeLat) / latRange) * this.gridRows);
    const col = Math.floor(((crimeLng - this.leftLongBound) / lngRange) * this.gridCols);

    return { row, col };
  }

  isCrimeWithinBounds(crimeLat: number, crimeLng: number): boolean {
    return (
      crimeLat >= this.botLatBound &&
      crimeLat <= this.topLatBound &&
      crimeLng >= this.leftLongBound &&
      crimeLng <= this.rightLongBound
    );
  }
  resetGrid() {
    this.grid = Array.from({ length: this.gridRows }, () => Array(this.gridCols).fill(0));
  }


  countCrimesInGrid() {
    this.resetGrid();
    this.crimes.forEach((crime) => {
      if (this.isCrimeWithinBounds(crime.lat, crime.lng)) {
        const { row, col } = this.mapCrimeToGridCell(crime.lat, crime.lng);
        this.grid[row][col] += crime.crimeNumber;
      }
    });

    this.mapsAPILoader.load().then(() => {
      this.createOrUpdateRectangles();
    });
  }

  createOrUpdateRectangles() {
    if (!this.map) {
      return;
    }

    for (let i = 0; i < this.gridRows; i++) {
      for (let j = 0; j < this.gridCols; j++) {
        const crimeCount = this.grid[i][j];
        const bounds = this.getRectangleBounds(i, j);
        const fillColor = this.getCellColor(crimeCount);

        if (crimeCount > 0) {
          if (!this.rectangles[i][j]) {
            const rectangle = new google.maps.Rectangle({
              map: this.map,
              bounds,
              fillColor,
              fillOpacity: 1,
              strokeWeight: 0,
            });
            this.rectangles[i][j] = rectangle;
          } else {
            this.rectangles[i][j].setOptions({ bounds, fillColor });
          }
        } else if (this.rectangles[i][j]) {
          this.rectangles[i][j].setMap(null);
          this.rectangles[i][j] = null;
        }
      }
    }
  }

  onMapReady(map: google.maps.Map) {
    this.map = map;
    this.sendGrpcRequest();
  }

  getRectangleBounds(row: number, col: number): google.maps.LatLngBoundsLiteral {
    const latRange = this.topLatBound - this.botLatBound;
    const lngRange = this.rightLongBound - this.leftLongBound;
    const latStep = latRange / this.gridRows;
    const lngStep = lngRange / this.gridCols;

    const north = this.topLatBound - row * latStep;
    const south = north - latStep;
    const west = this.leftLongBound + col * lngStep;
    const east = west + lngStep;

    return {
      north,
      south,
      east,
      west,
    };
  }

  getCellColor(crimeCount: number): string {
    if (crimeCount <= 3) {
      return 'rgba(76, 175, 80, 0.5)'; // Green
    } else if (crimeCount > 3 && crimeCount <= 7) {
      return 'rgba(255, 193, 7, 0.5)'; // Yellow
    } else {
      return 'rgba(244, 67, 54, 0.5)'; // Red
    }
  }


  onBoundsChange(bounds: google.maps.LatLngBounds) {
    if (!bounds) {
      return;
    }

    const ne = bounds.getNorthEast();
    const sw = bounds.getSouthWest();

    this.topLatBound = ne.lat();
    this.botLatBound = sw.lat();
    this.leftLongBound = sw.lng();
    this.rightLongBound = ne.lng();
    this.countCrimesInGrid();
  }

  onCityChange() {
    const tmp = this.getCityLatLng(this.selectedCity);
    this.lat = tmp['lat'];
    this.lng = tmp['lng'];
    this.zoom = 12;
    // call other functions or update variables based on the selected city
  }

  getCityLatLng(city: string) {
    switch (city) {
      case 'New York':
        return { lat: 40.712776, lng: -74.005974 };
      case 'Los Angeles':
        return { lat: 34.052235, lng: -118.243683 };
      // add more cases for other cities
      case 'Chicago':
        return { lat: 41.878113, lng: -87.629799 };
      case 'Austin':
        return { lat: 30.267153, lng: -97.743057 };
      case 'San Francisco':
        return { lat: 37.774929, lng: -122.419416 };
      case 'Seattle':
        return { lat: 47.606209, lng: -122.332071 };
      default:
        return { lat: 0, lng: 0 }; // default to (0, 0) if city not found
    }
  }

}