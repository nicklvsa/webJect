import { Component, OnInit, ViewChild, ElementRef, Directive } from '@angular/core';
import axios from 'axios';

@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.scss']
})
export class HomeComponent implements OnInit {

  pkgID = '';

  constructor() { }

  @ViewChild('appname') appName: ElementRef;

  getID() {
    const data = {
      app_path: this.appName.nativeElement.value
    };
    const jsonData = JSON.stringify(data);
    console.log(jsonData);
    axios.post('http://localhost:8081/tweak/identify', jsonData).then((response) => {
      console.log(response.data);
      const dataRes = JSON.parse(response.data);
      this.pkgID = dataRes;
    }).catch((err) => {
      console.log(err);
    });
  }

  addTweak() {}

  ngOnInit(): void {
  }

}
