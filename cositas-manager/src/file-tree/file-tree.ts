import {Component, effect, inject} from '@angular/core';
import {FormsModule} from '@angular/forms';
import {ApiService, FileTreeItem, MoveActionBody} from './api-service';
import {Observable} from 'rxjs';
import {AsyncPipe} from '@angular/common';

@Component({
  selector: 'file-tree',
  template: `
    <div>
      @if (fileTree$ | async; as fileTree) {
        @for (item of fileTree; track item.name) {
          <p>
            <input type="radio" name="selectedItem" id="{{item.name}}" [value]="item.name" [(ngModel)]=selectedItem/>
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

    <button box-="round" popovertarget="moveDialog">󰉒 Move to...</button>
    <dialog id="moveDialog" container-="fill" popover>
      <div box-="round">
        @if(selectedItem) {
          <p>Move {{selectedItem}} to ...<p>

          <div is-="separator"></div>

          <div box-="round" shear-="top">
            <div class="row">
                <span is-="badge">Target directory</span>
            </div>
            <input id="targetDir" type="text" [(ngModel)]="targetDir" style="width: 50ch;"/>
          </div>

          <div is-="separator"></div>

          <p>
            <select name="mediaTree" id="mediaTree" [(ngModel)]="mediaDir" style="width: 50ch;">
              @if (mediaTree$ | async; as mediaTree) {
                @for (item of mediaTree; track item) {
                  <option value="{{item}}">{{item}}</option>
                } @empty {
                  <option>There are no items.</option>
                }
              } @else {
                <p>Loading...</p>
              }
            </select>
            <button box-="square" (click)="useMediaDir()">Use</button>
          </p>

          <div is-="separator"></div>

          <div style="float: right;">
            <button box-="round" variant-="red" popovertarget="moveDialog" popovertargetaction="hide">Cancel</button>
            <button box-="round" variant-="green" popovertarget="moveDialog" popovertargetaction="hide" (click)="moveActionAPI()">Confirm</button>
          </div>
        } @else {
          Select an item first.
          <div style="float: right;">
            <button box-="round" variant-="red" popovertarget="moveDialog" popovertargetaction="hide">OK</button>
          </div>
        }
      </div>
    </dialog>
  `,
  imports: [AsyncPipe, FormsModule]
})
export class FileTree {
  private apiService = inject(ApiService);
  fileTree$!: Observable<FileTreeItem[]>;
  mediaTree$!: Observable<string[]>
  targetDir = '/media/Movies/'
  selectedItem: string = '';
  mediaDir: string = '';

  constructor() {
    effect(() => {
      this.fileTree$ = this.apiService.getFileTree();
      this.mediaTree$ = this.apiService.getMediaTree();
    });
  }

  refreshFileTree() {
    this.fileTree$ = this.apiService.getFileTree();
  }

  moveActionAPI() {
    let moveBody: MoveActionBody = {
      destinationDirectory: this.targetDir,
      sourceName: this.selectedItem
    }
    this.apiService.move(moveBody).subscribe((actionResponse) => {
      alert(actionResponse.commandOutput);
    });
  }

  useMediaDir() {
    this.targetDir = this.mediaDir
  }
}


