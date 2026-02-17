import {Component, inject, signal} from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { FileTree} from '../file-tree/file-tree';
import {ApiService} from '../file-tree/api-service';
import {HttpErrorResponse} from '@angular/common/http';
import {throwError} from 'rxjs';

@Component({
  selector: 'app-root',
  imports: [RouterOutlet, FileTree],
  templateUrl: './app.html',
  styleUrl: './app.css'
})
export class App {
  protected readonly title = signal('cositas-manager');
  private apiService = inject(ApiService);
  actionOutput: string = '';
  showActionOutputDialog = signal(false);

  hideActionOutputDialog() {
    this.showActionOutputDialog.set(false);
  }

  chmodActionAPI() {
    this.apiService.chmod().subscribe((actionResponse) => {
      this.actionOutput = actionResponse.commandOutput;
      this.showActionOutputDialog.set(true);
    });
  }

  _7zzip001ActionAPI() {
    this.apiService._7zzip001().subscribe((actionResponse) => {
      this.actionOutput = actionResponse.commandOutput;
      this.showActionOutputDialog.set(true);
    });
  }

  _7zzipActionAPI() {
    this.apiService._7zzip().subscribe((actionResponse) => {
      this.actionOutput = actionResponse.commandOutput;
      this.showActionOutputDialog.set(true);
    });
  }

  _7z7z001ActionAPI() {
    this.apiService._7z7z001().subscribe((actionResponse) => {
      this.actionOutput = actionResponse.commandOutput;
      this.showActionOutputDialog.set(true);
    });
  }

  rmzipActionAPI() {
    this.apiService.rmzip().subscribe((actionResponse) => {
      this.actionOutput = actionResponse.commandOutput;
      this.showActionOutputDialog.set(true);
    });
  }
}
