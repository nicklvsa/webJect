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

  @ViewChild('appName', { read: ElementRef }) appName: ElementRef;

  getID() {
    const data = {
      app_path: this.appName.nativeElement.value
    };
    axios.post('http://localhost:8081/tweak/identify', {data}).then((response) => {
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
