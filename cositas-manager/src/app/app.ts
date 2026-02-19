import {Component, effect, inject, signal, WritableSignal} from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { FileTree} from '../file-tree/file-tree';
import {ApiService, MoveActionBody} from '../file-tree/api-service';
import {catchError, startWith, switchMap, tap, timer} from 'rxjs';

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
  jobsList: WritableSignal<string[]> = signal([])

  constructor() {
    effect(() => {
      timer(0, 5000).pipe(
        switchMap(() => {
          return this.apiService.getJobs()
        }),
        tap((data) => this.jobsList.set(data)),
        catchError(error => {
          throw new Error(`Polling error: ${error}`);
        })
      ).subscribe();
    });
  }

  hideActionOutputDialog() {
    this.showActionOutputDialog.set(false);
    this.actionOutput = '';
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

  moveActionAPI() {
    this.loading.set(true);
    this.showActionOutputDialog.set(true);
    let moveBody: MoveActionBody = {
      destinationDirectory: '',
      sourceName: ''
    }
    this.apiService.move(moveBody).subscribe((actionResponse) => {
      this.loading.set(false);
      this.actionOutput = actionResponse.commandOutput;
    });
  }
}
