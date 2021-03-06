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
  @ViewChild('bundleid') bundleID: ElementRef;
  @ViewChild('zipfile') zipFile: ElementRef;

  getID() {
    const data = {
      app_path: this.appName.nativeElement.value
    };
    const jsonData = JSON.stringify(data);
    console.log(jsonData);
    axios.post('http://localhost:8081/tweak/identify', jsonData).then((response) => {
      console.log(response.data);
      this.pkgID = response.data['content_msg'] || response.data['err_msg'];
    }).catch((err) => {
      console.log(err);
    });
  }

  addTweak() {
    const file = this.zipFile.nativeElement.file.files[0];
    const id = this.bundleID.nativeElement.value;
    if (id !== '' && file != null) {
      axios.post('http://localhost:8081/tweak/add', file).then((response) => {
        console.log(response.data);
      }).catch((err) => {
        console.log(err);
      });
    }
  }

  ngOnInit(): void {
  }

}
