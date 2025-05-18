namespace RealSearch {
  type ScheduleItem = {
    item: Item.Remote;
    data: Raw[];
  };

  type Raw = {
    id: number;
    item_id: number;
    channel_id: number;
    channel_name: string;
    size: number;
    text: string;
    file_suffix: string;
    msg_id: number;
    supports_streaming: boolean;
    link: string;
    date: number;
  };
}
