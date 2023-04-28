import { TestBed } from '@angular/core/testing';

import { CrimeService } from './crime.service';

describe('CrimeService', () => {
  let service: CrimeService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(CrimeService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
