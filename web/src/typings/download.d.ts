namespace Download {
  type Task = {
    id: number;
    item_id: number;
    msg_id: number;
    text: string;
    size: number;
    date: number;
    priority: number;
    downloading: boolean;
    downloaded: boolean;
    fatal_error: boolean;
    error: string;
    error_at: number;
    file?: Telegram.File;
  };
}
