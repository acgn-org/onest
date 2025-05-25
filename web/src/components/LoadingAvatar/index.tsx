import { type FC, useState } from "react";

import { Avatar, Skeleton } from "@mantine/core";

export interface LoadingAvatarProps {
  src?: string;
  alt: string;
}

export const LoadingAvatar: FC<LoadingAvatarProps> = ({ src, alt }) => {
  const [isLoading, setIsLoading] = useState(true);
  const isAvatarVisible = src && !isLoading;
  return (
    <>
      <Avatar
        src={src}
        alt={alt}
        onLoad={() => setIsLoading(false)}
        style={{
          display: isAvatarVisible ? undefined : "none",
        }}
      />
      {!isAvatarVisible && <Skeleton height={38} circle />}
    </>
  );
};
export default LoadingAvatar;
