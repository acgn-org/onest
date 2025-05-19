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

  type MatchedRaw = Raw & {
    matched: boolean;
    matched_text: string;
  };

  type Rule = {
    id: number;
    name: string;
    regexp: string;
    cn_index: number;
    en_index: number;
    channel_id: number;
    channel_name: string;
  };
}
