import { Component, OnInit } from '@angular/core';
import { Location } from "@angular/common";
import { Http } from "@angular/http";

@Component({
  selector: 'app-create',
  templateUrl: './create.component.html',
  styleUrls: ['./create.component.css']
})
export class CreateComponent implements OnInit {

  public movie: any;

  constructor(private location: Location, private http: Http) {
    this.movie = {
      "name": "",
      "genre": "",
      "formats": {
        "digital": false,
        "bluray": false,
        "dvd": false
      }
    }
  }

  ngOnInit() {
  }

  public save() {
    if (this.movie.name && this.movie.genre) {
      this.http.post("http://localhost:3333/movies", JSON.stringify(this.movie))
        .subscribe(result => {
          this.location.back()
        });
    }
  }

}
