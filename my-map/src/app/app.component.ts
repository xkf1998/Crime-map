import { Component } from '@angular/core';
import { MapsAPILoader, AgmMarker } from '@agm/core';
import { GrpcService } from './grpc.service';
import { CrimeService } from './crime.service';
import { GetCrimesRequest } from './crime_generated/crime_pb';
import { ThemePalette } from '@angular/material/core';
import { FormsModule } from '@angular/forms';
import { ProgressSpinnerMode } from '@angular/material/progress-spinner';


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

  start_date: string;
  end_date: string;
  startTimestamp = 0;
  endTimestamp = 2000000000;

  topLatBound: number;
  botLatBound: number;
  leftLongBound: number;
  rightLongBound: number;

  crimes: any;
  crimes_dict: any;
  crimes_display: any;

  gridRows = 50;
  gridCols = 100;
  grid: any[][] = Array.from({ length: this.gridRows }, () => Array(this.gridCols).fill(0));

  cities = ['New York', 'Chicago', 'Austin', 'Los Angeles', 'San Francisco', 'Seattle'];
  selectedCity: string = '';

  public isHeatMapEnabled = true;

  selectedMarker: any;
  infoWindowOpen = false;
  isLoading: boolean = false;

  color: ThemePalette = 'warn';
  mode: ProgressSpinnerMode = 'indeterminate';

  constructor(private mapsAPILoader: MapsAPILoader, private grpcService: GrpcService, private crimeService: CrimeService) {
    this.rectangles = new Array(this.gridRows);
    for (let i = 0; i < this.gridRows; i++) {
      this.rectangles[i] = new Array(this.gridCols).fill(null);
    }
    this.crimes = [
      { lat: 51.678418, lng: 7.809007, crimeNumber: 5, description: "1" },
      { lat: 51.679418, lng: 7.810007, crimeNumber: 10, description: "2" },
      { lat: 51.680418, lng: 7.811007, crimeNumber: 3, description: "3" },
      { lat: 51.681418, lng: 7.812007, crimeNumber: 8, description: "4" },
      // require the initial data for the map, now hardcoded.
    ];

  }

  toggleMapMode() {
    this.isHeatMapEnabled = !this.isHeatMapEnabled;
  }

  fetchCrimeData(time_min: number, time_max: number, long_min: number, long_max: number, lat_min: number, lat_max: number) {
    this.isLoading = true;
    const request = new GetCrimesRequest();
    request.setTimeMin(time_min);
    request.setTimeMax(time_max);
    request.setLongitudeMin(long_min);
    request.setLongitudeMax(long_max);
    request.setLatitudeMin(lat_min);
    request.setLatitudeMax(lat_max);
    this.crimeService.getCrimes(request).subscribe(response => {
      this.crimes = [];
      this.crimes_dict = {};
      for (let i = 0; i < response['array'][0].length; i++) {
        var crime = {}
        crime['lat'] = response['array'][0][i][2];
        crime['lng'] = response['array'][0][i][1];
        crime['crimeNumber'] = 1;
        crime['description'] = response['array'][0][i][3];
        var timestamp = response['array'][0][i][0];
        const date = new Date(timestamp * 1000);
        crime['date'] = date.toISOString();
        this.crimes.push(crime);

        var key = String(crime['lat']) + " " + String(crime['lng'])
        if (key in this.crimes_dict) {
          this.crimes_dict[key].push(crime);
        }
        else {
          var tmp = [];
          tmp.push(crime);
          this.crimes_dict[key] = tmp;
        }
      }
      this.countCrimesInGrid();
    }, error => {
      console.error('Error:', error);
    }, () => {
      // 在请求完成时，将 isLoading 设置为 false
      this.isLoading = false;
    }
    );
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

  crimeCountToColor(m: number, M: number, k: number, crimeCount: number): string {
    return Math.floor(M - (M - m) * Math.exp(- k * crimeCount)).toString();
  }

  getCellColor(crimeCount: number): string {
    let k = 0.04;
    let r = this.crimeCountToColor(255, 255, k, crimeCount);
    let g = this.crimeCountToColor(255, 0, k, crimeCount);
    let b = this.crimeCountToColor(0, 0, k, crimeCount);
    let a = (0.7 - 0.7 * Math.exp(- 0.4 * crimeCount)).toString();
    return 'rgba(' + r + ', ' + g + ', ' + b + ', ' + a + ')';
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

  onBoundsChange2(bounds: google.maps.LatLngBounds) {
    if (!bounds) {
      return;
    }

    const ne = bounds.getNorthEast();
    const sw = bounds.getSouthWest();

    this.topLatBound = ne.lat();
    this.botLatBound = sw.lat();
    this.leftLongBound = sw.lng();
    this.rightLongBound = ne.lng();
    // this.createMarkers();
  }

  createMarkers() {
    for (const crime of this.crimes) {
      const marker = new AgmMarker(this.map);
      marker.latitude = crime.lat;
      marker.longitude = crime.lng;
      marker.label = crime.crimeNumber.toString();
      marker.markerClick.subscribe(() => {
        console.log('Marker clicked!');
      });
    }
  }


  onCityChange() {
    const tmp = this.getCityLatLng(this.selectedCity);
    this.lat = tmp['lat'];
    this.lng = tmp['lng'];
    this.zoom = 12;
    // this.fetchCrimeData(this.startTimestamp, this.endTimestamp, this.leftLongBound, this.rightLongBound, this.botLatBound, this.topLatBound);
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

  onDateChange() {
    if (this.start_date && this.end_date) {
      const start = new Date(this.start_date);
      const end = new Date(this.end_date);
      // Convert start and end dates to Unix timestamps
      this.startTimestamp = start.getTime() / 1000; // Divide by 1000 to get Unix timestamp in seconds
      this.endTimestamp = end.getTime() / 1000;

      if (this.startTimestamp < this.endTimestamp) {
        // this.fetchCrimeData(this.startTimestamp, this.endTimestamp, this.leftLongBound, this.rightLongBound, this.botLatBound, this.topLatBound);
        // this.countCrimesInGrid();
      } else {
        // Handle the case where start_date is after end_date
        this.start_date = null;
        this.end_date = null;
        alert("Error: start_date cannot be after end_date");
      }
    }
  }

  UpdateMap() {
    this.fetchCrimeData(this.startTimestamp, this.endTimestamp, this.leftLongBound, this.rightLongBound, this.botLatBound, this.topLatBound);
    this.countCrimesInGrid();
  }

  onMarkerClick(marker: any) {
    var key = String(marker['lat']) + " " + String(marker['lng'])
    this.crimes_display = this.crimes_dict[key];
    console.log(this.crimes_display);
    this.infoWindowOpen = true;
  }
  onCloseClick() {
    this.selectedMarker = null;
    this.infoWindowOpen = false;
  }
}