import { Component, OnInit, ViewChild, ElementRef, Directive } from '@angular/core';

@Component({
  selector: 'app-nav',
  templateUrl: './nav.component.html',
  styleUrls: ['./nav.component.scss']
})
export class NavComponent implements OnInit {

  constructor() {}

  @ViewChild('navBurger', { read: ElementRef }) navBurger: ElementRef;
  @ViewChild('navMenu', { read: ElementRef }) navMenu: ElementRef;

  toggleNav() {
    this.navBurger.nativeElement.classList.toggle('is-active');
    this.navMenu.nativeElement.classList.toggle('is-active');
  }

  ngOnInit(): void {}

}
