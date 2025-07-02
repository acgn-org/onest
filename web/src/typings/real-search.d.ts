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
    selected: boolean;
    match_content: string | null;
    matched: boolean;
    matched_target: string;
    matched_target_repeated: boolean;
    priority: number | undefined;
  };

  type Rule = {
    id: number;
    name: string;
    regexp: string;
    cn_index: number;
    en_index: number;
    sample_pattern: string;
    channel_id: number;
    channel_name: string;
  };
}
