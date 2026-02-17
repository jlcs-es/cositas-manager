import {Component, effect, inject} from '@angular/core';
import {ApiService, FileTreeItem} from './api-service';
import {Observable} from 'rxjs';
import {AsyncPipe} from '@angular/common';

@Component({
  selector: 'file-tree',
  template: `
    <ul marker-="tree">
      @if (fileTree$ | async; as fileTree) {
        @for (item of fileTree; track item.name) {
          <li>
            <input type="radio" name="selectedItem" />
            {{item.permissions}} {{item.name}}@if(item.isDirectory){/}
          </li>
        } @empty {
          <li>There are no items.</li>
        }
      } @else {
        <li>Loading...</li>
      }
    </ul>


    <button box-="round"> Refresh</button>
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
}


