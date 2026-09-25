"use client";

import SecondaryButton from "@/components/SecondaryButton";
import { useRouter } from "next/navigation";
import React, { useRef, useState } from "react";
import { Upload, Modal } from "antd";
import { BiExpand } from "react-icons/bi";
import { FaPlus } from "react-icons/fa6";
import styled from "styled-components";
import { FiX } from "react-icons/fi";

const MAX_GRID_IMAGES = 8;

const AssetImages = () => {
  const router = useRouter();
  const [images, setImages] = useState<string[]>([]);
  // Ref for the hidden file input trigger
  const uploadBtnRef = useRef<any>(null);

  // For preview modal
  const [previewOpen, setPreviewOpen] = useState(false);
  const [previewImage, setPreviewImage] = useState<string>("");
  const [previewTitle, setPreviewTitle] = useState<string>("");

  const handleAntdChange = ({ fileList }: { fileList: any[] }) => {
    const newImages = fileList
      .filter((file) => !file.url)
      .map((file) => URL.createObjectURL(file.originFileObj));
    setImages((prev) => [...prev, ...newImages]);
  };

  const handleAdd = () => {
    uploadBtnRef.current?.click();
  };

  const handleRemove = (idx: number) => {
    setImages((imgs) => imgs.filter((_, i) => i !== idx));
  };

  const handlePreview = (img: string, idx: number) => {
    setPreviewImage(img);
    setPreviewOpen(true);
    setPreviewTitle(`Image ${idx + 1}`);
  };

  const displayImages = images.slice(0, MAX_GRID_IMAGES);
  const extraCount = images.length - MAX_GRID_IMAGES;

  return (
    <>
      <HeadingContent>
        <Heading>Asset Images</Heading>

        <SecondaryButton
          onClick={handleAdd}
          buttonStyle={{
            width: "150px",
            display: "flex",
            alignItems: "center",
            gap: "2px",
            marginRight: "20px",
            marginBottom: "16px",
          }}
        >
          <FaPlus />
          Add Image
        </SecondaryButton>

        <Upload
          multiple
          accept="image/*"
          showUploadList={false}
          beforeUpload={() => false}
          onChange={handleAntdChange}
        >
          <button ref={uploadBtnRef} style={{ display: "none" }} />
        </Upload>
      </HeadingContent>
      <Wrapper>
        <HeaderContent>
          <HeaderText>
            <span> {images.length} </span>
            {images.length > 1 ? "images uploaded" : "image uploaded"}{" "}
          </HeaderText>
          <BiExpand
            onClick={() => router.push("/asset-images")}
            size={20}
            style={{ cursor: "pointer" }}
          />
        </HeaderContent>

        <Grid>
          {displayImages.map((img, idx) => (
            <ThumbBox key={idx}>
              <Img
                src={img}
                alt={`Asset ${idx}`}
                onClick={() => handlePreview(img, idx)}
                style={{ cursor: "pointer" }}
              />
              <RemoveBtn onClick={() => handleRemove(idx)} title="Remove">
                <FiX />
              </RemoveBtn>
              {extraCount > 0 && idx === MAX_GRID_IMAGES - 1 && (
                <Overlay>+{extraCount}</Overlay>
              )}
            </ThumbBox>
          ))}
        </Grid>
      </Wrapper>

      {/* AntD Image Preview Modal */}
      <Modal
        open={previewOpen}
        footer={null}
        onCancel={() => setPreviewOpen(false)}
        title={previewTitle}
        width={800}
        bodyStyle={{
          padding: 0,
          display: "flex",
          justifyContent: "center",
          alignItems: "center",
        }}
        centered
      >
        <img
          alt={previewTitle}
          src={previewImage}
          style={{
            width: "100%",
            maxHeight: "100vh",
            objectFit: "contain",
            borderRadius: 12,
          }}
        />
      </Modal>
    </>
  );
};

export default AssetImages;

const Wrapper = styled.section`
  padding: 24px;
  border-radius: 20px;
  border: 1px solid #e0e0e0;
  display: flex;
  flex-direction: column;
  gap: 20px;
`;

const Heading = styled.h1`
  font-size: 20px;
  font-weight: 600;
  line-height: 24.38px;
  color: #00225a;
`;

const HeadingContent = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
`;

const HeaderContent = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
  font-size: 12px;
`;

const HeaderText = styled.p`
  font-weight: 500;
  font-size: 14px;
  color: #00225a;
  display: flex;
  align-items: center;
  gap: 2px;
`;

const Grid = styled.div`
  display: grid;
  grid-template-columns: repeat(8, 1fr);
  gap: 10px;
  width: 100%;
  overflow: hidden;
  @media (max-width: 900px) {
    grid-template-columns: repeat(4, 1fr);
  }
  @media (max-width: 500px) {
    grid-template-columns: repeat(2, 1fr);
  }
`;

const ThumbBox = styled.div`
  position: relative;
  width: 100%;
  aspect-ratio: 1/1;
  border-radius: 10px;
  overflow: hidden;
  background: #f2f4f7;
  box-shadow: 0 1px 4px 0 rgba(20, 30, 38, 0.06);
  display: flex;
  align-items: center;
  justify-content: center;
`;

const Img = styled.img`
  width: 100%;
  height: 100%;
  object-fit: cover;
`;

const RemoveBtn = styled.button`
  position: absolute;
  top: 4px;
  right: 4px;
  background: #fff;
  border: none;
  border-radius: 50%;
  color: #004988;
  width: 18px;
  height: 18px;
  font-size: 13px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  z-index: 2;
`;

const Overlay = styled.div`
  position: absolute;
  inset: 0;
  background: rgba(30, 58, 138, 0.78);
  color: #fff;
  font-weight: 700;
  font-size: 18px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 10px;
  z-index: 3;
`;
