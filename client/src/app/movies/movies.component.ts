import { Component, OnInit } from '@angular/core';
import { Http } from '@angular/http';
import { Router } from '@angular/router';
import { Location } from '@angular/common';


@Component({
  selector: 'app-movies',
  templateUrl: './movies.component.html',
  styleUrls: ['./movies.component.css']
})
export class MoviesComponent implements OnInit {

  public movies: any;

  constructor(private http: Http, private router: Router, private location: Location) {
    this.movies = [];
  }

  ngOnInit() {
    this.location.subscribe(() => {
      this.refresh();
    });
    this.refresh();
  }

  private refresh() {
    this.http.get("http://localhost:3333/movies").subscribe((response: any) => {
      console.log("response", response)
      this.movies = JSON.parse(response._body)
    })
      
  }

  /**
   * search
   */
  public search(event :any) {
    let url = "http://localhost:3333/movies";
    if (event.target.value) {
      url = "http://localhost:3333/search/" + event.target.value;
    }
    this.http.get(url)
      .subscribe((response: any) => {
        this.movies = JSON.parse(response._body)
      })
    
  }

  public create() {
    this.router.navigate(["create"])
  }

}
