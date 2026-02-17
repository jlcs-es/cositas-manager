import {Component, effect, inject} from '@angular/core';
import {ApiService, FileTreeItem} from './api-service';
import {Observable} from 'rxjs';
import {AsyncPipe} from '@angular/common';

@Component({
  selector: 'file-tree',
  template: `
    <div>
      @if (fileTree$ | async; as fileTree) {
        @for (item of fileTree; track item.name) {
          <p>
            <input type="radio" name="selectedItem" />
            {{item.permissions}} {{item.size}} {{item.name}}@if(item.isDirectory){/}
          </p>
        } @empty {
          <p>There are no items.</p>
        }
      } @else {
        <p>Loading...</p>
      }
    </div>


    <button box-="round" (click)="refreshFileTree()"> Refresh</button>
    <button box-="round">󰉒 Move to...</button>
  `,
  imports: [AsyncPipe]
})
export class FileTree {
  private apiService = inject(ApiService);
  fileTree$!: Observable<FileTreeItem[]>;

  constructor() {
    effect(() => {
      this.fileTree$ = this.apiService.getFileTree();
    });
  }

  refreshFileTree() {
    this.fileTree$ = this.apiService.getFileTree();
  }
}


