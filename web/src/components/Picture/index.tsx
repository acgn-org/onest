import type { FC, CSSProperties } from "react";

export interface PictureProps {
  name: string;
  alt: string;
  dir?: string;
  defaultType?: string;
  sources?: PictureSource[];
  className?: string;
  aspectRatio?: number | string;
  imgStyle?: CSSProperties;
}

export interface PictureSource {
  fileType: string;
  mimeType: string;
}

const Picture: FC<PictureProps> = ({
  name,
  alt,
  defaultType = "png",
  sources = [
    { fileType: "avif", mimeType: "image/avif" },
    { fileType: "webp", mimeType: "image/webp" },
  ],
  className,
  imgStyle,
  aspectRatio,
}) => {
  const getImageUrl = (name: string, type: string) => {
    return new URL(`/src/assets/${name}.${type}`, import.meta.url).href;
  };
  return (
    <picture className={className} style={{ display: "flex" }}>
      {sources.map((source) => (
        <source
          key={JSON.stringify(source.fileType)}
          srcSet={getImageUrl(name, source.fileType)}
          type={source.mimeType}
        />
      ))}
      <img
        alt={alt}
        src={getImageUrl(name, defaultType)}
        style={{
          ...imgStyle,
          aspectRatio,
        }}
      />
    </picture>
  );
};
export default Picture;
