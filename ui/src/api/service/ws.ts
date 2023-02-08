import { ReplaySubject, Subject, filter, first, lastValueFrom, retry } from "rxjs";
import { WebSocketSubject } from "rxjs/webSocket";

const nextId = (() => {
  let id = 0;
  return () => {
    id++;
    if (id > 199999999) {
      id = 1;
    }
    return `${id}`;
  };
})();

export interface Request {
  action: string;
  params?: Array<any> | {};
  id?: string | number;
}

export interface Response {
  id?: string;
  data: any;
  msg: any;
  code: number;
}

export enum ConnectionState {
  CONNECTED = "Connected",
  CONNECTING = "Connecting",
  CLOSING = "Closing",
  DISCONNECTED = "Disconnected",
}

export class WebSocketService {
  private connectionState = new ReplaySubject<ConnectionState>(1);
  private socket: WebSocketSubject<ArrayBuffer | Object>;

  private messageObserver = new Subject<any>();
  private binaryObserver = new Subject<ArrayBuffer>();

  constructor() {
    this.socket = null as any;
  }

  init = (url: string) => {
    this.connectionState.next(ConnectionState.CONNECTING);

    this.socket = new WebSocketSubject({
      binaryType: "arraybuffer",
      url,
      openObserver: {
        next: () => this.connectionState.next(ConnectionState.CONNECTED),
      },
      closingObserver: {
        next: () => this.connectionState.next(ConnectionState.CLOSING),
      },
      closeObserver: {
        next: () => this.connectionState.next(ConnectionState.DISCONNECTED),
      },
      deserializer: (e: MessageEvent) => {
        try {
          if (e.data instanceof ArrayBuffer) {
            return e.data;
          } else {
            return JSON.parse(e.data);
          }
        } catch (e) {
          console.error(e);
          return null;
        }
      },
    });

    // message
    this.socket.pipe(
      retry(),
      filter((v: any) => !(v instanceof ArrayBuffer)),
    ).subscribe((message) => {
      this.messageObserver.next(message);
    });

    // binary message
    this.socket.pipe(
      retry(),
      filter((value: any) => value instanceof ArrayBuffer),
    ).subscribe((message) => {
      this.binaryObserver.next(message);
    });

    this.connectionState.subscribe((state) => {
      console.log(`WebSocket state ${state}`);
    });
  };

  request = async (request: Request): Promise<any> => {
    if (!request.id) {
      request.id = nextId();
    }
    if (!request.params) {
      request.params = [];
    }

    this.socket.next(request);

    const obs = this.messageObserver.pipe(
      filter((v: any) => request.id === v.id),
      first(),
    );

    return lastValueFrom(obs).then((message: Response) => {
      if (message.code !== 200) {
        throw message;
      }
      return message.data;
    });
  };
}


