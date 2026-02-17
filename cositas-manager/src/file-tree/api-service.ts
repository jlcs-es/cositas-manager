import {inject, Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs';

export interface FileTreeItem {
  permissions: string;
  name: string;
  isDirectory: boolean;
}

export interface ActionResponse {
  commandOutput: string;
}

@Injectable({providedIn: 'root'})
export class ApiService {
  private http = inject(HttpClient);

  getFileTree(): Observable<FileTreeItem[]> {
    return this.http.post<FileTreeItem[]>(`/api/info/listfiles`, null);
  }

  chmod(): Observable<ActionResponse> {
    return this.http.post<ActionResponse>(`/api/action/chmod`, null);
  }

  _7zzip001(): Observable<ActionResponse> {
    return this.http.post<ActionResponse>(`/api/action/7zzip001`, null);
  }

  _7zzip(): Observable<ActionResponse> {
    return this.http.post<ActionResponse>(`/api/action/7zzip`, null);
  }

  _7z7z001(): Observable<ActionResponse> {
    return this.http.post<ActionResponse>(`/api/action/7z7z001`, null);
  }

  rmzip(): Observable<ActionResponse> {
    return this.http.post<ActionResponse>(`/api/action/rmzip`, null);
  }

}
