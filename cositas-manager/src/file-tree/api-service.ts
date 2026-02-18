import {inject, Injectable} from '@angular/core';
import {HttpClient, HttpErrorResponse} from '@angular/common/http';
import {catchError, Observable, throwError} from 'rxjs';

export interface FileTreeItem {
  permissions: string;
  size: string;
  name: string;
  isDirectory: boolean;
}

export interface MoveActionBody {
  sourceName: string;
  destinationDirectory: string;
}

export interface ActionResponse {
  commandOutput: string;
}

@Injectable({providedIn: 'root'})
export class ApiService {
  private http = inject(HttpClient);

  getFileTree(): Observable<FileTreeItem[]> {
    return this.http.post<FileTreeItem[]>(`/api/info/listfiles`, null)
      .pipe(catchError(this.handleError));
  }

  getMediaTree(): Observable<string[]> {
    return this.http.post<string[]>(`/api/info/treemedia`, null)
      .pipe(catchError(this.handleError));
  }

  chmod(): Observable<ActionResponse> {
    return this.http.post<ActionResponse>(`/api/action/chmod`, null)
      .pipe(catchError(this.handleError));
  }

  _7zzip001(): Observable<ActionResponse> {
    return this.http.post<ActionResponse>(`/api/action/7zzip001`, null)
      .pipe(catchError(this.handleError));
  }

  _7zzip(): Observable<ActionResponse> {
    return this.http.post<ActionResponse>(`/api/action/7zzip`, null)
      .pipe(catchError(this.handleError));
  }

  _7z7z001(): Observable<ActionResponse> {
    return this.http.post<ActionResponse>(`/api/action/7z7z001`, null)
      .pipe(catchError(this.handleError));
  }

  rmzip(): Observable<ActionResponse> {
    return this.http.post<ActionResponse>(`/api/action/rmzip`, null)
      .pipe(catchError(this.handleError));
  }

  move(body: MoveActionBody): Observable<ActionResponse> {
    return this.http.post<ActionResponse>(`/api/action/move`, body)
      .pipe(catchError(this.handleError));
  }

  private handleError(error: HttpErrorResponse) {
    if (error.status === 0) {
      // A client-side or network error occurred. Handle it accordingly.
      alert(`An error occurred: ${JSON.stringify(error.error)}`);
    } else {
      // The backend returned an unsuccessful response code.
      // The response body may contain clues as to what went wrong.
      alert(
        `Backend returned code ${error.status}, body was:\n${JSON.stringify(error.error, null, 2)}`);
    }
    // Return an observable with a user-facing error message.
    return throwError(() => new Error('Something bad happened; please try again later.'));
  }
}
