namespace Item {
  type Local = {
    id: number;
    channel_id: number;
    name: string;
    regexp: string;
    pattern: string;
    match_pattern: string;
    match_content: string;
    date_start: number;
    date_end: number;
    process: number;
    priority: number;
    target_path: string;
  };

  type Remote = {
    id: number;
    rule_id: number;
    name: string;
    name_cn: string;
    name_en: string;
    date_start: number;
    date_end: number;
  };

  type ViewMode = "active" | "error" | "all";
}
