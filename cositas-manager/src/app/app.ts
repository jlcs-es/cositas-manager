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
  loading = signal(true);

  hideActionOutputDialog() {
    this.showActionOutputDialog.set(false);
  }

  chmodActionAPI() {
    this.loading.set(true)
    this.showActionOutputDialog.set(true);
    this.apiService.chmod().subscribe((actionResponse) => {
      this.loading.set(false)
      this.actionOutput = actionResponse.commandOutput;
    });
  }

  _7zzip001ActionAPI() {
    this.loading.set(true)
    this.showActionOutputDialog.set(true);
    this.apiService._7zzip001().subscribe((actionResponse) => {
      this.loading.set(false)
      this.actionOutput = actionResponse.commandOutput;
    });
  }

  _7zzipActionAPI() {
    this.loading.set(true)
    this.showActionOutputDialog.set(true);
    this.apiService._7zzip().subscribe((actionResponse) => {
      this.loading.set(false)
      this.actionOutput = actionResponse.commandOutput;
    });
  }

  _7z7z001ActionAPI() {
    this.loading.set(true)
    this.showActionOutputDialog.set(true);
    this.apiService._7z7z001().subscribe((actionResponse) => {
      this.loading.set(false)
      this.actionOutput = actionResponse.commandOutput;
    });
  }

  rmzipActionAPI() {
    this.loading.set(true)
    this.showActionOutputDialog.set(true);
    this.apiService.rmzip().subscribe((actionResponse) => {
      this.loading.set(false)
      this.actionOutput = actionResponse.commandOutput;
    });
  }
}
