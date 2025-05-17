namespace Item {
  type Item = {
    id: number;
    channel_id: number;
    name: string;
    regexp: string;
    pattern: string;
    date_start: number;
    date_end: number;
    process: number;
    priority: number;
    target_path: string;
  };
}
