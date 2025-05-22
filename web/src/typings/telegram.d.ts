namespace Telegram {
  type Meta = {
    "@type": string;
    "@extra": string;
    "@client_id": number;
  };

  type LocalFile = Meta & {
    can_be_downloaded: boolean;
    can_be_deleted: boolean;
    is_downloading_active: boolean;
    is_downloading_completed: boolean;
    download_offset: number;
    download_prefix_size: number;
    downloaded_size: number;
  };

  type RemoteFile = Meta & {
    id: string;
    unique_id: string;
    is_uploading_active: boolean;
    is_uploading_completed: boolean;
    uploaded_size: number;
  };

  type File = Meta & {
    id: number;
    size: number;
    expected_size: number;
    local: LocalFile;
    remote: RemoteFile;
  };

  type ChatType = Meta & {
    supergroup_id: number;
    is_channel: boolean;
  };

  type MiniThumbnail = Meta & {
    width: number;
    height: number;
    data: string;
  };

  type ChatPhotoInfo = Meta & {
    small: File;
    big: File;
    minithumbnail: MiniThumbnail;
    has_animation: boolean;
    is_personal: boolean;
  };

  type Chat = Meta & {
    id: number;
    chat_id: number;
    title: string;
    type: ChatType;
    photo: ChatPhotoInfo;
    data: number;
  };
}
